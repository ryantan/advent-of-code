package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
	"math"
	"strconv"
)

//var fileName = "../sample.txt"

//var fileName = "../sample2.txt"

var fileName = "../input.txt"

func followHead(head [2]int, tail [2]int) [2]int {
	if head == tail {
		return tail
	}

	diffX := head[0] - tail[0]
	diffY := head[1] - tail[1]

	diff := diffX*diffX + diffY*diffY
	//fmt.Printf("diff: %d (%d, %d)\n", diff, diffX, diffY)

	// 1 space straight.
	if diff == 1 {
		return tail
	}

	// 1 space diagonal.
	if diff == 2 {
		return tail
	}

	// 2 space straight.
	if diff == 4 {
		return [2]int{tail[0] + (diffX / 2), tail[1] + (diffY / 2)}
	}

	// 1 space in one dimension, 2 in the other.
	if diff == 5 {
		// Move diagonally 1 space.
		if math.Abs(float64(diffY)) > math.Abs(float64(diffX)) {
			// 1 space in x, 2 space in y
			return [2]int{tail[0] + diffX, tail[1] + (diffY / 2)}
		} else {
			// 2 space in x, 1 space in y
			return [2]int{tail[0] + (diffX / 2), tail[1] + diffY}
		}
	}

	// 2 space in both dimensions
	if diff == 8 {
		// Move diagonally 1 space.
		return [2]int{tail[0] + (diffX / 2), tail[1] + (diffY / 2)}
	}

	panic("Something went wrong.")
}

func a(numberOfKnots int) int {
	// Keeps track of knot positions.
	knotCoords := make([][2]int, numberOfKnots)

	// Keeps track of tail (last knot) positions.
	tailCoords := make([][2]int, 0)

	direction := ""
	moves := 0

	scanner := common.GetLineScanner(fileName)
	for scanner.Scan() {
		l := scanner.Text()
		_, err := fmt.Sscanf(l, "%s %d", &direction, &moves)
		if err != nil {
			panic("Could not parse input.")
		}
		//fmt.Printf("(%s %d)\n", direction, moves)

		// Update head, then update tail.

		for i := 0; i < moves; i++ {
			if direction == "L" {
				knotCoords[0][0]--
			} else if direction == "U" {
				knotCoords[0][1]--
			} else if direction == "R" {
				knotCoords[0][0]++
			} else if direction == "D" {
				knotCoords[0][1]++
			}
			//fmt.Printf("Knot %d: (%d, %d)\n", 0, knotCoords[0][0], knotCoords[0][1])

			for k := 1; k < numberOfKnots; k++ {
				knotCoords[k] = followHead(knotCoords[k-1], knotCoords[k])
				//fmt.Printf("Knot %d: (%d, %d)\n", k, knotCoords[k][0], knotCoords[k][1])
			}
			//fmt.Printf("Tail: (%d, %d)\n", currentTail[0], currentTail[1])
			tailCoords = append(tailCoords, knotCoords[numberOfKnots-1])
			//printKnots(knotCoords)
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

	//fmt.Printf("Part1: %d\n", totalPositions)
	return totalPositions
}

func printKnots(knots [][2]int) {
	minX, minY, maxX, maxY := 0, 0, 0, 0
	for _, knot := range knots {
		if knot[0] < minX {
			minX = knot[0]
		}
		if knot[1] < minY {
			minY = knot[1]
		}
		if knot[0] > maxX {
			maxX = knot[0]
		}
		if knot[1] > maxY {
			maxY = knot[1]
		}
	}
	minX, minY, maxX, maxY = minX-2, minY-2, maxX+2, maxY+2

	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			c := "."
			for k := len(knots) - 1; k >= 0; k-- {
				if knots[k][0] == x && knots[k][1] == y {
					if k == 0 {
						c = "H"
					} else {
						c = strconv.Itoa(k)
					}
				}
			}
			print(c)
		}
		print("\n")
	}
}

func main() {
	part1 := a(2)
	fmt.Printf("Part1: %d\n", part1)
	part2 := a(10)
	fmt.Printf("Part2: %d\n", part2)
}
