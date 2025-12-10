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

type Point struct {
	X int
	Y int
}

type Edge struct {
	A, B Point
}

func main() {
	filepath := "../testdata/day9.txt" // adjust

	verts, err := loadPoints(filepath)
	if err != nil {
		panic(err)
	}

	n := len(verts)
	fmt.Printf("N RED TILES = %d\n", n)

	// We have some polygon formed by red tiles
	edges := buildEdges(verts)

	part1 := maxRectangleAreaAllPairs(verts)
	part2 := maxRectangleAreaNoCut(verts, edges)

	fmt.Println("Part 1 (no constraints) :", part1)
	fmt.Println("Part 2 (only red/green) :", part2)
}

func loadPoints(path string) ([]Point, error) {
	var pts []Point

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
		coords := s.Split(line, ",")
		if len(coords) != 2 {
			return nil, fmt.Errorf("invalid line: %q", line)
		}
		x, err1 := strconv.Atoi(coords[0])
		y, err2 := strconv.Atoi(coords[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid integers on line: %q", line)
		}
		pts = append(pts, Point{X: x, Y: y})
	}
	return pts, scanner.Err()
}

func buildEdges(verts []Point) []Edge {
	n := len(verts)
	edges := make([]Edge, n)
	for i := 0; i < n; i++ {
		edges[i] = Edge{
			A: verts[i],
			B: verts[(i+1)%n],
		}
	}
	return edges
}

// More succinct that previous version 'brute force'
// This is all we need to solve part 1
func maxRectangleAreaAllPairs(verts []Point) int64 {
	n := len(verts)
	var best int64

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			a := verts[i]
			b := verts[j]
			w := abs(a.X-b.X) + 1
			h := abs(a.Y-b.Y) + 1
			area := int64(w) * int64(h)
			if area > best {
				best = area
			}
		}
	}
	return best
}

// Same as above, but if rectangle 'cuts' our polygon we discard it
func maxRectangleAreaNoCut(verts []Point, edges []Edge) int64 {
	n := len(verts)
	var best int64

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			p := verts[i]
			q := verts[j]

			left := min(p.X, q.X)
			right := max(p.X, q.X)
			bottom := min(p.Y, q.Y)
			top := max(p.Y, q.Y)

			w := right - left + 1
			h := top - bottom + 1
			area := int64(w) * int64(h)
			if area <= best {
				continue
			}

			if rectangleCutByPolygon(left, right, bottom, top, edges) {
				continue
			}

			best = area
		}
	}
	return best
}

// rectangleCutByPolygon returns true if ANY polygon edge passes through the
// interior of the rectangle [left,right] x [bottom,top]. Touching edges are OK.
func rectangleCutByPolygon(left, right, bottom, top int, edges []Edge) bool {
	for _, e := range edges {
		if rectEdgeIntersect(left, right, bottom, top, e) {
			return true
		}
	}
	return false
}

// rectEdgeIntersect checks whether the edges of a rectangle intersect our polygon
func rectEdgeIntersect(left, right, bottom, top int, e Edge) bool {
	lx1, lx2 := e.A.X, e.B.X
	if lx1 > lx2 {
		lx1, lx2 = lx2, lx1
	}
	ly1, ly2 := e.A.Y, e.B.Y
	if ly1 > ly2 {
		ly1, ly2 = ly2, ly1
	}

	// If the edge is completely to left/right/below/above (including just touching), it's "away".
	away :=
		max(lx1, lx2) <= left || // entirely on or left of left side
			min(lx1, lx2) >= right || // entirely on or right of right side
			max(ly1, ly2) <= bottom || // entirely on or below bottom
			min(ly1, ly2) >= top // entirely on or above top

	// If NOT away, then this edge passes through the rectangle interior.
	return !away
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
