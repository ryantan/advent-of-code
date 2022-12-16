package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

type coord [2]int

func (c coord) manhattanDistance(c2 coord) int {
	distanceX, distanceY := c2[0]-c[0], c2[1]-c[1]
	if distanceX < 0 {
		distanceX = -distanceX
	}
	if distanceY < 0 {
		distanceY = -distanceY
	}
	return distanceX + distanceY
}

type Sensor struct {
	pos           coord
	nearestBeacon coord
	radius        int
}

func newSensor(sensorX, sensorY, beaconX, beaconY int) *Sensor {
	sensor := Sensor{
		pos:           coord{sensorX, sensorY},
		nearestBeacon: coord{beaconX, beaconY},
	}

	sensor.radius = sensor.pos.manhattanDistance(sensor.nearestBeacon)
	return &sensor
}

func main() {
	//scanner := common.GetLineScanner("../sample.txt")
	scanner := common.GetLineScanner("../input.txt")

	var sensorX, sensorY, beaconX, beaconY int
	var sensors []*Sensor
	y := 2000000
	var sensorsInRangeForY []*Sensor
	beaconByCoord := make(map[coord]bool, 0)
	for scanner.Scan() {
		l := scanner.Text()
		_, _ = fmt.Sscanf(l, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &sensorX, &sensorY, &beaconX, &beaconY)

		sensor := newSensor(sensorX, sensorY, beaconX, beaconY)
		sensors = append(sensors, sensor)
		beaconByCoord[sensor.nearestBeacon] = true

		if sensor.pos[1]-sensor.radius <= y && sensor.pos[1]+sensor.radius >= y {
			sensorsInRangeForY = append(sensorsInRangeForY, sensor)
		}
	}

	minX, maxX := int(^uint(0)>>1), 0
	var ranges [][2]int
	for _, sensor := range sensorsInRangeForY {
		dy := sensor.pos[1] - y
		if dy < 0 {
			dy = -dy
		}
		dx := sensor.radius - dy
		x1 := sensor.pos[0] - dx
		x2 := sensor.pos[0] + dx

		ranges = append(ranges, [2]int{x1, x2})

		if x1 < minX {
			minX = x1
		}
		if x2 > maxX {
			maxX = x2
		}
	}

	definitelyNotHere := 0
	for i := minX; i <= maxX; i++ {
		inRange := 0
		for _, r := range ranges {
			if r[0] <= i && r[1] >= i {
				inRange = 1
				break
			}
		}
		if beaconByCoord[coord{i, y}] {
			// Beacon present, skip.
		} else {
			definitelyNotHere += inRange
		}
	}

	fmt.Printf("Part1: %d\n", definitelyNotHere)
}
