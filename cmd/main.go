package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	path := "../testdata/day5_test.txt"
	var freshIDRanges []string
	var items []string

	// Read in our data line by line
	dat, err := readLines(path)
	if err != nil {
		log.Printf("Problem reading input file")
	}

	listTracker := 0
	for _, line := range dat {
		if line == "" {
			listTracker = 1
		}
		if listTracker == 0 {
			freshIDRanges = append(freshIDRanges, line)
		} else {
			items = append(items, line)
		}
	}

	fmt.Println("Fresh ID Ranges:")
	for _, id := range freshIDRanges {
		fmt.Printf("%v\n", id)
	}
	fmt.Printf("\nItems:")
	for _, item := range items {
		fmt.Printf("%v\n", item)
	}

}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
