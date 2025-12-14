package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	s "strings"
)

type Point struct{ X, Y int }

type Shape struct {
	ID     int
	W, H   int
	Filled []Point // coords of '#'
}

type Region struct {
	W, H   int
	Counts []int
}

type Placement struct {
	Type    int   // shape type index
	Cells   []int // flattened indices covered in region grid
	VarName string
	UB      int // upper bound for this placement var (count for its type)
}

// Day 12 - Tetris =^.^=

func main() {
	filepath := "../testdata/day12.txt"

	shapes, regions, err := loadData(filepath)
	if err != nil {
		panic(err)
	}

	shapeArea := make([]int, len(shapes))
	for i := range shapes {
		shapeArea[i] = len(shapes[i].Filled)
	}

	sol1 := 0

regionsLoop:
	for _, r := range regions {
		xy := r.W * r.H
		s := 0

		for t, cnt := range r.Counts {
			s += cnt * shapeArea[t]
			if s > xy {
				continue regionsLoop
			}
		}
		sol1++
	}

	fmt.Println(sol1)
}

func loadData(path string) ([]Shape, []Region, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var shapes []Shape
	var regions []Region

	// Phase: keep reading shapes while we see "<int>:"
	// Once we see a line that starts with "<int>x<int>:" we switch to regions.
	inRegions := false

	for scanner.Scan() {
		line := s.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if !inRegions {
			// Region line starts the region section.
			if isRegionHeader(line) {
				inRegions = true
				// fallthrough to parse region below
			} else {
				// Otherwise it must be a shape header like "0:"
				if !s.HasSuffix(line, ":") {
					return nil, nil, fmt.Errorf("expected shape header like '0:' got %q", line)
				}
				idStr := s.TrimSuffix(line, ":")
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return nil, nil, fmt.Errorf("bad shape id %q: %w", idStr, err)
				}

				// Read shape grid rows until blank line OR next header (defensive)
				var grid []string
				for scanner.Scan() {
					row := s.TrimSpace(scanner.Text())
					if row == "" {
						break
					}
					// Defensive: if we accidentally hit a new header, stop (rare, but safe)
					if isRegionHeader(row) || s.HasSuffix(row, ":") {
						return nil, nil, fmt.Errorf("unexpected header %q while reading grid for shape %d", row, id)
					}
					grid = append(grid, row)
				}
				if err := scanner.Err(); err != nil {
					return nil, nil, err
				}
				if len(grid) == 0 {
					return nil, nil, fmt.Errorf("shape %d has no grid rows", id)
				}

				w := len(grid[0])
				h := len(grid)
				var filled []Point

				for y := 0; y < h; y++ {
					if len(grid[y]) != w {
						return nil, nil, fmt.Errorf("shape %d grid not rectangular", id)
					}
					for x := 0; x < w; x++ {
						switch grid[y][x] {
						case '#':
							filled = append(filled, Point{X: x, Y: y})
						case '.':
						default:
							return nil, nil, fmt.Errorf("shape %d invalid char %q", id, grid[y][x])
						}
					}
				}

				// Optional: enforce contiguous IDs starting at 0
				if id != len(shapes) {
					return nil, nil, fmt.Errorf("expected shape id %d but got %d", len(shapes), id)
				}

				shapes = append(shapes, Shape{ID: id, W: w, H: h, Filled: filled})
				continue
			}
		}

		// Regions section: line like "12x5: 1 0 1 0 3 2"
		region, err := parseRegionLine(line, len(shapes))
		if err != nil {
			return nil, nil, err
		}
		regions = append(regions, region)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return shapes, regions, nil
}

func isRegionHeader(line string) bool {
	// True for patterns like "12x5:" (possibly with trailing spaces already trimmed)
	// We do a simple parse attempt: "<int>x<int>:" at the start.
	colon := s.IndexByte(line, ':')
	if colon < 0 {
		return false
	}
	head := s.TrimSpace(line[:colon]) // "12x5"
	xPos := s.IndexByte(head, 'x')
	if xPos < 0 {
		return false
	}
	if xPos == 0 || xPos == len(head)-1 {
		return false
	}
	_, err1 := strconv.Atoi(head[:xPos])
	_, err2 := strconv.Atoi(head[xPos+1:])
	return err1 == nil && err2 == nil
}

func parseRegionLine(line string, numShapes int) (Region, error) {
	parts := s.Split(line, ":")
	if len(parts) != 2 {
		return Region{}, fmt.Errorf("bad region line %q", line)
	}

	dims := s.TrimSpace(parts[0])     // "12x5"
	countStr := s.TrimSpace(parts[1]) // "1 0 1 0 3 2"

	xPos := s.IndexByte(dims, 'x')
	if xPos < 0 {
		return Region{}, fmt.Errorf("bad dims %q", dims)
	}

	w, err := strconv.Atoi(dims[:xPos])
	if err != nil {
		return Region{}, fmt.Errorf("bad region width in %q", dims)
	}
	h, err := strconv.Atoi(dims[xPos+1:])
	if err != nil {
		return Region{}, fmt.Errorf("bad region height in %q", dims)
	}

	fields := s.Fields(countStr)
	if len(fields) != numShapes {
		return Region{}, fmt.Errorf("region %q has %d counts, expected %d", line, len(fields), numShapes)
	}

	counts := make([]int, numShapes)
	for i := range fields {
		n, err := strconv.Atoi(fields[i])
		if err != nil {
			return Region{}, fmt.Errorf("bad count %q in %q", fields[i], line)
		}
		counts[i] = n
	}

	return Region{W: w, H: h, Counts: counts}, nil
}
