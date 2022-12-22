package main

import (
	util "aoc2022/aoc"
	inp "aoc2022/input"
	"fmt"
	"regexp"
)

type marker int

const (
	Wrap marker = iota
	Wall
	Free
)

func markerFromByte(b byte) marker {
	switch b {
	case ' ':
		return Wrap
	case '#':
		return Wall
	case '.':
		return Free
	default:
		panic("Unknown marker")
	}
}

type facing int

const (
	North facing = iota
	West
	South
	East
	NUM_FACINGS
)

func (f facing) score() int {
	switch f {
	case North:
		return 3
	case West:
		return 2
	case South:
		return 1
	case East:
		return 0
	default:
		panic("Unknown facing")
	}
}

type instruction interface {
	execute(board [][]marker, transposedBoard [][]marker, pos util.Coord, f facing) (util.Coord, facing)
}

type turn int

const (
	Left turn = iota
	Right
)

func (t turn) execute(board [][]marker, transposedBoard [][]marker, pos util.Coord, f facing) (util.Coord, facing) {
	switch t {
	case Left:
		return pos, (f + 1) % NUM_FACINGS
	case Right:
		return pos, (f + NUM_FACINGS - 1) % NUM_FACINGS
	default:
		panic("Unknown turn")
	}
}

func turnFromByte(b byte) turn {
	switch b {
	case 'L':
		return Left
	case 'R':
		return Right
	default:
		panic("Unknown turn")
	}
}

type forward int

func (fw forward) execute(board [][]marker, transposedBoard [][]marker, pos util.Coord, f facing) (util.Coord, facing) {
	var cnt int
	var straightPos int
	var straight []marker
	switch f {
	case North:
		cnt = -int(fw)
		straightPos = pos.Y
		straight = transposedBoard[pos.X]
	case West:
		cnt = -int(fw)
		straightPos = pos.X
		straight = board[pos.Y]
	case South:
		cnt = int(fw)
		straightPos = pos.Y
		straight = transposedBoard[pos.X]
	case East:
		cnt = int(fw)
		straightPos = pos.X
		straight = board[pos.Y]
	}

	nextCoord := forwardInStraight(straight, straightPos, cnt)

	switch f {
	case North, South:
		pos.Y = nextCoord
	case West, East:
		pos.X = nextCoord
	}

	return pos, f
}

func wrapIndex(index int, length int) int {
	if mod := index % length; mod >= 0 {
		return mod
	} else {
		return length + mod
	}
}

func forwardInStraight(straight []marker, pos int, cnt int) int {
	straightLen := len(straight)
	var updatePos func(int) int
	if cnt >= 0 {
		updatePos = func(p int) int { return wrapIndex(p+1, straightLen) }
	} else {
		updatePos = func(p int) int { return wrapIndex(p-1, straightLen) }
	}
	lastPos := pos
	for i := 0; i < util.Abs(cnt); i++ {
		lastPos = pos
		pos = updatePos(pos)
		for straight[pos] == Wrap {
			pos = updatePos(pos)
		}
		if straight[wrapIndex(pos, straightLen)] == Wall {
			return lastPos
		}
	}
	return pos
}

var turnRegex = regexp.MustCompile(`R|L`)
var forwardRegex = regexp.MustCompile(`\d+`)

func parseInstructions(s string) []instruction {
	turns := util.Map(turnRegex.FindAllString(s, -1), func(m string) turn { return turnFromByte(m[0]) })
	forwards := util.Map(forwardRegex.FindAllString(s, -1), func(m string) forward { return forward(inp.MustAtoi(m)) })
	instructions := make([]instruction, len(turns)+len(forwards))
	for i := range instructions {
		if i%2 == 0 {
			instructions[i] = forwards[i/2]
		} else {
			instructions[i] = turns[i/2]
		}
	}
	return instructions
}

func startPos(board [][]marker) (util.Coord, facing) {
	for x, m := range board[0] {
		if m == Free {
			return util.Coord{X: x, Y: 0}, East
		}
	}
	panic("No suitable starting position")
}

func followInstructions(board [][]marker, transposedBoard [][]marker, startC util.Coord, startF facing, instructions []instruction) (util.Coord, facing) {
	pos := startC
	f := startF
	for _, instr := range instructions {
		pos, f = instr.execute(board, transposedBoard, pos, f)
	}
	return pos, f
}

func transpose[T interface{}](m [][]T) [][]T {
	transposed := make([][]T, len(m[0]))
	for col := range transposed {
		transposed[col] = make([]T, len(m))
		for row := range transposed[col] {
			transposed[col][row] = m[row][col]
		}
	}
	return transposed
}

func main() {
	lines := inp.ReadLines("input")
	instructions := parseInstructions(lines[len(lines)-1])
	boardLines := lines[:len(lines)-2]
	width := util.Max(util.Map(boardLines, func(l string) int { return len(l) })...)
	board := make([][]marker, len(boardLines))
	for i, l := range boardLines {
		board[i] = make([]marker, width)
		for j := range board[i] {
			if j < len(l) {
				board[i][j] = markerFromByte(l[j])
			} else {
				board[i][j] = Wrap
			}
		}
	}
	transposedBoard := transpose(board)

	startC, startF := startPos(board)
	pos, f := followInstructions(board, transposedBoard, startC, startF, instructions)
	fmt.Printf("Part 1: %d\n", 1000*(pos.Y+1)+4*(pos.X+1)+f.score())
}
