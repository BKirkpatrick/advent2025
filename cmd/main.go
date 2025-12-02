package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	path := "../testdata/day1.txt"
	var previousPosition int = 50
	var currentPosition int = 50 // initial condition
	var zeroCounter int = 0
	var fullTurnCounter int = 0
	var flipCounter int = 0
	var nudgeCounter int = 0

	// Read in our data line by line
	lines, err := readLines(path)
	if err != nil {
		log.Printf("Problem reading input file")
	}

	fmt.Printf("The dial starts by pointing at %v\n", currentPosition)

	// Go through each line and extract direction and distance of rotation
	for _, line := range lines {
		dir, dist, err := parseMessage(line)
		if err != nil {
			log.Printf("Error parsing message: %v\n", err)
		}
		// I don't just care about phase offset anymore = ^.^ = ffs
		// fullTurns will record quotient of division
		fullTurns := dist / 100
		fullTurnCounter += fullTurns
		partialTurn := dist % 100 // this is how much my position is going to change.

		if dir == "R" {
			currentPosition = previousPosition + partialTurn
		} else if dir == "L" {
			currentPosition = previousPosition - partialTurn
		}

		if currentPosition == 0 {
			zeroCounter += 1
		} else if currentPosition > 99 {
			if previousPosition != 0 {
				nudgeCounter += 1
			}
			currentPosition = currentPosition - 100
		} else if currentPosition < 0 {
			if previousPosition != 0 {
				nudgeCounter += 1
			}
			currentPosition += 100
		}

		fmt.Printf("The dial is rotated %v to point at %v, Full turns: %v, zeroes: %v, flips: %v, nudge: %v\n", line, currentPosition, fullTurns, zeroCounter, flipCounter, nudgeCounter)
		previousPosition = currentPosition
	}

	fmt.Printf("We have %v entries\n", len(lines))

	fmt.Printf("Arrow landed on ZERO %v times\n", zeroCounter)
	fmt.Printf("Arrow blew past ZERO %v times\n", fullTurnCounter)
	fmt.Printf("Arrow did a little flip over ZERO %v times\n", flipCounter)
	fmt.Printf("Arrow nudge over ZERO %v times\n", nudgeCounter)

	answer := zeroCounter + fullTurnCounter + flipCounter + nudgeCounter

	fmt.Printf("The code is: %v\n", answer)

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

func parseMessage(msg string) (string, int, error) {
	dir := string(msg[0])
	distStr := string(msg[1:])
	dist, err := strconv.Atoi(distStr)
	if err != nil {
		log.Printf("Could not convert distance to int, err: %v\n", err)
	}
	return dir, dist, err
}
