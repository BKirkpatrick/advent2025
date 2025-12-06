package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, _ := os.ReadFile("../testdata/day6.txt")

	rows := strings.Split(strings.TrimRight(string(input), "\n"), "\n")
	h := len(rows)

	ops := []rune{}
	for _, f := range strings.Fields(rows[h-1]) {
		ops = append(ops, rune(f[0]))
	}

	var nums [][]int
	for _, row := range rows[:h-1] {
		fs := strings.Fields(row)
		arr := make([]int, len(fs))
		for i, f := range fs {
			n, _ := strconv.Atoi(f)
			arr[i] = n
		}
		nums = append(nums, arr)
	}

	fmt.Println("part 1:", c(nums, ops))
	fmt.Println("part 2:", c2(rows[:h-1], rows[h-1]))
}

func c(nums [][]int, ops []rune) int {
	sum := 0
	for p := 0; p < len(nums[0]); p++ {
		acc := nums[0][p]
		for r := 1; r < len(nums); r++ {
			if ops[p] == '+' {
				acc += nums[r][p]
			} else {
				acc *= nums[r][p]
			}
		}
		sum += acc
	}
	return sum
}

func c2(lines []string, opLine string) int64 {
	var opPos []int
	for i, c := range opLine {
		if c == '+' || c == '*' {
			opPos = append(opPos, i)
		}
	}
	result := int64(0)

	for i := range opPos {
		start := opPos[i]
		end := len(opLine)
		if i+1 < len(opPos) {
			end = opPos[i+1]
		}

		var colText []string
		for _, line := range lines {
			if start < len(line) {
				if end > len(line) {
					end = len(line)
				}
				colText = append(colText, line[start:end])
			} else {
				colText = append(colText, "")
			}
		}

		nums := extractVertical(colText)
		op := rune(opLine[start])

		result += applyOp(nums, op)
	}

	return result
}

func extractVertical(lines []string) []int {
	maxH := len(lines)
	var out []int

	for col := 0; col < len(lines[0]); col++ {
		var b strings.Builder
		for row := range maxH {
			if col < len(lines[row]) && lines[row][col] != ' ' {
				b.WriteByte(lines[row][col])
			}
		}
		if b.Len() > 0 {
			n, _ := strconv.Atoi(b.String())
			out = append(out, n)
		}
	}
	return out
}

func applyOp(nums []int, op rune) int64 {
	acc := int64(nums[0])
	for _, n := range nums[1:] {
		if op == '+' {
			acc += int64(n)
		} else {
			acc *= int64(n)
		}
	}
	return acc
}
