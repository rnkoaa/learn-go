package main

import (
	"fmt"
	"learn-go/recipes/search"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var batchSize = 1000

func main() {
	// var mainWg = sync.WaitGroup{}
	// mainWg.Add(1)
	defer timeTrack(time.Now(), "Main")

	// count, _ := LineCount("data/recipeitems-latest.json")
	// fmt.Println(count)
	// db, _ := OpenDB()
	// defer db.Close()
	// defer SetupCloseHandler(db)
	// var recipeService = NewSqliteRecipeService(db)

	// recipeDto, _ := recipeService.FindRecipe("5160756b96cc62079cc2db15")
	// fmt.Println(recipeDto.Name)

	// recipeDtos, _ := recipeService.FindAllByRecipeName("Drop Biscuits and Sausage Gravy")
	// if len(recipeDtos) > 0 {
	// 	for idx, dto := range recipeDtos {
	// 		fmt.Printf("%d -> %s, %s\n", idx, dto.ID, dto.Name)
	// 	}
	// }
	// fmt.Println("Length: ", len(recipeDtos))

	// recipeDtos, _ = recipeService.FindAllRecipes()
	// fmt.Println("-------------------------------------")
	// if len(recipeDtos) > 0 {
	// 	for idx, dto := range recipeDtos {
	// 		fmt.Printf("%d -> %s, %s\n", idx, dto.ID, dto.Name)
	// 	}
	// }

	// IndexRecipes(recipeService)
	bleveSearch := search.NewBleveSearch("/tmp/recipes.bleve")
	// go indexDbRecipes(&mainWg, *bleveSearch, recipeService)
	SearchForItems("Steamed Green", bleveSearch)
	// mainWg.Wait()
}

// SearchForItems -
func SearchForItems(searchTerm string, bleveSearch *search.BleveSearch) {
	// We are looking to an Event with some string which match with dotGo
	query := bleve.NewMatchQuery(searchTerm)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = 20
	searchResult, err := bleveSearch.Search(searchRequest)
	if err != nil {
		fmt.Println("Something went wrong", err)
	}
	fmt.Printf("Total Results: %d\n", searchResult.Total)
	fmt.Printf("Request Time :%v\n", searchResult.Took)
	fmt.Printf("Search Result Size: %d\n", searchResult.Size())
	// for i := 0; i < len(searchResult.Hits); i++ {
	// 	hit := searchResult.Hits[i]

	// 	// fmt.Printf("Hit InternalId: %v, Hit Id: %v %v\n", hit.IndexInternalID[0],
	// 	// 	hit.ID,
	// 	// 	hit.Fields)
	// 	fmt.Printf("ID: %s\n", hit.ID)
	// 	fmt.Printf("IndexInternalID: %s\n", hit.IndexInternalID)

	// }
	fmt.Println("Done")
}

// IndexRecipes - Index all recipes in the db by paging through them
func IndexRecipes(recipeService RecipeService) {
	bleveSearch := search.BleveSearch{}

	_, err := bleveSearch.OpenIndex("/tmp/recipes.bleve")
	if err != nil {
		fmt.Println("Failed To Open Index Directory.")
	}

	var recipeID = "516c733296cc62548fd302d0"
	recipeDto, _ := recipeService.FindRecipe(recipeID)
	fmt.Println(recipeDto.Name)

	bleveSearch.IndexRecipe(recipeDto)

	// We are looking to an Event with some string which match with dotGo
	query := bleve.NewMatchQuery("PT45M")
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := bleveSearch.Search(searchRequest)
	if err != nil {
		fmt.Println("Something went wrong", err)
	}
	fmt.Printf("Total Results: %d\n", searchResult.Total)
	fmt.Println(searchResult.Took)
	for i := 0; i < len(searchResult.Hits); i++ {
		hit := searchResult.Hits[i]

		// fmt.Printf("Hit InternalId: %v, Hit Id: %v %v\n", hit.IndexInternalID[0],
		// 	hit.ID,
		// 	hit.Fields)
		fmt.Printf("ID: %s\n", hit.ID)
		fmt.Printf("IndexInternalID: %s\n", hit.IndexInternalID)

	}
	fmt.Println("Done")
}

func indexDbRecipes(wg *sync.WaitGroup, bleveSearch search.BleveSearch, recipeService RecipeService) {
	count, _ := recipeService.Count()
	fmt.Printf("Total Items in Table: %d\n", count)

	totalPages := int(count/1000) + 1
	fmt.Printf("Total Pages: %d\n", totalPages)
	for idx := 0; idx < totalPages; idx++ {
		currentPage := idx
		recipeDtos, _ := recipeService.FindAllByLimitAndOffset(1000, currentPage)
		fmt.Printf("Current Page: %d Selected: %d items For Index.\n",
			currentPage+1,
			len(recipeDtos))
		_, err := bleveSearch.BatchIndexRecipe(recipeDtos)
		if err != nil {
			fmt.Printf("Error Batch Indexing at %d for count: %d\n", currentPage,
				len(recipeDtos))
		}
	}
	wg.Done()
}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupCloseHandler(db *gorm.DB) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(db *gorm.DB) {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		// DeleteFiles()
		db.Close()
		os.Exit(0)
	}(db)
}

// func batchIndex(bleveIndex bleve.Index, recipeDtos []domain.RecipeDTO) {

// }

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
