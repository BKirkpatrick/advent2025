package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	s "strings"
)

func main() {
	path := "../testdata/day5.txt"
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

	howBadIsThis(freshIDRanges, items)

	freshItems := getMyFreshItems(items, freshIDRanges)
	nFreshItems := len(freshItems)

	// Answer
	fmt.Printf("\nANSWER: %v\n", nFreshItems)

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

func getBoundsFromRange(idRange string) (int, int) {
	limits := s.Split(idRange, "-")
	start, _ := strconv.Atoi(limits[0])
	end, _ := strconv.Atoi(limits[1])
	return start, end

}

func howBadIsThis(idRanges []string, items []string) {
	nItems := len(items)
	nIDRanges := len(idRanges)
	nIDs := 0
	for _, idRange := range idRanges {
		start, end := getBoundsFromRange(idRange)
		nIDs += (end - start)
	}
	fmt.Printf("We have to check %v items against %v fresh IDS split across %v ID Ranges\n", nItems, nIDs, nIDRanges)
}

func isInRanges(id int, idRanges []string) bool {
	for _, idRange := range idRanges {
		start, end := getBoundsFromRange(idRange)
		if id >= start && id <= end {
			return true
		}
	}
	return false
}

func getMyFreshItems(myItems []string, idRanges []string) []string {
	var myFreshItems []string
	for _, item := range myItems {
		id, _ := strconv.Atoi(item)
		if isInRanges(id, idRanges) {
			myFreshItems = append(myFreshItems, item)
		}
	}
	return myFreshItems
}
