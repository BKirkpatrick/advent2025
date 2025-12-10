package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	s "strings"
)

// Day 10
// For a display with N buttons we are looking for a point in N-dimensional space.
// You are either 0 or 1 distance along each of the N dimensions.
// If I transform the button circuits into N-dimensional vectors, can we look through
// the components of those vectors to see what combo (+ and -) give desired output?

// Data structure will have target vector (lights)
// and a list of 'basis' vectors (the button circuits)
// I should also store the joltage requirements - they are going to be needed for part 2...

type Data struct {
	lights   string
	buttons  [][]int
	joltages []int
}

func main() {
	filepath := "../testdata/day10_test.txt" // adjust

	data, err := loadData(filepath)
	if err != nil {
		panic(err)
	}

	n := len(data)
	fmt.Printf("N data = %v\n", n)

	for i, j := range data {
		fmt.Printf("Data[%v]: %v\n", i, j)
	}

}

func loadData(path string) ([]Data, error) {
	var data []Data

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := s.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// we have a valid line
		elems := s.Split(line, " ")
		n := len(elems)

		// sort out lights
		lights := elems[0]

		// sort out joltages
		var joltageList []int
		joltages := elems[n-1]

		joltages = s.Replace(joltages, "{", "", 1)
		joltages = s.Replace(joltages, "}", "", 1)

		joltagesSplit := s.Split(joltages, ",")

		for _, j := range joltagesSplit {
			jInt, _ := strconv.Atoi(j)
			joltageList = append(joltageList, jInt)
		}

		// sort out buttons
		buttons := elems[1 : n-1]

		var buttonList [][]int
		for _, j := range buttons {
			var button []int

			j = s.Replace(j, "(", "", 1)
			j = s.Replace(j, ")", "", 1)
			buttonElem := s.Split(j, ",")

			for _, q := range buttonElem {
				qInt, _ := strconv.Atoi(q)
				button = append(button, qInt)
			}
			buttonList = append(buttonList, button)
		}

		data = append(data, Data{lights: lights, buttons: buttonList, joltages: joltageList})
	}
	return data, scanner.Err()
}

// stupid wee helper functions that Go doesn't have
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
