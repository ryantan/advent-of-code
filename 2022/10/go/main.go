package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

func main() {
	cycle, value, opValue, totalSamplesValue := 0, 1, 0, 0

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
		if cycle%20 == 0 && cycle%40 != 0 {
			totalSamplesValue += cycle * value
		}
	}

	scanner := common.GetLineScanner("../input.txt")
	for scanner.Scan() {
		tick()
		if _, err := fmt.Sscanf(scanner.Text(), "addx %d", &opValue); err == nil {
			tick()
			value += opValue
		}
	}

	fmt.Printf("Part1: %d\n", totalSamplesValue)
}
