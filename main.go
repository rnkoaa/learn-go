package main

import "fmt"

// Vertex struct
type Vertex struct {
	X int
	Y int
}

func main() {
	v := Vertex{1, 2}
	p := &v
	p.X = 4
	fmt.Println(v.X)
	fmt.Println(p.X)
}
