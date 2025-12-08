package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	s "strings"
)

type Vec3 struct {
	X, Y, Z float64
}

func main() {
	filepath := "../testdata/day8_test.txt"
	input, _ := os.ReadFile(filepath)

	rows := s.Split(s.TrimRight(string(input), "\n"), "\n")
	h := len(rows)
	w := len(rows[0])

	coords, _ := loadVec3File(filepath)

	fmt.Printf("We have %v rows, %v columns\n\n", h, w)

	for i, row := range rows {
		fmt.Printf("%v: %v\n", i, row)
	}

	fmt.Println("")

	for i, row := range coords {
		fmt.Printf("%v: %v\n", i, row)
	}

}

func loadVec3File(path string) ([]Vec3, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var points []Vec3

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := s.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		coords := s.Split(line, ",")

		if len(coords) != 3 {
			return nil, fmt.Errorf("invalid line: %q", line)
		}

		x, err1 := strconv.Atoi(coords[0])
		y, err2 := strconv.Atoi(coords[1])
		z, err3 := strconv.Atoi(coords[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return nil, fmt.Errorf("invalid integers on line: %q", line)
		}

		points = append(points, Vec3{float64(x), float64(y), float64(z)})
	}

	return points, scanner.Err()
}
