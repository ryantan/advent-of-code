package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"
)

func add(worry int, value int) int {
	return worry + value
}

func multiply(worry int, value int) int {
	return worry * value
}

func square(worry int, _ int) int {
	return worry * worry
}

func double(worry int, _ int) int {
	return worry + worry
}

var big2 = big.NewInt(2)
var factoringCalled = 0

//var squareCalled = 0
//var mulCalled = 0
//var doubleCalled = 0
//var addCalled = 0

type Op struct {
	operand  string
	op       func(int, int) int
	label    string // For debugging
	value    int
	bigValue *big.Int
}

type Monkey struct {
	op               Op
	divisor          int
	destinationTrue  int
	destinationFalse int
	inspected        int
	items            []*Item
}

type Item struct {
	initialWorry   int
	worry          *big.Int
	opChain        []Op // Keeping full chain for debugging
	opChainPending []Op
	factors        map[int]bool
}

func (item *Item) divisibleBy(factor int) bool {
	_, exists := item.factors[factor]
	return exists
}

func (item *Item) factorize() *big.Int {
	factoringCalled++
	fmt.Printf("factoringCalled #%d\n", factoringCalled)

	opsPending := len(item.opChainPending)
	if opsPending == 0 {
		return item.worry
	}

	tempInt := big.NewInt(1)
	//// Assumption: there'll only be 1 addition in every opChainPending.
	//for i, op := range item.opChainPending {
	//	fmt.Printf("op type: %s\n", op.label)
	//	if op.label == "add" {
	//		if i != len(item.opChainPending)-1 {
	//			panic("Add found in middle of opChainPending")
	//		}
	//	}
	//}

	for _, op := range item.opChainPending {
		//fmt.Printf("op: %+v\n", op)
		switch op.label {
		case "square":
			//squareCalled++
			//fmt.Printf("squareCalled: #%d\n", squareCalled)

			item.worry.Mul(item.worry, tempInt)
			tempInt = big.NewInt(1)
			item.worry.Exp(item.worry, big2, nil)
			//item.worry.Mul(item.worry, item.worry)

		case "multiply":
			//mulCalled++
			//fmt.Printf("mulCalled: #%d\n", mulCalled)
			tempInt.Mul(tempInt, op.bigValue)
		case "double":
			//doubleCalled++
			//fmt.Printf("doubleCalled: #%d\n", doubleCalled)
			tempInt.Mul(tempInt, big2)
		case "add":
			//addCalled++
			//fmt.Printf("addCalled: #%d\n", addCalled)

			item.worry.Mul(item.worry, tempInt)
			tempInt = big.NewInt(1)
			item.worry.Add(item.worry, op.bigValue)
		}
		item.opChain = append(item.opChain, op)
	}
	item.opChainPending = make([]Op, 0)

	factors := make(map[int]bool, 0)
	for _, factor := range noteWorthyFactors {
		if big.NewInt(0).Mod(item.worry, big.NewInt(int64(factor))).Int64() == 0 {
			factors[factor] = true
		}
	}

	item.factors = factors

	return item.worry
}

// Inspects item and result which monkey to throw to.
func (m *Monkey) inspect(item *Item) int {
	item.opChainPending = append(item.opChainPending, m.op)

	switch m.op.label {
	case "square":
		// Doesn't affect factors.
		break
	case "multiply":
		if _, exists := noteWorthyFactorsHash[m.op.value]; exists {
			item.factors[m.op.value] = true
		}
	case "double":
		if _, exists := noteWorthyFactorsHash[2]; exists {
			item.factors[2] = true
		}
	case "add":
		// Factorize again.
		item.factorize()
		break
	}
	if item.divisibleBy(m.divisor) {
		return m.destinationTrue
	}
	return m.destinationFalse
}

var noteWorthyFactorsHash = make(map[int]bool, 0)
var noteWorthyFactors = make([]int, 0)

func main() {
	scanner := common.GetLineScanner("../sample.txt")
	//scanner := common.GetLineScanner("../input.txt")

	getInstruction := func() string {
		scanner.Scan()
		parts := strings.Split(scanner.Text(), ": ")
		return parts[1]
	}

	monkeys := make([]*Monkey, 0)
	items := make([]*Item, 0)

	for scanner.Scan() {
		// "  Starting items"
		itemsRaw := strings.Split(getInstruction(), ", ")
		monkeyItems := make([]*Item, 0)
		for _, worry := range itemsRaw {
			worry, _ := strconv.Atoi(worry)

			//factors := make(map[int]bool, 0)
			//for i := 1; i <= worry; i++ {
			//	if worry%i == 0 {
			//		factors[i] = true
			//	}
			//}

			item := Item{
				initialWorry: worry,
				//factors:      factors,
				factors:        make(map[int]bool, 0),
				worry:          big.NewInt(int64(worry)),
				opChainPending: make([]Op, 0),
				opChain:        make([]Op, 0),
			}
			items = append(items, &item)
			monkeyItems = append(monkeyItems, &item)
		}

		x, operand, y, divisor, destinationTrue, destinationFalse := "", "", "", 0, 0, 0

		// "  Operation"
		_, _ = fmt.Sscanf(getInstruction(), "new = %s %s %s", &x, &operand, &y)

		// "  Test"
		_, _ = fmt.Sscanf(getInstruction(), "divisible by %d", &divisor)

		// "    If true"
		_, _ = fmt.Sscanf(getInstruction(), "throw to monkey %d", &destinationTrue)

		// "    If false"
		_, _ = fmt.Sscanf(getInstruction(), "throw to monkey %d", &destinationFalse)

		// Just in case there might be value * old.
		if x != "old" {
			x, y = y, x
		}

		// We are assuming lhs is always old.
		rhs := 0
		//fmt.Printf("Worry: %d, x: %s, y: %s\n", worry, x, y)
		if v, err := strconv.Atoi(y); err == nil {
			rhs = v
		}

		bigRhs := big.NewInt(int64(rhs))

		var op Op
		if operand == "*" {
			if y == "old" {
				op = Op{operand: operand, op: square, label: "square", value: rhs, bigValue: bigRhs}
			} else {
				op = Op{operand: operand, op: multiply, label: "multiply", value: rhs, bigValue: bigRhs}
			}
		} else {
			if y == "old" {
				op = Op{operand: operand, op: double, label: "double", value: rhs, bigValue: bigRhs}
			} else {
				op = Op{operand: operand, op: add, label: "add", value: rhs, bigValue: bigRhs}
			}
		}

		noteWorthyFactors = append(noteWorthyFactors, divisor)
		noteWorthyFactorsHash[divisor] = true

		monkey := Monkey{
			op:               op,
			divisor:          divisor,
			destinationTrue:  destinationTrue,
			destinationFalse: destinationFalse,
			items:            monkeyItems,
		}
		monkeys = append(monkeys, &monkey)

		// Discard newline.
		scanner.Scan()
	}
	numberOfMonkeys := len(monkeys)
	monkeyInspected := make([]int, numberOfMonkeys)

	for _, item := range items {
		item.factorize()
	}

	fmt.Printf("== At the start ==\nmonkeyInspected: %+v\n\n", monkeyInspected)
	fmt.Printf("noteWorthyFactors: %+v\n\n", noteWorthyFactors)
	printMonkeyItems(monkeys)
	printItems(items)

	//rounds := 1
	//rounds := 2
	//rounds := 20
	rounds := 1000
	//rounds := 2000
	//rounds := 10000
	start := time.Now()
	for r := 1; r <= rounds; r++ {
		roundStart := time.Now()
		for m, monkey := range monkeys {
			tempItems := monkey.items
			//fmt.Printf("monkey %d items: %+v\n", m, items)
			for _, item := range tempItems {
				monkey.items = monkey.items[1:]
				destination := monkey.inspect(item)
				monkeyInspected[m]++
				//fmt.Printf("Throws %d to monkey %d\n", item.initialWorry, destination)
				monkeys[destination].items = append(monkeys[destination].items, item)
				//printMonkeyItems(monkeys)
			}
		}
		//fmt.Printf("== After round %d ==\nmonkeyInspected: %+v\n\n", r, monkeyInspected)
		//printMonkeyItems(monkeys)
		//printItems(items)
		elapsed := time.Since(roundStart)
		fmt.Printf("Round %d took %s\n\n", r, elapsed)
	}
	elapsed := time.Since(start)
	fmt.Printf("%d rounds took %s\n\n", rounds, elapsed)
	fmt.Printf("== After round %d ==\nmonkeyInspected: %+v\n\n", rounds, monkeyInspected)
	printMonkeyItems(monkeys)
	//printItems(items)

	sort.Ints(monkeyInspected)
	fmt.Printf("Sorted monkeyInspected: %+v\n\n", monkeyInspected)

	fmt.Printf("level of monkey business: %d\n", monkeyInspected[numberOfMonkeys-1]*monkeyInspected[numberOfMonkeys-2])

}

func printItems(items []*Item) {
	for i, item := range items {
		fmt.Printf("\nItem %d:\n", i)
		fmt.Printf("  Initial: %d\n", item.initialWorry)
		fmt.Printf("  Current: %s\n", item.worry.String())

		factors := make([]int, 0)
		for factor, _ := range item.factors {
			factors = append(factors, factor)
		}
		sort.Ints(factors)
		for f, factor := range factors {
			fmt.Printf("  f%d: %d\n", f, factor)
		}
	}
}

func printMonkeyItems(monkeys []*Monkey) {
	for i, monkey := range monkeys {
		itemWorryLevels := make([]int, 0)
		for _, item := range monkey.items {
			itemWorryLevels = append(itemWorryLevels, item.initialWorry)
		}
		fmt.Printf("Monkey %d: %v\n", i, itemWorryLevels)
	}
}
