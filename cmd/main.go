package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	s "strings"
)

func main() {
	path := "../testdata/day2.txt"
	var idList []string
	var invalid bool
	var nInvalidID int
	var invalidSum int64 = 0
	//var invalidIdList []string

	// Read in our data line by line
	dat, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Problem reading input file")
	}
	idRanges := string(dat)

	// Split that single string up, splitting at ","
	idList = s.Split(idRanges, ",")

	for _, idRange := range idList {
		fmt.Printf("ID RANGE: %v\n", idRange)
		firstID, lastID := getFirstLastID(idRange)
		startSearching, _ := strconv.Atoi(firstID)
		stopSearching, _ := strconv.Atoi(lastID)

		for id := startSearching; id <= stopSearching; id++ {
			//fmt.Printf("checking %v\n", id)
			idStr := strconv.Itoa(id)
			invalid = isInvalid(idStr)
			if invalid {
				nInvalidID++
				invalidSum += int64(id)
				// invalidIdList = append(invalidIdList, idStr)
				fmt.Printf("INVALID ID: %v\n", idStr)
			}
		}

		//fmt.Printf("%v: %v --> first:%v last:%v\n", i, idRange, firstID, lastID)
	}

	fmt.Printf("Found %v Invalid IDs\n", nInvalidID)
	fmt.Printf("SUM of Invalid IDs = %v\n", invalidSum)

}

func getFirstLastID(idRange string) (string, string) {
	ids := s.Split(idRange, "-")
	firstID := ids[0]
	lastID := ids[1]
	return firstID, lastID
}

func isInvalid(id string) bool {
	length := len(id)

	// Try all possible size of substring
	for size := 1; size <= length/2; size++ {
		if length%size != 0 {
			continue
		}
		sub := id[:size]
		match := true

		for i := size; i < length; i += size {
			if id[i:i+size] != sub {
				match = false
				break
			}
		}
		// Repeating pattern found; invalid ID
		if match {
			return true
		}
	}
	// No repeating pattern found → ID is valid
	return false
}
