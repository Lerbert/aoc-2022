package main

import (
	util "aoc2022/aoc"
	inp "aoc2022/input"
	"fmt"
	"strings"
)

type direction int

const (
	North direction = iota
	West
	South
	East
)

func directionFromByte(b byte) direction {
	switch b {
	case '^':
		return North
	case '<':
		return West
	case 'v':
		return South
	case '>':
		return East
	default:
		panic("Unknown direction")
	}
}

type blizzard struct {
	facing    direction
	pos       util.Coord
	mapWidth  int
	mapHeight int
}

func (b *blizzard) move() {
	switch b.facing {
	case North:
		b.pos = util.Coord{X: b.pos.X, Y: util.WrapIndex(b.pos.Y-1, b.mapHeight)}
	case West:
		b.pos = util.Coord{X: util.WrapIndex(b.pos.X-1, b.mapWidth), Y: b.pos.Y}
	case South:
		b.pos = util.Coord{X: b.pos.X, Y: util.WrapIndex(b.pos.Y+1, b.mapHeight)}
	case East:
		b.pos = util.Coord{X: util.WrapIndex(b.pos.X+1, b.mapWidth), Y: b.pos.Y}
	}
}

func moveBlizzards(blizzards []blizzard) {
	for i := range blizzards {
		(&blizzards[i]).move()
	}
}

func bfs(blizzards []blizzard, start util.Coord, end util.Coord) int {
	toVisit := map[util.Coord]struct{}{start: {}}
	return bfsRec(blizzards, toVisit, start, end)
}

func bfsRec(blizzards []blizzard, toVisit map[util.Coord]struct{}, start util.Coord, end util.Coord) int {
	for steps := 0; ; steps++ {
		if len(toVisit) == 0 {
			panic("End not reachable")
		}

		moveBlizzards(blizzards)
		blizzardPositions := util.Map(blizzards, func(b blizzard) util.Coord { return b.pos })

		mapWidth, mapHeight := blizzards[0].mapWidth, blizzards[0].mapHeight

		visitNext := make(map[util.Coord]struct{}) // Use map to avoid duplicate insertions
		for c := range toVisit {
			// Look at all neighbors + waiting
			consider := append(c.OrthogonalNeighbors(), c)
			for _, n := range consider {
				if n == end {
					// End never has a blizzard, so we can reach it in the next step
					// Quitting here leaves the blizzards in the right position
					return steps + 1
				}
				if (util.Between(n.X, 0, mapWidth-1) && util.Between(n.Y, 0, mapHeight-1)) || n == start {
					// Valid neighbor
					if !util.Contains(blizzardPositions, n) {
						visitNext[n] = struct{}{}
					}
				}
			}
		}
		toVisit = visitNext
	}
}

func parseBlizzards(lines []string) []blizzard {
	blizzards := make([]blizzard, 0)
	width := len(lines[0]) - 2 // Walls on all edges
	height := len(lines) - 2
	for y, l := range lines {
		for x, b := range []byte(l) {
			if b == '^' || b == '<' || b == 'v' || b == '>' {
				// Subtract one from coordinates so that (0, 0) is the top left corner of the area inside the walls
				blizzards = append(blizzards, blizzard{facing: directionFromByte(b), pos: util.Coord{X: x - 1, Y: y - 1}, mapWidth: width, mapHeight: height})
			}
		}
	}
	return blizzards
}

func main() {
	lines := inp.ReadLines("input")
	blizzards := parseBlizzards(lines)
	// Subtract one from coordinates so that (0, 0) is the top left corner of the area inside the walls
	start := util.Coord{Y: -1, X: strings.IndexByte(lines[0], '.') - 1}
	end := util.Coord{Y: len(lines) - 2, X: strings.IndexByte(lines[len(lines)-1], '.') - 1}

	minutesSE1 := bfs(blizzards, start, end)
	fmt.Printf("Part 1: %d\n", minutesSE1)

	minutesES := bfs(blizzards, end, start)
	minutesSE2 := bfs(blizzards, start, end)
	fmt.Printf("Part 2: %d\n", minutesSE1+minutesES+minutesSE2)
}
