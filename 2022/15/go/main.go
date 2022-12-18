package main

import (
	"fmt"
	"github/ryantan/advent-of-code/2022/common"
)

type coord [2]int

func (c coord) manhattanDistance(c2 coord) int {
	distanceX, distanceY := c2[0]-c[0], c2[1]-c[1]
	if distanceX < 0 {
		distanceX = -distanceX // abs without dealing with floats.
	}
	if distanceY < 0 {
		distanceY = -distanceY // abs without dealing with floats.
	}
	return distanceX + distanceY
}

type Sensor struct {
	pos           coord
	nearestBeacon coord
	radius        int
	minX          int
	maxX          int
	minY          int
	maxY          int
}

func newSensor(sensorX, sensorY, beaconX, beaconY int) *Sensor {
	sensor := Sensor{
		pos:           coord{sensorX, sensorY},
		nearestBeacon: coord{beaconX, beaconY},
	}

	sensor.radius = sensor.pos.manhattanDistance(sensor.nearestBeacon)
	sensor.minX = sensorX - sensor.radius
	sensor.maxX = sensorX + sensor.radius
	sensor.minY = sensorY - sensor.radius
	sensor.maxY = sensorY + sensor.radius

	return &sensor
}

var max = 4000000

//var max = 20

func main() {
	//scanner := common.GetLineScanner("../sample.txt")
	scanner := common.GetLineScanner("../input.txt")

	var sensorX, sensorY, beaconX, beaconY int
	var sensors []*Sensor
	beaconByCoord := make(map[coord]bool, 0)
	for scanner.Scan() {
		l := scanner.Text()
		_, _ = fmt.Sscanf(l, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &sensorX, &sensorY, &beaconX, &beaconY)
		sensor := newSensor(sensorX, sensorY, beaconX, beaconY)
		sensors = append(sensors, sensor)
		beaconByCoord[sensor.nearestBeacon] = true

	}

	// Part 1
	y := max / 2
	fmt.Printf("Part1: %d\n", impossibleLocationsByY(sensors, beaconByCoord, y))

	// Part 2
	if beacon := part2(sensors, beaconByCoord); beacon != nil {
		fmt.Printf("Part2: %d\n", beacon[0]*4000000+beacon[1])
	} else {
		fmt.Printf("Part2: cannot find!\n")
	}
}

func impossibleLocationsByY(sensors []*Sensor, beaconByCoord map[coord]bool, y int) int {
	minX, maxX := int(^uint(0)>>1), 0
	var ranges [][2]int
	for _, sensor := range sensors {
		if sensor.minY <= y && sensor.maxY >= y {
		} else {
			continue
		}

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

	return definitelyNotHere
}

var logLevel = 1

func possibleLocationsByY2(sensors []*Sensor, y, minX, maxX int) (int, coord) {
	var ranges [][2]int
	for s, sensor := range sensors {
		// Skip out-of-range sensors.
		if sensor.minY > y || sensor.maxY < y {
			if logLevel > 5 {
				fmt.Printf("Skipping sensor %d\n", s)
			}
			continue
		}

		// Calculate range (x1:x2) for y.
		dy := sensor.pos[1] - y
		if dy < 0 {
			dy = -dy
		}
		dx := sensor.radius - dy
		x1, x2 := sensor.pos[0]-dx, sensor.pos[0]+dx

		if (x1 < minX && x2 < minX) || (x1 > maxX && x2 > maxX) {
			// Skip if x1:x2 does not overlap with minX:maxX.
		} else {
			if x1 < minX {
				x1 = minX
			}
			if x2 > maxX {
				x2 = maxX
			}
			ranges = append(ranges, [2]int{x1, x2})
		}
	}
	if logLevel > 4 {
		fmt.Printf("Ranges: %v\n", ranges)
	}

	// Sort ranges
	for i := 0; i < len(ranges)-1; i++ {
		for j := 0; j < len(ranges)-i-1; j++ {
			if ranges[j][0] == ranges[j+1][0] {
				if ranges[j][1] > ranges[j+1][1] {
					ranges[j], ranges[j+1] = ranges[j+1], ranges[j]
				}
			}
			if ranges[j][0] > ranges[j+1][0] {
				ranges[j], ranges[j+1] = ranges[j+1], ranges[j]
			}
		}
	}
	if logLevel > 3 {
		fmt.Printf("Ranges after sort: %v\n", ranges)
	}

	flatRanges := [][2]int{ranges[0]}
	flatRangesIndex := 0
	for i := 1; i < len(ranges); i++ {
		if ranges[i][0] <= flatRanges[flatRangesIndex][1] {
			if ranges[i][1] <= flatRanges[flatRangesIndex][1] {
				// Skip if r2 is completely within r1
				continue
			}
			// Extend current flag range with new end x.
			flatRanges[flatRangesIndex][1] = ranges[i][1]
		} else {
			flatRangesIndex++
			flatRanges = append(flatRanges, ranges[i])
		}
	}
	if logLevel > 3 {
		fmt.Printf("Flattened ranges: %v\n", flatRanges)
	}

	// Find possible locations via iterating through ranges.
	//possibleLocations := 0
	//var possibleLocation coord
	//x := minX
	//for _, r := range flatRanges {
	//	if r[0] > x {
	//		possibleLocations += r[0] - x
	//		possibleLocation = coord{x, y}
	//	}
	//	// Move to end of this range, i.e. 1 after r[1]
	//	x = r[1] + 1
	//	if x > maxX {
	//		break
	//	}
	//}
	//if x < maxX {
	//	possibleLocations += maxX - x + 1
	//	possibleLocation = coord{x, y}
	//}
	//return possibleLocations, possibleLocation

	// Find possible locations by deduction, assuming there's only 1.
	if len(flatRanges) == 1 {
		// If there's only 1 range, the possible locations are at the ends.
		if flatRanges[0][0] != minX {
			return 1, coord{minX, y}
		} else if flatRanges[0][1] != maxX {
			return 1, coord{maxX, y}
		}
	} else if len(flatRanges) == 2 {
		// If there's 2 ranges, the possible locations are after 1st range.
		return 1, coord{flatRanges[0][1] + 1, y}
	}
	// 0 ranges, or range is mixX:maxX
	return 0, coord{}
}

func part2(sensors []*Sensor, beaconByCoord map[coord]bool) *coord {
	// Filter out sensors that are too far away
	// Note: No sensors were out of range in my input.
	// 1. Get diamond enclosing area of interest.
	// 2. Get radius of diamond.
	// 3. Check each sensor's distance to center of diamond.
	// 4. If any sensor's distance is more than sensor radius + diamond radius, rule them out.
	//var aoiSensors []*Sensor
	//aoiCenter, aoiRadius := coord{max / 2, max / 2}, max
	//for s, sensor := range sensors {
	//	if sensor.pos.manhattanDistance(aoiCenter) <= sensor.radius+aoiRadius {
	//		aoiSensors = append(aoiSensors, sensor)
	//	} else {
	//		fmt.Printf("filtered out sensor %d\n", s)
	//	}
	//}
	//fmt.Printf("filtered out %d sensors\n", len(sensors)-len(aoiSensors))
	//sensors = aoiSensors

	// Sort sensors
	// Bubble sort.
	for i := 0; i < len(sensors)-1; i++ {
		for j := 0; j < len(sensors)-i-1; j++ {
			if sensors[j].minY > sensors[j+1].minY {
				sensors[j], sensors[j+1] = sensors[j+1], sensors[j]
			}
		}
	}

	for i := 0; i < len(sensors)-1; i++ {
		for j := 0; j < len(sensors)-i-1; j++ {
			if sensors[j].minX > sensors[j+1].minX {
				sensors[j], sensors[j+1] = sensors[j+1], sensors[j]
			}
		}
	}

	// Review all sensors.
	if logLevel > 1 {
		for s, sensor := range sensors {
			fmt.Printf("Sensor %d: %v, b: %v, r: %d, x: %d-%d, y: %d-%d\n", s, sensor.pos, sensor.nearestBeacon, sensor.radius, sensor.minX, sensor.maxX, sensor.minY, sensor.maxY)
		}
	}

	// Approach #1
	for y := 0; y <= max; y++ {
		possibleLocationCount, possibleLocation := possibleLocationsByY2(sensors, y, 0, max)

		if logLevel > 1 {
			fmt.Printf("possibleLocationCount at %d: %d %v\n", y, possibleLocationCount, possibleLocation)
		}
		if possibleLocationCount > 0 {
			return &possibleLocation
		}
	}

	// Approach #2: Assuming answer lies just 1 pixel outside
	// perimeters of sensor ranges, we can do collision detection of
	// perimeters with diamonds.

	// Approach #3: Divide area of interest into diamonds of radius 1000, and
	// filter out diamonds that overlaps completely with sensor ranges based
	// on radius (if manhattan distance of centers is smaller than sum of
	// individual radius).
	// Perform approach #1 on left over diamonds.

	return nil
}
