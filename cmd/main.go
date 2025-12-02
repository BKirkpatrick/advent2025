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
			evenLength := hasEvenLength(idStr)
			if evenLength {
				invalid = isInvalid(idStr)
				if invalid {
					nInvalidID++
					invalidSum += int64(id)
					// invalidIdList = append(invalidIdList, idStr)
					fmt.Printf("INVALID ID: %v\n", idStr)
				}
			}

		}

		//fmt.Printf("%v: %v --> first:%v last:%v\n", i, idRange, firstID, lastID)
	}

	fmt.Printf("Found %v Invalid IDs\n", nInvalidID)
	fmt.Printf("SUM of Invalid IDs = %v\n", invalidSum)

	// fmt.Printf("The dial starts by pointing at %v\n", currentPosition)
	// // Go through each line and extract direction and distance of rotation
	// for _, line := range lines {
	// 	dir, dist, err := parseMessage(line)
	// 	if err != nil {
	// 		log.Printf("Error parsing message: %v\n", err)
	// 	}
	// 	// I don't just care about phase offset anymore = ^.^ = ffs
	// 	// fullTurns will record quotient of division
	// 	fullTurns := dist / 100
	// 	fullTurnCounter += fullTurns
	// 	partialTurn := dist % 100 // this is how much my position is going to change.

	// 	if dir == "R" {
	// 		currentPosition = previousPosition + partialTurn
	// 	} else if dir == "L" {
	// 		currentPosition = previousPosition - partialTurn
	// 	}

	// 	if currentPosition == 0 {
	// 		zeroCounter += 1
	// 	} else if currentPosition > 99 {
	// 		if previousPosition != 0 {
	// 			nudgeCounter += 1
	// 		}
	// 		currentPosition = currentPosition - 100
	// 	} else if currentPosition < 0 {
	// 		if previousPosition != 0 {
	// 			nudgeCounter += 1
	// 		}
	// 		currentPosition += 100
	// 	}

	// 	fmt.Printf("The dial is rotated %v to point at %v, Full turns: %v, zeroes: %v, flips: %v, nudge: %v\n", line, currentPosition, fullTurns, zeroCounter, flipCounter, nudgeCounter)
	// 	previousPosition = currentPosition
	// }

	// fmt.Printf("Arrow landed on ZERO %v times\n", zeroCounter)
	// fmt.Printf("Arrow blew past ZERO %v times\n", fullTurnCounter)
	// fmt.Printf("Arrow did a little flip over ZERO %v times\n", flipCounter)
	// fmt.Printf("Arrow nudge over ZERO %v times\n", nudgeCounter)

	// answer := zeroCounter + fullTurnCounter + flipCounter + nudgeCounter

	// fmt.Printf("The code is: %v\n", answer)

}

func getFirstLastID(idRange string) (string, string) {
	ids := s.Split(idRange, "-")
	firstID := ids[0]
	lastID := ids[1]
	return firstID, lastID
}

func hasEvenLength(id string) bool {
	var good bool
	length := len(id)
	if length == 0 {
		good = false
	} else if length%2 == 0 {
		good = true
	}
	return good
}

func isInvalid(id string) bool {
	var invalid bool
	length := len(id)
	l := id[:length/2]
	r := id[length/2:]
	if l == r {
		invalid = true
	} else {
		invalid = false
	}
	return invalid
}

// func parseMessage(msg string) (string, int, error) {
// 	dir := string(msg[0])
// 	distStr := string(msg[1:])
// 	dist, err := strconv.Atoi(distStr)
// 	if err != nil {
// 		log.Printf("Could not convert distance to int, err: %v\n", err)
// 	}
// 	return dir, dist, err
// }
