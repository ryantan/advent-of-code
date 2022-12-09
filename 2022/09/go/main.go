package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
	"math"
)

//var fileName = "../sample.txt"
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

	panic("Something went wrong.")
}

func a() {
	scanner := common.GetLineScanner(fileName)

	tailCoords := make([][2]int, 0)
	tailCoords = append(tailCoords, [2]int{0, 0})
	headCoords := make([][2]int, 0)
	headCoords = append(headCoords, [2]int{0, 0})

	currentTail := [2]int{0, 0}
	currentHead := [2]int{0, 0}

	for scanner.Scan() {
		l := scanner.Text()
		direction := ""
		moves := 0
		_, err := fmt.Sscanf(l, "%s %d", &direction, &moves)
		if err != nil {
			panic("Could not parse input.")
		}
		fmt.Printf("(%s %d)\n", direction, moves)

		// Update head, then update tail.

		for i := 0; i < moves; i++ {
			if direction == "L" {
				currentHead[0]--
			} else if direction == "U" {
				currentHead[1]--
			} else if direction == "R" {
				currentHead[0]++
			} else if direction == "D" {
				currentHead[1]++
			}
			headCoords = append(headCoords, currentHead)
			//fmt.Printf("Head: (%d, %d)\n", currentHead[0], currentHead[1])
			//fmt.Printf("headCoords: %+v\n", headCoords)

			currentTail = followHead(currentHead, currentTail)
			//fmt.Printf("Tail: (%d, %d)\n", currentTail[0], currentTail[1])
			tailCoords = append(tailCoords, currentTail)

			//fmt.Printf("tailCoords: %+v\n", tailCoords)
		}
	}

	positions := make(map[string]bool, 0)
	totalPositions := 0
	for _, coord := range tailCoords {
		position := fmt.Sprintf("%d,%d", coord[0], coord[1])
		if _, exists := positions[position]; !exists {
			positions[position] = true
			totalPositions++
		}
	}

	fmt.Printf("Part1: %d\n", totalPositions)
}

func main() {
	a()
}
