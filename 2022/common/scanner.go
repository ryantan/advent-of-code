package common

import (
	"bufio"
	"os"
)

// GetScanner opens a file and returns a bufio.Scanner split by lines.
func GetLineScanner(file string) *bufio.Scanner {
	f, err := os.Open(file)
	if err != nil {
		panic("Can not read input " + file)
	}

	scanner := bufio.NewScanner(f)
	return scanner
}

// GetCharScanner opens a file and returns a bufio.Scanner split by runes.
func GetCharScanner(file string) *bufio.Scanner {
	f, err := os.Open(file)
	if err != nil {
		panic("Can not read input " + file)
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanRunes)
	return scanner
}
