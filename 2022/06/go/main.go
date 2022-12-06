package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

//var fileName = "../sample.txt"
var fileName = "../input.txt"

// Check if there are repeat chars in the block.
func isDiff(block map[rune]int) bool {
	//fmt.Printf("block: %v\n", block)
	for _, i := range block {
		if i > 1 {
			return false
		}
	}
	return true
}

func findSequenceOfDistinct(s []rune, count int) int {
	if len(s) < count {
		return 0
	}

	// Check first block.
	buffer := make(map[rune]int, count)
	for i := 0; i < count; i++ {
		buffer[s[i]]++
	}
	//fmt.Printf("first buffer: %v\n", buffer)

	if isDiff(buffer) {
		return 0
	}

	// Start of last count chars.
	for x, i := 0, count; i < len(s); x, i = x+1, i+1 {
		buffer[s[i]]++
		buffer[s[x]]--
		if buffer[s[x]] == 0 {
			delete(buffer, s[x])
		}
		if isDiff(buffer) {
			return i + 1
		}
	}

	// Undefined response.
	return len(s)
}

func a() {
	scanner := common.GetLineScanner(fileName)
	for scanner.Scan() {
		l := []rune(scanner.Text())
		processed := findSequenceOfDistinct(l, 4)
		fmt.Printf("Part1: %d\n", processed)
	}

	//fmt.Printf("Part1 done\n")
}

func b() {
	scanner := common.GetLineScanner(fileName)
	for scanner.Scan() {
		l := []rune(scanner.Text())
		processed := findSequenceOfDistinct(l, 14)
		fmt.Printf("Part2: %d\n", processed)
	}

	//fmt.Printf("Part2 done\n")
}

func main() {
	a()
	b()
}
