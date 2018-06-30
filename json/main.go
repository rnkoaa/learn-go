package main

import (
	"encoding/json"
	"fmt"
)

// Bird struct to hold bird
// type Bird struct {
// 	Species     string `json:"species"`
// 	Description string `json:"description"`
// }
type Bird struct {
	Species     string
	Description string
}

func main() {
	birdJson := `{"species": "pigeon",
		"description": "likes to perch on rocks"}`
	fmt.Println(birdJson)

	var bird Bird
	json.Unmarshal([]byte(birdJson), &bird)
	fmt.Printf("Bird[Species: %s, Description: %s]\n", bird.Species, bird.Description)
}
