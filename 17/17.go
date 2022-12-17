package main

import (
	util "aoc2022/aoc"
	inp "aoc2022/input"
	"fmt"
)

const GRID_WIDTH = 7
const MAX_X = GRID_WIDTH - 1

type shape int

const (
	// ° = "Center"
	// #°##
	Minus shape = iota
	// .#.
	// #°#
	// .#.
	Plus
	// ..#
	// ..#
	// ##°
	Hook
	// #
	// °
	// #
	// #
	Bar
	// °#
	// ##
	Square
	NUM_SHAPES = iota
)

var SHAPES [NUM_SHAPES]shape = [NUM_SHAPES]shape{Minus, Plus, Hook, Bar, Square}

type piece struct {
	s      shape
	center util.Coord
}

func (p *piece) edges() (util.Coord, util.Coord, util.Coord, util.Coord) {
	// top, left, bottom, right (offset from center)
	switch p.s {
	case Minus:
		return util.Coord{X: 0, Y: 0}, util.Coord{X: -1, Y: 0}, util.Coord{X: 0, Y: 0}, util.Coord{X: 2, Y: 0}
	case Plus:
		return util.Coord{X: 0, Y: 1}, util.Coord{X: -1, Y: 0}, util.Coord{X: 0, Y: -1}, util.Coord{X: 1, Y: 0}
	case Hook:
		return util.Coord{X: 0, Y: 2}, util.Coord{X: -2, Y: 0}, util.Coord{X: 0, Y: 0}, util.Coord{X: 0, Y: 0}
	case Bar:
		return util.Coord{X: 0, Y: 1}, util.Coord{X: 0, Y: 0}, util.Coord{X: 0, Y: -2}, util.Coord{X: 0, Y: 0}
	case Square:
		return util.Coord{X: 0, Y: 0}, util.Coord{X: 0, Y: 0}, util.Coord{X: 1, Y: -1}, util.Coord{X: 1, Y: -1}
	default:
		panic("Unknown shape")
	}
}

func (p *piece) fall(grid [][GRID_WIDTH]bool) bool {
	ok := p.canFall(grid)
	if ok {
		p.center.Y -= 1
	}
	return ok
}

func (p *piece) right(grid [][GRID_WIDTH]bool) {
	if p.canRight(grid) {
		p.center.X += 1
	}
}

func (p *piece) left(grid [][GRID_WIDTH]bool) {
	if p.canLeft(grid) {
		p.center.X -= 1
	}
}

func (p *piece) settle(grid [][GRID_WIDTH]bool) {
	switch p.s {
	case Minus:
		minusSettle(grid, p.center)
	case Plus:
		plusSettle(grid, p.center)
	case Hook:
		hookSettle(grid, p.center)
	case Bar:
		barSettle(grid, p.center)
	case Square:
		squareSettle(grid, p.center)
	default:
		panic("Unknown shape")
	}
}

func minusSettle(grid [][GRID_WIDTH]bool, center util.Coord) {
	for x := center.X - 1; x < center.X+3; x++ {
		grid[center.Y][x] = true
	}
}

func plusSettle(grid [][GRID_WIDTH]bool, center util.Coord) {
	for x := center.X - 1; x < center.X+2; x++ {
		grid[center.Y][x] = true
	}
	for y := center.Y - 1; y < center.Y+2; y++ {
		grid[y][center.X] = true
	}
}

func hookSettle(grid [][GRID_WIDTH]bool, center util.Coord) {
	for x := center.X - 2; x < center.X+1; x++ {
		grid[center.Y][x] = true
	}
	for y := center.Y; y < center.Y+3; y++ {
		grid[y][center.X] = true
	}
}

func barSettle(grid [][GRID_WIDTH]bool, center util.Coord) {
	for y := center.Y - 2; y < center.Y+2; y++ {
		grid[y][center.X] = true
	}
}

func squareSettle(grid [][GRID_WIDTH]bool, center util.Coord) {
	for x := center.X; x < center.X+2; x++ {
		for y := center.Y - 1; y < center.Y+1; y++ {
			grid[y][x] = true
		}
	}
}

func (p *piece) canFall(grid [][GRID_WIDTH]bool) bool {
	switch p.s {
	case Minus:
		return minusCanFall(grid, p.center)
	case Plus:
		return plusCanFall(grid, p.center)
	case Hook:
		return hookCanFall(grid, p.center)
	case Bar:
		return barCanFall(grid, p.center)
	case Square:
		return squareCanFall(grid, p.center)
	default:
		panic("Unknown shape")
	}
}

func minusCanFall(grid [][GRID_WIDTH]bool, center util.Coord) bool {
	return center.Y > 0 && !util.Any(grid[center.Y-1][center.X-1:center.X+3])
}

func plusCanFall(grid [][GRID_WIDTH]bool, center util.Coord) bool {
	return center.Y > 1 && !grid[center.Y-1][center.X-1] && !grid[center.Y-1][center.X+1] && !grid[center.Y-2][center.X]
}

func hookCanFall(grid [][GRID_WIDTH]bool, center util.Coord) bool {
	return center.Y > 0 && !util.Any(grid[center.Y-1][center.X-2:center.X+1])
}

func barCanFall(grid [][GRID_WIDTH]bool, center util.Coord) bool {
	return center.Y > 2 && !grid[center.Y-3][center.X]
}

func squareCanFall(grid [][GRID_WIDTH]bool, center util.Coord) bool {
	return center.Y > 1 && !util.Any(grid[center.Y-2][center.X:center.X+2])
}

func (p *piece) canRight(grid [][GRID_WIDTH]bool) bool {
	switch p.s {
	case Minus:
		return minusCanRight(grid, p.center)
	case Plus:
		return plusCanRight(grid, p.center)
	case Hook:
		return hookCanRight(grid, p.center)
	case Bar:
		return barCanRight(grid, p.center)
	case Square:
		return squareCanRight(grid, p.center)
	default:
		panic("Unknown shape")
	}
}

func minusCanRight(grid [][GRID_WIDTH]bool, center util.Coord) bool {
	return center.X < MAX_X-2 && !grid[center.Y][center.X+3]
}

func plusCanRight(grid [][GRID_WIDTH]bool, center util.Coord) bool {
	return center.X < MAX_X-1 && !grid[center.Y][center.X+2] && !grid[center.Y+1][center.X+1] && !grid[center.Y-1][center.X+1]
}

func hookCanRight(grid [][GRID_WIDTH]bool, center util.Coord) bool {
	return center.X < MAX_X && !util.Any(util.Map(grid[center.Y:center.Y+3], func(s [GRID_WIDTH]bool) bool { return s[center.X+1] }))
}

func barCanRight(grid [][GRID_WIDTH]bool, center util.Coord) bool {
	return center.X < MAX_X && !util.Any(util.Map(grid[center.Y-2:center.Y+2], func(s [GRID_WIDTH]bool) bool { return s[center.X+1] }))
}

func squareCanRight(grid [][GRID_WIDTH]bool, center util.Coord) bool {
	return center.X < MAX_X-1 && !util.Any(util.Map(grid[center.Y-1:center.Y+1], func(s [GRID_WIDTH]bool) bool { return s[center.X+2] }))
}

func (p *piece) canLeft(grid [][GRID_WIDTH]bool) bool {
	switch p.s {
	case Minus:
		return minusCanLeft(grid, p.center)
	case Plus:
		return plusCanLeft(grid, p.center)
	case Hook:
		return hookCanLeft(grid, p.center)
	case Bar:
		return barCanLeft(grid, p.center)
	case Square:
		return squareCanLeft(grid, p.center)
	default:
		panic("Unknown shape")
	}
}

func minusCanLeft(grid [][GRID_WIDTH]bool, center util.Coord) bool {
	return center.X > 1 && !grid[center.Y][center.X-2]
}

func plusCanLeft(grid [][GRID_WIDTH]bool, center util.Coord) bool {
	return center.X > 1 && !grid[center.Y][center.X-2] && !grid[center.Y+1][center.X-1] && !grid[center.Y-1][center.X-1]
}

func hookCanLeft(grid [][GRID_WIDTH]bool, center util.Coord) bool {
	return center.X > 2 && !grid[center.Y][center.X-3] && !grid[center.Y+1][center.X-1] && !grid[center.Y+2][center.X-1]
}

func barCanLeft(grid [][GRID_WIDTH]bool, center util.Coord) bool {
	return center.X > 0 && !util.Any(util.Map(grid[center.Y-2:center.Y+2], func(s [GRID_WIDTH]bool) bool { return s[center.X-1] }))
}

func squareCanLeft(grid [][GRID_WIDTH]bool, center util.Coord) bool {
	return center.X > 0 && !util.Any(util.Map(grid[center.Y-1:center.Y+1], func(s [GRID_WIDTH]bool) bool { return s[center.X-1] }))
}

func printGrid(grid [][GRID_WIDTH]bool) {
	for i := range grid {
		line := grid[len(grid)-1-i]
		for _, b := range line {
			if b {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
}

func startPiece(shape int, highestPoint int, grid *[][GRID_WIDTH]bool) piece {
	piece := piece{s: SHAPES[shape%NUM_SHAPES]}
	top, left, bottom, _ := piece.edges()
	piece.center = util.Coord{X: 2 - left.X, Y: highestPoint + 4 - bottom.Y}
	for piece.center.Y+top.Y >= len(*grid) {
		*grid = append(*grid, [GRID_WIDTH]bool{})
	}
	return piece
}

func simulate(jets []byte, untilRocksStopped int) int {
	lenJets := len(jets)
	stoppedPieces := 0
	highestPoint := -1
	needNewPiece := true
	grid := make([][GRID_WIDTH]bool, 5)
	var currentPiece piece
	for step := 0; stoppedPieces < untilRocksStopped; step++ {
		if needNewPiece {
			currentPiece = startPiece(stoppedPieces, highestPoint, &grid)
			needNewPiece = false
		}
		if step%2 != 0 {
			// Fall
			couldFall := currentPiece.fall(grid)
			if !couldFall {
				currentPiece.settle(grid)
				stoppedPieces++
				needNewPiece = true
				top, _, _, _ := currentPiece.edges()
				highestPoint = util.Max(currentPiece.center.Y+top.Y, highestPoint)
			}
		} else {
			// Push
			direction := jets[(step/2)%lenJets]
			if direction == '>' {
				currentPiece.right(grid)
			} else if direction == '<' {
				currentPiece.left(grid)
			} else {
				panic("Unknown direction")
			}
		}
	}
	return highestPoint
}

func main() {
	lines := inp.ReadLines("input")

	highestPoint := simulate([]byte(lines[0]), 2022) + 1
	fmt.Printf("Part 1: %d\n", highestPoint)
}
