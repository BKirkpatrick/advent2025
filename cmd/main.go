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

	// Part 1: Count all paths from "you" to "out"
	pathCount := countPaths("you", "out", graph)
	fmt.Printf("\n[Part 1] Total paths from 'you' to 'out': %v\n", pathCount)

	// Part 2: Count paths from "svr" to "out" that visit both "dac" and "fft"
	memo := make(map[string]int)
	pathCount2 := countPathsWithRequiredMemo("svr", "out", graph, false, false, memo)
	fmt.Printf("[Part 2] Paths from 'svr' to 'out' visiting both 'dac' and 'fft': %v\n", pathCount2)
	fmt.Printf("Memo cache size: %v entries\n", len(memo))
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

// countPathsWithRequired counts paths from start to target that visit both required nodes
// the difference here is I need to track whether I have "seen" DAC and FFT nodes
// Only paths that have seen both are valid
// UPDATE - nah this takes waaaaay too long
func countPathsWithRequired(current, target string, graph map[string][]string, seenDac, seenFft bool) int {
	// Update flags if we're at a required node
	if current == "dac" {
		seenDac = true
	}
	if current == "fft" {
		seenFft = true
	}

	// Base case: we reached the target
	if current == target {
		// Only count this path if we've seen both required nodes
		if seenDac && seenFft {
			return 1
		}
		return 0
	}

	// Get neighbors of current node
	neighbors, exists := graph[current]
	if !exists {
		// Dead end - no path from here
		return 0
	}

	// Sum up paths from all neighbors, passing along our flags
	totalPaths := 0
	for _, neighbor := range neighbors {
		totalPaths += countPathsWithRequired(neighbor, target, graph, seenDac, seenFft)
	}

	return totalPaths
}

// countPathsWithRequiredMemo is the memoised version for performance
// What is the trick here?
// If we reach the same node with the same (seenDac, seenFft) state,
// the answer will be the same - so cache it and move on with our life.
func countPathsWithRequiredMemo(current, target string, graph map[string][]string, seenDac, seenFft bool, memo map[string]int) int {
	// Update flags if we're at a required node
	if current == "dac" {
		seenDac = true
	}
	if current == "fft" {
		seenFft = true
	}

	// Create cache key from state: (node, seenDac, seenFft)
	cacheKey := fmt.Sprintf("%s,%t,%t", current, seenDac, seenFft)

	// Check if we've already computed this
	if cached, found := memo[cacheKey]; found {
		return cached
	}

	// Base case: we reached the target
	if current == target {
		// Only count this path if we've seen both required nodes
		if seenDac && seenFft {
			return 1
		}
		return 0
	}

	// Get neighbors of current node
	neighbors, exists := graph[current]
	if !exists {
		// Dead end - no path from here
		memo[cacheKey] = 0
		return 0
	}

	// Sum up paths from all neighbors, passing along our flags
	totalPaths := 0
	for _, neighbor := range neighbors {
		totalPaths += countPathsWithRequiredMemo(neighbor, target, graph, seenDac, seenFft, memo)
	}

	// Cache the result before returning
	memo[cacheKey] = totalPaths
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
