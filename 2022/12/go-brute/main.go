package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
	"strings"
)

type coord [2]int

func (c coord) move(delta coord) coord {
	return coord{c[0] + delta[0], c[1] + delta[1]}
}

var neighborsPos = []coord{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}

type Node struct {
	id      int
	pos     coord
	label   string
	value   rune
	visited bool
	isStart bool
	isEnd   bool
}

type Edge struct {
	source *Node
	target *Node
	weight int
}

var nodes = make([]*Node, 0)
var nodesByCoord = make(map[coord]*Node, 0)

func a() int {
	heightMap := make([][]string, 0)
	//scanner := common.GetLineScanner("../sample.txt")
	//scanner := common.GetLineScanner("../input.txt")
	//scanner := common.GetLineScanner("../input2.txt")
	//scanner := common.GetLineScanner("../input3.txt")
	scanner := common.GetLineScanner("../input4.txt")
	//scanner := common.GetLineScanner("../input5.txt")
	//scanner := common.GetLineScanner("../input6.txt")
	//scanner := common.GetLineScanner("../input7.txt")

	var startNodeId, endNodeId int
	for scanner.Scan() {
		//l := []rune(scanner.Text())
		l := strings.Split(scanner.Text(), "")

		row := make([]string, 0)
		for _, char := range l {

			node := Node{
				id:    len(nodes),
				pos:   coord{len(row), len(heightMap)},
				label: char,
				value: []rune(char)[0],
			}

			if char == "E" {
				endNodeId = node.id
				node.value = 'z'
				//node.value = 'd'
				node.isEnd = true
			}
			if char == "S" {
				startNodeId = node.id
				node.value = 'a'
				node.isStart = true
			}
			nodes = append(nodes, &node)
			nodesByCoord[node.pos] = &node

			row = append(row, char)
		}
		heightMap = append(heightMap, row)
	}
	//fmt.Printf("map: %v\nsCoord: %v\neCoord: %v\n", heightMap, nodes[startNodeId].pos, nodes[endNodeId].pos)
	//fmt.Printf("nodesByCoord: %v\n", nodesByCoord)

	shortestPathSteps, shortestPath, shortestPathDir := findPath(heightMap, nodes[startNodeId], nodes[endNodeId])

	fmt.Printf("shortestPath: %d, shortestPathDir: %d\n", len(shortestPath), len(shortestPathDir))
	printPathOnMap(heightMap, shortestPath, shortestPathDir)

	return shortestPathSteps
}

func findPath(heightMap [][]string, start, end *Node) (int, []*Node, []int) {

	shortestPathSteps := len(heightMap) * len(heightMap[0])
	var shortestPath []*Node
	var shortestPathDir []int

	var internal func(node *Node, path []*Node, pathHash map[int]bool, pathDir []int, deadEnds map[int]bool, dir int, depth string) bool
	internal = func(node *Node, path []*Node, pathHash map[int]bool, pathDir []int, deadEnds map[int]bool, dir int, depth string) bool {

		//fmt.Printf("%sVisit %v %s (%d)\n", depth, node.pos, node.label, node.value)

		path = append(path, node)
		pathHash[node.id] = true
		if dir != -1 {
			// Ignore first.
			pathDir[len(pathDir)-1] = dir
		}
		pathDir = append(pathDir, dir)

		if node.id == end.id {
			//fmt.Printf("Reached the end!\n")
			pathDir[len(pathDir)-1] = -1
			printPath(path)
			if len(path) < shortestPathSteps {
				shortestPathSteps = len(path)
				shortestPath = path
				shortestPathDir = pathDir
			}
			return false
		}

		hasNonDeadEnds := false
		for d, move := range neighborsPos {
			neighbourCoord := node.pos.move(move)
			neighbour, nodeExists := nodesByCoord[neighbourCoord]
			if !nodeExists {
				//fmt.Printf("%s  neighbour at %v does not exist.\n", depth, neighbourCoord)
				continue
			}

			if pathHash[neighbour.id] == true {
				// Path intersected, dead end.
				//fmt.Printf("%s  neighbour at %v is already in path.\n", depth, neighbourCoord)
				continue
			}
			if deadEnds[neighbour.id] {
				// Dead end.
				//fmt.Printf("%s  neighbour at %v is a deadend.\n", depth, neighbourCoord)
				continue
			}

			// Check if reachable
			diff := neighbour.value - node.value
			if diff > 1 {
				// Cannot reach
				//fmt.Printf("%s  neighbour at %v %s cannot be reached, diff: %d\n", depth, neighbour.pos, neighbour.label, diff)
				continue
			}

			// Can reach
			//fmt.Printf("%s  neighbour at %v %s looks ok, diff: %d\n", depth, neighbour.pos, neighbour.label, diff)

			// Copy
			newPath := make([]*Node, len(path))
			copy(newPath, path)

			newPathDir := make([]int, len(pathDir))
			copy(newPathDir, pathDir)

			// Copy hash.
			newPathHash := map[int]bool{}
			for k, v := range pathHash {
				newPathHash[k] = v
			}

			// Copy dead ends.
			newDeadEnds := map[int]bool{}
			for k, v := range deadEnds {
				newDeadEnds[k] = v
			}

			isDeadEnd := internal(neighbour, newPath, newPathHash, newPathDir, newDeadEnds, d, depth+"  ")
			if !isDeadEnd {
				hasNonDeadEnds = true
			}
		}
		if !hasNonDeadEnds {
			deadEnds[node.id] = true
			// All neighbours are dead ends.
			return true
		}
		return false
	}

	internal(start, []*Node{}, make(map[int]bool, 0), []int{}, make(map[int]bool, 0), -1, "")

	return shortestPathSteps, shortestPath, shortestPathDir
}

func printPath(path []*Node) {
	fmt.Printf("%d steps ", len(path))
	for _, node := range path {
		fmt.Printf("-> %v (%s)", node.pos, node.label)
	}
	fmt.Printf("\n")
}

//var dirChar = []string{">", "V", "<", "^"}

var dirChar = []string{"<", "^", ">", "v"}

func printPathOnMap(heightMap [][]string, path []*Node, pathDir []int) {
	fmt.Printf("\n\n=== Path ===\n")
	printPath(path)

	var mapWithPath [][]string
	for y := 0; y < len(heightMap); y++ {
		row := make([]string, 0)
		for x := 0; x < len(heightMap[0]); x++ {
			row = append(row, ".")
		}
		mapWithPath = append(mapWithPath, row)
	}
	fmt.Printf("\nlen(path)=%d, len(pathDir)=%d\n\n", len(path), len(pathDir))
	for i, node := range path {
		//fmt.Printf("node.pos: %v, pathDir[i]: %v\n", node.pos, pathDir[i])
		if pathDir[i] != -1 {
			mapWithPath[node.pos[1]][node.pos[0]] = dirChar[pathDir[i]]
		} else {
			// End node.
			mapWithPath[node.pos[1]][node.pos[0]] = "E"
		}
	}
	printMap(mapWithPath)
}

func printMap(mapArray [][]string) {
	for y := 0; y < len(mapArray); y++ {
		for x := 0; x < len(mapArray[0]); x++ {
			fmt.Printf("%s", mapArray[y][x])
		}
		fmt.Printf("\n")
	}
}

func main() {
	part1 := a()
	fmt.Printf("Part1: %d\n", part1)
}
