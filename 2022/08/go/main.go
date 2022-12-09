package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
	"strconv"
)

//var fileName = "../sample.txt"
var fileName = "../input.txt"

func a() {
	scanner := common.GetLineScanner(fileName)

	trees := make([][]int, 0)

	for scanner.Scan() {
		l := []rune(scanner.Text())

		row := make([]int, 0)
		for _, ts := range l {
			t, err := strconv.Atoi(string(ts))
			if err != nil {
				panic("Cannot convert tree height to int.")
			}
			row = append(row, t)
		}
		trees = append(trees, row)

	}
	//fmt.Printf("trees: %+v\n", trees)

	h := len(trees)
	w := len(trees[0])
	//fmt.Printf("w: %d, h: %d\n", w, h)

	checkIsVisible := func(x int, y int) bool {
		//fmt.Printf("Checking %d %d\n", x, y)
		tree := trees[y][x]

		// Check left
		visibleLeft := true
		for l := x - 1; l >= 0; l-- {
			//fmt.Printf("Checking left %d\n", l)
			if trees[y][l] >= tree {
				visibleLeft = false
				break
			}
		}
		if visibleLeft {
			return true
		}

		// Check top
		visibleTop := true
		for t := y - 1; t >= 0; t-- {
			//fmt.Printf("Checking top %d\n", t)
			if trees[t][x] >= tree {
				visibleTop = false
				break
			}
		}
		if visibleTop {
			return true
		}

		// Check top
		visibleRight := true
		for r := x + 1; r < w; r++ {
			if trees[y][r] >= tree {
				visibleRight = false
				break
			}
		}
		if visibleRight {
			return true
		}

		// Check top
		visibleBottom := true
		for b := y + 1; b < h; b++ {
			if trees[b][x] >= tree {
				visibleBottom = false
				break
			}
		}
		if visibleBottom {
			return true
		}

		return false
	}

	totalVisible := 0
	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			if checkIsVisible(x, y) {
				totalVisible += 1
			}
		}
	}
	totalVisible += w + w + (h-2)*2

	fmt.Printf("Part1: %d\n", totalVisible)

	getScore := func(x int, y int) int {
		//fmt.Printf("=== Checking %d %d\n", x, y)
		tree := trees[y][x]

		// Check left
		visibleLeft := x
		for l := x - 1; l >= 0; l-- {
			//fmt.Printf("Checking left %d\n", l)
			if trees[y][l] >= tree {
				visibleLeft = x - l
				break
			}
		}
		//fmt.Printf("visibleLeft %d\n", visibleLeft)

		// Check top
		visibleTop := y
		for t := y - 1; t >= 0; t-- {
			//fmt.Printf("Checking top %d\n", t)
			if trees[t][x] >= tree {
				visibleTop = y - t
				break
			}
		}
		//fmt.Printf("visibleTop %d\n", visibleTop)

		// Check top
		visibleRight := w - x - 1
		for r := x + 1; r < w; r++ {
			if trees[y][r] >= tree {
				visibleRight = r - x
				break
			}
		}
		//fmt.Printf("visibleRight %d\n", visibleRight)

		// Check top
		visibleBottom := h - y - 1
		for b := y + 1; b < h; b++ {
			if trees[b][x] >= tree {
				visibleBottom = b - y
				break
			}
		}
		//fmt.Printf("visibleBottom %d\n", visibleBottom)

		return visibleLeft * visibleTop * visibleRight * visibleBottom
	}

	maxScore := 0
	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			score := getScore(x, y)
			if score > maxScore {
				maxScore = score
			}
		}
	}

	fmt.Printf("Part2: %d\n", maxScore)
}

func main() {
	a()
}
