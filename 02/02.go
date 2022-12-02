package main

import (
	inp "aoc2022/input"
	"fmt"
	"log"
)

type outcome byte

const (
	Lose outcome = 'X'
	Draw outcome = 'Y'
	Win  outcome = 'Z'
)

type symbol int

const (
	Rock        symbol = 0
	Paper       symbol = 1
	Scissors    symbol = 2
	NUM_SYMBOLS        = 3
)

func (s symbol) value() uint {
	switch s {
	case Rock:
		return 1
	case Paper:
		return 2
	case Scissors:
		return 3
	}
	log.Fatal("Unknown symbol", s)
	panic("")
}

func (s1 symbol) beats(s2 symbol) outcome {
	if s1 == s2 {
		return Draw
	}
	if (int(s1)+1)%NUM_SYMBOLS == int(s2) {
		return Lose
	} else {
		return Win
	}
}

func (s symbol) playToGet(desiredOutcome outcome) symbol {
	switch desiredOutcome {
	case Lose:
		// return symbol((int(s) - 1) % NUM_SYMBOLS) // doesn't work because in Go -1 % 3 == -1
		return symbol((int(s) + NUM_SYMBOLS - 1) % NUM_SYMBOLS)
	case Draw:
		return s
	case Win:
		return symbol((int(s) + 1) % NUM_SYMBOLS)
	}
	log.Fatal("Unknown symbol or outcome", s, desiredOutcome)
	panic("")
}

func symbolFromChar(b byte) symbol {
	switch b {
	case 'A', 'X':
		return Rock
	case 'B', 'Y':
		return Paper
	case 'C', 'Z':
		return Scissors
	}
	log.Fatal("Unknown character for symbol", b)
	panic("")
}

type strategy struct {
	opponent symbol
	player   symbol
}

func (s strategy) outcome() uint {
	switch s.player.beats(s.opponent) {
	case Lose:
		return 0
	case Draw:
		return 3
	case Win:
		return 6
	}
	log.Fatal("Result not known")
	panic("")
}

func (s strategy) score() uint {
	return s.outcome() + s.player.value()
}

func strategyFromLine(s string) strategy {
	bytes := []byte(s)
	return strategy{
		opponent: symbolFromChar(bytes[0]),
		player:   symbolFromChar(bytes[2]),
	}
}

func strategyFromLine2(s string) strategy {
	bytes := []byte(s)
	op := symbolFromChar(bytes[0])
	return strategy{
		opponent: op,
		player:   op.playToGet(outcome(bytes[2])),
	}
}

func main() {
	lines := inp.ReadLines("input")

	score1 := uint(0)
	for _, l := range lines {
		strategy := strategyFromLine(l)
		score1 += strategy.score()
	}
	fmt.Printf("Part 1: %d\n", score1)

	score2 := uint(0)
	for _, l := range lines {
		strategy := strategyFromLine2(l)
		score2 += strategy.score()
	}
	fmt.Printf("Part 2: %d\n", score2)
}
