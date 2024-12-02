package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
	"time"
)

//var actions = []int{
//	0
//}

// Actions:
// 0 - Do nothing
// 1 - Build ore robot
// 2 - Build clay robot
// 3 - Build obsidian robot
// 4 - Build geode robot

func newBluePrint(l string) BluePrint {
	var bluePrintId, oreForOreRobot, oreForClayRobot, oreForObsidianRobot, clayForObsidianRobot, oreForGeodeRobot, obsidianForGeodeRobot int
	_, _ = fmt.Sscanf(l, "Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.", &bluePrintId, &oreForOreRobot, &oreForClayRobot, &oreForObsidianRobot, &clayForObsidianRobot, &oreForGeodeRobot, &obsidianForGeodeRobot)

	bluePrint := BluePrint{
		bluePrintId:           bluePrintId,
		oreForOreRobot:        oreForOreRobot,
		oreForClayRobot:       oreForClayRobot,
		oreForObsidianRobot:   oreForObsidianRobot,
		clayForObsidianRobot:  clayForObsidianRobot,
		oreForGeodeRobot:      oreForGeodeRobot,
		obsidianForGeodeRobot: obsidianForGeodeRobot,
	}
	return bluePrint
}

type BluePrint struct {
	bluePrintId           int
	oreForOreRobot        int
	oreForClayRobot       int
	oreForObsidianRobot   int
	clayForObsidianRobot  int
	oreForGeodeRobot      int
	obsidianForGeodeRobot int
}

type State struct {
	ore           int
	clay          int
	obsidian      int
	geode         int
	oreRobot      int
	clayRobot     int
	obsidianRobot int
	geodeRobot    int
}

func (s State) buildOreRobot(bluePrint *BluePrint) State {
	s.ore -= bluePrint.oreForOreRobot
	s.oreRobot++
	return s
}

func (s State) buildClayRobot(bluePrint *BluePrint) State {
	s.ore -= bluePrint.oreForClayRobot
	s.clayRobot++
	return s
}

func (s State) buildObsidianRobot(bluePrint *BluePrint) State {
	s.ore -= bluePrint.oreForObsidianRobot
	s.clay -= bluePrint.clayForObsidianRobot
	s.obsidianRobot++
	return s
}

func (s State) buildGeodeRobot(bluePrint *BluePrint) State {
	s.ore -= bluePrint.oreForGeodeRobot
	s.obsidian -= bluePrint.obsidianForGeodeRobot
	s.geodeRobot++
	return s
}

type timeAndId struct {
	time int
	id   int
}

var earliestTimesForState = map[State]timeAndId{}
var nodeIdsToAbort = map[int]bool{}
var endNodesByState = map[State]int{}
var foundInEarliestTimesForState = 0

type Node struct {
	id       int
	parent   *Node
	action   int
	time     int
	state    State
	children map[int]*Node
}

func (n *Node) hasAncestor(id int) bool {
	currentNode := n
	for {
		if currentNode.id == id {
			return true
		}
		currentNode = currentNode.parent
		if currentNode == nil {
			break
		}
	}
	return false
}

var aborted = 0
var nodesEvaluated = 0
var nodesByTime = map[int]int{}
var endNodes = make([]*Node, 0)

func enoughTimeToBuildObsidianRobot(timeLeft int, state State, bluePrint *BluePrint) bool {
	if timeLeft <= 0 {
		return false
	}
	if state.ore >= bluePrint.oreForObsidianRobot && state.clay >= bluePrint.clayForObsidianRobot {
		return true
	}
	if state.ore+(timeLeft*state.oreRobot) < bluePrint.oreForObsidianRobot {
		return false
	}
	if state.clay+(timeLeft*state.clayRobot) < bluePrint.clayForObsidianRobot {
		return false
	}

	return true
}
func enoughTimeToBuildGeodeRobot(timeLeft int, state State, bluePrint *BluePrint) bool {
	if timeLeft <= 0 {
		return false
	}
	if state.ore >= bluePrint.oreForGeodeRobot && state.obsidianRobot >= bluePrint.obsidianForGeodeRobot {
		return true
	}
	if state.ore+(timeLeft*state.oreRobot) < bluePrint.oreForGeodeRobot {
		return false
	}
	if state.obsidian+(timeLeft*state.obsidianRobot) < bluePrint.obsidianForGeodeRobot {
		return false
	}

	return true
}

func expandNode(parent *Node, bluePrint *BluePrint, nodeId *int, maxTime int) {
	for id, _ := range nodeIdsToAbort {
		if parent.hasAncestor(id) {
			nodeIdsToAbort[id] = false
			return
		}
	}

	newTime := parent.time + 1

	addNode := func(action int, state State) {
		newNode := &Node{
			id:       *nodeId,
			action:   action,
			parent:   parent,
			state:    state,
			time:     newTime,
			children: map[int]*Node{},
		}
		*nodeId++
		nodesEvaluated++
		if _, exists := nodesByTime[newTime]; !exists {
			nodesByTime[newTime] = 1
		} else {
			nodesByTime[newTime]++
		}

		// If this timing was beaten before, kill this branch.
		if earliestTimeForState, exists := earliestTimesForState[state]; exists {
			foundInEarliestTimesForState++
			if newTime >= earliestTimeForState.time {
				aborted++
				return
			} else {
				// Replace
				nodeIdsToAbort[earliestTimeForState.id] = true
				earliestTimesForState[state] = timeAndId{time: newTime, id: newNode.id}
			}
		} else {
			earliestTimesForState[state] = timeAndId{time: newTime, id: newNode.id}
		}

		if newTime == maxTime {
			fmt.Printf("endNodes length: %d\n", len(endNodes))
			endNodes = append(endNodes, newNode)
		} else {
			expandNode(newNode, bluePrint, nodeId, maxTime)
		}
		fmt.Printf("newNode (%d): %+v\n", *nodeId, newNode)
		parent.children[newNode.action] = newNode
	}

	timeLeft := maxTime - newTime

	if timeLeft == 5 && parent.state.obsidianRobot == 0 && !enoughTimeToBuildObsidianRobot(timeLeft-2, parent.state, bluePrint) {
		// We can give up this branch.
		aborted++
		return
	}
	if timeLeft == 5 && parent.state.geodeRobot == 0 && !enoughTimeToBuildGeodeRobot(timeLeft, parent.state, bluePrint) {
		// We can give up this branch.
		aborted++
		return
	}
	if timeLeft == 4 && parent.state.geodeRobot == 0 && parent.state.obsidianRobot == 0 {
		// We can give up this branch.
		aborted++
		return
	}
	if timeLeft == 2 && parent.state.geodeRobot == 0 && parent.state.obsidian < bluePrint.obsidianForGeodeRobot {
		// We can give up this branch.
		aborted++
		return
	}
	if timeLeft == 1 && parent.state.geodeRobot == 0 {
		// We can give up this branch.
		aborted++
		return
	}

	// Collect resources
	state := parent.state
	state.ore += parent.state.oreRobot
	state.clay += parent.state.clayRobot
	state.obsidian += parent.state.obsidianRobot
	state.geode += parent.state.geodeRobot

	// When there's only 1 logical option.
	if parent.state.obsidianRobot == 0 && parent.state.ore >= bluePrint.oreForObsidianRobot && parent.state.clay >= bluePrint.clayForObsidianRobot {
		addNode(3, state.buildObsidianRobot(bluePrint))
		return
	}

	// When there's only 1 logical option.
	if parent.state.geodeRobot == 0 && parent.state.ore >= bluePrint.oreForGeodeRobot && parent.state.obsidian >= bluePrint.obsidianForGeodeRobot {
		addNode(3, state.buildGeodeRobot(bluePrint))
		return
	}

	// When there's only 1 logical option.
	if parent.state.clayRobot == 0 && parent.state.ore >= bluePrint.oreForClayRobot {
		addNode(3, state.buildClayRobot(bluePrint))
		return
	}

	func() {
		addNode(0, state)
	}()

	// Last step should always be collection, not building.
	// Because there's no point to build anymore.
	if timeLeft == 0 {
		// Skip any building.
		return
	}

	if parent.state.ore >= bluePrint.oreForOreRobot {
		if parent.state.ore > timeLeft*bluePrint.oreForGeodeRobot && parent.state.ore > timeLeft*bluePrint.oreForObsidianRobot {
			// We have more ore than we ever gonna use!
		} else {
			addNode(1, state.buildOreRobot(bluePrint))
		}
	}
	if maxTime-parent.time > 3 && parent.state.ore >= bluePrint.oreForClayRobot {
		if parent.state.clay > timeLeft*bluePrint.clayForObsidianRobot {
			// We have more clay than we ever gonna use!
		} else {
			func() {
				newState := state
				newState.ore -= bluePrint.oreForClayRobot
				newState.clayRobot++
				addNode(2, newState)
			}()
		}
	}
	if parent.state.ore >= bluePrint.oreForObsidianRobot && parent.state.clay >= bluePrint.clayForObsidianRobot {
		if parent.state.obsidian > timeLeft*bluePrint.obsidianForGeodeRobot {
			// We have more obsidian than we ever gonna use!
		} else {
			addNode(3, state.buildObsidianRobot(bluePrint))
		}
	}
	if parent.state.ore >= bluePrint.oreForGeodeRobot && parent.state.obsidian >= bluePrint.obsidianForGeodeRobot {
		func() {
			newState := state
			newState.ore -= bluePrint.oreForGeodeRobot
			newState.obsidian -= bluePrint.obsidianForGeodeRobot
			newState.geodeRobot++
			addNode(4, newState)
		}()
	}

	//fmt.Printf("parent.children: %+v\n", parent.children)
}

func findMaxGeodes(l string) int {
	start := time.Now()

	var maxTime = 24
	var nodeId = 1

	bluePrint := newBluePrint(l)

	root := Node{
		id: 0,
		state: State{
			oreRobot: 1,
		},
		children: make(map[int]*Node),
	}
	expandNode(&root, &bluePrint, &nodeId, maxTime)

	fmt.Printf("Time taken: %s\n", time.Since(start))
	fmt.Printf("endNodes length: %d\n", len(endNodes))

	maxGeodes := 0
	for _, node := range endNodes {
		if _, exists := endNodesByState[node.state]; !exists {
			endNodesByState[node.state] = 1
		} else {
			endNodesByState[node.state]++
		}
		//fmt.Printf("Node: %+v\n", node)
		if node.state.geode > maxGeodes {
			maxGeodes = node.state.geode
		}
	}

	//fmt.Printf("endNodesByState: %+v\n", endNodesByState)

	return maxGeodes
}

func main() {
	scanner := common.GetLineScanner("../sample.txt")

	bluePrintId := 1
	totalQualityLevel := 0
	for scanner.Scan() {
		l := scanner.Text()
		geodes := findMaxGeodes(l)
		fmt.Printf("######## Blue print %d makes %d geodes\n", bluePrintId, geodes)
		totalQualityLevel += bluePrintId * geodes
		if bluePrintId == 1 {
			break
		}
		bluePrintId++
	}

	fmt.Printf("aborted: %d\n", aborted)
	fmt.Printf("nodesEvaluated: %d\n", nodesEvaluated)
	fmt.Printf("nodesByTime: %+v\n", nodesByTime)
	fmt.Printf("foundInEarliestTimesForState: %d\n", foundInEarliestTimesForState)
	fmt.Printf("Part 1: %d\n", totalQualityLevel)
}
