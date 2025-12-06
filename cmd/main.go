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
	path := "../testdata/day6.txt"
	var mathsSheet [][]string
	//var operators []string

	// Read in our data line by line
	dat, err := readLines(path)
	if err != nil {
		log.Printf("Problem reading input file")
	}

	mathsSheet = buildMathsStr(dat)

	ans := doMyHomework(mathsSheet)

	//print result
	// for _, row := range mathsSheet {
	// 	fmt.Printf("%v\n", row)
	// }
	// fmt.Printf("%v\n", len(mathsSheet))

	fmt.Printf("ANSWER %v\n", ans)

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

func buildMathsStr(dat []string) [][]string {
	var mathsSheet [][]string
	for _, line := range dat {
		fields := s.Fields(line) // splits on any whitespace, skips empties
		if len(fields) == 0 {
			continue
		}
		mathsSheet = append(mathsSheet, fields)
	}
	return mathsSheet
}

func doMathsOperation(mathSheetCol []string) int {
	ans := 0
	n := len(mathSheetCol)
	numbers := mathSheetCol[:(n - 1)]
	operator := mathSheetCol[n-1]
	fmt.Printf("numbers: %v\n", numbers)
	fmt.Printf("operator: %v\n", operator)
	if operator == "*" {
		fmt.Printf("Performing multiplication\n")
		ans = 1
		for _, j := range numbers {
			intJ, _ := strconv.Atoi(string(j))
			//fmt.Printf("Multiplying %v by %v...\n", ans, intJ)
			ans *= intJ
			//fmt.Printf("Got %v\n", ans)
		}
	} else if operator == "+" {
		fmt.Printf("Performing addition\n")
		ans = 0
		for _, j := range numbers {
			intJ, _ := strconv.Atoi(string(j))
			ans += intJ
		}
	} else {
		fmt.Printf("Unknown Operation: %v\n", operator)
	}
	//fmt.Printf("About to return answer: %v\n", ans)
	return ans
}

func column(sheet [][]string, col int) []string {
	out := make([]string, len(sheet))
	for i := range sheet {
		out[i] = sheet[i][col]
	}
	return out
}

func doMyHomework(mathsSheet [][]string) int {
	ans := 0
	nCols := len(mathsSheet[0])
	fmt.Printf("We have %v operations to do\n", nCols)

	for i := 0; i < nCols; i++ {
		col := column(mathsSheet, i)
		fmt.Printf("\nColumn: %v\n", col)
		ansCol := doMathsOperation(col)
		ans += ansCol
	}
	return ans
}
