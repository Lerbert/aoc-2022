package main

import (
	inp "aoc2022/input"
	"fmt"
	"log"
)

func priority(b byte) uint {
	if b <= 'Z' && b >= 'A' {
		return uint(b-'A') + 27
	} else if b <= 'z' && b >= 'a' {
		return uint(b-'a') + 1
	} else {
		log.Fatal("Non alpha symbol", b)
		panic("")
	}
}

type rucksack []byte

func (r rucksack) compartment1() []byte {
	return r[:len(r)/2]
}

func (r rucksack) compartment2() []byte {
	return r[len(r)/2:]
}

func (r rucksack) findDuplicate() byte {
	itemsInC1 := make(map[byte]struct{})
	for _, item := range r.compartment1() {
		itemsInC1[item] = struct{}{}
	}
	for _, item := range r.compartment2() {
		if _, ok := itemsInC1[item]; ok {
			return item
		}
	}
	log.Fatal("No duplicate item found")
	panic("")
}

func (r rucksack) uniqueItems() []byte {
	m := make(map[byte]struct{})
	for _, item := range r {
		m[item] = struct{}{}
	}
	keys := make([]byte, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func rucksackFromLine(s string) rucksack {
	return []byte(s)
}

func groupBy[T interface{}](s []T, n int) [][]T {
	grouped := make([][]T, 0, len(s)/n)
	groupIndex := -1
	for i, e := range s {
		if i%n == 0 {
			groupIndex++
			grouped = append(grouped, make([]T, 0, n))
		}
		grouped[groupIndex] = append(grouped[groupIndex], e)
	}
	return grouped
}

func findRucksackBadge(rs []rucksack) byte {
	foundItems := make(map[byte]int)
	numRucksacks := len(rs)
	for _, r := range rs {
		for _, b := range r.uniqueItems() {
			val, ok := foundItems[b]
			if ok {
				val++
				if val == numRucksacks {
					return b
				} else {
					foundItems[b] = val
				}
			} else {
				foundItems[b] = 1
			}
		}
	}
	log.Fatal("No duplicate item found")
	panic("")
}

func main() {
	lines := inp.ReadLines("input")

	rucksacks := make([]rucksack, len(lines))
	for i, l := range lines {
		rucksacks[i] = rucksackFromLine(l)
	}

	prioSum1 := uint(0)
	for _, r := range rucksacks {
		duplicate := r.findDuplicate()
		prioSum1 += priority(duplicate)
	}
	fmt.Printf("Part 1: %d\n", prioSum1)

	groupedElves := groupBy(rucksacks, 3)
	prioSum2 := uint(0)
	for _, g := range groupedElves {
		badge := findRucksackBadge(g)
		prioSum2 += priority(badge)
	}
	fmt.Printf("Part 2: %d\n", prioSum2)
}
