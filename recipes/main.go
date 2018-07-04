package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Id -
type Id struct {
	Oid string `json:"$oid"`
}

// NewId -
func NewId() Id {
	return Id{}
}

// Timestamp -
type Timestamp struct {
	Date int64 `json:"$date"`
}

// NewTimestamp
func NewTimestamp() Timestamp {
	return Timestamp{}
}

// Recipe - json object
type Recipe struct {
	RecipeID      Id `json:"_id,omitempty"`
	Name          string
	Ingredients   string
	URL           string `json:"url"`
	Image         string
	CookTime      string
	Source        string
	RecipeYied    string
	DatePublished string
	PrepTime      string
	Description   string
	Created       Timestamp `json:"ts"`
}

// ToString -
func (rp Recipe) ToString() string {
	result := ""
	result += "[Name: " + rp.Name + "]"
	return result
}

// NewRecipe -
func NewRecipe() Recipe {
	return Recipe{}
}

// FromJSON - converts recipe string to json
func FromJSON(recipeJSON string) Recipe {
	recipe := NewRecipe()
	err := json.Unmarshal([]byte(recipeJSON), &recipe)
	if err != nil {
		fmt.Printf("Failed to Process Recipe: %v\n", err)
	}
	return recipe
}

func readLines() (string, error) {
	file, err := os.Open("data/recipeitems-latest.json")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	linesCount := 0
	for scanner.Scan() {
		linesCount++
	}
	fmt.Printf("Total Lines: %v\n", linesCount)
	return "", nil
}

func readRecipeFile(ch chan<- Recipe, doneCh chan<- struct{}, wg *sync.WaitGroup) {
	file, err := os.Open("data/recipeitems-latest.json")
	if err != nil {
		fmt.Println("An error occurred opening file, ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	linesCount := 0
	for scanner.Scan() {
		ch <- FromJSON(scanner.Text())
		linesCount++
	}
	doneCh <- struct{}{}
	wg.Done()
}

var wg = sync.WaitGroup{}

func main() {
	wg.Add(2)

	var recipeCh = make(chan Recipe, 50)
	var doneCh = make(chan struct{})

	// spawn a goroutine to read all files and convert to recipe structs.
	go readRecipeFile(recipeCh, doneCh, &wg)

	// spawn another go routine to read the channel and persist this object into
	// a database, postgres
	go func(ch <-chan Recipe) {
		for {
			select {
			case entry := <-recipeCh:
				fmt.Println(entry.Name)
				break
			case <-doneCh:
				wg.Done()
				break
			}
		}
	}(recipeCh)
	wg.Wait()
}
