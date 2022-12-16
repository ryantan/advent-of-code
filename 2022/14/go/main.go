package main

import (
	"fmt"
	"strings"
)

import "github/ryantan/advent-of-code/2022/common"

type coord [2]int

func newCoord(s string) *coord {
	x, y := 0, 0
	_, _ = fmt.Sscanf(s, "%d,%d", &x, &y)
	c := coord{x, y}
	return &c
}

type Sand struct {
	pos          coord
	lowestRockY  int
	floor        int
	rocks        []*coord
	stuffByCoord map[coord]bool
}

// newSand produces sand at 500, 0
func newSand(rocks []*coord, stuffByCoord map[coord]bool, maxY int, floor int) *Sand {
	sand := Sand{
		pos:          coord{500, 0},
		lowestRockY:  maxY,
		floor:        floor,
		rocks:        rocks,
		stuffByCoord: stuffByCoord,
	}
	return &sand
}

// drop moves the sand down till it settles, and returns true when it can still move.
// Returns false, false when settled, or false, true when it went past the lowest rock.
func (s *Sand) drop() (bool, bool) {
	if s.stuffByCoord[s.pos] {
		return false, true
	}
	oldPos := s.pos
	//fmt.Printf("Original pos: %v\n", s.pos)
	s.pos[1]++

	if s.pos[1] >= s.floor {
		s.pos[1]--
	} else if s.stuffByCoord[s.pos] {
		// Something is below s.

		// Try left down.
		s.pos[0]--
		if !s.stuffByCoord[s.pos] {
			// Can go left down, stay.
		} else {
			// Cannot go left down, try right down.
			s.pos[0] += 2
			if !s.stuffByCoord[s.pos] {
				// Can go right down, stay.
			} else {
				// Cannot go right down, back.
				s.pos[0]--
				s.pos[1]--
				fmt.Printf("Cannot go right down, back, s: %v\n", s.pos)
			}
		}

	} else {
		// Can go down.
	}

	lost := false
	// Went pass lowest rock? Lost!
	if s.pos[1] >= s.lowestRockY {
		lost = true
	}

	if oldPos == s.pos {
		// Didn't move? Settled!
		return false, lost
	}

	// Moved!
	return true, lost
}

func main() {
	//scanner := common.GetLineScanner("../sample.txt")
	scanner := common.GetLineScanner("../input.txt")

	rocks := make([]*coord, 0)

	for scanner.Scan() {
		l := strings.Split(scanner.Text(), " -> ")
		//fmt.Printf("l: %v\n", l)
		for i := 0; i < len(l)-1; i++ {
			start := newCoord(l[i])
			end := newCoord(l[i+1])

			dx, dy := end[0]-start[0], end[1]-start[1]
			//fmt.Printf("dx: %d, dy: %d\n", dx, dy)
			if dx*dx > dy*dy {
				startX, endX := start[0], start[0]+dx
				if startX > endX {
					startX, endX = endX, startX
				}
				for x := startX; x <= endX; x++ {
					rocks = append(rocks, &coord{x, start[1]})
				}
			} else {
				startY, endY := start[1], start[1]+dy
				if startY > endY {
					startY, endY = endY, startY
				}
				for y := startY; y <= endY; y++ {
					rocks = append(rocks, &coord{start[0], y})
				}
			}
			//fmt.Printf("Next rock line\n")
		}
		//fmt.Printf("Next rock formation\n")
	}

	stuffByCoord := make(map[coord]bool, 0)
	minX, maxX, maxY := 10000000, 0, 0
	for _, rock := range rocks {
		if rock[1] > maxY {
			maxY = rock[1]
		}
		if rock[0] > maxX {
			maxX = rock[0]
		}
		if rock[0] < minX {
			minX = rock[0]
		}
		stuffByCoord[*rock] = true
	}

	part1 := 0
	sand := make([]coord, 0)
	for {
		s := newSand(rocks, stuffByCoord, maxY, maxY+2)
		for j := 0; j < maxY+2; j++ {
			canMove, lost := s.drop()
			if !canMove {
				if lost && part1 == 0 {
					part1 = len(sand)
				}
				sand = append(sand, s.pos)
				stuffByCoord[s.pos] = true
				break
			}
		}
		if s.stuffByCoord[coord{500, 0}] {
			break
		}
	}
	//fmt.Printf("lostCount: %d\n", lostCount)
	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", len(sand))

	printMap(rocks, sand, minX, maxX, maxY)
}

//func printRocks(rocks []*coord) {
//	fmt.Printf("%d Rocks: %v\n", len(rocks), rocks)
//	for _, rock := range rocks {
//		fmt.Printf("Rock: %v\n", rock)
//	}
//}

func printMap(rocks []*coord, sand []coord, minX, maxX, maxY int) {
	var grid [][]string

	// Expand min max X with sand.
	for _, s := range sand {
		if s[0] > maxX {
			maxX = s[0]
		}
		if s[0] < minX {
			minX = s[0]
		}
	}
	minX, maxX = minX-2, maxX+2

	// Fill grid with air.
	for y := 0; y <= maxY; y++ {
		var row []string
		for x := minX; x <= maxX; x++ {
			row = append(row, ".")
		}
		grid = append(grid, row)
	}

	// Add Sand.
	fmt.Printf("x: %d-%d, y: 0-%d\n", minX, maxX, maxY)
	for _, s := range sand {
		if s[1] < maxY && s[0] > minX && s[0] < maxX {
			grid[s[1]][s[0]-minX] = "o"
		}
	}

	// Add rocks.
	for _, r := range rocks {
		grid[r[1]][r[0]-minX] = "#"
	}

	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX-minX; x++ {
			print(grid[y][x])
		}
		print("\n")
	}
	print("\n")
}
