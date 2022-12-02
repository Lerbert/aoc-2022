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

func (s1 symbol) beats(s2 symbol) int {
	if s1 == s2 {
		return 0
	}
	switch s1 {
	case Rock:
		switch s2 {
		case Paper:
			return -1
		case Scissors:
			return 1
		}
	case Paper:
		switch s2 {
		case Scissors:
			return -1
		case Rock:
			return 1
		}
	case Scissors:
		switch s2 {
		case Rock:
			return -1
		case Paper:
			return 1
		}
	}
	log.Fatal("Unknown symbols", s1, s2)
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
	case -1:
		return 0
	case 0:
		return 3
	case 1:
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

func main() {
	lines := inp.ReadLines("input")
	score := uint(0)
	for _, l := range lines {
		strategy := strategyFromLine(l)
		score += strategy.score()
	}
	fmt.Println(score)
}
