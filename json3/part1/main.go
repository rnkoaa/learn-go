package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Dog - dog object
type Dog struct {
	ID     int
	Name   string
	Breed  string
	BornAt time.Time
}

// JSONDog -
type JSONDog struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Breed  string `json:"breed"`
	BornAt int64  `json:"born_at"`
}

// NewJSONDog -
func NewJSONDog(dog Dog) JSONDog {
	return JSONDog{
		dog.ID,
		dog.Name,
		dog.Breed,
		dog.BornAt.Unix(),
	}
}

// Dog - decode json to dog
func (jd JSONDog) Dog() Dog {
	return Dog{
		jd.ID,
		jd.Name,
		jd.Breed,
		time.Unix(jd.BornAt, 0),
	}
}

// MarshalJSON implementation
func (d Dog) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewJSONDog(d))
}

// UnmarshalJSON implementation
func (d *Dog) UnmarshalJSON(data []byte) error {
	var jd JSONDog
	if err := json.Unmarshal(data, &jd); err != nil {
		return err
	}
	*d = jd.Dog()
	return nil
}

func operate1() {
	dog := Dog{1, "bowser", "husky", time.Now()}
	b, err := json.Marshal(NewJSONDog(dog))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func operate2() {
	b := []byte(`{
		"id":1,
		"name":"bowser",
		"breed":"husky",
		"born_at":1480979203}`)
	var jsonDog JSONDog
	json.Unmarshal(b, &jsonDog)
	fmt.Println(jsonDog.Dog())
}

func operate3() {
	dog := Dog{1, "bowser", "husky", time.Now()}
	b, err := json.Marshal(dog)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	b = []byte(`{
    "id":1,
    "name":"bowser",
    "breed":"husky",
    "born_at":1480979203}`)
	dog = Dog{}
	json.Unmarshal(b, &dog)
	fmt.Println(dog)
}

func main() {
	// fmt.Println("Hello, Json 3")
	operate1()

	operate2()

	operate3()
}
