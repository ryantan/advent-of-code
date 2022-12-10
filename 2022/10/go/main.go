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
	//opValue, opCountDown := 0, 0
	nextSampleCycle := 20
	//sampledValues := make([]int, 0)
	totalSamplesValue := 0
	value := 1

	scanner := common.GetLineScanner(fileName)
	for scanner.Scan() {
		l := strings.Split(scanner.Text(), " ")

		valueAfterSample := value
		if l[0] == "noop" {
			cycle++
		} else {
			v, err := strconv.Atoi(l[1])
			if err != nil {
				panic("Could not parse input")
			}
			//fmt.Printf("v: %d\n", v)
			valueAfterSample = value + v

			cycle += 2
		}

		if cycle >= nextSampleCycle {
			//sampledValues = append(sampledValues, value)
			totalSamplesValue += nextSampleCycle * value
			fmt.Printf("=== %d: %d\n", cycle-1, value)
			nextSampleCycle += 40
		}
		value = valueAfterSample

		fmt.Printf("%d: %d\n", cycle, value)

	}

	fmt.Printf("Part1: %d\n", totalSamplesValue)
}

func main() {
	a()
}
