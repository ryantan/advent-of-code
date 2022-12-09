package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

//var fileName = "../sample.txt"

//var fileName = "../sample2.txt"

var fileName = "../input.txt"

func followHead(head [2]int, tail [2]int) [2]int {
	diffX, diffY := head[0]-tail[0], head[1]-tail[1]
	diffX2, diffY2 := diffX*diffX, diffY*diffY
	diff := diffX2 + diffY2
	//fmt.Printf("diff: %d (%d, %d)\n", diff, diffX, diffY)

	// 0: same position, 1: 1 space straight, or 2: 1 space diagonal.
	if diff == 0 || diff == 1 || diff == 2 {
		// No change.
		return tail
	}

	// 4: 2 space straight left/up/right/down,
	// 8: 2 space diagonally
	if diff == 4 || diff == 8 {
		// Move 1 left/up/right/down or 1 diagonally.
		return [2]int{tail[0] + (diffX / 2), tail[1] + (diffY / 2)}
	}

	// Last possible case: 5: 1 space in one dimension, 2 in the other.
	if diffY2 > diffX2 {
		// 1 space in x, 2 space in y
		return [2]int{tail[0] + diffX, tail[1] + (diffY / 2)}
	} else {
		// 2 space in x, 1 space in y
		return [2]int{tail[0] + (diffX / 2), tail[1] + diffY}
	}
}

func a(numberOfKnots int) int {
	// Keeps track of knot positions.
	knotCoords := make([][2]int, numberOfKnots)

	// Keeps track of tail (the last knot) positions.
	tailCoords := make([][2]int, 0)

	direction, moves := "", 0
	scanner := common.GetLineScanner(fileName)
	for scanner.Scan() {
		l := scanner.Text()
		_, err := fmt.Sscanf(l, "%s %d", &direction, &moves)
		if err != nil {
			panic("Could not parse input.")
		}

		for i := 0; i < moves; i++ {
			// Move head.
			if direction == "L" {
				knotCoords[0][0]--
			} else if direction == "U" {
				knotCoords[0][1]--
			} else if direction == "R" {
				knotCoords[0][0]++
			} else if direction == "D" {
				knotCoords[0][1]++
			}

			// Move knots.
			for k := 1; k < numberOfKnots; k++ {
				knotCoords[k] = followHead(knotCoords[k-1], knotCoords[k])
			}

			// Track tail positions.
			tailCoords = append(tailCoords, knotCoords[numberOfKnots-1])
			//debug.PrintKnots(knotCoords)
		}
		//fmt.Printf("tailCoords: %+v\n", tailCoords)
	}

	// Count unique coords.
	positions := make(map[string]bool, 0)
	totalPositions := 0
	for _, coord := range tailCoords {
		position := fmt.Sprintf("%d,%d", coord[0], coord[1])
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
