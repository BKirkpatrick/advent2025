package main

import (
	"fmt"
	"log"
	"os"
	s "strings"
)

var splitCounter int

func main() {
	input, _ := os.ReadFile("../testdata/day7.txt")
	splitCounter = 0

	rows := s.Split(s.TrimRight(string(input), "\n"), "\n")
	h := len(rows)
	w := len(rows[0])

	fmt.Printf("We have %v rows, %v columns\n\n", h, w)

	for i, row := range rows {
		fmt.Printf("%v: %v\n", i, row)
	}

	fmt.Println("")

	timeLines := countTimelines(rows, w)
	nTimeLines := 0

	for _, v := range timeLines {
		nTimeLines += v
	}

	fmt.Printf("\nTimelines: %v\n", nTimeLines)

}

func countTimelines(finalField []string, w int) []int {
	timeLineCounter := make([]int, w)
	for i, row := range finalField {
		b := []byte(row)
		if i == 0 {
			//special case
			for j, char := range b {
				if char == 'S' {
					timeLineCounter[j] = 1
				}
			}
		} else {
			for j, char := range b {
				if char == '^' {
					timeLineCounter[j-1] += timeLineCounter[j]
					timeLineCounter[j+1] += timeLineCounter[j]
					timeLineCounter[j] = 0
				}
			}
		}
		fmt.Printf("%v\n", timeLineCounter)
	}
	return timeLineCounter
}

// advanceTime uses row1 to update row2 to row2Updated
func advanceTime(row1 string, row2 string, w int) string {
	row2Updated := s.Repeat(".", w)
	b := []byte(row2Updated)
	for i, char := range row1 {
		if char == '.' {
			// row 2 should be unchanged
			continue
		} else if char == '^' {
			// row 2 should be unchanged - shelter behind splitter
			continue
		} else if char == '|' {
			// check what we are about to hit
			if row2[i] == '.' {
				// Here the beam hits empty space and continues
				b[i] = '|'
			} else if row2[i] == '^' {
				// Beam has hit splitter
				splitCounter++
				// Splitter stays
				b[i] = '^'
				// and positions adjacent to splitter become beams
				b[i-1] = '|'
				b[i+1] = '|'
			} else {
				log.Fatalf("Unknown character: %v\n", char)
			}
		}
	}
	row2Updated = string(b)
	return row2Updated
}

func initiateManifold(startRow string, w int) string {
	betterStartRow := s.Repeat(".", w)
	b := []byte(betterStartRow)
	for i, char := range startRow {
		if char == '.' {
			// do nothing
			continue
		} else if char == 'S' {
			// this is a tachyon source
			b[i] = '|'
		} else {
			log.Fatalf("Unknow character: %v\n", char)
		}
	}
	betterStartRow = string(b)
	return betterStartRow
}
