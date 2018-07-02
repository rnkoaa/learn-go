package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Id struct {
	Oid string `json:"$oid"`
}

type Timestamp struct {
	Date string `json:"$date"`
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
	RecipeYied    int64
	DatePublished int64
	PrepTime      int64
	Description   int64
	Created       Timestamp `json:"ts"`
}

func readRecipesFile() (string, error) {
	file, err := os.Open("data/recipeitems-latest.json")
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println("bytes read: ", bytesread)
	return string(buffer), nil
}

func readLines() (string, error) {
	file, err := os.Open("data/recipeitems-latest.json")
	if err != nil {
		return "", err
	}
	defer file.Close()
	// _, err := file.Stat()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return "", err
	// }

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// linesCount := make([]int64, 0)
	linesCount := 0
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		linesCount++
	}
	fmt.Printf("Total Lines: %v\n", linesCount)
	return "", nil
}

var concurrency = 1

func main() {
	// Read all json files in one shot
	// data, err := readRecipesFile()
	// if err != nil {
	// 	fmt.Println("Unable to Open File.", err)
	// }
	// fmt.Println(data)
	// readLines()

	// This channel has no buffer, so it only accepts input when something is ready
	// to take it out. This keeps the reading from getting ahead of the writers.
	workQueue := make(chan string)

	// We need to know when everyone is done so we can exit.
	complete := make(chan bool)

	// Read the lines into the work queue.
	go func() {
		file, err := os.Open("data/recipeitems-latest.json")
		if err != nil {
			log.Fatal(err)
		}

		// Close when the functin returns
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			workQueue <- scanner.Text()
		}

		// Close the channel so everyone reading from it knows we're done.
		close(workQueue)
	}()

	// Now read them all off, concurrently.
	for i := 0; i < concurrency; i++ {
		go startWorking(workQueue, complete)
	}

	// Wait for everyone to finish.
	for i := 0; i < concurrency; i++ {
		<-complete
	}

	fmt.Println("In Recipes.")
}

func startWorking(queue chan string, complete chan bool) {
	var count = 0
	for range queue {
		// Do the work with the line.
		count++
	}
	fmt.Printf("Total Count: %v\n", count)

	// Let the main process know we're done.
	complete <- true
}
