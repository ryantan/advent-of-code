package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
	"strconv"
)

type Packet struct {
	content         []rune
	originalContent []rune
	head            int
	currentElement  int // -1 for '[', -2 for ']'
}

func NewPacket(content string) *Packet {
	p := Packet{
		originalContent: []rune(content),
	}
	return &p
}

// convertToList converts an int to a list.
func (p *Packet) convertToList() {
	// Update underlying string to close the list.
	originalString := fmt.Sprintf("[%d]", p.currentElement)
	aIndex := p.head - (len(originalString) - 2)

	//fmt.Printf("p.content before part a: %s\n", string(p.content))
	newContent := append([]rune{}, p.content[:aIndex]...)
	newContent = append(newContent, []rune(originalString)...)
	newContent = append(newContent, p.content[p.head:]...)
	p.content = newContent
	p.head++ // Preserve position after adding the new '['.
	//fmt.Printf("p.content after replace: %s\n", string(p.content))
}

func (p *Packet) reset() {
	p.content = p.originalContent
	p.head = 0
	p.currentElement = 0
}

// readNext reads the next element, either a '[', ']' or int. Returns false if end of content.
func (p *Packet) readNext() bool {
	if p.head >= len(p.content) {
		// No more.
		return false
	}

	c := p.content[p.head]
	if c == '[' {
		p.currentElement = -1
		p.head++
	} else if c == ']' {
		p.currentElement = -2
		p.head++
	} else {
		if c == ',' {
			// Ignore
			c = p.content[p.head]
			p.head++
		}

		if c >= 48 && c <= 57 {
			// Digit, read all subsequent digits as long as they are not "," or "]"
			var number []rune
			for ; p.head < len(p.content); p.head++ {
				nextChar := p.content[p.head]
				if nextChar == ']' || nextChar == ',' {
					break
				}
				number = append(number, nextChar)
			}

			element, err := strconv.Atoi(string(number))
			if err != nil {
				panic("Could not parse number")
			}
			p.currentElement = element
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

	for {
		hasNext1 := p.readNext()
		hasNext2 := pRight.readNext()
		if !hasNext1 && !hasNext2 {
			//fmt.Printf("Both ran out\n")
			return false
		} else if !hasNext1 && hasNext2 {
			// Left ran out.
			//fmt.Printf("Left ran out\n")
			return true
		} else if hasNext1 && !hasNext2 {
			// Right ran out.
			//fmt.Printf("Right ran out\n")
			return false
		}

		// Repeatedly convert to list till we are ready to compare ints.
		for {
			//fmt.Printf("Current elements: %d | %d\n", p.currentElement, pRight.currentElement)
			if p.currentElement == -1 && pRight.currentElement >= 0 {
				// Left is a list, right is an int. Convert right to list.
				//fmt.Printf("Convert right to list\n")
				pRight.convertToList()
				p.readNext()
			} else if p.currentElement >= 0 && pRight.currentElement == -1 {
				// Right is a list, left is an int. Convert left to list.
				//fmt.Printf("Convert left to list\n")
				p.convertToList()
				pRight.readNext()
			} else {
				break
			}
		}

		if p.currentElement == -2 && pRight.currentElement != -2 {
			// Left ends the list first.
			//fmt.Printf("Left list ran out\n")
			return true
		} else if p.currentElement != -2 && pRight.currentElement == -2 {
			// Right ends the list first.
			//fmt.Printf("Right list ran out\n")
			return false
		} else if p.currentElement < pRight.currentElement {
			// Left is smaller.
			//fmt.Printf("Left is smaller\n")
			return true
		} else if p.currentElement > pRight.currentElement {
			// Right is smaller.
			//fmt.Printf("Right is smaller\n")
			return false
		}
	}
}

func main() {
	//scanner := common.GetLineScanner("../sample.txt")
	scanner := common.GetLineScanner("../input.txt")

	// Read all packets.
	packets := make([]*Packet, 0)
	for scanner.Scan() {
		l := scanner.Text()
		if len(l) == 0 {
			continue
		}
		packets = append(packets, NewPacket(l))
	}

	// Part 1.
	sum := 0
	for i := 0; i < len(packets)/2; i++ {
		if packets[i*2].compare(packets[i*2+1]) {
			sum += i + 1
		}
	}
	fmt.Printf("Part 1: %d\n", sum)

	// Part 2
	// Add divider packets.
	divider1 := NewPacket("[[2]]")
	divider2 := NewPacket("[[6]]")
	packets = append(packets, divider1, divider2)

	// Bubble sort.
	for i := 0; i < len(packets)-1; i++ {
		for j := 0; j < len(packets)-i-1; j++ {
			if !packets[j].compare(packets[j+1]) {
				packets[j], packets[j+1] = packets[j+1], packets[j]
			}
		}
	}

	// Look for pos.
	// We can keep track of positions during the sort, but at higher cost.
	decoderKey := 1
	for i, packet := range packets {
		if packet == divider1 || packet == divider2 {
			decoderKey *= i + 1
		}
	}

	fmt.Printf("Part 2: %d\n", decoderKey)
}
