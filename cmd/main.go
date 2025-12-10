package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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

	fmt.Println("")
	for i := range n {
		fmt.Printf("TARGET for %v = %v\n", data[i].lights, lightsTargetVector(data[i].lights))
	}

	fmt.Println("")
	for i := range n {
		fmt.Printf("BUTTON VECTORS for %v = %v\n", data[i].buttons, buttonVectors(data[i].buttons, len(data[i].lights)))
	}
	bv1 := []int{0, 1, 0}
	bv2 := []int{1, 1, 1}
	bv3 := addWrapButtonVectors(bv1, bv2)
	fmt.Printf("Wrapped Sum of %v and %v = %v\n", bv1, bv2, bv3)

}

func calculateNPushes(data Data, buttonsPushed []int) int {
	n := len(buttonsPushed)
	return n
}

func calculateJoltage(data Data, buttonsPushed []int) int {
	joltage := 0
	for _, j := range buttonsPushed {
		joltage += data.joltages[j]
	}
	return joltage
}

func addWrapButtonVectors(bv1 []int, bv2 []int) []int {
	var bv3 []int
	var elem int
	for i := range len(bv1) {
		// add components
		elem = bv1[i] + bv2[i]
		// wrap
		elem = elem % 2
		bv3 = append(bv3, elem)
	}
	return bv3
}

func buttonVectors(buttons [][]int, nDim int) [][]int {
	var vectors [][]int
	for _, button := range buttons {
		bv := buttonVector(button, nDim)
		vectors = append(vectors, bv)
	}
	return vectors
}

func buttonVector(button []int, nDim int) []int {
	var vector []int
	for i := range nDim - 2 {
		if slices.Contains(button, i) {
			vector = append(vector, 1)
		} else {
			vector = append(vector, 0)
		}
	}
	return vector
}

func lightsTargetVector(lights string) []int {
	var target []int

	for _, j := range lights {
		if j == '.' {
			target = append(target, 0)
		} else if j == '#' {
			target = append(target, 1)
		}
	}
	return target
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
