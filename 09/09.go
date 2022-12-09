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
	head position
	tail position
}

func (r *rope) move(m motion) []position {
	// Move head
	switch m.d {
	case Right:
		r.head.x += m.l
	case Up:
		r.head.y += m.l
	case Left:
		r.head.x -= m.l
	case Down:
		r.head.y -= m.l
	}

	return r.pullTail()
}

func (r *rope) pullTail() []position {
	if r.tail.adjacent(r.head) {
		return []position{r.tail}
	}

	// diagonal if necessary
	if abs(r.head.x-r.tail.x) == 1 {
		r.tail.x = r.head.x
	} else if abs(r.head.y-r.tail.y) == 1 {
		r.tail.y = r.head.y
	}

	return r.pullTailRec()
}

func (r *rope) pullTailRec() []position {
	if r.tail.adjacent(r.head) {
		return []position{r.tail}
	}

	if r.head.x == r.tail.x {
		if r.tail.y > r.head.y {
			r.tail.y -= 1
		} else {
			r.tail.y += 1
		}
	} else if r.head.y == r.tail.y {
		if r.tail.x > r.head.x {
			r.tail.x -= 1
		} else {
			r.tail.x += 1
		}
	} else {
		panic("unreachable")
	}
	tailPos := r.tail
	return append(r.pullTailRec(), tailPos)
}

func main() {
	lines := inp.ReadLines("input")
	motions := make([]motion, len(lines))
	for i, l := range lines {
		motions[i] = motionFromLine(l)
	}
	rope := rope{
		head: position{x: 0, y: 0},
		tail: position{x: 0, y: 0},
	}

	tailVisited := make(map[position]struct{})
	tailVisited[rope.tail] = struct{}{}
	for _, m := range motions {
		positions := rope.move(m)
		for _, p := range positions {
			tailVisited[p] = struct{}{}
		}
	}
	fmt.Printf("Part 1: %d\n", len(tailVisited))

	// bestScore := max(scenicScore)
	// fmt.Printf("Part 2: %d\n", bestScore)
}
