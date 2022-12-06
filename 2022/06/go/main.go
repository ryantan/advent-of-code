package main

import (
	"bufio"
	"fmt"
	"os"
)

func isDiff(last4 map[rune]int) bool {
	//fmt.Printf("last4: %v\n", last4)
	for _, i := range last4 {
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

	// Check first 4.
	buffer := make(map[rune]int, count)
	for i := 0; i < count; i++ {
		buffer[s[i]]++
	}
	//fmt.Printf("first buffer: %v\n", buffer)

	if isDiff(buffer) {
		return 0
	}

	// Start of last count chars.
	x := 0
	for i := count; i < len(s); i++ {
		buffer[s[i]]++
		buffer[s[x]]--
		if buffer[s[x]] == 0 {
			delete(buffer, s[x])
		}
		if isDiff(buffer) {
			return i + 1
		}
		x++
	}

	return len(s)
}

func a() {
	//f, err := os.Open("../sample.txt")
	f, err := os.Open("../input.txt")
	if err != nil {
		panic("Can't read input")
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		l := []rune(scanner.Text())
		processed := findSequenceOfDistinct(l, 4)
		fmt.Printf("Part1: %d\n", processed)
	}

	fmt.Printf("Part1 done\n")
}

func b() {
	//f, err := os.Open("../sample.txt")
	f, err := os.Open("../input.txt")
	if err != nil {
		panic("Can't read input")
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		l := []rune(scanner.Text())
		processed := findSequenceOfDistinct(l, 14)
		fmt.Printf("Part2: %d\n", processed)
	}

	fmt.Printf("Part2 done\n")
}

func main() {
	a()
	b()
}
