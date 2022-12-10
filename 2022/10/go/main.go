package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
	"strconv"
	"strings"
)

//var fileName = "../sample.txt"

var fileName = "../input.txt"

func a() {
	cycle := 0
	nextSampleCycle := 20
	totalSamplesValue := 0
	value := 1

	takeSample := func() {
		cycle++
		fmt.Printf("%d: %d\n", cycle, value)
		if cycle == nextSampleCycle {
			totalSamplesValue += nextSampleCycle * value
			fmt.Printf("=== %d: %d (%d)\n", cycle, value, nextSampleCycle*value)
			nextSampleCycle += 40
		}
	}

	scanner := common.GetLineScanner(fileName)
	for scanner.Scan() {
		l := strings.Split(scanner.Text(), " ")

		if l[0] == "noop" {
			takeSample()
			continue
		}

		opValue, err := strconv.Atoi(l[1])
		if err != nil {
			panic("Could not parse input")
		}

		takeSample()
		takeSample()
		value += opValue
	}

	fmt.Printf("Part1: %d\n", totalSamplesValue)
}

func main() {
	a()
}
