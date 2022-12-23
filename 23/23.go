package main

import (
	util "aoc2022/aoc"
	inp "aoc2022/input"
	"fmt"
	"math"
)

type direction int

const (
	North direction = iota
	South
	West
	East
	NUM_DIRECTIONS int = iota
)

func neighbors(c util.Coord) []util.Coord {
	n := c.Neighbors()
	return append(n, n[0])
}

var dirOffsets = map[direction]int{North: 0, East: 1, South: 2, West: 3}

func neighborsInDirection(neighbors []util.Coord, dir direction) []util.Coord {
	offset := dirOffsets[dir]
	return neighbors[2*offset : 2*offset+3]
}

func propose(current util.Coord, elves map[util.Coord]struct{}, startDirection direction) util.Coord {
	n := neighbors(current)
	if util.All(util.Map(n, func(c util.Coord) bool { return !util.SetIn(elves, c) })) {
		return current
	}
	for i := 0; i < 4; i++ {
		nID := neighborsInDirection(n, direction((int(startDirection)+i)%NUM_DIRECTIONS))
		if !util.SetIn(elves, nID[0]) && !util.SetIn(elves, nID[1]) && !util.SetIn(elves, nID[2]) {
			return nID[1]
		}
	}
	return current
}

func move(elf util.Coord, elves map[util.Coord]struct{}, proposals map[util.Coord]util.Coord) {
	dest, ok := proposals[elf]
	if ok {
		delete(elves, elf)
		elves[dest] = struct{}{}
	}
}

func round(elves map[util.Coord]struct{}, startDirection direction) {
	proposals := make(map[util.Coord]util.Coord)
	destinationCnt := make(map[util.Coord]int)

	// Get proposals
	for elf := range elves {
		dest := propose(elf, elves, startDirection)
		proposals[elf] = dest
		_, ok := destinationCnt[dest]
		if !ok {
			destinationCnt[dest] = 0
		}
		destinationCnt[dest]++
	}

	// Remove duplicate proposals
	for elf, dest := range proposals {
		if destinationCnt[dest] > 1 {
			delete(proposals, elf)
		}
	}

	// Copy for safe iteration
	elvesIter := make([]util.Coord, len(elves))
	i := 0
	for elf := range elves {
		elvesIter[i] = elf
		i++
	}

	// Move elves
	for _, elf := range elvesIter {
		move(elf, elves, proposals)
	}
}

func simulate(elves map[util.Coord]struct{}, limit int) {
	startingDirection := North
	for r := 0; r < limit; r++ {
		round(elves, startingDirection)
		startingDirection = (startingDirection + 1) % direction(NUM_DIRECTIONS)
	}
}

func boundingBox(elves map[util.Coord]struct{}) (util.Coord, util.Coord) {
	minX, maxX := math.MaxInt, math.MinInt
	minY, maxY := math.MaxInt, math.MinInt
	for elf := range elves {
		if elf.X < minX {
			minX = elf.X
		}
		if elf.X > maxX {
			maxX = elf.X
		}
		if elf.Y < minY {
			minY = elf.Y
		}
		if elf.Y > maxY {
			maxY = elf.Y
		}
	}
	return util.Coord{X: minX, Y: minY}, util.Coord{X: maxX, Y: maxY}
}

func emptyInBB(elves map[util.Coord]struct{}) int {
	bbMin, bbMax := boundingBox(elves)
	return (bbMax.X-bbMin.X+1)*(bbMax.Y-bbMin.Y+1) - len(elves)
}

func printElves(elves map[util.Coord]struct{}) {
	bbMin, bbMax := boundingBox(elves)
	for y := bbMin.Y; y <= bbMax.Y; y++ {
		for x := bbMin.X; x <= bbMax.X; x++ {
			if _, ok := elves[util.Coord{X: x, Y: y}]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
}

func parseElves(lines []string) map[util.Coord]struct{} {
	elves := make(map[util.Coord]struct{})
	for y, l := range lines {
		for x, b := range []byte(l) {
			if b == '#' {
				elves[util.Coord{X: x, Y: y}] = struct{}{}
			}
		}
	}
	return elves
}

func main() {
	lines := inp.ReadLines("input")
	elves := parseElves(lines)

	simulate(elves, 10)
	empytTiles := emptyInBB(elves)
	fmt.Printf("Part 1: %d\n", empytTiles)
}
