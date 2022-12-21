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

func main() {
	lines := inp.ReadLines("input")
	monkeys := make(map[string]monkey)
	for _, l := range lines {
		m := monkeyFromLine(l)
		monkeys[m.name] = m
	}

	rootNumber := findMonkeyNumber(monkeys, "root")
	fmt.Printf("Part 1: %d\n", rootNumber)
}
