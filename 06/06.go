package main

import (
	inp "aoc2022/input"
	"fmt"
)

func hasDuplicates[T comparable](s []T) bool {
	m := make(map[T]struct{})
	for _, e := range s {
		_, ok := m[e]
		if ok {
			return true
		} else {
			m[e] = struct{}{}
		}
	}
	return false
}

func findDistincts(message []byte, n int) int {
	for i := 0; i < len(message)-n; i++ {
		if !hasDuplicates(message[i : i+n]) {
			return i + n
		}
	}
	return -1
}

func findStartOfPacket(message []byte) int {
	return findDistincts(message, 4)
}

func findStartOfMessage(message []byte) int {
	return findDistincts(message, 14)
}

func main() {
	lines := inp.ReadLines("input")
	message := []byte(lines[0])

	startOfPacket := findStartOfPacket(message)
	fmt.Printf("Part 1: %d\n", startOfPacket)

	startOfMessage := findStartOfMessage(message)
	fmt.Printf("Part 2: %d\n", startOfMessage)
}
