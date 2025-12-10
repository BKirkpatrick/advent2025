package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	s "strings"
	"time"
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

type Calculation struct {
	target        []int
	buttonVectors [][]int
	joltages      []int
}

type Result struct {
	target        []int
	buttonVectors [][]int
	buttonsPushed []int
	buttonScore   int
	joltageScore  int
}

func main() {
	filepath := "../testdata/day10_test.txt" // adjust

	data, err := loadData(filepath)
	if err != nil {
		panic(err)
	}

	n := len(data)
	fmt.Printf("N data = %v\n", n)

	calcs := prepareCalculations(data)

	for i, j := range calcs {
		if i == 0 {
			fmt.Printf("%v: %v\n", i, j)
		}
	}

	myCalc := calcs[0]

	nPresses := areWeThereYet(myCalc)

	fmt.Printf("ANS = %v\n", nPresses)

	// results, answer := doAllCalculations(calcs)

	// for i, j := range results {
	// 	fmt.Printf("%v: %v\n", i, j)
	// }

	// fmt.Printf("\nANSWER = %v\n", answer)

}

func doAllCalculations(calcs []Calculation) ([]Result, int) {
	nPushes := 0
	var results []Result
	for _, calc := range calcs {
		res := doCalculation(calc)
		myResult := res[0]
		results = append(results, myResult)
		nPushes += myResult.buttonScore
	}
	return results, nPushes
}

func vectorsEqual(v1 []int, v2 []int) bool {
	var ans bool
	for i := range len(v1) {
		if v1[i] != v2[i] {
			// vectors not equal
			ans = false
			break
		} else {
			ans = true
			continue
		}
	}
	return ans
}

func doCalculation(calc Calculation) []Result {
	var results []Result
	totalButtons := len(calc.buttonVectors)
	// fmt.Printf("\nI am about to tackle this calculation\n:%v\n", calc)
	// fmt.Printf("Calculation involves %v buttons\n", totalButtons)
	tar := calc.target
	// first you need to loop over how many buttons you want to push
	var buttonTests [][]int

	// Generate all combinations: 1 button, 2 buttons, ..., up to totalButtons
	for numButtons := 1; numButtons <= totalButtons; numButtons++ {
		combos := generateCombinations(totalButtons, numButtons)
		buttonTests = append(buttonTests, combos...)
	}

	for _, combo := range buttonTests {
		ans := make([]int, len(tar))
		for _, ind := range combo {
			ans = addWrapButtonVectors(ans, calc.buttonVectors[ind])
		}
		if vectorsEqual(ans, tar) {
			fmt.Printf("This works! %v\n", combo)
			//ind := findSliceIndex(calc.buttonVectors, combo)
			jolts := calculateJoltage(calc, tar)
			results = append(results, Result{tar, calc.buttonVectors, combo, len(combo), jolts})
			break
		}
	}
	return results
}

func areWeThereYet(calc Calculation) int {
	//var results []Result
	var fronteir [][]int
	bvs := calc.buttonVectors
	// now we are targetting the joltage
	tar := calc.joltages
	nLights := len(tar)
	start := make([]int, nLights)
	fronteir = append(fronteir, start)
	// So my idea is just to test all first steps, then all seconds from those first, etc
	// And test whether I have arrived at my destination
	nPresses := 0
	for x := 0; x < 12; x++ {
		// We press the button
		nPresses++
		var holder [][]int
		fmt.Printf("HERE IS OUR FRONTEIR - %v\n", fronteir)
		for _, pos := range fronteir {
			// Every one of these is a position in our fronteir
			// fmt.Printf("%v\n", pos)
			for _, bv := range bvs {
				// fmt.Printf("BV - %v\n", bv)
				res := addVectors(pos, bv)
				holder = append(holder, res)
				// fmt.Printf("%v\n", res)

				time.Sleep(time.Millisecond * 500)

				if vectorsEqual(res, tar) {
					// We reached our destination
					fmt.Printf("We have arrived at %v after %v button presses", res, nPresses)
					break
				}

			}

		}
		fronteir = append(fronteir, holder...)
	}
	return nPresses
}

func prepareCalculations(data []Data) []Calculation {
	var calcs []Calculation

	for _, j := range data {
		target := lightsTargetVector(j.lights)
		nDim := len(j.lights)
		bvs := buttonVectors(j.buttons, nDim)
		jolts := j.joltages
		calcs = append(calcs, Calculation{target, bvs, jolts})
	}
	return calcs
}

func calculateNPushes(data Data, buttonsPushed []int) int {
	n := len(buttonsPushed)
	return n
}

func calculateJoltage(calc Calculation, target []int) int {
	joltage := 0
	for i, j := range target {
		if j == 1 {
			joltage += calc.joltages[i]
		}
	}
	return joltage
}

func addVectors(bv1 []int, bv2 []int) []int {
	var bv3 []int
	var elem int
	for i := range len(bv1) {
		elem = bv1[i] + bv2[i]
		bv3 = append(bv3, elem)
	}
	return bv3
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

// generateCombinations returns all combinations of k buttons from n total buttons (0-indexed)
func generateCombinations(n, k int) [][]int {
	var result [][]int
	var current []int

	var backtrack func(start int)
	backtrack = func(start int) {
		if len(current) == k {
			combo := make([]int, k)
			copy(combo, current)
			result = append(result, combo)
			return
		}

		for i := start; i < n; i++ {
			current = append(current, i)
			backtrack(i + 1)
			current = current[:len(current)-1]
		}
	}

	backtrack(0)
	return result
}
