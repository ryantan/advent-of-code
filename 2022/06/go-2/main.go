package main

import (
	"bufio"
	"fmt"
	"os"
)

//var fileName = "../sample.txt"
var fileName = "../input.txt"

// Check if there are repeat chars in the block.
func isDiff(block map[string]int) bool {
	//fmt.Printf("block: %v\n", block)
	for _, i := range block {
		if i > 1 {
			return false
		}
	}
	return true
}

func findSequenceOfDistinct(scanner *bufio.Scanner, scanner2 *bufio.Scanner, blockSize int) int {
	buffer := make(map[string]int, blockSize)

	// Check first block.
	for i := 0; i < blockSize; i++ {
		if !scanner.Scan() {
			return 0
		}
		c := scanner.Text()
		buffer[c]++
	}
	if isDiff(buffer) {
		return 0
	}

	// Start of last blockSize chars.
	for i := blockSize; scanner.Scan(); i++ {
		scanner2.Scan()
		startOfBlock := scanner2.Text()
		latest := scanner.Text()
		buffer[latest]++
		buffer[startOfBlock]--
		if buffer[startOfBlock] == 0 {
			delete(buffer, startOfBlock)
		}
		if isDiff(buffer) {
			return i + 1
		}
	}

	// Undefined response.
	return 0
}

// getScanner opens a file and returns a bufio.Scanner.
func getScanner(file string) *bufio.Scanner {
	f, err := os.Open(file)
	if err != nil {
		panic("Can not read input " + file)
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanRunes)
	return scanner
}

func a() {
	scanner := getScanner(fileName)
	scanner2 := getScanner(fileName)
	processed := findSequenceOfDistinct(scanner, scanner2, 4)
	fmt.Printf("Part1: %d\n", processed)
}

func b() {
	scanner := getScanner(fileName)
	scanner2 := getScanner(fileName)
	processed := findSequenceOfDistinct(scanner, scanner2, 14)
	fmt.Printf("Part2: %d\n", processed)
}

func main() {
	a()
	b()
}
