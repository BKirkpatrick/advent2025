package main

import (
	"bufio"
	"fmt"
	"os"
	s "strings"
)

// Day 11 - Counting paths - OMG can I build a DAG? =^.^=

func main() {
	filepath := "../testdata/day11.txt"

	graph, err := loadGraph(filepath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded graph with %v nodes\n", len(graph))

	// Count all paths from "you" to "out"
	pathCount := countPaths("you", "out", graph)

	fmt.Printf("\nTotal paths from 'you' to 'out': %v\n", pathCount)
}

// countPaths uses DFS to count all paths from start to end
// DFS = depth first search, so I'm guessing it should work
// well for this type of structure?
func countPaths(current, target string, graph map[string][]string) int {
	// Base case: we reached the target
	if current == target {
		return 1
	}

	// Get neighbors of current node
	neighbors, exists := graph[current]
	if !exists {
		// Dead end - no path from here
		return 0
	}

	// Sum up paths from all neighbors
	totalPaths := 0
	for _, neighbor := range neighbors {
		totalPaths += countPaths(neighbor, target, graph)
	}

	return totalPaths
}

// loadGraph parses the input file and builds an adjacency list
func loadGraph(path string) (map[string][]string, error) {
	graph := make(map[string][]string)

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

		// Parse line: "device: output1 output2 output3"
		parts := s.Split(line, ": ")
		if len(parts) != 2 {
			continue
		}

		device := parts[0]
		outputs := s.Fields(parts[1]) // Split by whitespace

		graph[device] = outputs
	}

	return graph, scanner.Err()
}
