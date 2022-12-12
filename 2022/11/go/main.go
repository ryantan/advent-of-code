package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
	"math"
	"sort"
	"strconv"
	"strings"
)

func main() {
	//scanner := common.GetLineScanner("../sample.txt")
	scanner := common.GetLineScanner("../input.txt")

	getInstruction := func() string {
		scanner.Scan()
		parts := strings.Split(scanner.Text(), ": ")
		return parts[1]
	}

	monkeys := make([]func(worry int) (int, int), 0)
	monkeyItems := make([][]int, 0)
	ops := map[string]func(lhs int, rhs int) int{
		"*": func(lhs int, rhs int) int { return lhs * rhs },
		"+": func(lhs int, rhs int) int { return lhs + rhs },
	}

	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "Mon") {

			// "  Starting items"
			itemsRaw := strings.Split(getInstruction(), ", ")
			items := make([]int, len(itemsRaw))
			for i, worry := range itemsRaw {
				worry, _ := strconv.Atoi(worry)
				items[i] = worry
			}
			monkeyItems = append(monkeyItems, items)

			x, op, y, divisor, destinationTrue, destinationFalse := "", "", "", 0, 0, 0

			// "  Operation"
			_, _ = fmt.Sscanf(getInstruction(), "new = %s %s %s", &x, &op, &y)

			// "  Test"
			_, _ = fmt.Sscanf(getInstruction(), "divisible by %d", &divisor)

			// "    If true"
			_, _ = fmt.Sscanf(getInstruction(), "throw to monkey %d", &destinationTrue)

			// "    If false"
			_, _ = fmt.Sscanf(getInstruction(), "throw to monkey %d", &destinationFalse)

			monkeys = append(monkeys, func(worry int) (int, int) {
				lhs, rhs := worry, worry
				//fmt.Printf("Worry: %d, x: %s, y: %s\n", worry, x, y)
				if v, err := strconv.Atoi(x); err == nil {
					lhs = v
				}
				if v, err := strconv.Atoi(y); err == nil {
					rhs = v
				}
				//fmt.Printf("Worry: %d, lhs: %d, rhs: %d\n", worry, lhs, rhs)

				worry = ops[op](lhs, rhs)
				//fmt.Printf("Worry is now %d\n", worry)
				worry = int(math.Floor(float64(worry) / 3))
				//fmt.Printf("Worry is now %d\n", worry)

				if worry%divisor == 0 {
					return worry, destinationTrue
				}
				return worry, destinationFalse
			})
		}
	}
	numberOfMonkeys := len(monkeyItems)
	monkeyInspected := make([]int, numberOfMonkeys)

	fmt.Printf("monkeyItems: %+v\n", monkeyItems)
	rounds := 20
	//rounds := 10000
	//rounds := 1000
	for i := 1; i <= rounds; i++ {
		for m, monkey := range monkeys {
			items := monkeyItems[m]
			//fmt.Printf("monkey %d items: %+v\n", m, items)
			for _, worry := range items {
				//fmt.Printf("worry: %d\n", worry)
				monkeyItems[m] = monkeyItems[m][1:]
				newWorry, destination := monkey(worry)
				monkeyInspected[m]++
				//fmt.Printf("Throws %d to monkey %d\n", newWorry, destination)
				monkeyItems[destination] = append(monkeyItems[destination], newWorry)
				//fmt.Printf("monkeyItems: %+v\n\n", monkeyItems)
			}
		}
		fmt.Printf("Round end monkeyItems: %+v\n", monkeyItems)
		fmt.Printf("== After round %d ==\nmonkeyInspected: %+v\n\n", i, monkeyInspected)
	}

	sort.Ints(monkeyInspected)
	fmt.Printf("Sorted nmonkeyInspected: %+v\n\n", monkeyInspected)

	fmt.Printf("level of monkey business: %d\n", monkeyInspected[numberOfMonkeys-1]*monkeyInspected[numberOfMonkeys-2])

}
