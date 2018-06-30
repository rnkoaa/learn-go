package main

import (
	"fmt"
	"sync"
)

var wg = sync.WaitGroup{}

func main() {
	wg.Add(1)

	// fmt.Println("Hello, World!")
	go func() {
		fmt.Println("Inside Go Routine")
		wg.Done()
	}()

	fmt.Println("Inside Main Go-Routine")
	wg.Wait()
}
