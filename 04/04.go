package main

import (
	inp "aoc2022/input"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type idRange struct {
	lower int
	upper int
}

func (r idRange) contains(number int) bool {
	return r.lower <= number && r.upper >= number
}

func (r idRange) includes(other idRange) bool {
	return other.lower >= r.lower && other.upper <= r.upper
}

func (r idRange) overlaps(other idRange) bool {
	return r.contains(other.lower) || r.contains(other.upper) || other.includes(r)
}

func idRangeFromString(s string) idRange {
	bounds := strings.Split(s, "-")
	lower, errL := strconv.Atoi(bounds[0])
	upper, errU := strconv.Atoi(bounds[1])
	if errL != nil || errU != nil {
		log.Fatal("Could not convert string to int", errL, errU)
	}
	return idRange{
		lower,
		upper,
	}
}

type elfPair struct {
	first  idRange
	second idRange
}

func elfPairFromLine(s string) elfPair {
	ranges := strings.Split(s, ",")
	return elfPair{
		first:  idRangeFromString(ranges[0]),
		second: idRangeFromString(ranges[1]),
	}
}

func main() {
	lines := inp.ReadLines("input")

	pairs := make([]elfPair, len(lines))
	for i, l := range lines {
		pairs[i] = elfPairFromLine(l)
	}

	numIncluded := 0
	for _, p := range pairs {
		if p.first.includes(p.second) || p.second.includes(p.first) {
			numIncluded++
		}
	}
	fmt.Printf("Part 1: %d\n", numIncluded)

	numOverlapping := 0
	for _, p := range pairs {
		if p.first.overlaps(p.second) {
			numOverlapping++
		}
	}
	fmt.Printf("Part 2: %d\n", numOverlapping)
}
