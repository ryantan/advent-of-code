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

func newBluePrint(l string) *BluePrint {
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
	return &bluePrint
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

func (s State) isBetter(otherState State) bool {
	//return s.ore > otherState.ore && s.oreRobot > otherState.oreRobot && s.clay > otherState.clay && s.clayRobot > otherState.clayRobot && s.obsidian > otherState.obsidian && s.obsidianRobot > otherState.obsidianRobot && s.geode > otherState.geode && s.geodeRobot > otherState.geodeRobot
	return s.oreRobot > otherState.oreRobot && s.clayRobot > otherState.clayRobot && s.obsidianRobot > otherState.obsidianRobot && s.geodeRobot > otherState.geodeRobot
}

// tooMuchOre checks if we have more ore than we would ever use!
func (s State) tooMuchOre(timeLeft int, bluePrint *BluePrint) bool {
	if bluePrint.oreForGeodeRobot > bluePrint.oreForObsidianRobot {
		return s.ore > timeLeft*bluePrint.oreForGeodeRobot
	}
	return s.ore > timeLeft*bluePrint.oreForObsidianRobot
}

// tooMuchClay checks if we have more clay than we would ever use!
func (s State) tooMuchClay(timeLeft int, bluePrint *BluePrint) bool {
	return s.clay > timeLeft*bluePrint.clayForObsidianRobot
}

// tooMuchObsidian checks if we have more obsidian than we would ever use!
func (s State) tooMuchObsidian(timeLeft int, bluePrint *BluePrint) bool {
	return s.obsidian > timeLeft*bluePrint.obsidianForGeodeRobot
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

type Node struct {
	id       int
	parent   *Node
	action   int
	time     int
	state    State
	children map[int]*Node
}

//var endNodes = make([]*Node, 0)

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

func stateIsViableWithBluePrint(state State, timeLeft int, bluePrint *BluePrint) bool {

	//if timeLeft == 5 && state.obsidianRobot == 0 && !enoughTimeToBuildObsidianRobot(timeLeft-1, state, bluePrint) {
	//	// We can give up this branch.
	//	return false
	//}
	//if timeLeft == 4 && state.geodeRobot == 0 && state.obsidianRobot == 0 {
	//	// We can give up this branch.
	//	return false
	//}
	if timeLeft == 2 && state.geodeRobot == 0 && !enoughTimeToBuildGeodeRobot(timeLeft-1, state, bluePrint) {
		// We can give up this branch.
		return false
	}
	if timeLeft == 2 && state.geodeRobot == 0 && state.obsidian < bluePrint.obsidianForGeodeRobot {
		// We can give up this branch.
		return false
	}
	if timeLeft == 1 && state.geodeRobot == 0 {
		// We can give up this branch.
		return false
	}

	return true
}

//func expandNode(parent *Node, bluePrint *BluePrint, nodeId *int, maxTime int) {
//
//	newTime := parent.time + 1
//
//	addNode := func(action int, state State) {
//		newNode := &Node{
//			id:       *nodeId,
//			action:   action,
//			parent:   parent,
//			state:    state,
//			time:     newTime,
//			children: map[int]*Node{},
//		}
//		*nodeId++
//		nodesEvaluated++
//		if _, exists := nodesByTime[newTime]; !exists {
//			nodesByTime[newTime] = 1
//		} else {
//			nodesByTime[newTime]++
//		}
//
//		// If this timing was beaten before, kill this branch.
//		if earliestTimeForState, exists := earliestTimesForState[state]; exists {
//			foundInEarliestTimesForState++
//			if newTime >= earliestTimeForState.time {
//				aborted++
//				return
//			} else {
//				// Replace
//				nodeIdsToAbort[earliestTimeForState.id] = true
//				earliestTimesForState[state] = timeAndId{time: newTime, id: newNode.id}
//			}
//		} else {
//			earliestTimesForState[state] = timeAndId{time: newTime, id: newNode.id}
//		}
//
//		if oldCount, exists := statesSeen[state]; exists {
//			statesSeen[state] = oldCount + 1
//		} else {
//			statesSeen[state] = 1
//			distinctStatesSeen++
//		}
//
//		if _, exists := statesSeenByTime[newTime]; !exists {
//			statesSeenByTime[newTime] = map[State]int{}
//		}
//		if _, exists := statesSeenByTime[newTime][state]; exists {
//			statesSeenByTime[newTime][state]++
//		} else {
//			statesSeenByTime[newTime][state] = 1
//		}
//		if _, exists := distinctStatesSeenByTime[newTime]; exists {
//			distinctStatesSeenByTime[newTime]++
//		} else {
//			distinctStatesSeenByTime[newTime] = 1
//		}
//
//		if newTime == maxTime {
//			fmt.Printf("endNodes length: %d\n", len(endNodes))
//			endNodes = append(endNodes, newNode)
//		} else {
//			expandNode(newNode, bluePrint, nodeId, maxTime)
//		}
//		fmt.Printf("newNode (%d): %+v\n", *nodeId, newNode)
//		parent.children[newNode.action] = newNode
//	}
//
//	timeLeft := maxTime - newTime
//
//	if !stateIsViableWithBluePrint(parent.state, timeLeft, bluePrint) {
//		return
//	}
//
//	// Collect resources
//	state := parent.state
//	state.ore += parent.state.oreRobot
//	state.clay += parent.state.clayRobot
//	state.obsidian += parent.state.obsidianRobot
//	state.geode += parent.state.geodeRobot
//
//	// When there's only 1 logical option.
//	if parent.state.obsidianRobot == 0 && parent.state.ore >= bluePrint.oreForObsidianRobot && parent.state.clay >= bluePrint.clayForObsidianRobot {
//		addNode(3, state.buildObsidianRobot(bluePrint))
//		return
//	}
//
//	// When there's only 1 logical option.
//	if parent.state.geodeRobot == 0 && parent.state.ore >= bluePrint.oreForGeodeRobot && parent.state.obsidian >= bluePrint.obsidianForGeodeRobot {
//		addNode(3, state.buildGeodeRobot(bluePrint))
//		return
//	}
//
//	// When there's only 1 logical option.
//	if parent.state.clayRobot == 0 && parent.state.ore >= bluePrint.oreForClayRobot {
//		addNode(3, state.buildClayRobot(bluePrint))
//		return
//	}
//
//	func() {
//		addNode(0, state)
//	}()
//
//	// Last step should always be collection, not building.
//	// Because there's no point to build anymore.
//	if timeLeft == 0 {
//		// Skip any building.
//		return
//	}
//
//	if parent.state.ore >= bluePrint.oreForOreRobot {
//		if !parent.state.tooMuchOre(timeLeft, bluePrint) {
//			addNode(1, state.buildOreRobot(bluePrint))
//		}
//	}
//	if maxTime-parent.time > 3 && parent.state.ore >= bluePrint.oreForClayRobot {
//		if !parent.state.tooMuchClay(timeLeft, bluePrint) {
//			addNode(3, state.buildClayRobot(bluePrint))
//		}
//	}
//	if parent.state.ore >= bluePrint.oreForObsidianRobot && parent.state.clay >= bluePrint.clayForObsidianRobot {
//		if !parent.state.tooMuchObsidian(timeLeft, bluePrint) {
//			addNode(3, state.buildObsidianRobot(bluePrint))
//		}
//	}
//	if parent.state.ore >= bluePrint.oreForGeodeRobot && parent.state.obsidian >= bluePrint.obsidianForGeodeRobot {
//		func() {
//			addNode(4, state.buildGeodeRobot(bluePrint))
//		}()
//	}
//
//	//fmt.Printf("parent.children: %+v\n", parent.children)
//}

func findMaxGeodes(l string) int {
	start := time.Now()

	var maxTime = 24
	//var nodeId = 1
	var queue = make([]*Node, 0)

	bluePrint := newBluePrint(l)

	root := Node{
		id: 0,
		state: State{
			oreRobot: 1,
		},
		children: make(map[int]*Node),
	}
	queue = append(queue, &root)

	var currentNode *Node
	endNodes := make([]*Node, 0)
	var earliestTimesForState = map[State]int{}
	var foundInEarliestTimesForState = 0
	var aborted = 0
	var nodesEvaluated = 0
	var nodesByTime = map[int]int{}
	var endNodesByState = map[State]int{}

	addNode := func(action int, state State, parent *Node, newTime int) {
		newNode := &Node{
			id:       0,
			action:   action,
			parent:   parent,
			state:    state,
			time:     newTime,
			children: map[int]*Node{},
		}
		nodesEvaluated++

		// If this timing was beaten before, kill this branch.
		if earliestTimeForState, exists := earliestTimesForState[state]; exists {
			foundInEarliestTimesForState++
			if newTime >= earliestTimeForState {
				aborted++
				return
			} else {
				// Replace
				earliestTimesForState[state] = newTime
			}
		} else {
			earliestTimesForState[state] = newTime
		}

		if _, exists := nodesByTime[newTime]; !exists {
			nodesByTime[newTime] = 1
		} else {
			nodesByTime[newTime]++
		}

		//fmt.Printf("Queued 1 node\n")
		queue = append(queue, newNode)
	}

	maxGeodeRobotsOnLastRound := 0

	//currentProcessingTime := 0
	//startForThisTime := time.Now()
	for len(queue) > 0 {
		currentNode = queue[0]

		// Dequeue first item.
		//fmt.Printf("Deqeued 1 node\n")
		queue = queue[1:]

		newTime := currentNode.time + 1

		//// If we want to log processing duration per time left.
		//if newTime > currentProcessingTime {
		//	fmt.Printf("last processing time: %s for %d\n", time.Since(startForThisTime), currentProcessingTime)
		//	currentProcessingTime = newTime
		//	startForThisTime = time.Now()
		//}

		if currentNode.time == maxTime {
			endNodes = append(endNodes, currentNode)
			continue
		}

		timeLeft := maxTime - newTime

		if timeLeft == 0 {
			if currentNode.state.geode+currentNode.state.geodeRobot > maxGeodeRobotsOnLastRound {
				// Consider.
				maxGeodeRobotsOnLastRound = currentNode.state.geode + currentNode.state.geodeRobot
				//fmt.Printf("currentNode.state: %+v\n", currentNode.state)
				//fmt.Printf("maxGeodeRobotsOnLastRound: %d\n", maxGeodeRobotsOnLastRound)
			} else {
				// Don't consider.
				continue
			}
		}

		//if !stateIsViableWithBluePrint(currentNode.state, timeLeft, bluePrint) {
		//	aborted++
		//	continue
		//}

		// Collect resources
		state := currentNode.state
		state.ore += currentNode.state.oreRobot
		state.clay += currentNode.state.clayRobot
		state.obsidian += currentNode.state.obsidianRobot
		state.geode += currentNode.state.geodeRobot

		// When there's only 1 logical option.
		//if bluePrint.oreForOreRobot > bluePrint.oreForClayRobot {
		//	if currentNode.state.clayRobot == 0 && currentNode.state.ore >= bluePrint.oreForClayRobot {
		//		addNode(2, state.buildClayRobot(bluePrint), currentNode, newTime)
		//		continue
		//	}
		//	if currentNode.state.oreRobot == 1 && currentNode.state.ore >= bluePrint.oreForOreRobot {
		//		addNode(1, state.buildOreRobot(bluePrint), currentNode, newTime)
		//		continue
		//	}
		//} else {
		//	if currentNode.state.oreRobot == 1 && currentNode.state.ore >= bluePrint.oreForOreRobot {
		//		addNode(1, state.buildOreRobot(bluePrint), currentNode, newTime)
		//		continue
		//	}
		//	if currentNode.state.clayRobot == 0 && currentNode.state.ore >= bluePrint.oreForClayRobot {
		//		addNode(2, state.buildClayRobot(bluePrint), currentNode, newTime)
		//		continue
		//	}
		//}
		//if currentNode.state.oreRobot == 1 && currentNode.state.ore >= bluePrint.oreForOreRobot {
		//	addNode(1, state.buildOreRobot(bluePrint), currentNode, newTime)
		//	continue
		//}
		//if currentNode.state.clayRobot == 0 && currentNode.state.ore >= bluePrint.oreForClayRobot {
		//	addNode(2, state.buildClayRobot(bluePrint), currentNode, newTime)
		//	continue
		//}
		if currentNode.state.obsidianRobot == 0 && currentNode.state.ore >= bluePrint.oreForObsidianRobot && currentNode.state.clay >= bluePrint.clayForObsidianRobot {
			addNode(3, state.buildObsidianRobot(bluePrint), currentNode, newTime)
			continue
		}
		if currentNode.state.geodeRobot == 0 && currentNode.state.ore >= bluePrint.oreForGeodeRobot && currentNode.state.obsidian >= bluePrint.obsidianForGeodeRobot {
			addNode(4, state.buildGeodeRobot(bluePrint), currentNode, newTime)
			continue
		}

		// First 2 steps should always be collection, not building, because there's not enough resources to build anything.
		// Last step should always be collection, not building, because there's no point to build anymore.
		if timeLeft == 0 || newTime == 1 || newTime == 2 {
			addNode(0, state, currentNode, newTime)
			continue
		}

		// Collect.
		addNode(0, state, currentNode, newTime)

		// Don't build ore robot in last 4 minutes
		if timeLeft > 3 && currentNode.state.ore >= bluePrint.oreForOreRobot {
			//if !currentNode.state.tooMuchOre(timeLeft, bluePrint) {
			addNode(1, state.buildOreRobot(bluePrint), currentNode, newTime)
			//}
		}

		// Don't build clay robot in last 4 minutes
		if timeLeft > 3 && currentNode.state.ore >= bluePrint.oreForClayRobot {
			//if !currentNode.state.tooMuchClay(timeLeft, bluePrint) {
			addNode(2, state.buildClayRobot(bluePrint), currentNode, newTime)
			//}
		}
		if currentNode.state.ore >= bluePrint.oreForObsidianRobot && currentNode.state.clay >= bluePrint.clayForObsidianRobot {
			//if !currentNode.state.tooMuchObsidian(timeLeft, bluePrint) {
			addNode(3, state.buildObsidianRobot(bluePrint), currentNode, newTime)
			//}
		}
		if currentNode.state.ore >= bluePrint.oreForGeodeRobot && currentNode.state.obsidian >= bluePrint.obsidianForGeodeRobot {
			addNode(4, state.buildGeodeRobot(bluePrint), currentNode, newTime)
		}

		// TODO: Filter out clearly inferior states.
	}

	fmt.Printf("Time taken: %s\n", time.Since(start))
	fmt.Printf("endNodes length: %d\n", len(endNodes))
	fmt.Printf("aborted: %d\n", aborted)
	fmt.Printf("nodesEvaluated: %d\n", nodesEvaluated)
	fmt.Printf("nodesByTime: %+v\n", nodesByTime)
	fmt.Printf("foundInEarliestTimesForState: %d\n", foundInEarliestTimesForState)
	//fmt.Printf("foundBetterInStatesAtTime: %d\n", foundBetterInStatesAtTime)
	//fmt.Printf("foundInBestStatesAtTime: %d\n", foundInBestStatesAtTime)
	//fmt.Printf("distinctStatesSeen: %d\n", distinctStatesSeen)
	//fmt.Printf("distinctStatesSeenByTime: %d\n", distinctStatesSeenByTime)

	maxGeodes := 0

	if len(endNodes) > 0 {
		var nodeWithMaxGeode *Node = endNodes[0]
		for _, node := range endNodes {
			if _, exists := endNodesByState[node.state]; !exists {
				endNodesByState[node.state] = 1
			} else {
				endNodesByState[node.state]++
			}
			//fmt.Printf("Node: %+v\n", node)
			if node.state.geode > maxGeodes {
				maxGeodes = node.state.geode
				nodeWithMaxGeode = node
			}
		}
		//fmt.Printf("endNodesByState: %+v\n", endNodesByState)

		printChain(nodeWithMaxGeode)

		//for i, node := range endNodes {
		//	fmt.Printf("Endnode %d state: %+v\n", i, node.state)
		//}
	}

	return maxGeodes
}

func printChain(n *Node) {
	if n == nil {
		fmt.Printf("No chain (no end nodes).\n")
		return
	}

	var chain []*Node
	for {
		chain = append(chain, n)
		n = n.parent
		if n.time == 0 {
			break
		}
	}

	for i := len(chain) - 1; i >= 0; i-- {
		node := chain[i]
		fmt.Printf("== Minute %d ==\n", node.time)
		if node.action == 0 {

		} else if node.action == 1 {
			fmt.Printf("Building Ore robot.\n")
		} else if node.action == 2 {
			fmt.Printf("Building Clay robot.\n")
		} else if node.action == 3 {
			fmt.Printf("Building Obsidian robot.\n")
		} else if node.action == 4 {
			fmt.Printf("Building Geode robot.\n")
		}
		fmt.Printf("State: %+v\n", node.state)
	}
}

func main() {
	//scanner := common.GetLineScanner("../sample.txt")
	scanner := common.GetLineScanner("../input.txt")

	bluePrintId := 1
	totalQualityLevel := 0
	start := time.Now()
	for scanner.Scan() {
		l := scanner.Text()
		geodes := findMaxGeodes(l)
		fmt.Printf("######## Blue print %d makes %d geodes\n", bluePrintId, geodes)
		totalQualityLevel += bluePrintId * geodes
		//if bluePrintId == 1 {
		//	break
		//}
		bluePrintId++
	}

	fmt.Printf("Time taken: %s\n", time.Since(start))
	fmt.Printf("Part 1: %d\n", totalQualityLevel)

	// Wrong answers
	// 1319 - Too low
}
