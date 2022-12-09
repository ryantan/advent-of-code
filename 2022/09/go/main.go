package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

var fileName = "../input.txt"

type coord [2]int

var directionalMoves = map[string][]int{"L": {-1, 0}, "U": {0, -1}, "R": {1, 0}, "D": {0, 1}}

func (c *coord) move(value ...int) *coord {
	c[0] += value[0]
	c[1] += value[1]
	return c
}

func (c *coord) follow(head coord) *coord {
	dX, dY := head[0]-c[0], head[1]-c[1]
	dX2, dY2 := dX*dX, dY*dY
	diff := dX2 + dY2

	// 4: 2 spaces straight, 8: 2 spaces diagonally
	if diff == 4 || diff == 8 {
		// Move 1 straight or diagonally.
		return c.move(dX/2, dY/2)
	}

	// 5: 1 space in one dimension, 2 in the other.
	if diff == 5 {
		if dY2 > dX2 {
			// Move 1 in x, 2 in y
			return c.move(dX, dY/2)
		} else {
			// Move 2 in x, 1 in y
			return c.move(dX/2, dY)
		}
	}

	// Remaining possibilities:
	// 0: same position, 1: 1 space straight, or 2: 1 space diagonal.
	return c // No change.
}

func findUniquePositions(numberOfKnots int) int {
	knots := make([]coord, numberOfKnots)
	tail := &knots[numberOfKnots-1]
	tailPositions, uniquePositions := make(map[coord]bool, 0), 0

	direction, moves := "", 0
	scanner := common.GetLineScanner(fileName)
	for scanner.Scan() {
		if _, err := fmt.Sscanf(scanner.Text(), "%s %d", &direction, &moves); err != nil {
			panic("Could not parse input.")
		}
		for i := 0; i < moves; i++ {
			knots[0].move(directionalMoves[direction]...)
			for k := 1; k < numberOfKnots; k++ {
				knots[k].follow(knots[k-1])
			}
			if _, exists := tailPositions[*tail]; !exists {
				tailPositions[*tail] = true
				uniquePositions++
			}
		}
	}
	return uniquePositions
}

func main() {
	fmt.Printf("Part1: %d\n", findUniquePositions(2))
	fmt.Printf("Part2: %d\n", findUniquePositions(10))
}
