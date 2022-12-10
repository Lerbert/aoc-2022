package main

import (
	inp "aoc2022/input"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type executer interface {
	duration() int
	execute(*cpu)
}

type cpu struct {
	pc int
	x  int
}

func (c *cpu) executeProgram(prog []executer, recordInterval int) map[int]int {
	xVals := make(map[int]int)
	c.pc = 1
	c.x = 1
	for _, instr := range prog {
		oldPc := c.pc
		c.pc += instr.duration()
		if oldPc/recordInterval < c.pc/recordInterval {
			xVals[c.pc/recordInterval*recordInterval] = c.x
		}
		instr.execute(c)
		if c.pc%recordInterval == 0 {
			xVals[c.pc] = c.x
		}
	}
	return xVals
}

type noop struct{}

func (noop) duration() int {
	return 1
}

func (noop) execute(c *cpu) {
	// intentionally left blank
}

type addx struct {
	incr int
}

func (addx) duration() int {
	return 2
}

func (a addx) execute(c *cpu) {
	c.x += a.incr
}

func parseInstruction(s string) executer {
	if s == "noop" {
		return noop{}
	} else if strings.HasPrefix(s, "addx ") {
		incr, err := strconv.Atoi(s[len("addx "):])
		if err != nil {
			log.Fatal("Could not parse increment", err)
		}
		return addx{incr}
	} else {
		log.Fatal("Unknown instruction", s)
		panic("")
	}
}

func main() {
	lines := inp.ReadLines("input")
	instructions := make([]executer, len(lines))
	for i, l := range lines {
		instructions[i] = parseInstruction(l)
	}
	c := cpu{pc: 0, x: 1}

	xVals := c.executeProgram(instructions, 20)
	signalStrength := 0
	for cycle, x := range xVals {
		if cycle%20 != 0 {
			log.Fatal("Error in recording")
		}
		if cycle <= 220 && (cycle-20)%40 == 0 {
			signalStrength += cycle * x
		}
	}
	fmt.Printf("Part 1: %d\n", signalStrength)
}
