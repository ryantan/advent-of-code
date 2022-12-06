package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

//var fileName = "../sample.txt"
var fileName = "../input.txt"

func priority(s string) int {
	if len(s) == 0 {
		return 0
	}

	charCode := int([]rune(s)[0])
	if charCode > 96 {
		return charCode - 96
	}

	return charCode - 38
}

func a() int {
	totalPriorities := 0

	scanner := common.GetLineScanner(fileName)
	for scanner.Scan() {
		contents := []rune(scanner.Text())

		var dup = ""

		size := len(contents) / 2
		hash := make(map[string]int, size)
		for i := 0; i < size; i++ {
			if hash[string(contents[i])] == 0 {
				hash[string(contents[i])] += 1
			}
		}
		for i := size; i < size*2; i++ {
			if hash[string(contents[i])] == 1 {
				dup = string(contents[i])
				break
			}
		}

		//fmt.Printf("dup: %v\n", dup)
		//fmt.Printf("priority: %v\n", priority(dup))

		totalPriorities += priority(dup)
	}

	fmt.Printf("totalPriorities: %d\n", totalPriorities)
	return totalPriorities
}

func b() int {
	totalPriorities := 0

	scanner := common.GetLineScanner(fileName)
	for scanner.Scan() {
		contents1 := []rune(scanner.Text())

		scanner.Scan()
		contents2 := []rune(scanner.Text())

		scanner.Scan()
		contents3 := []rune(scanner.Text())

		var dup = ""

		hash := make(map[string]int, len(contents1))

		for i := 0; i < len(contents1); i++ {
			if hash[string(contents1[i])] == 0 {
				hash[string(contents1[i])] += 1
			}
		}
		for i := 0; i < len(contents2); i++ {
			if hash[string(contents2[i])] == 1 {
				hash[string(contents2[i])] += 1
			}
		}
		for i := 0; i < len(contents3); i++ {
			if hash[string(contents3[i])] == 2 {
				dup = string(contents3[i])
				break
			}
		}

		//fmt.Printf("dup: %v\n", dup)
		//fmt.Printf("priority: %v\n", priority(dup))

		totalPriorities += priority(dup)
	}

	fmt.Printf("totalPriorities: %d\n", totalPriorities)
	return totalPriorities
}

func main() {
	a()
	b()
}
