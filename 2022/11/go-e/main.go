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
	opChainPending []Op
	factors        map[int]bool
	remainders     map[int]int // remainder for each factor
}

func printChain(chain []Op) string {
	output := make([]string, 0)
	for _, op := range chain {
		output = append(output, fmt.Sprintf("%s %d", op.label, op.value))
	}
	return strings.Join(output, ", ")
}

var chainRemainderCalled = 0
var chainRemainderHits = 0
var chainRemainderMisses = 0

//var memoizedChainRemainder = make(map[string]int, 0)

// factor, value, chain
var memoizedChainRemainder2 = make(map[int]map[int]map[string]int, 0)

func getChainString(chain []Op) string {
	output := make([]string, 0)
	for _, op := range chain {
		output = append(output, fmt.Sprintf("%s %d", op.label, op.value))
	}
	return strings.Join(output, ",")
}

func getKey(factor, initial int, chain []Op) string {
	return fmt.Sprintf("%d|%d|%s", factor, initial, getChainString(chain))
}

func (item *Item) divisibleBy(factor int) bool {
	remainder := item.remainders[factor]
	return remainder == 0
}

//func (item *Item) squash() *big.Int {
//	opsPending := len(item.opChainPending)
//	if opsPending == 0 {
//		return item.worry
//	}
//
//	for _, op := range item.opChainPending {
//		//fmt.Printf("op: %+v\n", op)
//		switch op.label {
//		case "square":
//			//squareCalled++
//			//fmt.Printf("squareCalled: #%d\n", squareCalled)
//			item.worry.Mul(item.worry, item.worry)
//		case "multiply":
//			//mulCalled++
//			//fmt.Printf("mulCalled: #%d\n", mulCalled)
//			item.worry.Mul(item.worry, op.bigValue)
//		case "add":
//			//addCalled++
//			//fmt.Printf("addCalled: #%d\n", addCalled)
//			item.worry.Add(item.worry, op.bigValue)
//		}
//	}
//	// Clear pending chain.
//	item.opChainPending = make([]Op, 0)
//
//	return item.worry
//}

// Inspects item and result which monkey to throw to.
func (m *Monkey) inspect(item *Item) int {

	//// if last item is same op, merge them.
	//opsCount := len(item.opChainPending)
	//merged := false
	//if opsCount > 0 {
	//	lastOp := item.opChainPending[opsCount-1]
	//
	//	if m.op.label == lastOp.label {
	//		if m.op.label == "multiply" {
	//			merged = true
	//			mergedValue := lastOp.value * m.op.value
	//			item.opChainPending[opsCount-1] = Op{
	//				operand:  lastOp.operand,
	//				op:       lastOp.op,
	//				label:    lastOp.label,
	//				value:    mergedValue,
	//				bigValue: big.NewInt(int64(mergedValue)),
	//			}
	//		} else if m.op.label == "add" {
	//			merged = true
	//			mergedValue := lastOp.value + m.op.value
	//			item.opChainPending[opsCount-1] = Op{
	//				operand:  lastOp.operand,
	//				op:       lastOp.op,
	//				label:    lastOp.label,
	//				value:    mergedValue,
	//				bigValue: big.NewInt(int64(mergedValue)),
	//			}
	//		}
	//	}
	//}
	//if !merged {
	//	item.opChainPending = append(item.opChainPending, m.op)
	//}

	if item.remainders == nil {
		newRemainders := make(map[int]int, 0)
		for _, factor := range noteWorthyFactors {
			newRemainders[factor] = item.initialWorry % factor
		}
		item.remainders = newRemainders
	}

	item.opChainPending = append(item.opChainPending, m.op)

	lastOp := m.op
	newRemainders := make(map[int]int, 0)
	for _, factor := range noteWorthyFactors {
		r := item.remainders[factor]
		switch lastOp.label {
		case "square":
			r = (r * r) % factor
			break
		case "multiply":
			r = ((lastOp.value % factor) * r) % factor
			break
		//case "double":
		//	r = ((2 % factor) * chainRemainder(remainingChain)) % factor
		//	break
		case "add":
			r = ((lastOp.value % factor) + r) % factor
			break
		default:
			panic("No such label")
		}
		newRemainders[factor] = r
	}
	item.remainders = newRemainders

	//if len(item.opChainPending) > 10 && item.worry.Int64() < 100000 {
	//	item.squash()
	//}

	if item.divisibleBy(m.divisor) {
		return m.destinationTrue
	}
	return m.destinationFalse
}

var noteWorthyFactorsHash = make(map[int]bool, 0)
var noteWorthyFactors = make([]int, 0)

func main() {
	//scanner := common.GetLineScanner("../sample.txt")
	scanner := common.GetLineScanner("../input.txt")

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
			item := Item{
				initialWorry:   worry,
				factors:        make(map[int]bool, 0),
				worry:          big.NewInt(int64(worry)),
				opChainPending: make([]Op, 0),
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
				//op = Op{operand: operand, op: double, label: "double", value: rhs, bigValue: bigRhs}
				op = Op{operand: operand, op: multiply, label: "multiply", value: 2, bigValue: big2}
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

	//for _, item := range items {
	//	item.factorize()
	//}

	fmt.Printf("== At the start ==\nmonkeyInspected: %+v\n\n", monkeyInspected)
	fmt.Printf("noteWorthyFactors: %+v\n\n", noteWorthyFactors)
	printMonkeyItems(monkeys)
	printItems(items)

	//rounds := 1
	//rounds := 2
	//rounds := 20
	//rounds := 1000
	//rounds := 2000
	rounds := 10000
	start := time.Now()
	for r := 1; r <= rounds; r++ {
		//roundStart := time.Now()
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

		//elapsed := time.Since(roundStart)
		//fmt.Printf("Round %d took %s\n\n", r, elapsed)

		//if r%1000 == 0 || r == 40 || r == 20 || r == 1 {
		if r%1000 == 0 || r == 20 || r == 1 {
			fmt.Printf("== After round %d ==\nmonkeyInspected: %+v\n\n", r, monkeyInspected)
			//printMonkeyItems(monkeys)
			//printItems(items)
			elapsed := time.Since(start)
			fmt.Printf("Elapsed: %s\n\n", elapsed)
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("%d rounds took %s\n\n", rounds, elapsed)
	fmt.Printf("== After round %d ==\nmonkeyInspected: %+v\n\n", rounds, monkeyInspected)
	printMonkeyItems(monkeys)
	//printItems(items)

	sort.Ints(monkeyInspected)
	fmt.Printf("Sorted monkeyInspected: %+v\n\n", monkeyInspected)

	fmt.Printf("level of monkey business: %d\n", monkeyInspected[numberOfMonkeys-1]*monkeyInspected[numberOfMonkeys-2])

	fmt.Printf("chainRemainderCalled: %d\n", chainRemainderCalled)
	fmt.Printf("chainRemainderHits: %d\n", chainRemainderHits)
	fmt.Printf("chainRemainderMisses: %d\n", chainRemainderMisses)
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
