package main

import (
	inp "aoc2022/input"
	"fmt"
	"regexp"
	"strconv"
)

type crate byte

type Stack[T interface{}] struct {
	s []T
}

func StackFromSlice[T interface{}](es []T) Stack[T] {
	return Stack[T]{
		s: es,
	}
}

func (s *Stack[T]) Length() int {
	return len(s.s)
}

func (s *Stack[T]) Push(es ...T) {
	s.s = append(s.s, es...)
}

func (s *Stack[T]) Pop(n int) []T {
	popped := s.s[len(s.s)-n:]
	for i, j := 0, len(popped)-1; i < j; i, j = i+1, j-1 {
		popped[i], popped[j] = popped[j], popped[i]
	}
	s.s = s.s[:len(s.s)-n]
	return popped
}

func (s *Stack[T]) Peek() T {
	return s.s[len(s.s)-1]
}

type instruction struct {
	from int
	to   int
	cnt  int
}

func (instr instruction) execute(stacks *[]Stack[crate]) {
	(*stacks)[instr.to-1].Push((*stacks)[instr.from-1].Pop(instr.cnt)...)
}

func instructionFromLine(l string) instruction {
	instructionRe := regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)
	matches := instructionRe.FindStringSubmatch(l)
	from, _ := strconv.Atoi(matches[2])
	to, _ := strconv.Atoi(matches[3])
	cnt, _ := strconv.Atoi(matches[1])
	return instruction{
		from,
		to,
		cnt,
	}
}

func splitAt[T comparable](l []T, split T) ([]T, []T) {
	for i, e := range l {
		if e == split {
			return l[:i], l[i+1:]
		}
	}
	return l, make([]T, 0)
}

func topElems(stacks []Stack[crate]) []crate {
	tops := make([]crate, len(stacks))
	for i, s := range stacks {
		tops[i] = s.Peek()
	}
	return tops
}

func main() {
	// Rewrite stacks to lines in input by hand
	lines := inp.ReadLines("input")
	stackLines, instructionLines := splitAt(lines, "")

	stacks := make([]Stack[crate], len(stackLines))
	for i, l := range stackLines {
		stacks[i] = StackFromSlice([]crate(l))
	}

	instructions := make([]instruction, len(instructionLines))
	for i, l := range instructionLines {
		instructions[i] = instructionFromLine(l)
	}

	for _, instr := range instructions {
		instr.execute(&stacks)
		// fmt.Println(stacks)
	}
	fmt.Printf("Part 1: %s\n", topElems(stacks))
}
