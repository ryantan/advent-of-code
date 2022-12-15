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

var nodes = make([]*Node, 0)
var nodesByCoord = make(map[coord]*Node, 0)
var maxSteps = 0

func a() int {
	heightMap := make([][]string, 0)
	//scanner := common.GetLineScanner("../sample.txt")
	scanner := common.GetLineScanner("../input.txt")

	var startNodeId, endNodeId int
	for scanner.Scan() {
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
	maxSteps = len(heightMap) * len(heightMap[0])

	shortestPathSteps, shortestPath, shortestPathDir := findPath(nodes[endNodeId], nodes[startNodeId])

	fmt.Printf("shortestPath: %d, shortestPathDir: %d\n", len(shortestPath), len(shortestPathDir))
	printPathOnMap(heightMap, shortestPath, shortestPathDir)

	return shortestPathSteps
}

func findPath(start, end *Node) (int, []*Node, []int) {
	var shortestPath []*Node
	var shortestPathDir []int
	shortestPathSteps := maxSteps

	queue := make([]*Node, 0)
	var currentNode *Node
	distance := make(map[int]int, 0)
	for _, node := range nodes {
		distance[node.id] = maxSteps
	}
	parent := make(map[int]*Node, 0)
	directionTo := make(map[int]int, 0)

	//printParents := func() {
	//	for id, node := range parent {
	//		fmt.Printf("%d > %d, ", node.id, id)
	//	}
	//	fmt.Printf("%d parents found.\n", len(parent))
	//}

	getPath := func(endNode *Node) ([]*Node, []int) {
		var path []*Node
		var pathDir []int

		currentNode = endNode
		for {
			//fmt.Printf("Looking for parent of %d\n", currentNode.id)
			if parentNode, exists := parent[currentNode.id]; exists {
				//fmt.Printf("%d > %d\n", parentNode.id, currentNode.id)
				path = append([]*Node{parentNode}, path...)
				pathDir = append([]int{directionTo[currentNode.id]}, pathDir...)
				currentNode = parentNode
			} else {
				//fmt.Printf("?? > %d\n", currentNode.id)
				break
			}
		}
		return path, pathDir
	}

	checkPath := func(endNode *Node) {
		path, pathDir := getPath(endNode)
		fmt.Printf("path: %v\n", path)
		if len(path) >= shortestPathSteps {
			return
		}
		shortestPath = path
		shortestPathDir = pathDir
		shortestPathSteps = len(path)
	}

	visitCount := 0

	queue = append(queue, start)
	distance[start.id] = 0

	for {
		//fmt.Printf("Queue length: %d\n", len(queue))
		if len(queue) == 0 {
			//fmt.Printf("Finished queue!\n")
			break
		}
		visitCount++

		currentNode = queue[0]
		//fmt.Printf("Visiting: %v %s\n", currentNode.pos, currentNode.label)

		//if currentNode.id == end.id {
		//	fmt.Printf("currentNode.value: %v\n", currentNode.value)
		//	//fmt.Printf("Reached end!\n")
		//	checkPath()
		//	break
		//}

		if currentNode.value == 'a' {
			//fmt.Printf("Reached an 'a'!\n")
			checkPath(currentNode)
			break
		}

		queue = queue[1:]
		//printQueue("After popping", queue)

		for d, move := range neighborsPos {
			neighbourCoord := currentNode.pos.move(move)
			neighbour, nodeExists := nodesByCoord[neighbourCoord]
			if !nodeExists {
				//fmt.Printf("%s  neighbour at %v does not exist.\n", depth, neighbourCoord)
				continue
			}

			// Check if reachable
			diff := currentNode.value - neighbour.value
			if diff > 1 {
				// Cannot reach
				continue
			}

			cost := distance[currentNode.id] + 1

			if cost < distance[neighbour.id] {
				distance[neighbour.id] = cost
				queue = append(queue, neighbour)
				parent[neighbour.id] = currentNode
				directionTo[neighbour.id] = d
				//printParents()
			}
		}
		//printQueue("After queuing", queue)
	}
	//printParents()
	//fmt.Printf("visitCount: %d\n", visitCount)

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
		//mapWithPath[node.pos[1]][node.pos[0]] = "#"
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

func printQueue(label string, queue []*Node) {
	//fmt.Printf("queue: %v\n", queue)
	fmt.Printf("%s: %d in queue: ", label, len(queue))
	for _, node := range queue {
		fmt.Printf("-> %v (%s)", node.pos, node.label)
	}
	fmt.Printf("\n")
}

func main() {
	part1 := a()
	fmt.Printf("Part1: %d\n", part1)
}
