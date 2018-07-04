package main

import (
	"fmt"
	"learn-go/libraries/stringutil"
)

func main() {
	name := "Richard Amoako Agyei"

	name = stringutil.Reverse(name)
	fmt.Println(name)
}
