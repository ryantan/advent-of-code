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
	id    int
	pos   coord
	label string
	value rune
}

var height = 0
var width = 0
var nodes = make([]*Node, 0)
var nodesByCoord = make(map[coord]*Node, 0)
var maxSteps = 0

func main() {
	//scanner := common.GetLineScanner("../sample.txt")
	scanner := common.GetLineScanner("../input.txt")

	var startNodeId, endNodeId int
	for scanner.Scan() {
		l := strings.Split(scanner.Text(), "")

		for x, char := range l {
			node := Node{
				id:    len(nodes),
				pos:   coord{x, height},
				label: char,
				value: []rune(char)[0],
			}
			if char == "E" {
				endNodeId = node.id
				node.value = 'z'
			}
			if char == "S" {
				startNodeId = node.id
				node.value = 'a'
			}
			nodes = append(nodes, &node)
			nodesByCoord[node.pos] = &node
		}
		height++
	}
	maxSteps = len(nodes)
	width = maxSteps / height

	shortestPathToS, shortestPathToSDir, shortestPathToA, shortestPathToADir := findPath(nodes[endNodeId], nodes[startNodeId])

	printPathOnMap(height, width, shortestPathToS, shortestPathToSDir)
	printPathOnMap(height, width, shortestPathToA, shortestPathToADir)

	fmt.Printf("Part1: %d\n", len(shortestPathToS))
	fmt.Printf("Part2: %d\n", len(shortestPathToA))
}

func findPath(start, end *Node) ([]*Node, []int, []*Node, []int) {
	var shortestPathToS []*Node
	var shortestPathToSDir []int // Only required for nice path rendering.
	var shortestPath []*Node
	var shortestPathDir []int // Only required for nice path rendering.
	shortestPathSteps := maxSteps

	queue := make([]*Node, 0)
	var currentNode *Node

	// Initialize distance to some large value. We use max steps in this case.
	distance := make(map[int]int, 0)
	for _, node := range nodes {
		distance[node.id] = maxSteps
	}
	parent := make(map[int]*Node, 0)
	directionTo := make(map[int]int, 0) // Only required for nice path rendering.

	getPath := func(endNode *Node) ([]*Node, []int) {
		var path []*Node
		var pathDir []int
		for parentNode := parent[endNode.id]; parentNode != nil; parentNode = parent[parentNode.id] {
			path = append([]*Node{parentNode}, path...)
			pathDir = append([]int{directionTo[parentNode.id]}, pathDir...)
		}
		return path, pathDir
	}

	checkPath := func(node *Node) {
		path, pathDir := getPath(node)

		// Keep answer for part 1.
		if node.id == end.id {
			shortestPathToS = path
			shortestPathToSDir = pathDir // Only required for nice path rendering.
		}

		// Continue for part 2.
		if len(path) >= shortestPathSteps {
			return
		}
		shortestPath = path
		shortestPathDir = pathDir // Only required for nice path rendering.
		shortestPathSteps = len(path)
	}

	queue = append(queue, start)
	distance[start.id] = 0

	for len(queue) > 0 {
		currentNode = queue[0]
		if currentNode.value == 'a' {
			checkPath(currentNode)
			// We can technically break if we are only interested in part 1.
			//break
		}

		// Dequeue first item.
		queue = queue[1:]

		// Check neighbours.
		for d, move := range neighborsPos {
			// Check if node exists on map.
			neighbour, nodeExists := nodesByCoord[currentNode.pos.move(move)]
			if !nodeExists {
				continue
			}

			// Check if neighbour can be reached.
			// Take note the logic is reversed as we're going from End to start.
			if diff := currentNode.value - neighbour.value; diff > 1 {
				continue
			}

			cost := distance[currentNode.id] + 1
			if cost < distance[neighbour.id] {
				distance[neighbour.id] = cost
				queue = append(queue, neighbour)
				parent[neighbour.id] = currentNode
				directionTo[neighbour.id] = d // Only required for nice path rendering.
			}
		}
	}
	return shortestPathToS, shortestPathToSDir, shortestPath, shortestPathDir
}

func printPath(path []*Node) {
	fmt.Printf("%d steps ", len(path))
	for _, node := range path {
		fmt.Printf("-> %v (%s)", node.pos, node.label)
	}
	fmt.Printf("\n")
}

var dirChar = []string{"<", "^", ">", "v"}

func printPathOnMap(height, width int, path []*Node, pathDir []int) {
	fmt.Printf("\n\n=== Path ===\n")
	//printPath(path)

	var mapWithPath [][]string
	for y := 0; y < height; y++ {
		row := make([]string, 0)
		for x := 0; x < width; x++ {
			row = append(row, ".")
		}
		mapWithPath = append(mapWithPath, row)
	}
	for i, node := range path {
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
