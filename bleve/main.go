package main

import (
	"fmt"
	"os"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var bleveIdx bleve.Index

// Bleve connect or create index in persistence
func Bleve(indexPath string) (bleve.Index, error) {

	//with bleveIdx isn't set...
	if bleveIdx == nil {
		var err error

		// try to open the persistence file
		bleveIdx, err = bleve.Open(indexPath)

		//if it doesn't exists or something goes wrong
		if err != nil {
			// create a new mapping file and create a new index
			mapping := bleve.NewIndexMapping()
			bleveIdx, err = bleve.New(indexPath, mapping)

			if err != nil {
				return nil, err
			}
		}
	}
	return bleveIdx, nil

}

func idxDestroy() {
	os.RemoveAll(testIdx)
}

// Event - struct to hold an event object
type Event struct {
	ID          int
	Name        string
	Description string
	Local       string
	Website     string
	Start       time.Time
	End         time.Time
}

// Index - Index an Event
func (event *Event) Index(index bleve.Index) error {
	err := index.Index(string(event.ID), event)
	return err
}

const (
	checkMark = "\u2713"
	ballotX   = "\u2717"
	testIdx   = "test.bleve"
	dbFile    = "test.sqlite3.db"
)

func idxCreate() (bleve.Index, error) {
	idx, err := Bleve(testIdx)
	return idx, err
}

// create a SQLite3 database file, create an events table and fill with some data.
func dbCreate() (*gorm.DB, []Event) {
	db, _ := gorm.Open("sqlite3", dbFile)
	db.DropTableIfExists(&Event{})
	db.CreateTable(&Event{})
	eventList := fillDatabase(db)
	return db, eventList
}

// indexEvents add the eventList to the index
func indexEvents(idx bleve.Index, eventList []Event) {
	for _, event := range eventList {
		event.Index(idx)
	}
}

// fill the database with some data
func fillDatabase(db *gorm.DB) []Event {
	eventList := []Event{
		{1, "dotGo 2015", "The European Go conference", "Paris", "http://www.dotgo.eu/", time.Date(2015, 11, 19, 9, 0, 0, 0, time.UTC), time.Date(2015, 11, 19, 18, 30, 0, 0, time.UTC)},

		{2, "GopherCon INDIA 2016", "The Go Conference in India", "Bengaluru", "http://www.gophercon.in/", time.Date(2016, 2, 19, 0, 0, 0, 0, time.UTC), time.Date(2016, 2, 20, 23, 59, 0, 0, time.UTC)},

		{3, "GopherCon 2016", "GopherCon, It is the largest event in the world dedicated solely to the Go programming language. It's attended by the best and the brightest of the Go team and community.", "Denver", "http://gophercon.com/", time.Date(2016, 7, 11, 0, 0, 0, 0, time.UTC), time.Date(2016, 7, 13, 23, 59, 0, 0, time.UTC)},
	}

	// inserting the events
	for _, event := range eventList {
		db.Create(&event)
	}

	return eventList
}

func dbDestroy() {
	os.RemoveAll(dbFile)
}

func main() {
	idxDestroy()
	idx, err := idxCreate()
	if err != nil {
		fmt.Println("Failed to open bleve Index")
	}

	if idx == nil {
		fmt.Println("Bleve Index is null")
	}

	_, eventList := dbCreate()

	event := eventList[0]
	err = event.Index(idx)
	if err != nil {
		fmt.Println("Error Index event", err)
	}
	// fmt.Println(event)

	indexEvents(idx, eventList)
	// We are looking to an Event with some string which match with dotGo
	query := bleve.NewMatchQuery("GopherCon")
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := idx.Search(searchRequest)
	if err != nil {
		fmt.Println("Something went wrong", err)
	}
	fmt.Printf("Total Results: %d\n", searchResult.Total)
	fmt.Println(searchResult.Took)
	fmt.Println(searchResult.Hits[0].ID)
	for i := 0; i < len(searchResult.Hits); i++ {
		hit := searchResult.Hits[i]

		// fmt.Printf("Hit InternalId: %v, Hit Id: %v %v\n", hit.IndexInternalID[0],
		// 	hit.ID,
		// 	hit.Fields)
		fmt.Println(hit)
	}
	fmt.Println("Done")
}
