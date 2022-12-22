package main

import (
	util "aoc2022/aoc"
	inp "aoc2022/input"
	"fmt"
	"log"
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
	execute2(board [][]marker, transposedBoard [][]marker, pos util.Coord, f facing) (util.Coord, facing)
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

func (t turn) execute2(board [][]marker, transposedBoard [][]marker, pos util.Coord, f facing) (util.Coord, facing) {
	return t.execute(board, transposedBoard, pos, f)
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

func (fw forward) execute2(board [][]marker, transposedBoard [][]marker, pos util.Coord, f facing) (util.Coord, facing) {
	cnt := int(fw)
	var lastPos util.Coord
	var lastF facing
	for i := 0; i < cnt; i++ {
		lastPos = pos
		lastF = f
		pos = incrementCoord(pos, f)
		pos, f = wrapCube(board, pos, f, lastPos)
		if board[pos.Y][pos.X] == Wall {
			return lastPos, lastF
		}
	}
	return pos, f
}

func incrementCoord(pos util.Coord, f facing) util.Coord {
	switch f {
	case North:
		return util.Coord{X: pos.X, Y: pos.Y - 1}
	case West:
		return util.Coord{X: pos.X - 1, Y: pos.Y}
	case South:
		return util.Coord{X: pos.X, Y: pos.Y + 1}
	case East:
		return util.Coord{X: pos.X + 1, Y: pos.Y}
	default:
		panic("Unknown direction")
	}
}

func wrapCube(board [][]marker, pos util.Coord, f facing, lastPos util.Coord) (util.Coord, facing) {
	width, height := len(board[0]), len(board)
	if pos.X < 0 || pos.X >= width || pos.Y < 0 || pos.Y >= height || board[pos.Y][pos.X] == Wrap {
		return leave(currentFace(lastPos), lastPos, f)
	} else {
		return pos, f
	}
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

type cubeFace int

const (
	Top cubeFace = iota
	Bottom
	Front
	Back
	LeftFace
	RightFace
	NUM_FACES int = iota
)

// Net
//
//	TE
//	S
//
// WB
// N
var (
	TOP_FACE     = util.Coord{X: 50, Y: 0}
	BOTTOM_FACE  = util.Coord{X: 50, Y: 100}
	FRONT_FACE   = util.Coord{X: 50, Y: 50} // South
	BACK_FACE    = util.Coord{X: 0, Y: 150} // North
	LEFT_FACE    = util.Coord{X: 0, Y: 100} // West
	RIGHT_FACE   = util.Coord{X: 100, Y: 0} // East
	FACE_CORNERS = [NUM_FACES]util.Coord{TOP_FACE, BOTTOM_FACE, FRONT_FACE, BACK_FACE, LEFT_FACE, RIGHT_FACE}
)

const CUBE_SIZE = 50

func currentFace(pos util.Coord) cubeFace {
	if TOP_FACE.X <= pos.X && pos.X < TOP_FACE.X+CUBE_SIZE && TOP_FACE.Y <= pos.Y && pos.Y < TOP_FACE.Y+CUBE_SIZE {
		return Top
	} else if BOTTOM_FACE.X <= pos.X && pos.X < BOTTOM_FACE.X+CUBE_SIZE && BOTTOM_FACE.Y <= pos.Y && pos.Y < BOTTOM_FACE.Y+CUBE_SIZE {
		return Bottom
	} else if FRONT_FACE.X <= pos.X && pos.X < FRONT_FACE.X+CUBE_SIZE && FRONT_FACE.Y <= pos.Y && pos.Y < FRONT_FACE.Y+CUBE_SIZE {
		return Front
	} else if BACK_FACE.X <= pos.X && pos.X < BACK_FACE.X+CUBE_SIZE && BACK_FACE.Y <= pos.Y && pos.Y < BACK_FACE.Y+CUBE_SIZE {
		return Back
	} else if LEFT_FACE.X <= pos.X && pos.X < LEFT_FACE.X+CUBE_SIZE && LEFT_FACE.Y <= pos.Y && pos.Y < LEFT_FACE.Y+CUBE_SIZE {
		return LeftFace
	} else if RIGHT_FACE.X <= pos.X && pos.X < RIGHT_FACE.X+CUBE_SIZE && RIGHT_FACE.Y <= pos.Y && pos.Y < RIGHT_FACE.Y+CUBE_SIZE {
		return RightFace
	} else {
		log.Fatal("Position is not on cube ", pos)
		panic("")
	}
}

func relativeCoord(face cubeFace, pos util.Coord) util.Coord {
	return util.Coord{X: pos.X - FACE_CORNERS[face].X, Y: pos.Y - FACE_CORNERS[face].Y}
}

// func normal(from cubeFace, to cubeFace, pos util.Coord) util.Coord {
// 	relative := relativeCoord(from, pos)
// 	return util.Coord{X: FACE_CORNERS[to].X + relative.X, Y: FACE_CORNERS[to].Y + relative.Y}
// }

func crossover(from cubeFace, to cubeFace, pos util.Coord) util.Coord {
	relative := relativeCoord(from, pos)
	return util.Coord{X: FACE_CORNERS[to].X + relative.Y, Y: FACE_CORNERS[to].Y + relative.X}
}

// func invertedX(from cubeFace, to cubeFace, pos util.Coord) util.Coord {
// 	relative := relativeCoord(from, pos)
// 	return util.Coord{X: FACE_CORNERS[to].X + CUBE_SIZE - 1 - relative.X, Y: FACE_CORNERS[to].Y + relative.Y}
// }

func invertedY(from cubeFace, to cubeFace, pos util.Coord) util.Coord {
	relative := relativeCoord(from, pos)
	return util.Coord{X: FACE_CORNERS[to].X + relative.X, Y: FACE_CORNERS[to].Y + CUBE_SIZE - 1 - relative.Y}
}

func leave(face cubeFace, pos util.Coord, f facing) (util.Coord, facing) {
	switch face {
	case Top:
		switch f {
		case West:
			return invertedY(Top, LeftFace, pos), East
		case North:
			return crossover(Top, Back, pos), East
		}
	case Bottom:
		switch f {
		case South:
			return crossover(Bottom, Back, pos), West
		case East:
			return invertedY(Bottom, RightFace, pos), West
		}
	case Front:
		switch f {
		case West:
			return crossover(Front, LeftFace, pos), South
		case East:
			return crossover(Front, RightFace, pos), North
		}
	case Back:
		switch f {
		case East:
			return crossover(Back, Bottom, pos), North
		case South:
			return invertedY(Back, RightFace, pos), South
		case West:
			return crossover(Back, Top, pos), South
		}
	case LeftFace:
		switch f {
		case North:
			return crossover(LeftFace, Front, pos), East
		case West:
			return invertedY(LeftFace, Top, pos), East
		}
	case RightFace:
		switch f {
		case North:
			return invertedY(RightFace, Back, pos), North
		case East:
			return invertedY(RightFace, Bottom, pos), West
		case South:
			return crossover(RightFace, Front, pos), West
		}
	}
	log.Fatal("Transition should not be possible ", face, " ", f)
	panic("")
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

func followInstructions2(board [][]marker, transposedBoard [][]marker, startC util.Coord, startF facing, instructions []instruction) (util.Coord, facing) {
	pos := startC
	f := startF
	for _, instr := range instructions {
		pos, f = instr.execute2(board, transposedBoard, pos, f)
	}
	return pos, f
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

func password(pos util.Coord, f facing) int {
	return 1000*(pos.Y+1) + 4*(pos.X+1) + f.score()
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
	fmt.Printf("Part 1: %d\n", password(pos, f))

	pos, f = followInstructions2(board, transposedBoard, startC, startF, instructions)
	fmt.Printf("Part 2: %d\n", password(pos, f))
}
