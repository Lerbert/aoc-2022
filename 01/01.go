package main

import (
	inp "aoc2022/input"
	"fmt"
	"log"
	"math"
	"strconv"
)

const ESTIMATED_INPUTS = 2000
const ESTIMATED_ITEMS_PER_ELF = 8

func main() {
	lines := inp.ReadLines("input")
	elves := splitElves(lines)
	calories_per_elf := make([]int, len(elves))
	for i, elf := range elves {
		calories_per_elf[i] = sum(elf)
	}
	fmt.Println(max(calories_per_elf))
	fmt.Println(sum(maxN(calories_per_elf, 3)))
}

func splitElves(lines []string) [][]int {
	elves := make([][]int, 1, ESTIMATED_INPUTS/ESTIMATED_ITEMS_PER_ELF)
	curElf := 0
	for _, line := range lines {
		if line == "" {
			elves = append(elves, make([]int, 0, ESTIMATED_ITEMS_PER_ELF))
			curElf++
		} else {
			calories, err := strconv.Atoi(line)
			if err != nil {
				log.Fatal("Could not parse int", err)
			}
			elves[curElf] = append(elves[curElf], calories)
		}
	}
	return elves
}

func sum(s []int) int {
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

func max(s []int) int {
	max := math.MinInt
	for _, v := range s {
		if v > max {
			max = v
		}
	}
	return max
}

func maxN(s []int, n int) []int {
	// index i is the top i+1-th element
	maxs := make([]int, n)
	for i := range maxs {
		maxs[i] = math.MinInt
	}
	for _, v := range s {
		comp := v
		swapAlways := false
		for i := range maxs {
			if swapAlways || comp > maxs[i] {
				comp, maxs[i] = maxs[i], comp
				// when we swapped once, we can alwas swap because maxs is sorted --> find correct spot for v and then shift the remaining elements to the right
				swapAlways = true
			}
		}
	}
	return maxs
}
