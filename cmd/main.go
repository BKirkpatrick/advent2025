package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	s "strings"
)

var splitCounter int
var timeLineCounter int

func main() {
	input, _ := os.ReadFile("../testdata/day7_test.txt")
	splitCounter = 0
	timeLineCounter = 0

	rows := s.Split(s.TrimRight(string(input), "\n"), "\n")
	h := len(rows)
	w := len(rows[0])

	fmt.Printf("We have %v rows, %v columns\n\n", h, w)

	for i, row := range rows {
		fmt.Printf("%v: %v\n", i, row)
	}

	fmt.Println("")

	var finalField []string
	var thisRow string
	for i := range h - 1 {
		if i == 0 {
			finalField = append(finalField, rows[0])
			thisRow = initiateManifold(rows[0], w)
		}
		nextRow := advanceTime(thisRow, rows[i+1], w)
		finalField = append(finalField, nextRow)
		thisRow = nextRow
	}

	for i, row := range finalField {
		fmt.Printf("%v: %v\n", i, row)
	}

	fmt.Printf("\nTotal Splits: %v\n", splitCounter)
	fmt.Printf("\nTimelines: %v\n", timeLineCounter)

	// var nums [][]int
	// for _, row := range rows[:h-1] {
	// 	fs := strings.Fields(row)
	// 	arr := make([]int, len(fs))
	// 	for i, f := range fs {
	// 		n, _ := strconv.Atoi(f)
	// 		arr[i] = n
	// 	}
	// 	nums = append(nums, arr)
	// }

	// fmt.Println("part 1:", c(nums, ops))
	// fmt.Println("part 2:", c2(rows[:h-1], rows[h-1]))
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
				timeLineCounter += 2
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

func c(nums [][]int, ops []rune) int {
	sum := 0
	for p := 0; p < len(nums[0]); p++ {
		acc := nums[0][p]
		for r := 1; r < len(nums); r++ {
			if ops[p] == '+' {
				acc += nums[r][p]
			} else {
				acc *= nums[r][p]
			}
		}
		sum += acc
	}
	return sum
}

func c2(lines []string, opLine string) int64 {
	var opPos []int
	for i, c := range opLine {
		if c == '+' || c == '*' {
			opPos = append(opPos, i)
		}
	}
	result := int64(0)

	for i := range opPos {
		start := opPos[i]
		end := len(opLine)
		if i+1 < len(opPos) {
			end = opPos[i+1]
		}

		var colText []string
		for _, line := range lines {
			if start < len(line) {
				if end > len(line) {
					end = len(line)
				}
				colText = append(colText, line[start:end])
			} else {
				colText = append(colText, "")
			}
		}

		nums := extractVertical(colText)
		op := rune(opLine[start])

		result += applyOp(nums, op)
	}

	return result
}

func extractVertical(lines []string) []int {
	maxH := len(lines)
	var out []int

	for col := 0; col < len(lines[0]); col++ {
		var b strings.Builder
		for row := range maxH {
			if col < len(lines[row]) && lines[row][col] != ' ' {
				b.WriteByte(lines[row][col])
			}
		}
		if b.Len() > 0 {
			n, _ := strconv.Atoi(b.String())
			out = append(out, n)
		}
	}
	return out
}

func applyOp(nums []int, op rune) int64 {
	acc := int64(nums[0])
	for _, n := range nums[1:] {
		if op == '+' {
			acc += int64(n)
		} else {
			acc *= int64(n)
		}
	}
	return acc
}
