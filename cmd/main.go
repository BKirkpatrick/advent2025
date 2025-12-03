package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
)

func main() {
	path := "../testdata/day3.txt"
	var bankSlice []int
	var highestJoltage int
	var totalJoltage int = 0

	// Read in our data line by line
	batteries, err := readLines(path)
	if err != nil {
		log.Printf("Problem reading input file")
	}

	for _, bank := range batteries {
		//fmt.Printf("B: %v\n", bank)
		bankSlice, _ = batteryBankToSlice(bank)
		highestJoltage = findJuiciestBatteries(bankSlice, 12)

		totalJoltage += highestJoltage

		fmt.Printf("B: %v --> %v\n", bank, highestJoltage)
	}
	fmt.Printf("Number of banks: %v\n", len(batteries))
	fmt.Printf("TOTAL JOLTAGE = %v\n", totalJoltage)
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

func batteryBankToSlice(batteryBank string) ([]int, error) {
	var bankSlice []int
	for _, j := range batteryBank {
		j_int, err := strconv.Atoi(string(j))
		if err != nil {
			return nil, err
		}
		bankSlice = append(bankSlice, j_int)
	}
	return bankSlice, nil
}

func sliceToInt(s []int) int {
	res := 0
	op := 1
	for i := len(s) - 1; i >= 0; i-- {
		res += s[i] * op
		op *= 10
	}
	return res
}

func findJuiciestBatteries(bankSlice []int, nBatteries int) int {
	var searchingSlice []int
	var innerSlice []int
	var joltageDigits []int
	var juiciestBattery int
	remainingBatteries := nBatteries
	joltage := 0

	n := len(bankSlice)
	searchingSlice = bankSlice // initialise
	for remainingBatteries > 0 {
		// Keep getting the next best battery until I have found all that I need
		//fmt.Printf("Slice is %v long, I still have %v batteries to find, so we can search %v values\n", n, remainingBatteries, (n - remainingBatteries + 1))
		innerSlice = searchingSlice[:(n - remainingBatteries + 1)]
		//fmt.Printf("Searching: %v for juiciest battery\n", innerSlice)
		juiciestBattery = slices.Max(innerSlice)
		//fmt.Printf("Juiciest Battery: %v\n", juiciestBattery)
		joltageDigits = append(joltageDigits, juiciestBattery)
		// Now I need to find where that juiciest battery was and make a new search slice from there
		for i, j := range searchingSlice {
			if j == juiciestBattery {
				// you found the juiciest battery
				// make new searching slice
				searchingSlice = searchingSlice[i+1:]
				break
			}
		}
		// now I found the best bettery and added it to digits recorder
		// we now have to look for the next digit in another slice but we dont have its length yet
		n = len(searchingSlice)
		remainingBatteries-- // decrement number of remaining batteries
	}
	// We should now have the largest joltage, as a slice
	joltage = sliceToInt(joltageDigits)
	return joltage
}
