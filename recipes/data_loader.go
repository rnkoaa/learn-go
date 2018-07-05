package main

import (
	"bufio"
	"fmt"
	"learn-go/recipes/domain"
	"learn-go/recipes/repo"
	"log"
	"os"
	"sync"

	"github.com/jinzhu/gorm"
)

var wg = sync.WaitGroup{}

func OpenDB() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "/tmp/recipes.db")
	if err != nil {
		log.Printf("Error opening file")
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	return db, err
}

func Load() {
	db, _ := OpenDB()
	defer db.Close()
	recipeRepository := repo.NewSqliteRecipeRepository(db)
	recipeRepository.CreateTableIFNotExists()

	wg.Add(2)

	var recipeCh = make(chan domain.Recipe, 50)
	var doneCh = make(chan struct{})

	// spawn a goroutine to read all files and convert to recipe structs.
	go readRecipeFile(recipeCh, doneCh, &wg)

	// spawn another go routine to read the channel and persist this object into
	// a database, postgres
	go processJSONLines(recipeCh, doneCh, recipeRepository, &wg, db)

	wg.Wait()
}

func LineCount(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	linesCount := 0
	for scanner.Scan() {
		linesCount++
	}
	fmt.Printf("Total Lines: %v\n", linesCount)
	return linesCount, nil
}

func readRecipeFile(ch chan<- domain.Recipe, doneCh chan<- struct{}, wg *sync.WaitGroup) {
	file, err := os.Open("data/recipeitems-latest.json")
	if err != nil {
		fmt.Println("An error occurred opening file, ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	linesCount := 0
	for scanner.Scan() {
		ch <- domain.FromJSON(scanner.Text())
		linesCount++
	}
	doneCh <- struct{}{}
	wg.Done()
}

func processJSONLines(recipeCh <-chan domain.Recipe,
	doneCh <-chan struct{},
	recipeRepository repo.RecipeRepository,
	wg *sync.WaitGroup,
	db *gorm.DB) {
	recipeDtos := []domain.RecipeDTO{}
	for {
		select {
		case entry := <-recipeCh:
			recipeDto := domain.FromRecipe(entry)
			recipeDtos = append(recipeDtos, recipeDto)
			if len(recipeDtos) == 1000 {
				log.Printf("Got 1000, recipes to insert. working...")
				recipeRepository.SaveAll(recipeDtos)
				recipeDtos = recipeDtos[:0]
			}
			break
		case <-doneCh:
			log.Printf("Got %d items to insert for final", len(recipeDtos))
			recipeRepository.SaveAll(recipeDtos)
			recipeDtos = recipeDtos[:0]
			wg.Done()
			break
		}
	}
}
