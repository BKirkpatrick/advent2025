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
	path := "../testdata/day3_test.txt"
	var bankSlice []int
	var highestJoltage int
	//var totalJoltage int = 0

	// Read in our data line by line
	batteries, err := readLines(path)
	if err != nil {
		log.Printf("Problem reading input file")
	}

	fmt.Printf("Number of banks: %v\n", len(batteries))
	fmt.Printf("Battery banks:\n%v\n", batteries)

	for _, bank := range batteries {
		bankSlice, _ = batteryBankToSlice(bank)
		highestJoltage = findMaxJoltage(bankSlice, 12)
		fmt.Printf("BANK: %v; JOLT: %v\n", bankSlice, highestJoltage)

		// fmt.Printf("Bank:\n%v\n", bank)

		// highestJoltage = findHightestJoltageInBank(bankSlice)
		// fmt.Printf("Bank: %v --> Highest Joltage: %v\n", bankSlice, highestJoltage)
		// totalJoltage += highestJoltage
	}
	//fmt.Printf("TOTAL JOLTAGE = %v\n", totalJoltage)
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

func findMaxJoltage(bankSlice []int, nBatteries int) int {
	highestJoltage := 0
	joltage := 0
	n := len(bankSlice)
	// Okay so I don't care about finding big numbers if they are in the last
	// (nBatteries - 1) digits, because I cannot form an integer that is nBatteries long
	searchingSlice := bankSlice[:n-nBatteries+1]
	biggestFirstDigit := slices.Max(searchingSlice) // this is the biggest first digit
	// the only problem is there could be several of them
	// Everytime you find one of these, use its index to look up the whole bankSlice and
	// return the nBatteries-digit number it defines
	for i, j := range searchingSlice {
		if j == biggestFirstDigit {
			joltage = sliceToInt(bankSlice[i : i+nBatteries])
			if joltage > highestJoltage {
				highestJoltage = joltage
			}
		}
	}
	return highestJoltage
}

func findHightestJoltageInBank(bankSlice []int) int {
	highestNumber := 0
	joltage := 0
	n := len(bankSlice)
	// We can find the highest joltage quickly
	biggestFirstDigit := slices.Max(bankSlice)

	for m := range biggestFirstDigit {
		fmt.Printf("Searching (%v) for battery with value %v\n", m, biggestFirstDigit)
		for i, j := range bankSlice {
			//fmt.Printf("i: %v, j: %v\n", i, j)
			if j == biggestFirstDigit { // our first digit is as big as could be
				// look through remaining digits
				// first make sure that this is not the last digit in the slice
				fmt.Printf("Located largest first digit\n")
				if i == n-1 {
					// If you are here you reached the end of the list and you dont have a joltage
					// Start again with lower first digit
					fmt.Printf("Largest first digit is at the end of the list!\nWill search for next highest first digit instead...\n")
					biggestFirstDigit--
					// Bonus - this number must be the largest second digit we could find, so no need to keep searching for it!
					// biggestSecondDigit = j
					continue
				} else {
					fmt.Printf("Now looking for second digit...\n")
					biggestSecondDigit := slices.Max(bankSlice[i+1:])
					fmt.Printf("Largest second digit = %v\n", biggestSecondDigit)
					joltage = 10*j + biggestSecondDigit
					if joltage > highestNumber {
						// the joltage we just found is highest so far
						highestNumber = joltage
					}
					break
				}
			}
			continue // there could be multiple batteries that tie for highest joltage
		}
		if highestNumber >= 10 {
			// we search from highest to lowest, and if we have found one we must be done
			break
		}
	}
	return highestNumber

}
