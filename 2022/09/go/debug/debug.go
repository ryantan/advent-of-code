package debug

import (
	"fmt"
	"strconv"
)

// Illustrates knots in a grid like in the samples given.
func PrintKnots(knots [][2]int) {
	minX, minY, maxX, maxY := 0, 0, 0, 0
	for _, knot := range knots {
		if knot[0] < minX {
			minX = knot[0]
		}
		if knot[1] < minY {
			minY = knot[1]
		}
		if knot[0] > maxX {
			maxX = knot[0]
		}
		if knot[1] > maxY {
			maxY = knot[1]
		}
	}
	minX, minY, maxX, maxY = minX-2, minY-2, maxX+2, maxY+2

	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			c := "."
			for k := len(knots) - 1; k >= 0; k-- {
				if knots[k][0] == x && knots[k][1] == y {
					if k == 0 {
						c = "H"
					} else {
						c = strconv.Itoa(k)
					}
				}
			}
			print(c)
		}
		print("\n")
	}
}

func PrintPosition(prefix string, position [2]int) {
	fmt.Printf("%s: (%d, %d)\n", prefix, position[0], position[1])
}
