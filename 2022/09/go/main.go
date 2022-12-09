package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

//var fileName = "../sample.txt"

//var fileName = "../sample2.txt"

var fileName = "../input.txt"

type coord [2]int

var directionalMoves = map[string]coord{
	"L": {-1, 0},
	"U": {0, -1},
	"R": {1, 0},
	"D": {0, 1},
}

func (c *coord) Move(value coord) {
	c[0] += value[0]
	c[1] += value[1]
}

func (c *coord) Add(x int, y int) *coord {
	c[0] += x
	c[1] += y
	return c
}

func (c *coord) follow(head coord) *coord {
	diffX, diffY := head[0]-c[0], head[1]-c[1]
	diffX2, diffY2 := diffX*diffX, diffY*diffY
	diff := diffX2 + diffY2
	//fmt.Printf("diff: %d (%d, %d)\n", diff, diffX, diffY)

	// 0: same position, 1: 1 space straight, or 2: 1 space diagonal.
	if diff == 0 || diff == 1 || diff == 2 {
		// No change.
		return c
	}

	// 4: 2 space straight left/up/right/down,
	// 8: 2 space diagonally
	if diff == 4 || diff == 8 {
		// Move 1 left/up/right/down or 1 diagonally.
		return c.Add(diffX/2, diffY/2)
	}

	// Last possible case: 5: 1 space in one dimension, 2 in the other.
	if diffY2 > diffX2 {
		// 1 space in x, 2 space in y
		return c.Add(diffX, diffY/2)
	} else {
		// 2 space in x, 1 space in y
		return c.Add(diffX/2, diffY)
	}
}

func a(numberOfKnots int) int {
	// Keeps track of knot positions.
	knotCoords := make([]coord, numberOfKnots)

	// Keeps track of tail (the last knot) positions.
	tailCoords := make([]coord, 0)

	direction, moves := "", 0
	scanner := common.GetLineScanner(fileName)
	for scanner.Scan() {
		_, err := fmt.Sscanf(scanner.Text(), "%s %d", &direction, &moves)
		if err != nil {
			panic("Could not parse input.")
		}

		for i := 0; i < moves; i++ {
			// Move head.
			knotCoords[0].Move(directionalMoves[direction])

			// Move knots.
			for k := 1; k < numberOfKnots; k++ {
				knotCoords[k].follow(knotCoords[k-1])
			}

			// Track tail positions.
			tailCoords = append(tailCoords, knotCoords[numberOfKnots-1])
			//debug.PrintKnots(knotCoords)
		}
		//fmt.Printf("tailCoords: %+v\n", tailCoords)
	}

	// Count unique coords.
	positions := make(map[coord]bool, 0)
	totalPositions := 0
	for _, position := range tailCoords {
		if _, exists := positions[position]; !exists {
			positions[position] = true
			totalPositions++
		}
	}

	return totalPositions
}

func main() {
	part1 := a(2)
	fmt.Printf("Part1: %d\n", part1)
	part2 := a(10)
	fmt.Printf("Part2: %d\n", part2)
}
