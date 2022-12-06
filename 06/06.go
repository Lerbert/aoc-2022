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

func findStartOfPacket(message []byte) int {
	for i := 0; i < len(message)-4; i++ {
		if !hasDuplicates(message[i : i+4]) {
			return i + 4
		}
	}
	return -1
}

func main() {
	lines := inp.ReadLines("input")
	message := []byte(lines[0])

	startOfPacket := findStartOfPacket(message)
	fmt.Printf("Part 1: %d\n", startOfPacket)
}
