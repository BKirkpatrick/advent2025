package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	path := "../testdata/day1_test.txt"
	var currentPosition int = 50 // initial condition
	var zeroCounter int = 0

	// Read in our data line by line
	lines, err := readLines(path)
	if err != nil {
		log.Printf("Problem reading input file")
	}

	fmt.Printf("We have %v entries\n", len(lines))

	fmt.Printf("The dial starts by pointing at %v\n", currentPosition)

	// Go through each line and extract direction and distance of rotation
	for _, line := range lines {
		dir, dist, err := parseMessage(line)
		if err != nil {
			log.Printf("Error parsing message: %v\n", err)
		}
		// I only care about phase offset, so modulo 100
		dist = dist % 100
		if dir == "R" {
			currentPosition += dist
		} else if dir == "L" {
			currentPosition += (100 - dist)
		}
		// Make sure current position also stays <= 100
		currentPosition = currentPosition % 100
		fmt.Printf("The dial is rotated %v to point at %v\n", line, currentPosition)
		// Keep count of the zeros
		if currentPosition == 0 {
			zeroCounter += 1
		}
	}
	fmt.Printf("Arrow landed on ZERO %v times\n", zeroCounter)

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
