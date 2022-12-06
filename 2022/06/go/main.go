package main

import (
	"bufio"
	"fmt"
	"os"
)

//func isDiff(c1 rune,c2 rune,c3 rune,c4 rune) bool {
//	return c1 != c2 && c1 != c3 && c1 != c4 && (c2 != c3 && c2 != c4) && (c3 != c4)
//}

func isDiff(last4 map[rune]int) bool {
	fmt.Printf("last4: %v\n", last4)
	for _, i := range last4 {
		if i > 1 {
			return false
		}
	}
	return true
}

func findToken(s []rune) int {
	if len(s) < 4 {
		return 0
	}

	// Check first 4.
	last4 := make(map[rune]int, 4)
	last4[s[0]]++
	last4[s[1]]++
	last4[s[2]]++
	last4[s[3]]++
	fmt.Printf("first last4: %v\n", last4)

	if isDiff(last4) {
		return 0
	}

	// Start of last 4 chars.
	x := 0
	for i := 4; i < len(s); i++ {
		last4[s[i]]++
		last4[s[x]]--
		if last4[s[x]] == 0 {
			delete(last4, s[x])
		}
		if isDiff(last4) {
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
		processed := findToken(l)
		fmt.Printf("Part1: %d\n", processed)
	}

	fmt.Printf("Part1 done\n")
}

func main() {
	a()
}
