package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Tentacle a character from Day of Tentacles
type Tentacle struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (t Tentacle) toString() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(bytes)
}

func getTentacles() []Tentacle {
	rawString := `[{"name":"Purple Tentacle","description":"A mutant monster and lab assistant created by mad scientist Dr. Fred Edison."},{"name":"Green Tentacle","description":"Harmless and friendly brother of Purple Tentacle."},{"name":"Bernard Bernoulli","description":"Green Tentacle's friend, he's a nerd with glasses."}]`
	var tentacles []Tentacle
	json.Unmarshal([]byte(rawString), &tentacles)
	return tentacles
}

func main() {
	tentacles := getTentacles()
	fmt.Println(tentacles)
	for _, te := range tentacles {
		fmt.Println(te.toString())
	}
}
