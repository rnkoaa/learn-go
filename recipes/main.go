package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
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

var wg = sync.WaitGroup{}

func main() {
	wg.Add(2)

	var logCh = make(chan string, 50)
	var doneCh = make(chan struct{})
	go func(ch chan<- string) {
		for i := 0; i < 100; i++ {
			ch <- fmt.Sprintf("%d,Hello, World", i)
		}
		// logCh <- logEntry{time.Now(), logInfo, "App is Starting"}
		// logCh <- logEntry{time.Now(), logInfo, "App is Shutting Down"}
		// time.Sleep(100 * time.Millisecond)
		doneCh <- struct{}{}
		wg.Done()
	}(logCh)
	go func(ch <-chan string) {
		for {
			select {
			case entry := <-logCh:
				// fmt.Printf("%v - [%v] %v\n", entry.time.Format("2006-01-02T15:04:05"),
				// 	entry.severity,
				// 	entry.message)
				fmt.Println(entry)
				break
			case <-doneCh:
				wg.Done()
				break
			}
		}
	}(logCh)
	wg.Wait()
}

// func main() {
// 	wg.Add(2)
// 	// Read all json files in one shot
// 	// data, err := readRecipesFile()
// 	// if err != nil {
// 	// 	fmt.Println("Unable to Open File.", err)
// 	// }
// 	// fmt.Println(data)
// 	// readLines()

// 	// This channel has no buffer, so it only accepts input when something is ready
// 	// to take it out. This keeps the reading from getting ahead of the writers.
// 	workQueue := make(chan string)

// 	// We need to know when everyone is done so we can exit.
// 	// complete := make(chan bool)

// 	// Read the lines into the work queue.
// 	go func() {
// 		file, err := os.Open("data/recipeitems-latest.json")
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		// Close when the function returns
// 		defer file.Close()

// 		scanner := bufio.NewScanner(file)

// 		for scanner.Scan() {
// 			workQueue <- scanner.Text()
// 		}

// 		// Close the channel so everyone reading from it knows we're done.
// 		// close(workQueue)
// 		wg.Done()
// 	}()

// 	// Now read them all off, concurrently.
// 	for i := 0; i < concurrency; i++ {
// 		go startWorking(workQueue)
// 	}

// 	// Wait for everyone to finish.
// 	// for i := 0; i < concurrency; i++ {
// 	// 	<-complete
// 	// }

// 	fmt.Println("In Recipes.")
// 	wg.Wait()
// }

// func startWorking(queue chan string) {
// 	var count = 0
// 	for range queue {
// 		// Do the work with the line.
// 		count++
// 	}
// 	fmt.Printf("Total Count: %v\n", count)

// 	// Let the main process know we're done.
// 	// complete <- true
// 	wg.Done()
// }

// func startWorking(queue chan string, complete chan bool) {
// 	var count = 0
// 	for range queue {
// 		// Do the work with the line.
// 		count++
// 	}
// 	fmt.Printf("Total Count: %v\n", count)

// 	// Let the main process know we're done.
// 	complete <- true
// }
