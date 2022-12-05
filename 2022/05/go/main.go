package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func addCrates(stacksPtr *[][]string, l string) bool {
	stackCount := len(*stacksPtr)
	length := len(l)
	line := []rune(l)

	if line[1] == '1' {
		// Finish
		return true
	}

	stackIndex := 0
	for true {
		indexOfCrate := (stackIndex * 4) + 1
		if length < indexOfCrate {
			return false
		}
		crate := string(line[indexOfCrate])
		if crate == " " {
			// ignore
		} else {
			if stackIndex >= stackCount {
				*stacksPtr = append(*stacksPtr, []string{})
			}
			(*stacksPtr)[stackIndex] = append([]string{crate}, (*stacksPtr)[stackIndex]...)
		}
		stackIndex++
	}

	return false
}
func pop(stack *[]string) string {
	l := len(*stack)
	last := (*stack)[l-1]
	*stack = (*stack)[:l-1]
	return last
}

func a() {
	//f, err := os.Open("../sample.txt")
	f, err := os.Open("../input.txt")
	if err != nil {
		panic("Can't read input")
	}

	stacks := make([][]string, 3)
	fmt.Printf("stacks: %v\n", stacks)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		isDone := addCrates(&stacks, scanner.Text())
		fmt.Printf("stacks: %v\n", stacks)
		if isDone {
			fmt.Printf("done stacks: %v\n", stacks)
			break
		}
	}
	// empty line
	scanner.Scan()

	for scanner.Scan() {
		crates, from, to := getStep(scanner.Text())

		for i := 0; i < crates; i++ {
			crate := pop(&stacks[from])
			stacks[to] = append(stacks[to], crate)
			fmt.Printf("stacks: %v\n", stacks)
		}
	}

	lastCrates := ""
	for _, stack := range stacks {
		lastCrates = lastCrates + pop(&stack)
	}

	fmt.Printf("Part1: %s\n", lastCrates)
}

func getStep(l string) (int, int, int) {
	parts := strings.Split(l, " ")
	cratesCount, err := strconv.Atoi(parts[1])
	if err != nil {
		panic("Could not get number of crates to move.")
	}

	from, err := strconv.Atoi(parts[3])
	if err != nil {
		panic("Could not get from.")
	}

	to, err := strconv.Atoi(parts[5])
	if err != nil {
		panic("Could not get to.")
	}

	fmt.Printf("Move %d from %d to %d\n", cratesCount, from, to)
	return cratesCount, from - 1, to - 1
}

func main() {
	a()
}
