package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

func main() {
	cycle, x, opValue, totalSamplesValue := 0, 1, 0, 0

	tick := func() {
		c40 := cycle % 40
		if c40 == x-1 || c40 == x || c40 == x+1 {
			print("#")
		} else {
			print(".")
		}
		if x == 39 {
			print("\n")
		}

		cycle++
		if cycle%20 == 0 && cycle%40 != 0 {
			totalSamplesValue += cycle * x
		}
	}

	scanner := common.GetLineScanner("../input.txt")
	for scanner.Scan() {
		tick()
		if _, err := fmt.Sscanf(scanner.Text(), "addx %d", &opValue); err == nil {
			tick()
			x += opValue
		}
	}

	fmt.Printf("Part1: %d\n", totalSamplesValue)
}
