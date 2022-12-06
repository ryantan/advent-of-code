package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

//var fileName = "../sample.txt"
var fileName = "../input.txt"

func getRanges(l string) (int, int, int, int) {
	var from1, to1, from2, to2 int
	_, err := fmt.Sscanf(l, "%d-%d,%d-%d", &from1, &to1, &from2, &to2)
	if err != nil {
		panic("cannot parse input")
	}

	if to1 > to2 {
		return from2, to2, from1, to1
	}

	return from1, to1, from2, to2
}

func a() {
	total := 0

	scanner := common.GetLineScanner(fileName)
	for scanner.Scan() {
		from1, to1, from2, to2 := getRanges(scanner.Text())
		//fmt.Printf("%d - %d, %d - %d\n", from1, to1, from2, to2)

		if to1 == to2 && from1 == from2 {
			//fmt.Println("ranges are same")
			total++
			continue
		}

		if to2-from2 > to1-from1 {
			if to2 >= to1 && from2 <= from1 {
				total++
				continue
			}
		} else {
			if to1 >= to2 && from1 <= from2 {
				total++
				continue
			}
		}
	}

	fmt.Printf("part1: %d\n", total)
}

func b() {
	total := 0

	scanner := common.GetLineScanner(fileName)
	for scanner.Scan() {
		from1, to1, from2, to2 := getRanges(scanner.Text())
		//fmt.Printf("%d - %d, %d - %d\n", from1, to1, from2, to2)

		if to1 == to2 && from1 == from2 {
			total++
			continue
		}

		if to1 >= from2 {
			total++
			continue
		}
	}

	fmt.Printf("part2: %d\n", total)
}

func main() {
	a()
	b()
}
