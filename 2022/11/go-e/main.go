package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Item struct {
	initialWorry int
	remainders   map[int]int // remainder for each factor
}

func (item *Item) initialRemainders() {
	if item.remainders != nil {
		return
	}
	newRemainders := make(map[int]int, 0)
	for _, factor := range divisors {
		newRemainders[factor] = item.initialWorry % factor
	}
	item.remainders = newRemainders
}

// updateRemainders pre-computes remainders for all divisors, based on new op and value.
func (item *Item) updateRemainders(op string, value int) {
	newRemainders := make(map[int]int, 0)
	for _, factor := range divisors {
		r := item.remainders[factor]
		switch op {
		case "square":
			r = (r * r) % factor
			break
		case "multiply":
			r = ((value % factor) * r) % factor
			break
		case "add":
			r = ((value % factor) + r) % factor
			break
		default:
			panic("No such label")
		}
		newRemainders[factor] = r
	}
	item.remainders = newRemainders
}

func (item *Item) divisibleBy(factor int) bool {
	remainder := item.remainders[factor]
	return remainder == 0
}

type Monkey struct {
	op               string
	opValue          int
	divisor          int
	destinationTrue  int
	destinationFalse int
	inspected        int
	items            []*Item
}

// inspect the item and output which monkey to throw to.
func (m *Monkey) inspect(item *Item) int {
	// Inspect
	item.updateRemainders(m.op, m.opValue)

	// Route to next monkey.
	if item.divisibleBy(m.divisor) {
		return m.destinationTrue
	}
	return m.destinationFalse
}

func (m *Monkey) throw(otherMonkey *Monkey) {
	item := m.items[0]
	m.items = m.items[1:]
	otherMonkey.items = append(otherMonkey.items, item)
}

// Divisors from all monkeys.
var divisors = make([]int, 0)

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
				initialWorry: worry,
			}
			items = append(items, &item)
			monkeyItems = append(monkeyItems, &item)
		}

		x, operand, y, divisor, destinationTrue, destinationFalse := "", "", "", 0, 0, 0

		// "  Operation"
		_, _ = fmt.Sscanf(getInstruction(), "new = %s %s %s", &x, &operand, &y)
		// Just in case there might be value * old.
		if x != "old" {
			x, y = y, x
		}
		// We are assuming lhs is always old so we only care about rhs.
		rhs := 0
		if v, err := strconv.Atoi(y); err == nil {
			rhs = v
		}

		op := ""
		if operand == "*" {
			if y == "old" {
				op = "square"
			} else {
				op = "multiply"
			}
		} else {
			if y == "old" {
				rhs = 2
				op = "multiply"
			} else {
				op = "add"
			}
		}

		// "  Test"
		_, _ = fmt.Sscanf(getInstruction(), "divisible by %d", &divisor)
		divisors = append(divisors, divisor)

		// "    If true"
		_, _ = fmt.Sscanf(getInstruction(), "throw to monkey %d", &destinationTrue)
		// "    If false"
		_, _ = fmt.Sscanf(getInstruction(), "throw to monkey %d", &destinationFalse)

		monkey := Monkey{
			items:            monkeyItems,
			op:               op,
			opValue:          rhs,
			divisor:          divisor,
			destinationTrue:  destinationTrue,
			destinationFalse: destinationFalse,
		}
		monkeys = append(monkeys, &monkey)

		// Discard newline.
		scanner.Scan()
	}
	for _, item := range items {
		item.initialRemainders()
	}
	numberOfMonkeys := len(monkeys)
	monkeyInspected := make([]int, numberOfMonkeys)
	sort.Ints(divisors)
	fmt.Printf("== At the start ==\nmonkeyInspected: %+v\ndivisors: %+v\n\n", monkeyInspected, divisors)
	printMonkeyItems(monkeys)
	printItems(items)

	//rounds := 1000
	rounds := 10000
	start := time.Now()
	for r := 1; r <= rounds; r++ {
		for m, monkey := range monkeys {
			for _, item := range monkey.items {
				destination := monkey.inspect(item)
				monkeyInspected[m]++
				monkey.throw(monkeys[destination])
			}
		}

		if r%1000 == 0 || r == 20 || r == 1 {
			elapsed := time.Since(start)
			fmt.Printf("== After round %d (%s) ==\nmonkeyInspected: %+v\n\n", r, elapsed, monkeyInspected)
			//printMonkeyItems(monkeys)
			//printItems(items)
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("== After %d rounds (%s) ==\nmonkeyInspected: %+v\n\n", rounds, elapsed, monkeyInspected)
	printMonkeyItems(monkeys)
	printItems(items)

	sort.Ints(monkeyInspected)
	fmt.Printf("level of monkey business: %d\n", monkeyInspected[numberOfMonkeys-1]*monkeyInspected[numberOfMonkeys-2])
}

func printItems(items []*Item) {
	for i, item := range items {
		fmt.Printf("\nItem %d (%d):\n", i, item.initialWorry)
		for factor, remainder := range item.remainders {
			fmt.Printf("  %% %d: %d\n", factor, remainder)
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
