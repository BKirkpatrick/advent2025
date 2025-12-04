package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	path := "../testdata/day4.txt"
	var paperGrid [][]int
	var nnGrid [][]int
	var mask [][]int
	var result [][]int
	var nRolls int
	threshold := 4

	// Read in our data line by line
	dat, err := readLines(path)
	if err != nil {
		log.Printf("Problem reading input file")
	}
	paperGrid = buildGridFromDat(dat)

	nnGrid = calculateNNGrid(paperGrid)

	mask = makeMask(nnGrid, threshold)

	nRolls = calculateVulnerableRolls(paperGrid, mask)

	result = calculateFinalConfiguration(paperGrid, mask)

	before := renderIntGrid(paperGrid)
	after := renderIntGrid(result)

	fmt.Println("BEFORE")
	for i, row := range before {
		fmt.Printf("%v: %v\n", i, row)
	}

	// fmt.Println("")

	// fmt.Println("MASK")
	// for i, row := range mask {
	// 	fmt.Printf("%v: %v\n", i, row)
	// }

	fmt.Println("")

	fmt.Println("AFTER")
	for i, row := range after {
		fmt.Printf("%v: %v\n", i, row)
	}

	fmt.Println("")

	fmt.Printf("ANSWER: %v\n", nRolls)

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

func symbolToInt(symbol string) int {
	var out int
	if symbol == "@" {
		out = 1
	} else if symbol == "." {
		out = 0
	} else {
		log.Fatalf("Invalid character in paper grid")
	}
	return out
}

func intToSymbol(gridInt int) string {
	var out string
	if gridInt == 1 {
		out = "@"
	} else if gridInt == 0 {
		out = "."
	} else {
		log.Fatalf("Invalid integer in paper grid")
	}
	return out
}

func buildGridFromDat(dat []string) [][]int {
	var grid [][]int
	for _, row := range dat {
		var newRow []int
		for _, pigeonHole := range row {
			newRow = append(newRow, symbolToInt(string(pigeonHole)))
		}
		grid = append(grid, newRow)
	}
	return grid
}

func calculateNNGrid(grid [][]int) [][]int {
	var nnGrid [][]int
	var nn int
	nRows := len(grid)
	nCols := len(grid[0])
	for i, row := range grid {
		var newRow []int
		for j := range row {
			// x is the value inside the "pigeonhole" at grid[i][j]
			// I
			if i == 0 && j == 0 {
				// top left corner
				nn = grid[i][j+1] + grid[i+1][j] + grid[i+1][j+1]
			} else if i == 0 && j == nCols-1 {
				// top right corner
				nn = grid[i][j-1] + grid[i+1][j] + grid[i+1][j-1]
			} else if i == nRows-1 && j == 0 {
				// bottom left column
				nn = grid[i][j+1] + grid[i-1][j] + grid[i-1][j+1]
			} else if i == nRows-1 && j == nCols-1 {
				// bottom right corner
				nn = grid[i][j-1] + grid[i-1][j] + grid[i-1][j-1]
			} else if i == 0 {
				// top row, not a corner
				nn = grid[i][j-1] + grid[i][j+1] + grid[i+1][j-1] + grid[i+1][j] + grid[i+1][j+1]
			} else if i == nRows-1 {
				// bottom row, not a corner
				nn = grid[i][j-1] + grid[i][j+1] + grid[i-1][j-1] + grid[i-1][j] + grid[i-1][j+1]
			} else if j == 0 {
				// left column, not a corner
				nn = grid[i-1][j] + grid[i+1][j] + grid[i-1][j+1] + grid[i][j+1] + grid[i+1][j+1]
			} else if j == nCols-1 {
				// right column, not a corner
				nn = grid[i-1][j] + grid[i+1][j] + grid[i-1][j-1] + grid[i][j-1] + grid[i+1][j-1]
			} else {
				// pigeonhole with 8NNs
				nn = grid[i+1][j] + grid[i+1][j+1] + grid[i][j+1] + grid[i-1][j+1] + grid[i-1][j] + grid[i-1][j-1] + grid[i][j-1] + grid[i+1][j-1]
			}
			newRow = append(newRow, nn)
		}
		nnGrid = append(nnGrid, newRow)
	}
	return nnGrid
}

func makeMask(nnGrid [][]int, threshold int) [][]int {
	var mask [][]int
	for _, row := range nnGrid {
		var newRow []int
		for _, x := range row {
			if x < threshold {
				x = 1
			} else {
				x = 0
			}
			newRow = append(newRow, x)
		}
		mask = append(mask, newRow)
	}
	return mask
}

func calculateVulnerableRolls(paperGrid [][]int, mask [][]int) int {
	nRolls := 0
	for i, row := range paperGrid {
		for j := range row {
			nRolls += (paperGrid[i][j] * mask[i][j])
		}
	}
	return nRolls
}

func calculateFinalConfiguration(paperGrid [][]int, mask [][]int) [][]int {
	var result [][]int
	var x int
	for i, row := range paperGrid {
		var newRow []int
		for j, p := range row {
			if p == 1 {
				// There is a roll of paper at this location to begin with
				// We need to see whether it can go or not
				if mask[i][j] == 1 {
					// This roll should go
					x = 0
				} else if mask[i][j] == 0 {
					// This roll should stay
					x = p
				} else {
					log.Fatalf("Invalid value in mask: %v\n", mask[i][j])
				}
			} else {
				// There was no paper at this location, so it should remain empty
				x = 0
			}
			newRow = append(newRow, x)
		}
		result = append(result, newRow)
	}
	return result
}

func renderIntGrid(grid [][]int) [][]string {
	var renderedGrid [][]string
	for _, row := range grid {
		var newRow []string
		for _, x := range row {
			newRow = append(newRow, intToSymbol(x))
		}
		renderedGrid = append(renderedGrid, newRow)
	}
	return renderedGrid
}
