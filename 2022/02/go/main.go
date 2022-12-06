package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

//var fileName = "../sample.txt"
var fileName = "../input.txt"

func a() int {
	totalScore := 0

	// A = Rock
	// B = Paper
	// C = Scissors
	// X = Rock
	// Y = Paper
	// Z = Scissors
	scores := map[string]int{
		"A X": 4, // 3 + 1
		"A Y": 8, // 6 + 2
		"A Z": 3, // 0 + 3
		"B X": 1, // 0 + 1
		"B Y": 5, // 3 + 2
		"B Z": 9, // 6 + 3
		"C X": 7, // 6 + 1
		"C Y": 2, // 0 + 2
		"C Z": 6, // 3 + 3
	}

	scanner := common.GetLineScanner(fileName)
	for scanner.Scan() {
		moves := scanner.Text()
		score := scores[moves]
		//fmt.Printf("score: %d\n", score)
		totalScore += score
	}

	fmt.Println(totalScore)
	return totalScore
}

func b() int {
	totalScore := 0

	// A = Rock
	// B = Paper
	// C = Scissors
	// X = Lose
	// Y = Draw
	// Z = Win
	scores := map[string]int{
		"A X": 3, // 0 + 3 // Rock x Scissors
		"A Y": 4, // 3 + 1 // Rock x Rock
		"A Z": 8, // 6 + 2 // Rock x Paper
		"B X": 1, // 0 + 1 // Paper x Rock
		"B Y": 5, // 3 + 2 // Paper x Paper
		"B Z": 9, // 6 + 3 // Paper x Scissors
		"C X": 2, // 0 + 2 // Scissors x Paper
		"C Y": 6, // 3 + 3 // Scissors x Scissors
		"C Z": 7, // 6 + 1 // Scissors x Rock
	}

	scanner := common.GetLineScanner(fileName)
	for scanner.Scan() {
		moves := scanner.Text()
		score := scores[moves]
		//fmt.Printf("score: %d\n", score)
		totalScore += score
	}

	fmt.Println(totalScore)
	return totalScore
}

func main() {
	a()
	b()
}
