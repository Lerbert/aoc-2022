package main

import (
	util "aoc2022/aoc"
	inp "aoc2022/input"
	"fmt"
)

type uniqueNumber struct {
	id int
	n  int
}

func wrap(index int, length int) int {
	if mod := index % length; mod >= 0 {
		return mod
	} else {
		return length + mod
	}
}

func mix(numbers []uniqueNumber, rounds int) ([]uniqueNumber, map[uniqueNumber]int) {
	numbersLen := len(numbers)
	mixedIndices := make(map[uniqueNumber]int)
	mixed := make([]uniqueNumber, numbersLen)
	copy(mixed, numbers)
	for i, x := range numbers {
		mixedIndices[x] = i
	}

	for i := 0; i < rounds; i++ {
		mixingRound(numbers, mixed, mixedIndices)
	}

	return mixed, mixedIndices
}

func mixingRound(numbers []uniqueNumber, mixed []uniqueNumber, mixedIndices map[uniqueNumber]int) {
	numbersLen := len(numbers)
	for _, x := range numbers {
		xIndex := mixedIndices[x]
		offset := x.n % (numbersLen - 1)
		if offset >= 0 {
			for i := xIndex; i < xIndex+offset; i++ {
				fst, snd := mixed[wrap(i, numbersLen)], mixed[wrap(i+1, numbersLen)]
				mixed[wrap(i, numbersLen)], mixed[wrap(i+1, numbersLen)] = snd, fst
				mixedIndices[fst], mixedIndices[snd] = mixedIndices[snd], mixedIndices[fst]
			}
		} else {
			for i := xIndex; i > xIndex+offset; i-- {
				fst, snd := mixed[wrap(i, numbersLen)], mixed[wrap(i-1, numbersLen)]
				mixed[wrap(i, numbersLen)], mixed[wrap(i-1, numbersLen)] = snd, fst
				mixedIndices[fst], mixedIndices[snd] = mixedIndices[snd], mixedIndices[fst]
			}
		}
	}
}

func coordinateSum(mixed []uniqueNumber, mixedIndices map[uniqueNumber]int) int {
	coordinates := 0
	for i := 1; i <= 3; i++ {
		coordinates += mixed[wrap(mixedIndices[uniqueNumber{id: 0, n: 0}]+i*1000, len(mixed))].n
	}
	return coordinates
}

const DECRYPTION_KEY = 811589153
const ROUNDS = 10

func main() {
	lines := inp.ReadLines("input")
	numbers := make([]uniqueNumber, len(lines))
	for i, l := range lines {
		n := inp.MustAtoi(l)
		id := i + 1
		if n == 0 {
			id = 0
		}
		numbers[i] = uniqueNumber{id: id, n: n}
	}

	mixed, mixedIndices := mix(numbers, 1)
	coordinates := coordinateSum(mixed, mixedIndices)
	fmt.Printf("Part 1: %d\n", coordinates)

	// Apply decryption key
	numbers = util.Map(numbers, func(u uniqueNumber) uniqueNumber { return uniqueNumber{id: u.id, n: u.n * DECRYPTION_KEY} })

	mixed, mixedIndices = mix(numbers, ROUNDS)
	coordinates = coordinateSum(mixed, mixedIndices)
	fmt.Printf("Part 2: %d\n", coordinates)
}
