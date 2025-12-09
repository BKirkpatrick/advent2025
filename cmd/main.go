package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	s "strings"
)

type RedTile struct {
	X  float64 // x coord
	Y  float64 // y coord
	TL float64 // 'Top left' score
	TR float64 // 'Top right' score
	BL float64 // 'Bottom left' score
	BR float64 // 'Bottom right' score
}

type Tile struct {
	X float64
	Y float64
	C string
}

func main() {
	filepath := "../testdata/day9_test.txt"

	redTiles, _ := loadTiles(filepath)
	nRedTiles := len(redTiles)

	fmt.Printf("\nN RED TILES = %v\n\n", nRedTiles)

	polygon := buildPolygon(redTiles)

	for i, tile := range polygon {
		fmt.Printf("TILE %v = %v\n", i, tile)
	}

}

func buildPolygon(redTiles []Tile) []Tile {
	var vertices []Tile
	nTiles := len(redTiles)
	for i := range nTiles - 1 {
		p1 := redTiles[i]
		p2 := redTiles[i+1]
		line := connect2Points(p1, p2)
		vertices = append(vertices, line...)
	}
	// Now we have to add the one where it loops back
	p1 := redTiles[nTiles-1]
	p2 := redTiles[0]
	line := connect2Points(p1, p2)
	vertices = append(vertices, line...)

	return vertices
}

func connect2Points(p1 Tile, p2 Tile) []Tile {
	// We move from p1 --> p2
	var vertices []Tile
	if p1.X == p2.X {
		// line must be vertical
		if p1.Y > p2.Y {
			// we are moving upwards
			for i := p1.Y; i > p2.Y-1; i-- {
				vertices = append(vertices, Tile{p1.X, i, "green"})
			}
		} else {
			// we are moving downwards
			for i := p1.Y; i < p2.Y+1; i++ {
				vertices = append(vertices, Tile{p1.X, i, "green"})
			}
		}
	} else {
		// line must be horizontal
		if p1.X > p2.X {
			// we are moving left
			for i := p1.X; i > p2.X-1; i-- {
				vertices = append(vertices, Tile{i, p1.Y, "green"})
			}
		} else {
			// we are moving right
			for i := p1.X; i < p2.X+1; i++ {
				vertices = append(vertices, Tile{i, p1.Y, "green"})
			}
		}
	}
	vertices[0].C = "red"
	n := len(vertices)
	vertices = vertices[:n-1]
	return vertices
}

func buildEmptyGrid(h int, w int) []string {
	var grid []string
	var i int
	row := s.Repeat(".", w)
	for i = range h {
		grid = append(grid, row)
	}
	fmt.Printf("Loaded %v rows\n", i)

	return grid
}

func bruteForce(tiles []RedTile) float64 {
	var biggestArea float64
	var w float64
	var h float64
	nTiles := len(tiles)
	for i := range nTiles {
		x1 := tiles[i].X
		y1 := tiles[i].Y
		for j := range nTiles - i - 1 {
			x2 := tiles[i+1+j].X
			y2 := tiles[i+1+j].Y
			if x2 > x1 {
				w = x2 - x1 + 1
			} else {
				w = x1 - x2 + 1
			}
			if y2 > y1 {
				h = y2 - y1 + 1
			} else {
				h = y1 - y2 + 1
			}
			area := w * h
			if area > biggestArea {
				biggestArea = area
			}
		}
	}
	return biggestArea
}

func loadTiles(path string) ([]Tile, error) {
	var vertices []Tile

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
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

		vertices = append(vertices, Tile{float64(x), float64(y), "red"})
	}
	return vertices, scanner.Err()

}

func loadCoordsFile(path string) ([]RedTile, map[string]RedTile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var tiles []RedTile
	bestMap := make(map[string]RedTile)
	bestTL := 0
	bestTR := 0
	bestBL := 0
	bestBR := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := s.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		coords := s.Split(line, ",")

		if len(coords) != 2 {
			return nil, nil, fmt.Errorf("invalid line: %q", line)
		}

		x, err1 := strconv.Atoi(coords[0])
		y, err2 := strconv.Atoi(coords[1])

		if err1 != nil || err2 != nil {
			return nil, nil, fmt.Errorf("invalid integers on line: %q", line)
		}

		// Scores - higher the better
		tl := 1000000000000 - (x + y) // hilarious hack
		tr := x - y
		bl := y - x
		br := x + y

		redTile := RedTile{float64(x), float64(y), float64(tl), float64(tr), float64(bl), float64(br)}

		if tl > bestTL {
			bestTL = tl
			bestMap["TL"] = redTile
		}
		if tr > bestTR {
			bestTR = tr
			bestMap["TR"] = redTile
		}
		if bl > bestBL {
			bestBL = bl
			bestMap["BL"] = redTile
		}
		if br > bestBR {
			bestBR = br
			bestMap["BR"] = redTile
		}

		tiles = append(tiles, redTile)
	}

	return tiles, bestMap, scanner.Err()
}
