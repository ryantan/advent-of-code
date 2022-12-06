package main

import (
	"bufio"
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

//var fileName = "../sample.txt"
var fileName = "../input.txt"

// Check if there are repeat chars in the block.
func isDiff(block map[string]int, blockSize int) bool {
	ones := 0
	for _, i := range block {
		// Optional early return.
		if i > 1 {
			return false
		}
		ones += i
		if ones == blockSize {
			return true
		}
	}
	return false
}

func findSequenceOfDistinct(scanner *bufio.Scanner, scanner2 *bufio.Scanner, blockSize int) int {
	buffer := make(map[string]int, blockSize)
	for i := 0; scanner.Scan(); i++ {
		// Advance start of block scanner if we're past the first block.
		if i >= blockSize {
			scanner2.Scan()
			buffer[scanner2.Text()]--
		}
		buffer[scanner.Text()]++
		if isDiff(buffer, blockSize) {
			return i + 1
		}
	}

	// Undefined response.
	return 0
}

func a() {
	scanner := common.GetCharScanner(fileName)
	scanner2 := common.GetCharScanner(fileName)
	processed := findSequenceOfDistinct(scanner, scanner2, 4)
	fmt.Printf("Part1: %d\n", processed)
}

func b() {
	scanner := common.GetCharScanner(fileName)
	scanner2 := common.GetCharScanner(fileName)
	processed := findSequenceOfDistinct(scanner, scanner2, 14)
	fmt.Printf("Part2: %d\n", processed)
}

func main() {
	a()
	b()
}
