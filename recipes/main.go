package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	defer timeTrack(time.Now(), "Main")

	// count, _ := LineCount("data/recipeitems-latest.json")
	// fmt.Println(count)
	db, _ := OpenDB()
	defer db.Close()
	var recipeService = NewSqliteRecipeService(db)

	count, _ := recipeService.Count()
	fmt.Printf("Total Items in Table: %d\n", count)

	recipeDto, _ := recipeService.FindRecipe("5160756b96cc62079cc2db15")
	fmt.Println(recipeDto.Name)

	recipeDtos, _ := recipeService.FindAllByRecipeName("Drop Biscuits and Sausage Gravy")
	if len(recipeDtos) > 0 {
		for idx, dto := range recipeDtos {
			fmt.Printf("%d -> %s, %s\n", idx, dto.ID, dto.Name)
		}
	}
	fmt.Println("Length: ", len(recipeDtos))
}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		// DeleteFiles()
		os.Exit(0)
	}()
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
