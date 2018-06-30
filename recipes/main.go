package main

import (
	"bufio"
	"fmt"
	"os"
)

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

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return "", nil
}
func main() {
	// Read all json files in one shot
	// data, err := readRecipesFile()
	// if err != nil {
	// 	fmt.Println("Unable to Open File.", err)
	// }
	// fmt.Println(data)
	readLines()
}
