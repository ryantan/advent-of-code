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

const typesOfBlocks = 5

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
var patternsToSkip = 0
var highestYBeforeSkipping int

// How many times we should see the patterns before we consider it the actual pattern.
const minRepetitionOfPattern = 5

// How many generations to skip before starting to check for patterns.
// First few generations don't seem to be stable and are never part of patterns.
const earliestPatternStartGeneration = 3

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

	//generationLength := 1 * windDirectionsCount
	generationLength := typesOfBlocks * windDirectionsCount
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
	var part2Answer int
	//maxBlocks := 3
	//maxBlocks := 12
	//maxBlocks := 2022
	//maxBlocks := generationLength * 40
	maxBlocks := 1_000_000_000_000
	start := time.Now()
	blocksToSimulate := 0

blockSpawnLoop:
	for i := 0; i < maxBlocks; i++ {
		// region Count down to finish simulation.
		if patternsToSkip > 0 {

			if blocksToSimulate == 0 {
				part2Answer = highestY + 1 + patternsToSkip*heightDiffInPattenLength
				break blockSpawnLoop
			}
			blocksToSimulate--
		}
		// endregion

		startingY := highestY + 4
		//fmt.Printf("startingY: %d, highestY: %d, len(settledBlocks): %d\n", startingY, highestY, len(settledBlocks))

		// Add in additional rows if required.
		for j := len(settledBlocks); j <= startingY+4; j++ {
			settledBlocks = append(settledBlocks, [7]bool{})
		}

		b := newBlock(i%5, startingY)

		// region Periodic timer
		if i%1_000_000_000 == 0 {
			//if i%1_000_000 == 0 {
			fmt.Printf("%d: (%s) h=%d ah=%d\n", i, time.Since(start), len(settledBlocks), len(settledBlocksArchive))
			start = time.Now()
			//printSettledBlocks(settledBlocks, settledBlocksArchive, b)
		}
		// endregion

		// region Optional optimization - Check next 3 wind directions if they can be collapsed.
		ahead1 := windDirections[windIndex%windDirectionsCount]
		ahead2 := windDirections[(windIndex+1)%windDirectionsCount]
		ahead3 := windDirections[(windIndex+2)%windDirectionsCount]
		aheadS := int(ahead1 + ahead2*2 + ahead3*4)
		dX := collapseDirections[aheadS][b.blockType]
		b.y -= 3
		b.x += dX
		windIndex += 3
		// endregion

		// Let wind blow till block is settled.
		for {
			//printSettledBlocks(settledBlocks, settledBlocksArchive, b)
			windDirection := windDirections[windIndex%windDirectionsCount]
			windIndex++
			if windDirection == '>' {
				if b.canMoveRight() {
					b.moveRight()
				}
			} else if windDirection == '<' {
				if b.canMoveLeft() {
					b.moveLeft()
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

		// Capture our answer for part 1.
		if i == 2021 {
			part1Answer = highestY + 1 + len(settledBlocksArchive)
		}

		if i%generationLength == 0 && patternsToSkip == 0 {
			generationHeight := highestY - lastHighestY
			generationHeights = append(generationHeights, generationHeight)
			lastHighestY = highestY
			//fmt.Printf("generationHeights: %v\n", generationHeights)

			// Find pattern. Which seems to start after 1~2 generations
			if generation < earliestPatternStartGeneration {
				// Wait for a few generations before pattern starts.
			} else if generation == earliestPatternStartGeneration {
				// Define start of pattern.
				patternStartIndex = generation
				patternStartHeight = generationHeight
				fmt.Printf("patternStartIndex: %d, patternStartHeight: %d\n", patternStartIndex, patternStartHeight)
			} else if generationHeight == patternStartHeight && patternSeenCount == 0 {
				// We see the start of pattern again.
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

				if patternSeenCount < minRepetitionOfPattern {
					fmt.Printf("Waiting for more patterns to be seen.\n")
				} else {
					// Pretty confident we have found the patten.
					// But wait till we get answer for part 1.
					if i < 2022 {
						fmt.Printf("Continue for part 1 for now.\n")
					} else {
						fmt.Printf("i: %d\n", i)
						blocksLeft := maxBlocks - i - 1
						fmt.Printf("blocksLeft: %d\n", blocksLeft)
						fmt.Printf("patternLength: %d\n", patternLength)
						fmt.Printf("generationLength: %d\n", generationLength)
						fmt.Printf("heightDiffInPattenLength: %d\n", heightDiffInPattenLength)
						highestYBeforeSkipping = highestY
						fmt.Printf("highestYBeforeSkipping: %d\n", highestYBeforeSkipping)

						// If we're lucky, we don't need to simulate any remaining blocks when blocksToSimulate=0.
						blocksToSimulate = blocksLeft % (generationLength * patternLength)
						fmt.Printf("blocksToSimulate: %d\n", blocksToSimulate)

						blocksToSkip := blocksLeft - blocksToSimulate
						fmt.Printf("blocksToSkip: %d\n", blocksToSkip)

						patternsToSkip = blocksToSkip / generationLength / patternLength
						fmt.Printf("patternsToSkip: %d\n", patternsToSkip)
					}
				}
			}

			lastGenerationHeight = generationHeight
			generation++
		}
	}

	fmt.Printf("Total time spent: %s\n", time.Since(start))

	fmt.Printf("Part 1: %d\n", part1Answer)
	fmt.Printf("Part 2: %d\n", part2Answer)

	// Wrong answers:
	// 1540634005739 is too low.
	// 1540634005740 is too low.
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
