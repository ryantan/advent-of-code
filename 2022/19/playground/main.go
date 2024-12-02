package main

import (
	"fmt"
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

func main() {

	var queue = make([]*Node, 0)

	queue = append(queue, &Node{
		state: State{
			ore:       1,
			clay:      2,
			obsidian:  3,
			geode:     4,
			oreRobot:  5,
			clayRobot: 6,
		},
	})
	queue = append(queue, &Node{
		state: State{
			ore:           1,
			clay:          2,
			obsidian:      4,
			geode:         4,
			oreRobot:      5,
			clayRobot:     6,
			obsidianRobot: 7,
			geodeRobot:    8,
		},
	})
	queue = append(queue, &Node{
		state: State{
			ore:           1,
			clay:          2,
			obsidian:      3,
			geode:         4,
			oreRobot:      5,
			clayRobot:     6,
			obsidianRobot: 0,
			geodeRobot:    0,
		},
	})

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
		fmt.Printf("%d nodes has same state: %+v\n", count, state)
	}
}
