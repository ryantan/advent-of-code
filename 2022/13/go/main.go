package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
	"strconv"
)

type Packet struct {
	content         []rune
	originalContent []rune
	stack           []int
	head            int
	depth           int
	currentElement  int
}

func (p *Packet) convertToList() {
	// Pop.
	originalElement := p.pop()

	// Add list.
	p.stack = append(p.stack, -1)

	// Add back original element
	p.stack = append(p.stack, originalElement)

	// Update underlying string to close the list.
	originalString := fmt.Sprintf("[%d]", originalElement)
	aIndex := p.head - (len(originalString) - 2)
	bIndex := p.head
	b := p.content[bIndex:]

	newContent := make([]rune, 0)
	//fmt.Printf("p.content before part a: %s\n", string(p.content))
	newContent = append(newContent, p.content[:aIndex]...)
	newContent = append(newContent, []rune(originalString)...)
	newContent = append(newContent, b...)
	p.content = newContent
	p.head++
	//fmt.Printf("p.content after replace: %s\n", string(p.content))
}

func (p *Packet) pop() int {
	lastElement := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return lastElement
}

func (p *Packet) reset() {
	p.content = p.originalContent
	p.head = 0
	p.stack = nil
	p.depth = 0
	p.currentElement = 0
}
func (p *Packet) readNext() bool {
	if p.head >= len(p.content) {
		// No more.
		return false
	}

	c := p.content[p.head]
	if c == '[' {
		p.currentElement = -1
		p.stack = append(p.stack, -1)
		p.depth++
		p.head++
	} else if c == ']' {
		p.currentElement = -2
		p.pop()
		p.depth--
		p.head++
	} else {
		if c == ',' {
			// Ignore
			c = p.content[p.head]
			p.head++
		}

		if c >= 48 && c <= 57 {
			// Digit, read till next , or ] to terminate digit.
			number := []rune{c}
			p.head++

			// Read all subsequent digits as long as they are not "," or "]"
			for {
				if p.head >= len(p.content) {
					break
				}
				nextChar := p.content[p.head]
				if nextChar == ']' || nextChar == ',' {
					break
				}
				number = append(number, nextChar)
				p.head++
			}

			element, err := strconv.Atoi(string(number))
			if err != nil {
				panic("Could not parse number")
			}
			p.currentElement = element
			p.stack = append(p.stack, element)
		}
	}

	return true
}

// compare evaluates p as left packet and pRight as right packet and return true if order is correct.
func (p *Packet) compare(pRight *Packet) bool {
	// Reset
	p.reset()
	pRight.reset()

	//fmt.Println(string(p.originalContent))
	//fmt.Println(string(pRight.originalContent))

	correctOrder := false
	for i := 0; i < 10000; i++ {
		hasNext1 := p.readNext()
		hasNext2 := pRight.readNext()

		//fmt.Printf("Current elements: %d | %d\n", p.currentElement, pRight.currentElement)

		if !hasNext1 && !hasNext2 {
			//fmt.Printf("Both ran out\n")
			break
		} else if !hasNext1 && hasNext2 {
			// Left ran out.
			//fmt.Printf("Left ran out\n")
			correctOrder = true
			break
		} else if hasNext1 && !hasNext2 {
			// Right ran out.
			//fmt.Printf("Right ran out\n")
			break
		}

		//doneAllConvertingWeNeedToDo := false
		for {
			if p.currentElement == -1 && pRight.currentElement != -1 && pRight.currentElement != -2 {
				// Left is a list, right is not a list (and not closing a list either.
				// Convert right to list.
				//fmt.Printf("Convert right to list\n")
				pRight.convertToList()
				p.readNext()
				//fmt.Printf("Current elements: %d | %d\n", p.currentElement, pRight.currentElement)
			} else if p.currentElement != -1 && pRight.currentElement == -1 && p.currentElement != -2 {
				// Right is a list, left is not a list (and not closing a list either.
				// Convert left to list.
				//fmt.Printf("Convert left to list\n")
				p.convertToList()
				pRight.readNext()
				//fmt.Printf("Current elements: %d | %d\n", p.currentElement, pRight.currentElement)
			} else {
				//doneAllConvertingWeNeedToDo = true
				break
			}
		}

		if p.currentElement == -2 && pRight.currentElement != -2 {
			// Left ends the list first.
			//fmt.Printf("Left list ran out\n")
			correctOrder = true
			break
		} else if p.currentElement != -2 && pRight.currentElement == -2 {
			// Right ends the list first.
			//fmt.Printf("Right list ran out\n")
			break
		} else if p.currentElement == -2 && pRight.currentElement == -2 {
			// Both lists ended.
			//fmt.Printf("Both lists ended\n")
			continue
		} else if p.currentElement < pRight.currentElement {
			// Left is smaller.
			//fmt.Printf("Left is smaller\n")
			correctOrder = true
			break
		} else if p.currentElement > pRight.currentElement {
			// Right is smaller.
			//fmt.Printf("Right is smaller\n")
			break
		}
	}
	return correctOrder
}

func main() {
	//scanner := common.GetLineScanner("../sample.txt")
	//scanner := common.GetLineScanner("../sample2.txt")
	scanner := common.GetLineScanner("../input.txt")

	packets := make([]*Packet, 0)
	for scanner.Scan() {
		l1 := []rune(scanner.Text())
		if len(l1) == 0 {
			continue
		}
		packet := &Packet{
			originalContent: l1,
			content:         l1,
		}
		packets = append(packets, packet)
	}

	// Part 1.
	sum := 0
	for i := 0; i < len(packets)/2; i++ {
		//fmt.Printf("\n===== Pair %d:\n", i)
		if correctOrder := packets[i*2].compare(packets[i*2+1]); correctOrder {
			//fmt.Printf("Correct order!\n\n")
			sum += i + 1
		} else {
			//fmt.Printf("Wrong order!\n\n")
		}
	}

	fmt.Printf("Part 1: %d\n", sum)

	// Part 2
	// Add divider packets.
	divider1 := &Packet{
		originalContent: []rune("[[2]]"),
	}
	divider2 := &Packet{
		originalContent: []rune("[[6]]"),
	}
	packets = append(packets, divider1, divider2)

	for i := 0; i < len(packets)-1; i++ {
		for j := 0; j < len(packets)-i-1; j++ {
			if !packets[j].compare(packets[j+1]) {
				packets[j], packets[j+1] = packets[j+1], packets[j]
			}
		}
	}

	pos1, pos2 := 0, 0
	for i, packet := range packets {
		if packet == divider1 {
			pos1 = i + 1
		} else if packet == divider2 {
			pos2 = i + 1
		}
	}

	fmt.Printf("Part 2: %d\n", pos1*pos2)
}
