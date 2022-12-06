package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
	"sort"
	"strconv"
)

//var fileName = "../sample.txt"
var fileName = "../input.txt"

func readCalories() []int {
	caloriesByElf := make([]int, 0)
	caloriesNow := 0

	scanner := common.GetLineScanner(fileName)
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
	a()
	b()
}
