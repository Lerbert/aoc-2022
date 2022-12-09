package main

import (
	inp "aoc2022/input"
	"fmt"
	"log"
	"strconv"
)

type direction int

const (
	Right = iota
	Up
	Left
	Down
)

func directionFromChar(c byte) direction {
	switch c {
	case 'R':
		return Right
	case 'U':
		return Up
	case 'L':
		return Left
	case 'D':
		return Down
	}
	log.Fatal("Could not parse direction")
	panic("")
}

type motion struct {
	d direction
	l int
}

func motionFromLine(line string) motion {
	d := directionFromChar(line[0])
	l, err := strconv.Atoi(line[2:])
	if err != nil {
		log.Fatal("Could not convert motion length", err)
	}
	return motion{
		d,
		l,
	}
}

type position struct {
	x int
	y int
}

func (p position) adjacent(other position) bool {
	return abs(p.x-other.x) < 2 && abs(p.y-other.y) < 2
}

func abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

type rope struct {
	knots []position
}

func (r *rope) move(m motion) []position {
	if m.l == 0 {
		return []position{r.knots[len(r.knots)-1]}
	}

	// Move head
	switch m.d {
	case Right:
		r.knots[0].x++
	case Up:
		r.knots[0].y++
	case Left:
		r.knots[0].x--
	case Down:
		r.knots[0].y--
	}

	for i := 1; i < len(r.knots); i++ {
		r.pullKnot(i)
	}

	tailPos := r.knots[len(r.knots)-1]
	return append(r.move(motion{d: m.d, l: m.l - 1}), tailPos)
}

func (r *rope) pullKnot(i int) {
	head := &r.knots[i-1]
	tail := &r.knots[i]
	if !tail.adjacent(*head) {
		tail.x += sgn(head.x - tail.x)
		tail.y += sgn(head.y - tail.y)
	}
}

func sgn(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	} else {
		return 0
	}
}

func simulateRope(motions []motion, ropeLength int) int {
	initPos := position{x: 0, y: 0}
	knots := make([]position, ropeLength)
	for i := range knots {
		knots[i] = initPos
	}
	rope := rope{knots}

	tailVisited := make(map[position]struct{})
	tailVisited[initPos] = struct{}{}
	for _, m := range motions {
		positions := rope.move(m)
		for _, p := range positions {
			tailVisited[p] = struct{}{}
		}
	}
	return len(tailVisited)
}

func main() {
	lines := inp.ReadLines("input")
	motions := make([]motion, len(lines))
	for i, l := range lines {
		motions[i] = motionFromLine(l)
	}

	tailVisited := simulateRope(motions, 2)
	fmt.Printf("Part 1: %d\n", tailVisited)

	tailVisited = simulateRope(motions, 10)
	fmt.Printf("Part 2: %d\n", tailVisited)
}
