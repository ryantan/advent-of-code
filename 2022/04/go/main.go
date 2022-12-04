package main

import (
	"bufio"
	"fmt"
	"os"
)

//func getRangeInInt(s string) (int, int) {
//	elf := strings.Split(s, "-")
//
//	from, err := strconv.Atoi(elf[0])
//	if err != nil {
//		panic("can't parse from")
//	}
//
//	to, err := strconv.Atoi(elf[1])
//	if err != nil {
//		panic("can't parse to")
//	}
//
//	return from, to
//}
//
//func getRanges(l string) (int, int, int, int) {
//	elves := strings.Split(l, ",")
//	from1, to1 := getRangeInInt(elves[0])
//	from2, to2 := getRangeInInt(elves[1])
//	return from1, to1, from2, to2
//}

func getRanges(l string) (int, int, int, int) {
	var from1, to1, from2, to2 int
	_, err := fmt.Sscanf(l, "%d-%d,%d-%d", &from1, &to1, &from2, &to2)
	if err != nil {
		panic("cannot parse input")
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

		diff1 := to1 - from1
		diff2 := to2 - from2

		if diff2 > diff1 {
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

	fmt.Println(total)
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
			//fmt.Println("ranges are same")
			total++
			continue
		}

		if to2 >= from1 && to2 <= to1 {
			//fmt.Println("to2 is in to1-from1")
			total++
			continue
		} else if to1 >= from2 && to1 <= to2 {
			//fmt.Println("to1 is in to2-from2")
			total++
			continue
		}
	}

	fmt.Println(total)
}

func main() {
	a()
	b()
}
