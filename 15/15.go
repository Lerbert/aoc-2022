package main

import (
	util "aoc2022/aoc"
	inp "aoc2022/input"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type sensor struct {
	pos              util.Coord
	beacon           util.Coord
	distanceToBeacon int
}

func (s sensor) coveredAtHeight(height int) (util.Range, bool) {
	covered := s.distanceToBeacon - util.Abs(s.pos.Y-height)
	return util.Range{Lower: s.pos.X - covered, Upper: s.pos.X + covered}, covered > 0
}

func parseCoord(s string) util.Coord {
	split := strings.Split(s, ", ")
	x, errX := strconv.Atoi(strings.TrimPrefix(split[0], "x="))
	y, errY := strconv.Atoi(strings.TrimPrefix(split[1], "y="))
	if errX != nil || errY != nil {
		log.Fatal("Could not parse coords")
	}
	return util.Coord{X: x, Y: y}
}

func sensorfromline(s string) sensor {
	split := strings.Split(s, ": ")
	pos := parseCoord(strings.TrimPrefix(split[0], "Sensor at "))
	beacon := parseCoord(strings.TrimPrefix(split[1], "closest beacon is at "))
	return sensor{pos: pos, beacon: beacon, distanceToBeacon: pos.ManhattanDistance(beacon)}
}

func removeIndex[T interface{}](s []T, i int) []T {
	if i == len(s)-1 {
		return s[:i]
	} else {
		return append(s[:i], s[i+1:]...)
	}
}

func mergeRanges(ranges []util.Range, newRange util.Range) []util.Range {
	for i, r := range ranges {
		if newRange.Upper < r.Lower {
			ranges = append(ranges[:i+1], ranges[i:]...)
			ranges[i] = newRange
			return ranges
		} else if newRange.Includes(r) {
			return mergeRanges(removeIndex(ranges, i), newRange)
		} else if r.Includes(newRange) {
			return ranges
		} else if r.Contains(newRange.Lower) {
			lower := ranges[i].Lower
			return mergeRanges(removeIndex(ranges, i), util.Range{Lower: lower, Upper: newRange.Upper})
		} else if r.Contains(newRange.Upper) {
			upper := ranges[i].Upper
			return mergeRanges(removeIndex(ranges, i), util.Range{Lower: newRange.Lower, Upper: upper})
		}
	}
	return append(ranges, newRange)
}

func coveredBeacons(ranges []util.Range, beacons map[util.Coord]struct{}, height int) int {
	sum := 0
	for _, r := range ranges {
		for b := range beacons {
			if b.Y == height && r.Contains(b.X) {
				sum++
			}
		}
	}
	return sum
}

func rangesForHeight(sensors []sensor, height int) []util.Range {
	ranges := make([]util.Range, 0)
	for _, s := range sensors {
		newRange, ok := s.coveredAtHeight(height)
		if ok {
			ranges = mergeRanges(ranges, newRange)
		}
	}
	return ranges
}

func findUncoveredSpot(sensors []sensor, maxSearch int) (int, int) {
	var y int
	for y = 0; y <= maxSearch; y++ {
		ranges := rangesForHeight(sensors, y)
		inclusion := util.Map(ranges, func(r util.Range) bool { return r.Includes(util.Range{Lower: 0, Upper: maxSearch}) })
		if !util.Any(inclusion) {
			for _, r := range ranges {
				if r.Contains(0) {
					// There is exactly one uncovered spot between 0 and maxSearch, so it must be one after the range that contains zero
					return r.Upper + 1, y
				}
			}
			// If no range contains zero, zero must be uncovered
			return 0, y
		}
	}
	panic("No uncovered spot")
}

const HEIGHT = 2_000_000
const SEARCH = 4_000_000
const X_FACTOR = 4_000_000

func main() {
	lines := inp.ReadLines("input")
	sensors := util.Map(lines, sensorfromline)
	beacons := make(map[util.Coord]struct{})
	for _, s := range sensors {
		beacons[s.beacon] = struct{}{}
	}

	ranges := rangesForHeight(sensors, HEIGHT)
	coveredSpots := util.Sum(util.Map(ranges, func(r util.Range) int { return r.Upper - r.Lower + 1 }))
	correction := coveredBeacons(ranges, beacons, HEIGHT)
	fmt.Printf("Part 1: %d\n", coveredSpots-correction)

	x, y := findUncoveredSpot(sensors, SEARCH)
	fmt.Printf("Part 2: %d\n", x*X_FACTOR+y)
}
