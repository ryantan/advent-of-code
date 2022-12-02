package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func readCalories() []int {

	f, err := os.Open("../input.txt")
	if err != nil {
		panic("Can't read input")
	}

	caloriesByElf := make([]int, 0)

	scanner := bufio.NewScanner(f)
	caloriesNow := 0
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)

		//
		if len(line) == 0 {
			caloriesByElf = append(caloriesByElf, -caloriesNow)
			caloriesNow = 0
			//fmt.Println("new line")
			continue
		}

		x, err := strconv.Atoi(line)
		if err != nil {
			panic("wrong input")
		}

		//fmt.Println(x)
		caloriesNow += x
	}
	// Add last elf.
	caloriesByElf = append(caloriesByElf, -caloriesNow)

	//fmt.Println(caloriesByElf)
	return caloriesByElf
}

func a() {
	caloriesByElf := readCalories()

	max := 0
	for _, calories := range caloriesByElf {
		if calories < max {
			max = calories
		}
	}
	fmt.Println(-max)
}

func b() {
	caloriesByElf := readCalories()
	sort.Ints(caloriesByElf)
	fmt.Println(-(caloriesByElf[0] + caloriesByElf[1] + caloriesByElf[2]))
}

func main() {
	//a()
	b()
}
