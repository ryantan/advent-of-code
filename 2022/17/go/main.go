package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

type coord [2]int
type settledBlockArray [][7]bool

var blockWidths = []int{4, 3, 3, 1, 2}
var blockHeights = []int{1, 3, 3, 4, 2}

type block struct {
	blockType     int
	x             int // x of left most '#" in block.
	y             int // y of bottom most '#" in block.
	settledBlocks *settledBlockArray
}

func newBlock(blockType int, bottomY int, settledBlocks *settledBlockArray) *block {
	return &block{
		blockType:     blockType,
		x:             2,
		y:             bottomY,
		settledBlocks: settledBlocks,
	}
}

func (b *block) moveLeft() {
	//fmt.Println("Moving left.")
	b.x--
}

func (b *block) moveRight() {
	//fmt.Println("Moving right.")
	b.x++
}

func (b *block) moveDown() {
	//fmt.Println("Moving down.")
	b.y--
}

var blockParts = [][]coord{
	{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
	{{1, 0}, {0, 1}, {1, 1}, {2, 1}, {1, 2}},
	{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}},
	{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
	{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
}

// settle adds the block to settledBlocks and return new highest y.
func (b *block) settle() int {
	settledBlocks := *b.settledBlocks
	parts := blockParts[b.blockType]
	for _, part := range parts {
		settledBlocks[part[1]+b.y][part[0]+b.x] = true
	}
	for y := len(settledBlocks) - 1; y >= 0; y-- {
		for x := 0; x < 7; x++ {
			if settledBlocks[y][x] {
				return y
			}
		}
	}
	return 0
}

var leftChecks = [][]coord{
	{{-1, 0}},
	{{0, 0}, {-1, 1}, {0, 2}},
	{{-1, 0}, {1, 1}, {1, 2}},
	{{-1, 0}, {-1, 1}, {-1, 2}, {-1, 3}},
	{{-1, 0}, {-1, 1}},
}

func (b *block) canMoveLeft() bool {
	if b.x == 0 {
		return false
	}

	settledBlocks := *b.settledBlocks
	hasBlocking := false
	for _, c := range leftChecks[b.blockType] {
		if settledBlocks[b.y+c[1]][b.x+c[0]] {
			hasBlocking = true
			break
		}
	}
	return !hasBlocking
}

var rightChecks = [][]coord{
	{{4, 0}},
	{{2, 2}, {3, 1}, {2, 0}},
	{{3, 2}, {3, 1}, {3, 0}},
	{{1, 0}, {1, 1}, {1, 2}, {1, 3}},
	{{2, 0}, {2, 1}},
}

func (b *block) canMoveRight() bool {
	if b.x+b.width() > 6 {
		return false
	}

	settledBlocks := *b.settledBlocks
	hasBlocking := false
	for _, c := range rightChecks[b.blockType] {
		if settledBlocks[b.y+c[1]][b.x+c[0]] {
			hasBlocking = true
			break
		}
	}
	return !hasBlocking
}

var downChecks = [][]coord{
	{{0, -1}, {1, -1}, {2, -1}, {3, -1}},
	{{0, 0}, {1, -1}, {2, 0}},
	{{0, -1}, {1, -1}, {2, -1}},
	{{0, -1}},
	{{0, -1}, {1, -1}},
}

func (b *block) canMoveDown() bool {
	if b.y == 0 {
		return false
	}

	settledBlocks := *b.settledBlocks
	hasBlocking := false
	for _, c := range downChecks[b.blockType] {
		if settledBlocks[b.y+c[1]][b.x+c[0]] {
			hasBlocking = true
			break
		}
	}
	return !hasBlocking
}

func (b *block) height() int {
	return blockHeights[b.blockType]
}

func (b *block) width() int {
	return blockWidths[b.blockType]
}

func main() {
	//scanner := common.GetLineScanner("../sample.txt")
	scanner := common.GetLineScanner("../input.txt")

	highestY := -1
	settledBlocks := settledBlockArray{}
	scanner.Scan()
	windIndex := 0
	windDirections := []rune(scanner.Text())
	windDirectionsCount := len(windDirections)

	maxBlocks := 2022
	//maxBlocks := 3
	//maxBlocks := 12
	for i := 0; i < maxBlocks; i++ {
		startingY := highestY + 4

		// Add in additional rows if required.
		for j := len(settledBlocks); j <= startingY+4; j++ {
			settledBlocks = append(settledBlocks, [7]bool{})
		}

		b := newBlock(i%5, startingY, &settledBlocks)

		for {
			//printSettledBlocks(settledBlocks, b)

			windDirection := windDirections[windIndex%windDirectionsCount]
			windIndex++
			if windDirection == '>' {
				if b.canMoveRight() {
					b.moveRight()
				} else {
					//fmt.Printf("Cannot move right\n")
				}
			} else if windDirection == '<' {
				if b.canMoveLeft() {
					b.moveLeft()
				} else {
					//fmt.Printf("Cannot move left\n")
				}
			}

			if b.canMoveDown() {
				b.moveDown()
			} else {
				//fmt.Printf("Settling\n")
				highestY = b.settle()
				break
			}
		}

	}
	fmt.Printf("Part 1: %d\n", highestY+1)
}

func printSettledBlocks(blocksOriginal settledBlockArray, b *block) {
	blocks := make([][]string, 0)
	for _, row := range blocksOriginal {
		newRow := make([]string, 7)
		for x, cell := range row {
			if cell {
				newRow[x] = "#"
			} else {
				newRow[x] = "."
			}
		}
		blocks = append(blocks, newRow)
	}

	parts := blockParts[b.blockType]
	for _, part := range parts {
		blocks[part[1]+b.y][part[0]+b.x] = "@"
	}

	for y := len(blocks) - 1; y >= 0; y-- {
		for i := 0; i < 7; i++ {
			fmt.Print(blocks[y][i])
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}
