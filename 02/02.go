package main

import (
	inp "aoc2022/input"
	"fmt"
	"log"
)

type symbol int

const (
	Rock symbol = iota
	Paper
	Scissors
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

func (s1 symbol) beats(s2 symbol) byte {
	if s1 == s2 {
		return 'Y'
	}
	switch s1 {
	case Rock:
		switch s2 {
		case Paper:
			return 'X'
		case Scissors:
			return 'Z'
		}
	case Paper:
		switch s2 {
		case Scissors:
			return 'X'
		case Rock:
			return 'Z'
		}
	case Scissors:
		switch s2 {
		case Rock:
			return 'X'
		case Paper:
			return 'Z'
		}
	}
	log.Fatal("Unknown symbols", s1, s2)
	panic("")
}

func (s symbol) playToGet(desiredOutcome byte) symbol {
	switch desiredOutcome {
	case 'X':
		switch s {
		case Rock:
			return Scissors
		case Paper:
			return Rock
		case Scissors:
			return Paper
		}
	case 'Y':
		return s
	case 'Z':
		switch s {
		case Rock:
			return Paper
		case Paper:
			return Scissors
		case Scissors:
			return Rock
		}
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
	case 'X':
		return 0
	case 'Y':
		return 3
	case 'Z':
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
		player:   op.playToGet(bytes[2]),
	}
}

func main() {
	lines := inp.ReadLines("input")
	score := uint(0)
	for _, l := range lines {
		strategy := strategyFromLine2(l)
		score += strategy.score()
	}
	fmt.Println(score)
}
