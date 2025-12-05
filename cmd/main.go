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

	nFreshItems := countUniqueInRanges(freshIDRanges)

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

func countUniqueInRanges(idRanges []string) int {
	// Parse all ranges
	type Range struct {
		start, end int
	}
	var ranges []Range
	for _, idRange := range idRanges {
		start, end := getBoundsFromRange(idRange)
		ranges = append(ranges, Range{start, end})
	}

	// Sort ranges by start position
	// Use a simple bubble sort because I know the name of that sorting algorithm
	for i := 0; i < len(ranges); i++ {
		for j := i + 1; j < len(ranges); j++ {
			if ranges[j].start < ranges[i].start {
				ranges[i], ranges[j] = ranges[j], ranges[i]
			}
		}
	}

	// Merge overlapping ranges and count
	if len(ranges) == 0 {
		return 0
	}

	count := 0
	current := ranges[0]

	for i := 1; i < len(ranges); i++ {
		if ranges[i].start <= current.end+1 {
			// Overlapping or adjacent - we should merge this
			if ranges[i].end > current.end {
				current.end = ranges[i].end
			}
		} else {
			// No overlap - count current range and move to next
			count += current.end - current.start + 1
			current = ranges[i]
		}
	}
	// Don't forget the last range, idiot!
	count += current.end - current.start + 1

	return count
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

func findMaxID(idRanges []string) int {
	maxID := 0
	for _, idRange := range idRanges {
		_, end := getBoundsFromRange(idRange)
		if end > maxID {
			maxID = end
		}
	}
	return maxID
}

func makeList(n int) []string {
	result := make([]string, 0, n)
	for i := 1; i <= n; i++ {
		result = append(result, strconv.Itoa(i))
	}
	return result
}
