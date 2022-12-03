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

func rucksackFromLine(s string) rucksack {
	return []byte(s)
}

func main() {
	lines := inp.ReadLines("input")

	prioSum := uint(0)
	for _, l := range lines {
		rucksack := rucksackFromLine(l)
		duplicate := rucksack.findDuplicate()
		prioSum += priority(duplicate)
	}
	fmt.Printf("Part 1: %d\n", prioSum)
}
