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

type crt [6][40]bool

func (c *crt) print() {
	for _, row := range c {
		for _, lit := range row {
			if lit {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}

type cpu struct {
	cycle int
	pc    int
	x     int
}

func (c *cpu) executeProgram(prog []executer, screen *crt, recordInterval int) map[int]int {
	xVals := make(map[int]int)
	c.cycle = 0
	c.pc = 0
	c.x = 1

	fetch := 0
	var instr executer = noop{}

	for c.pc < len(prog) {
		if c.cycle%recordInterval == 0 {
			xVals[c.cycle] = c.x
		}
		if c.cycle > 0 {
			row := (c.cycle - 1) / 40
			col := (c.cycle - 1) % 40
			(*screen)[row][col] = c.x-col <= 1 && c.x-col >= -1
		}
		if fetch == 0 {
			instr.execute(c)
			instr = prog[c.pc]
			fetch = instr.duration()
			c.pc++
		}
		fetch--
		c.cycle++
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
	var screen crt

	xVals := c.executeProgram(instructions, &screen, 20)
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
	fmt.Println("Part 2:")
	screen.print()
}
