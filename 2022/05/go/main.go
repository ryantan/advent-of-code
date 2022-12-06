package main

import (
	"bufio"
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

//var fileName = "../sample.txt"
var fileName = "../input.txt"
var columnWidth = 4

func addCrates(stacks *[][]string, line []rune) bool {
	// Stop at the ` 1   2   3 ...` line.
	if line[1] == '1' {
		return true
	}

	for i := 0; i*columnWidth < len(line); i++ {
		// Append new stack when required.
		if i >= len(*stacks) {
			*stacks = append(*stacks, []string{})
		}

		// Prepend crate to stack.
		if crate := string(line[(i*columnWidth)+1]); crate != " " {
			(*stacks)[i] = append([]string{crate}, (*stacks)[i]...)
		}
	}

	return false
}

// We're assuming length is always valid!
func popAFew(stack *[]string, count int) []string {
	l := len(*stack) - count
	last := (*stack)[l:]
	*stack = (*stack)[:l]
	return last
}

func scanStacks(scanner *bufio.Scanner) [][]string {
	stacks := make([][]string, 3)
	for scanner.Scan() {
		if isDone := addCrates(&stacks, []rune(scanner.Text())); isDone {
			break
		}
	}
	// Empty line separating stacks and moves.
	scanner.Scan()
	return stacks
}

func getStep(l string) (int, int, int) {
	var cratesCount, from, to int
	_, err := fmt.Sscanf(l, "move %d from %d to %d", &cratesCount, &from, &to)
	if err != nil {
		panic("Could not parse step input " + l)
	}
	//fmt.Printf("Move %d from %d to %d\n", cratesCount, from, to)
	return cratesCount, from - 1, to - 1
}

func a(oneByOne bool) {
	scanner := common.GetLineScanner(fileName)

	// Scan stacks.
	stacks := scanStacks(scanner)

	// Scan moves.
	for scanner.Scan() {
		cratesCount, from, to := getStep(scanner.Text())

		intermediateSteps := 1
		if oneByOne {
			intermediateSteps, cratesCount = cratesCount, intermediateSteps
		}

		for i := 0; i < intermediateSteps; i++ {
			crates := popAFew(&stacks[from], cratesCount)
			stacks[to] = append(stacks[to], crates...)
			//fmt.Printf("stacks: %v\n", stacks)
		}
	}

	lastCrates := ""
	for _, stack := range stacks {
		lastCrates = lastCrates + popAFew(&stack, 1)[0]
	}

	fmt.Printf("Top crates: %s\n", lastCrates)
}

func main() {
	a(true)
	a(false)
}
