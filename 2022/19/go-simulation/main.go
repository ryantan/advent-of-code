package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

func simulateBluePrint(l string) int {
	var bluePrintId, oreForOreRobot, oreForClayRobot, oreForObsidianRobot, clayForObsidianRobot, oreForGeodeRobot, obsidianForGeodeRobot int
	_, _ = fmt.Sscanf(l, "Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.", &bluePrintId, &oreForOreRobot, &oreForClayRobot, &oreForObsidianRobot, &clayForObsidianRobot, &oreForGeodeRobot, &obsidianForGeodeRobot)

	var maxTime = 24
	var timeLeft = maxTime
	var ore, clay, obsidian, geode int
	oreRobot, clayRobot, obsidianRobot, geodeRobot := 1, 0, 0, 0
	currentlyBuilding := false
	currentlyBuildingEffect := func() {}

	suitableForGeodeRobot := func() bool {
		if obsidianRobot == 0 {
			return false
		}

		return (ore+(2*oreRobot) >= oreForGeodeRobot) && (obsidian+(2*obsidianRobot) >= obsidianForGeodeRobot)
	}

	suitableForObsidianRobot := func() bool {
		if clayRobot == 0 {
			return false
		}

		// Check production rate or check stock proportion?

		return (ore+(2*oreRobot) >= oreForObsidianRobot) && (clay+(2*clayRobot) >= clayForObsidianRobot)
	}

	suitableForClayRobot := func() bool {
		if float64(clay)/float64(ore) >= float64(clayForObsidianRobot)/float64(oreForObsidianRobot) {
			return false
		}
		if float64(clay)/float64(ore) >= float64(obsidianForGeodeRobot)/float64(oreForGeodeRobot) {
			return false
		}
		return ore > oreForClayRobot
	}

	suitableForOreRobot := func() bool {
		// We know ceiling of how much ore we need.
		//  TODO: Find max ore required among robots.
		maxOreRequired := oreForGeodeRobot * timeLeft
		if ore >= maxOreRequired {
			return false
		}

		return ore > oreForOreRobot
	}

	simulateMinute := func() {
		// Reset.
		currentlyBuilding = false
		currentlyBuildingEffect = func() {}

		if suitableForGeodeRobot() {
			fmt.Printf("It's a right time to build geode robot\n")
			// If there are enough resources for GeodeRobot, build one.
			if ore >= oreForGeodeRobot && obsidian >= obsidianForGeodeRobot {
				fmt.Printf("You have enough ore (%d) and obsidian (%d) to build geode robot\n", oreForGeodeRobot, obsidianForGeodeRobot)
				ore -= oreForGeodeRobot
				obsidian -= obsidianForGeodeRobot
				currentlyBuilding = true
				currentlyBuildingEffect = func() {
					geodeRobot++
				}
			} else {
				//fmt.Printf("You do not have enough ore (%d) and obsidianRobot (%d) to build geode robot\n", oreForGeodeRobot, obsidianRobot)
			}
		} else if suitableForObsidianRobot() {
			fmt.Printf("It's a right time to build obsidian robot\n")
			// If there are enough resources for ObsidianRobot, build one.
			if ore >= oreForObsidianRobot && clay >= clayForObsidianRobot {
				fmt.Printf("You have enough ore (%d) and clay (%d) to build obsidian robot\n", oreForObsidianRobot, clayForObsidianRobot)
				ore -= oreForObsidianRobot
				clay -= clayForObsidianRobot
				currentlyBuilding = true
				currentlyBuildingEffect = func() {
					obsidianRobot++
				}
			} else {
				fmt.Printf("You do not have enough ore (%d) and clayRobot (%d) to build obsidian robot\n", oreForObsidianRobot, clayRobot)
			}
		} else if suitableForClayRobot() {
			fmt.Printf("It's a right time to build clay robot\n")
			// If there are enough resources for ClayRobot, build one.
			if ore >= oreForClayRobot {
				fmt.Printf("You have enough ore (%d) to build clay robot\n", oreForClayRobot)

				ore -= oreForClayRobot
				currentlyBuilding = true
				fmt.Printf("Spend %d ore to build clay robot\n", oreForClayRobot)
				currentlyBuildingEffect = func() {
					clayRobot++
				}
			}
		} else if suitableForOreRobot() {
			if ore >= oreForOreRobot {

				// If there are enough resources for ClayRobot, build one.
				ore -= oreForOreRobot
				currentlyBuilding = true
				fmt.Printf("Spend %d ore to build ore robot\n", oreForOreRobot)
				currentlyBuildingEffect = func() {
					oreRobot++
				}
			}
		}

		// Collect resources
		ore += oreRobot
		clay += clayRobot
		obsidian += obsidianRobot
		geode += geodeRobot

		if currentlyBuilding {
			currentlyBuildingEffect()
		}
	}

	for m := 1; m <= maxTime; m++ {
		fmt.Printf("=== Minute %d\n", m)
		simulateMinute()
		timeLeft--
		fmt.Printf("%d ore, %d clay, %d obsidian, %d geode\n", ore, clay, obsidian, geode)
		fmt.Printf("%d ore robot, %d clay robot, %d obsidian robot, %d geode robot\n", oreRobot, clayRobot, obsidianRobot, geodeRobot)
	}

	return geode
}

func main() {
	scanner := common.GetLineScanner("../sample.txt")

	bluePrintId := 1
	totalQualityLevel := 0
	for scanner.Scan() {
		l := scanner.Text()
		geodes := simulateBluePrint(l)
		fmt.Printf("######## Blue print %d makes %d geodes\n", bluePrintId, geodes)
		totalQualityLevel += bluePrintId * geodes
		bluePrintId++
		//if bluePrintId == 2 {
		//	break
		//}
	}

	fmt.Printf("Part 1: %d\n", totalQualityLevel)
}
