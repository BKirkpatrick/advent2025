package main

import (
	"bufio"
	"fmt"
	"math"
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
	filepath := "../testdata/day10.txt" // adjust

	data, err := loadData(filepath)
	if err != nil {
		panic(err)
	}

	total := 0
	for _, d := range data {
		total += solveMachinePart2(d.buttons, d.joltages)
	}
	fmt.Printf("PART 2 ANSWER = %d\n", total)

}

type cand struct {
	dec []int // integer decrement vector d(S)
	k   int   // |S| (buttons pressed once in this parity-fixing step)
}

func solveMachinePart2(buttons [][]int, target []int) int {
	// Assumption: number of counters <= 64
	// If ever exceed 64, switch parity keys from uint64 to a byte-slice key.
	n := len(target)
	m := len(buttons)

	// Build per-button parity mask over counters.
	btnMask := make([]uint64, m)
	for i, btn := range buttons {
		var mask uint64
		for _, idx := range btn {
			mask ^= (1 << uint(idx)) // toggle bit for parity
		}
		btnMask[i] = mask
	}

	// Precompute all subsets S, grouped by parityKey(S).
	// For each subset, store:
	// - integer decrement vector d(S) (counts overlaps)
	// - k = popcount(S)
	buckets := make(map[uint64][]cand, 1<<minInt(m, 16))

	limit := 1 << uint(m)
	for subset := 0; subset < limit; subset++ {
		par := uint64(0)
		dec := make([]int, n)
		k := 0

		for i := 0; i < m; i++ {
			if (subset>>uint(i))&1 == 0 {
				continue
			}
			k++
			par ^= btnMask[i]
			for _, idx := range buttons[i] {
				dec[idx]++
			}
		}

		buckets[par] = append(buckets[par], cand{dec: dec, k: k})
	}

	// Memoized recursion f(b): minimum presses to reach b exactly.
	memo := make(map[string]int, 1024)

	var f func(b []int) int
	f = func(b []int) int {
		key := vecKey(b)
		if v, ok := memo[key]; ok {
			return v
		}

		// base case
		allZero := true
		for _, x := range b {
			if x != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			memo[key] = 0
			return 0
		}

		// parity key p = b mod 2 as bitmask
		var p uint64
		for i, x := range b {
			if x&1 == 1 {
				p |= (1 << uint(i))
			}
		}

		cands := buckets[p]
		if len(cands) == 0 {
			memo[key] = inf()
			return memo[key]
		}

		best := inf()
		for _, c := range cands {
			// r = b - d(S)
			// must be nonnegative; then child = r/2
			child := make([]int, n)
			ok := true
			for i := 0; i < n; i++ {
				r := b[i] - c.dec[i]
				if r < 0 {
					ok = false
					break
				}
				// Since parity matches, r should be even; integer division is safe.
				child[i] = r / 2
			}
			if !ok {
				continue
			}

			sub := f(child)
			if sub == inf() {
				continue
			}

			cost := c.k + 2*sub
			if cost < best {
				best = cost
			}
		}

		memo[key] = best
		return best
	}

	ans := f(target)
	if ans == inf() {
		panic(fmt.Sprintf("no solution for machine target=%v", target))
	}
	return ans
}

func vecKey(v []int) string {
	// Compact-ish string key for memoization.
	var bld s.Builder
	for i, x := range v {
		if i > 0 {
			bld.WriteByte(',')
		}
		bld.WriteString(strconv.Itoa(x))
	}
	return bld.String()
}

func inf() int { return math.MaxInt / 4 }

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
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
