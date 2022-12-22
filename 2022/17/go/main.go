package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
	"time"
)

type coord [2]int
type settledBlockArray [][7]bool

var blockWidths = []int{4, 3, 3, 1, 2}
var blockHeights = []int{1, 3, 3, 4, 2}

type block struct {
	blockType int
	x         int // x of left most '#" in block.
	y         int // y of bottom most '#" in block.
}

func newBlock(blockType int, bottomY int) *block {
	return &block{
		blockType: blockType,
		x:         2,
		y:         bottomY,
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
func (b *block) settle() {
	parts := blockParts[b.blockType]
	for _, part := range parts {
		settledBlocks[part[1]+b.y][part[0]+b.x] = true
	}
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

//var collapseDirections = map[int][5]int{
//	180: {-2, -2, -2, -2, -2},
//	182: {-1, -1, -1, -1, -1},
//	184: {1, 1, 1, 1, 1},
//	186: {1, 2, 2, 3, 3},
//}

var collapseDirections = map[int][5]int{
	420: {-2, -2, -2, -2, -2}, // <<< // 60 + 2*60 + 4*60
	428: {-1, -1, -1, -1, -1}, // <<> // 60 + 2*60 + 4*62
	424: {-1, -1, -1, -1, -1}, // <>< // 60 + 2*62 + 4*60
	432: {1, 1, 1, 1, 1},      // <>> // 60 + 2*62 + 4*62
	422: {-1, -1, -1, -1, -1}, // ><< // 62 + 2*60 + 4*60
	430: {1, 1, 1, 1, 1},      // ><> // 62 + 2*60 + 4*62
	426: {0, 1, 1, 1, 1},      // >>< // 62 + 2*62 + 4*60
	434: {1, 2, 2, 3, 3},      // >>> // 62 + 2*62 + 4*62
}

var settledBlocks = settledBlockArray{}
var lastGenerationHeight = 0
var lastHighestY = 0
var generationHeights []int
var generation int
var patternStartIndex = 0
var patternStartHeight = 0
var patternEndIndex = 0
var patternEndHeight = 0
var patternSeenIndex = 0
var prevPatternSeenIndex = 0
var patternSeenCount = 0
var patternLength = 0
var heightDiffInPattenLength = 0
var extrapolatedHeight float64

func isAMatch(startIndex, testIndex, length int) bool {
	matches := true
	heightDiffInPattenLength = 0
	for p := 0; p < patternLength; p++ {
		heightDiffInPattenLength += generationHeights[startIndex+p]
		if generationHeights[startIndex+p] != generationHeights[testIndex+p] {
			matches = false
		}
	}
	fmt.Printf("heightDiffInPattenLength: %d\n", heightDiffInPattenLength)
	return matches
}

func main() {
	//scanner := common.GetLineScanner("../sample.txt")
	scanner := common.GetLineScanner("../input.txt")

	highestY := -1
	//settledBlocks := settledBlockArray{}
	settledBlocksArchive := settledBlockArray{}
	scanner.Scan()
	windIndex := 0
	windDirections := []rune(scanner.Text())
	windDirectionsCount := len(windDirections)
	fmt.Printf("windDirectionsCount: %d\n", windDirectionsCount)

	generationLength := 5 * windDirectionsCount
	fmt.Printf("generationLength: %d\n", generationLength)

	updateHighestY := func(minY int) {
		for j := len(settledBlocks) - 1; j >= minY; j-- {
			for x := 0; x < 7; x++ {
				if settledBlocks[j][x] {
					highestY = j
					//fmt.Printf("highestY updated to %d\n", highestY)
					return
				}
			}
		}
	}

	var part1Answer int
	//maxBlocks := 2022
	//maxBlocks := generationLength * 40
	maxBlocks := 1_000_000_000_000
	earliestPatternStartGeneration := 3
	//earliestPatternStartGeneration := 100
	//maxBlocks := 3
	//maxBlocks := 12
	start := time.Now()

blockSpawnLoop:
	for i := 0; i < maxBlocks; i++ {
		startingY := highestY + 4
		//fmt.Printf("startingY: %d, highestY: %d, len(settledBlocks): %d\n", startingY, highestY, len(settledBlocks))

		// Add in additional rows if required.
		for j := len(settledBlocks); j <= startingY+4; j++ {
			settledBlocks = append(settledBlocks, [7]bool{})
		}

		b := newBlock(i%5, startingY)

		if i%1_000_000 == 0 {
			//if i%1_000_000 == 0 {
			fmt.Printf("%d: (%s) h=%d ah=%d\n", i, time.Since(start), len(settledBlocks), len(settledBlocksArchive))
			start = time.Now()
			//printSettledBlocks(settledBlocks, settledBlocksArchive, b)
		}

		// Check next 3 wind directions if they can be collapsed.
		//ahead := string(windDirections[windIndex%windDirectionsCount : (windIndex+3)%windDirectionsCount])
		ahead1 := windDirections[windIndex%windDirectionsCount]
		ahead2 := windDirections[(windIndex+1)%windDirectionsCount]
		ahead3 := windDirections[(windIndex+2)%windDirectionsCount]
		aheadS := int(ahead1 + ahead2*2 + ahead3*4)
		dX := collapseDirections[aheadS][b.blockType]
		b.y -= 3
		b.x += dX
		windIndex += 3

		for {
			//printSettledBlocks(settledBlocks, settledBlocksArchive, b)
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
				b.settle()
				updateHighestY(b.y)
				break
			}
		}
		if i == 2021 {
			part1Answer = highestY + 1 + len(settledBlocksArchive)
		}

		if i%generationLength == 0 {
			generationHeight := highestY - lastHighestY
			generationHeights = append(generationHeights, generationHeight)
			lastHighestY = highestY
			//fmt.Printf("generationHeights: %v\n", generationHeights)

			// Find pattern. Which seems to start after 1~2 generations
			if generation < earliestPatternStartGeneration {
				// Wait for a few generation before pattern starts.
			} else if generation == earliestPatternStartGeneration {
				patternStartIndex = generation
				patternStartHeight = generationHeight
				fmt.Printf("patternStartIndex: %d, patternStartHeight: %d\n", patternStartIndex, patternStartHeight)
			} else if generationHeight == patternStartHeight && patternSeenCount == 0 {
				if patternSeenIndex == 0 {
					patternSeenIndex = generation
					patternLength = patternSeenIndex - patternStartIndex
					fmt.Printf("patternSeenIndex: %d\n", patternSeenIndex)
				} else {
					// Verifying if there is a match.
					patternLength = patternSeenIndex - patternStartIndex
					currentLength := generation - patternSeenIndex
					if currentLength == patternLength {
						if isAMatch(patternStartIndex, patternSeenIndex, patternLength) {
							patternEndIndex = patternSeenIndex
							patternEndHeight = lastGenerationHeight
							patternSeenIndex = generation
							patternSeenCount = 2
							fmt.Printf("Pattern found at %d and repeated at %d\n", patternStartIndex, patternEndIndex)
							fmt.Printf("patternStartHeight: %d patternEndHeight: %d\n", patternStartHeight, patternEndHeight)
						} else {
							// Unset patternSeenIndex.
							fmt.Printf("%d to %d is not a patten\n", patternStartIndex, patternSeenIndex)
							patternSeenIndex = generation
							fmt.Printf("Set patternSeenIndex to current generation: %d\n", patternSeenIndex)
						}
					} else if currentLength > patternLength {
						// Unset patternSeenIndex.
						fmt.Printf("%d to %d is not a patten and no chance a next token will make it a pattern\n", patternStartIndex, patternSeenIndex)
						patternSeenIndex = prevPatternSeenIndex
						prevPatternSeenIndex = generation
						fmt.Printf("Set patternSeenIndex to prevPatternSeenIndex: %d\n", patternSeenIndex)
						fmt.Printf("New prevPatternSeenIndex: %d\n", prevPatternSeenIndex)
					} else {
						// Ignore.
						fmt.Printf("%d to %d is not a patten but there's a chance a next token will make it a pattern\n", patternStartIndex, patternSeenIndex)
						prevPatternSeenIndex = generation
						fmt.Printf("New prevPatternSeenIndex: %d\n", prevPatternSeenIndex)
						fmt.Printf("New patternSeenIndex: %d\n", patternSeenIndex)
					}
				}
			} else if generationHeight == patternEndHeight && patternSeenCount > 0 && generation-patternSeenIndex+1 == patternLength {
				// Measuring confidence of pattern.
				// We define pattern as high confidence when patternSeenCount > 5
				if isAMatch(patternStartIndex, patternSeenIndex, patternLength) {
					patternSeenCount++
					fmt.Printf("Pattern found again at %d and repeated %d times\n", patternSeenIndex, patternSeenCount)
					patternSeenIndex = generation + 1
				} else {
					// Reset
					patternEndIndex = 0
					patternEndHeight = 0
					patternSeenIndex = 0
					patternSeenCount = 0
					patternLength = 0
				}
				if patternSeenCount >= 5 {
					// Pretty confident we have found the patten.
					// Extrapolate results now.
					blocksLeftToSimulate := maxBlocks - i
					fmt.Printf("blocksLeftToSimulate: %d\n", blocksLeftToSimulate)
					fmt.Printf("generationLength: %d\n", generationLength)
					currentHeight := highestY + 1
					fmt.Printf("currentHeight: %d\n", currentHeight)
					fmt.Printf("heightDiffInPattenLength: %d\n", heightDiffInPattenLength)
					generationsLeftToSimulate := float64(blocksLeftToSimulate) / float64(generationLength)
					fmt.Printf("generationsLeftToSimulate: %f\n", generationsLeftToSimulate)
					fmt.Printf("patternLength: %d\n", patternLength)
					patternsLeftToSimulate := float64(generationsLeftToSimulate) / float64(patternLength)
					fmt.Printf("patternsLeftToSimulate: %f\n", patternsLeftToSimulate)
					extrapolatedHeight = float64(highestY) + (patternsLeftToSimulate * float64(heightDiffInPattenLength))
					fmt.Printf("extrapolatedHeight: %f\n", extrapolatedHeight)

					// Wait till we get answer for part 1.
					if i >= 2022 {
						break blockSpawnLoop
					}
					fmt.Printf("=== Continuing for part 1\n")
				}
			}

			lastGenerationHeight = generationHeight
			generation++
		}

		// Settled.

		//if i%1000 == 0 {
		//	//archiveStart := time.Now()
		//	unbrokenLineY := 0
		//	// Check if there's an unbroken line formed, if yes, ignore
		//	// them in subsequent simulations.
		//	// Only look around y + height of b.
		//	for j := len(settledBlocks) - 2; j > 0; j-- {
		//		hasUnbrokenLine := true
		//
		//		//fmt.Printf("settledBlocks[j]: %v\n", settledBlocks[j])
		//		// Check 2 rows at once.
		//		for k := 0; k < 7; k++ {
		//			if !settledBlocks[j][k] && !settledBlocks[j+1][k] {
		//				hasUnbrokenLine = false
		//				break
		//			}
		//		}
		//		if hasUnbrokenLine {
		//			unbrokenLineY = j
		//			break
		//		}
		//	}
		//	if unbrokenLineY > 0 {
		//		//fmt.Printf("unbrokenLineY: %d\n", unbrokenLineY)
		//		// Move 0 to unbrokenLineY-1 into archive
		//		toMove := settledBlocks[0:unbrokenLineY]
		//		settledBlocksArchive = append(settledBlocksArchive, toMove...)
		//		settledBlocks = settledBlocks[unbrokenLineY:]
		//		updateHighestY(0)
		//	}
		//	//fmt.Printf("Archiving %d lines took: %s, new archive len=%d\n", unbrokenLineY, time.Since(archiveStart), len(settledBlocksArchive))
		//}

	}
	fmt.Printf("Part 1: %d\n", part1Answer)
	fmt.Printf("Part 2: %f\n", extrapolatedHeight)
}

func printSettledBlocks(blocksOriginal settledBlockArray, archive settledBlockArray, b *block) {
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
