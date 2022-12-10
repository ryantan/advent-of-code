package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

var fileName = "../input.txt"

func main() {
	cycle, value, opValue := 0, 1, 0
	nextSampleCycle, totalSamplesValue := 20, 0

	tick := func() {
		x := cycle % 40
		if x == value-1 || x == value || x == value+1 {
			print("#")
		} else {
			print(".")
		}
		if x == 39 {
			print("\n")
		}

		cycle++
		if cycle == nextSampleCycle {
			totalSamplesValue += nextSampleCycle * value
			nextSampleCycle += 40
		}
	}

	scanner := common.GetLineScanner(fileName)
	for scanner.Scan() {
		tick()
		if _, err := fmt.Sscanf(scanner.Text(), "addx %d", &opValue); err == nil {
			tick()
			value += opValue
		}
	}

	fmt.Printf("Part1: %d\n", totalSamplesValue)
}
