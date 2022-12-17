package main

import (
	util "aoc2022/aoc"
	inp "aoc2022/input"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

type rockFormation struct {
	nodes []util.Coord
}

func (r *rockFormation) blocks(c util.Coord) bool {
	for i := 0; i < len(r.nodes)-1; i++ {
		start, end := r.nodes[i], r.nodes[i+1]
		if (start.X == end.X && start.X == c.X && util.Between(c.Y, start.Y, end.Y)) || (start.Y == end.Y && start.Y == c.Y && util.Between(c.X, start.X, end.X)) {
			return true
		}
	}
	return false
}

func (r *rockFormation) deepest() int {
	return util.Max(util.Map(r.nodes, func(c util.Coord) int { return c.Y }))
}

func rockFromLine(l string) rockFormation {
	points := strings.Split(l, " -> ")
	nodes := util.Map(points, func(s string) util.Coord {
		split := strings.Split(s, ",")
		x, errX := strconv.Atoi(split[0])
		y, errY := strconv.Atoi(split[1])
		if errX != nil || errY != nil {
			log.Fatal("Could not parse coordinates")
		}
		return util.Coord{X: x, Y: y}
	})
	return rockFormation{nodes}
}

func addSand(rocks *[]rockFormation, bottom int, sand *map[util.Coord]struct{}, pos util.Coord) bool {
	if pos.Y > bottom {
		return false
	}

	newPositions := []util.Coord{{X: pos.X, Y: pos.Y + 1}, {X: pos.X - 1, Y: pos.Y + 1}, {X: pos.X + 1, Y: pos.Y + 1}}
	for _, newPos := range newPositions {
		blocked := util.Map(*rocks, func(r rockFormation) bool { return r.blocks(newPos) })
		_, hasSand := (*sand)[newPos]
		if !util.Any(blocked) && !hasSand {
			return addSand(rocks, bottom, sand, newPos)
		}
	}
	(*sand)[pos] = struct{}{}
	return true
}

func fillWithSand(rocks *[]rockFormation, start util.Coord) int {
	deepest := util.Max(util.Map(*rocks, func(r rockFormation) int { return r.deepest() }))
	sand := make(map[util.Coord]struct{})
	sandGrains := 0
	for ; ; sandGrains++ {
		_, sourceBlocked := sand[start]
		if sourceBlocked {
			break
		}
		finiteFall := addSand(rocks, deepest, &sand, start)
		if !finiteFall {
			break
		}
	}
	return sandGrains
}

func main() {
	lines := inp.ReadLines("input")
	rocks := util.Map(lines, rockFromLine)

	sandComeToRest := fillWithSand(&rocks, util.Coord{X: 500, Y: 0})
	fmt.Printf("Part 1: %d\n", sandComeToRest)

	deepest := util.Max(util.Map(rocks, func(r rockFormation) int { return r.deepest() }))
	rocks = append(rocks, rockFormation{nodes: []util.Coord{{X: math.MinInt, Y: deepest + 2}, {X: math.MaxInt, Y: deepest + 2}}})
	sandComeToRest = fillWithSand(&rocks, util.Coord{X: 500, Y: 0})
	fmt.Printf("Part 2: %d\n", sandComeToRest)
}
