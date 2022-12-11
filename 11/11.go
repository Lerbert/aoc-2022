package main

import (
	util "aoc2022/aoc"
	inp "aoc2022/input"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

type operation int

const (
	Add operation = iota
	Mul
)

type monkey struct {
	items          []int
	op             operation
	rhs            int // negative number means old
	test           int
	onTrue         int
	onFalse        int
	inspectedItems int
}

func (m *monkey) takeTurn(monkeys *[]monkey) {
	for _, item := range m.items {
		var opRhs int
		if m.rhs < 0 {
			opRhs = item
		} else {
			opRhs = m.rhs
		}
		var worry int
		if m.op == Mul {
			worry = item * opRhs
		} else if m.op == Add {
			worry = item + opRhs
		}
		worry /= 3
		if worry%m.test == 0 {
			(*monkeys)[m.onTrue].items = append((*monkeys)[m.onTrue].items, worry)
		} else {
			(*monkeys)[m.onFalse].items = append((*monkeys)[m.onFalse].items, worry)
		}
	}
	m.inspectedItems += len(m.items)
	m.items = make([]int, 0)
}

func playRound(monkeys *[]monkey) {
	for i := range *monkeys {
		(*monkeys)[i].takeTurn(monkeys)
	}
}

func monkeyBusiness(monkeys *[]monkey) int {
	inspectedItems := util.Map(monkeys, func(m monkey) int { return m.inspectedItems })
	businessScore := util.MaxN(inspectedItems, 2)
	return businessScore[0] * businessScore[1]
}

func parseMonkey(lines []string) monkey {
	itemsS := strings.Split(strings.TrimPrefix(strings.TrimSpace(lines[1]), "Starting items: "), ", ")
	items := make([]int, len(itemsS))
	for i, s := range itemsS {
		item, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal("Could not parse item", err)
		}
		items[i] = item
	}

	opS := strings.TrimPrefix(strings.TrimSpace(lines[2]), "Operation: new = old ")
	var op operation
	if opS[0] == '*' {
		op = Mul
	} else if opS[0] == '+' {
		op = Add
	}
	var rhs int
	if opS[2:] == "old" {
		rhs = -1
	} else {
		rhsP, err := strconv.Atoi(opS[2:])
		if err != nil {
			log.Fatal("Could not parse operation", err)
		}
		rhs = rhsP
	}

	testS := strings.TrimPrefix(strings.TrimSpace(lines[3]), "Test: divisible by ")
	test, err := strconv.Atoi(testS)
	if err != nil {
		log.Fatal("Could not parse test", err)
	}

	onTrueS := strings.TrimPrefix(strings.TrimSpace(lines[4]), "If true: throw to monkey ")
	onTrue, err := strconv.Atoi(onTrueS)
	if err != nil {
		log.Fatal("Could not parse onTrue", err)
	}

	onFalseS := strings.TrimPrefix(strings.TrimSpace(lines[5]), "If false: throw to monkey ")
	onFalse, err := strconv.Atoi(onFalseS)
	if err != nil {
		log.Fatal("Could not parse onFalse", err)
	}

	return monkey{items, op, rhs, test, onTrue, onFalse, 0}
}

func main() {
	lines := inp.ReadLines("input")
	monkeys := make([]monkey, int(math.Ceil(float64(len(lines))/7)))
	for i := range monkeys {
		monkeys[i] = parseMonkey(lines[i*7 : (i+1)*7])
	}

	for round := 0; round < 20; round++ {
		playRound(&monkeys)
	}
	businessScore := monkeyBusiness(&monkeys)
	fmt.Printf("Part 1: %d\n", businessScore)
}
