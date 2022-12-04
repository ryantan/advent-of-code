package main

import (
	"bufio"
	"fmt"
	"os"
)

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
	//f, err := os.Open("../sample.txt")
	f, err := os.Open("../input.txt")
	if err != nil {
		panic("Can't read input")
	}

	total := 0

	scanner := bufio.NewScanner(f)
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
	//f, err := os.Open("../sample.txt")
	f, err := os.Open("../input.txt")
	if err != nil {
		panic("Can't read input")
	}

	total := 0

	scanner := bufio.NewScanner(f)
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
