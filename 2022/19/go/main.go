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

//func (s State) score() int {
//	return (s.geode * 10000) + (s.geodeRobot * 1000) + (s.obsidian * 100) + (s.obsidianRobot * 10) + (s.clayRobot * 5) + (s.clay * 5) + s.ore + s.oreRobot
//}

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
	bluePrint *BluePrint
	parent    *Node
	action    int
	maxTime   int
	time      int
	state     State
	children  map[int]*Node
}

func (n *Node) score() int {
	timeLeft := n.maxTime - n.time
	geode := n.state.geode + (n.state.geodeRobot * timeLeft)
	//obsidian := n.state.obsidian + (n.bluePrint.obsidianForGeodeRobot * n.state.obsidianRobot)) * 1000
	obsidian := n.state.obsidian + (n.state.obsidianRobot * timeLeft)
	//clay := (n.state.clay + (n.bluePrint.clayForObsidianRobot * n.state.obsidianRobot)
	clay := n.state.clay + (n.state.clayRobot * timeLeft)
	//ore := n.state.ore + (n.state.oreRobot * timeLeft)
	//return (geode * 1000000) +
	//	(n.state.geodeRobot * 100000) +
	//	(obsidian * 10000) +
	//	(n.state.obsidianRobot * 1000) +
	//	(clay * 100) +
	//	(n.state.clayRobot * 10) +
	//	ore +
	//	(n.state.oreRobot)
	//return (geode * 1000000) +
	//	//(n.state.geodeRobot * 100000) +
	//	(obsidian * 10000) +
	//	//(n.state.obsidianRobot * 1000) +
	//	(clay * 100) +
	//	//(n.state.clayRobot * 10) +
	//	ore
	return (geode * 1000000) +
		(obsidian * 10000) +
		(clay * 100)
}

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

func findMaxGeodes(l string) int {
	start := time.Now()

	var maxTime = 24
	//maxTime = 32
	var queue = make([]*Node, 0)

	bluePrint := newBluePrint(l)

	root := Node{
		bluePrint: bluePrint,
		maxTime:   maxTime,
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
			bluePrint: bluePrint,
			action:    action,
			parent:    parent,
			state:     state,
			maxTime:   maxTime,
			time:      newTime,
			children:  map[int]*Node{},
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

	currentProcessingTime := 0
	startForThisTime := time.Now()
	for len(queue) > 0 {
		currentNode = queue[0]

		newTime := currentNode.time + 1

		// If we want to log processing duration per time left.
		if newTime > currentProcessingTime {
			fmt.Printf("last processing time: %s for %d\n", time.Since(startForThisTime), currentProcessingTime)
			currentProcessingTime = newTime
			startForThisTime = time.Now()

			if currentProcessingTime > 14 {

				if len(queue) > 0 {
					//// Check if there are nodes with diff time.
					//testTime := queue[0].time
					//nodesThatHaveDiffTime := 0
					//for _, n := range queue {
					//	if n.time != testTime {
					//		nodesThatHaveDiffTime++
					//	}
					//}
					//fmt.Printf("nodesThatHaveDiffTime: %d\n", nodesThatHaveDiffTime)
					//if nodesThatHaveDiffTime > 0 {
					//	fmt.Printf("Some nodes have different time!\n")
					//}

					if len(queue) > 6000 {

						// Get distinct nodes
						stateHashed := map[State]int{}
						var distinctQueue []*Node
						for _, node := range queue {
							if _, exists := stateHashed[node.state]; !exists {
								distinctQueue = append(distinctQueue, node)
								stateHashed[node.state] = 0
							}
							stateHashed[node.state]++
						}

						for state, count := range stateHashed {
							if count > 1 {
								fmt.Printf("at %d: %d nodes has same state: %+v\n", newTime, count, state)
							}
						}

						if len(distinctQueue) != len(queue) {
							fmt.Printf("queue is %d while distinctQueue is %d\n", len(queue), len(distinctQueue))
						}

						//// Sort by node score.
						//for i := 0; i < len(queue)-1; i++ {
						//	for j := 0; j < len(queue)-i-1; j++ {
						//		if queue[j].score() < queue[j+1].score() {
						//			queue[j], queue[j+1] = queue[j+1], queue[j]
						//		}
						//	}
						//}

						//statesWithSameScores := map[int][]State{}
						//for _, node := range queue {
						//	score := node.score()
						//	if _, exists := statesWithSameScores[score]; !exists {
						//		statesWithSameScores[score] = []State{}
						//	}
						//	statesWithSameScores[score] = append(statesWithSameScores[score], node.state)
						//}

						//for score, states := range statesWithSameScores {
						//	if len(states) > 2 {
						//		for _, state := range states {
						//			fmt.Printf("at %d: %d nodes score %d: %+v\n", newTime, len(states), score, state)
						//		}
						//	}
						//}

						//queue = distinctQueue

						// Sort by node score.
						for i := 0; i < len(queue)-1; i++ {
							for j := 0; j < len(queue)-i-1; j++ {
								if queue[j].score() < queue[j+1].score() {
									queue[j], queue[j+1] = queue[j+1], queue[j]
								}
							}
						}

						//// Check scores.
						//if len(queue) > 20000 {
						//	for i := 0; i < 20; i++ {
						//		fmt.Printf("score of state #%d: %d\n", i, queue[i].score())
						//	}
						//}

						// Keep only top 1000 nodes.
						if len(queue) > 3000 {
							queue = queue[:3000]
						}
					}
				}
			}
		}

		// Dequeue first item.
		//fmt.Printf("Deqeued 1 node\n")
		queue = queue[1:]

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

		fmt.Printf("nodeWithMaxGeode: %+v\n", nodeWithMaxGeode)
		//printChain(nodeWithMaxGeode)

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
		//if bluePrintId == 3 {
		//	break
		//}
		bluePrintId++
	}

	fmt.Printf("Time taken: %s\n", time.Since(start))
	fmt.Printf("Part 1: %d\n", totalQualityLevel)

	// Wrong answers
	// 1319 - Too low
}
