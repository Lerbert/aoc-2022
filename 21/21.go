package main

import (
	inp "aoc2022/input"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type operation int

const (
	Add operation = iota
	Sub
	Mul
	Div
	Immediate
)

type monkey struct {
	name      string
	op        operation
	lhs       string
	rhs       string
	immediate int
}

var opRegex = regexp.MustCompile(`([a-z]+) (\+|\*|-|/) ([a-z]+)`)

func monkeyFromLine(l string) monkey {
	split := strings.Split(l, ": ")
	name := split[0]
	immediate, err := strconv.Atoi(split[1])
	op := Immediate
	lhs := ""
	rhs := ""
	if err != nil {
		immediate = 0
		matches := opRegex.FindStringSubmatch(split[1])
		lhs = matches[1]
		rhs = matches[3]
		switch matches[2] {
		case "+":
			op = Add
		case "*":
			op = Mul
		case "-":
			op = Sub
		case "/":
			op = Div
		}
	}
	return monkey{name: name, op: op, lhs: lhs, rhs: rhs, immediate: immediate}
}

func findMonkeyNumber(monkeys map[string]monkey, name string) int {
	memo := make(map[string]int)
	return findMonkeyNumberRec(monkeys, name, memo)
}

func findMonkeyNumberRec(monkeys map[string]monkey, currentMonkeyName string, memo map[string]int) int {
	n, ok := memo[currentMonkeyName]
	if ok {
		return n
	}

	currentMonkey := monkeys[currentMonkeyName]

	switch currentMonkey.op {
	case Immediate:
		n = currentMonkey.immediate
	case Add:
		n = findMonkeyNumberRec(monkeys, currentMonkey.lhs, memo) + findMonkeyNumberRec(monkeys, currentMonkey.rhs, memo)
	case Mul:
		n = findMonkeyNumberRec(monkeys, currentMonkey.lhs, memo) * findMonkeyNumberRec(monkeys, currentMonkey.rhs, memo)
	case Sub:
		n = findMonkeyNumberRec(monkeys, currentMonkey.lhs, memo) - findMonkeyNumberRec(monkeys, currentMonkey.rhs, memo)
	case Div:
		n = findMonkeyNumberRec(monkeys, currentMonkey.lhs, memo) / findMonkeyNumberRec(monkeys, currentMonkey.rhs, memo)
	}

	memo[currentMonkeyName] = n
	return n
}

// func hasHumn(monkeys map[string]monkey, name string) bool {
// 	memo := make(map[string]bool)
// 	return hasHumnRec(monkeys, name, memo)
// }

func hasHumnRec(monkeys map[string]monkey, currentMonkeyName string, memo map[string]bool) bool {
	b, ok := memo[currentMonkeyName]
	if ok {
		return b
	}

	currentMonkey := monkeys[currentMonkeyName]
	if currentMonkey.op == Immediate {
		b = currentMonkeyName == HUMN
	} else {
		b = hasHumnRec(monkeys, currentMonkey.lhs, memo) || hasHumnRec(monkeys, currentMonkey.rhs, memo)
	}

	memo[currentMonkeyName] = b
	return b
}

func findHumnNumber(monkeys map[string]monkey, start string) int {
	memoNumbers := make(map[string]int)
	memoHasHumn := make(map[string]bool)
	startMonkey := monkeys[start]
	var matchNumber int
	var monkey string
	if hasHumnRec(monkeys, startMonkey.lhs, memoHasHumn) {
		matchNumber = findMonkeyNumberRec(monkeys, startMonkey.rhs, memoNumbers)
		monkey = startMonkey.lhs
	} else {
		matchNumber = findMonkeyNumberRec(monkeys, startMonkey.lhs, memoNumbers)
		monkey = startMonkey.rhs
	}
	return findHumnNumberRec(monkeys, monkey, matchNumber, memoNumbers, memoHasHumn)
}

func findHumnNumberRec(monkeys map[string]monkey, currentMonkeyName string, matchNumber int, memoNumbers map[string]int, memoHasHumn map[string]bool) int {
	currentMonkey := monkeys[currentMonkeyName]

	if currentMonkey.op == Immediate {
		if currentMonkeyName == HUMN {
			return matchNumber
		} else {
			return currentMonkey.immediate
		}
	}

	// humn is queried exactly once, so it must appear in exactly one of lhs and rhs
	humn, noHumn := currentMonkey.lhs, currentMonkey.rhs
	if hasHumnRec(monkeys, noHumn, memoHasHumn) {
		humn, noHumn = noHumn, humn
	}

	noHumnNumber := findMonkeyNumberRec(monkeys, noHumn, memoNumbers)
	var nextMatchNumber int
	switch currentMonkey.op {
	case Add:
		nextMatchNumber = matchNumber - noHumnNumber
	case Mul:
		nextMatchNumber = matchNumber / noHumnNumber
	case Sub:
		if humn == currentMonkey.lhs {
			nextMatchNumber = matchNumber + noHumnNumber
		} else {
			nextMatchNumber = noHumnNumber - matchNumber
		}
	case Div:
		if humn == currentMonkey.lhs {
			nextMatchNumber = matchNumber * noHumnNumber
		} else {
			nextMatchNumber = noHumnNumber / matchNumber
		}
	}
	return findHumnNumberRec(monkeys, humn, nextMatchNumber, memoNumbers, memoHasHumn)
}

const (
	ROOT = "root"
	HUMN = "humn"
)

func main() {
	lines := inp.ReadLines("input")
	monkeys := make(map[string]monkey)
	for _, l := range lines {
		m := monkeyFromLine(l)
		monkeys[m.name] = m
	}

	rootNumber := findMonkeyNumber(monkeys, ROOT)
	fmt.Printf("Part 1: %d\n", rootNumber)

	humnNumber := findHumnNumber(monkeys, ROOT)
	fmt.Printf("Part 2: %d\n", humnNumber)
}
