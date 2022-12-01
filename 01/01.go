package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

const ESTIMATED_INPUTS = 2000
const ESTIMATED_ITEMS_PER_ELF = 8

func main() {
	lines := readLines("input")
	elves := splitElves(lines)
	calories_per_elf := make([]int, len(elves))
	for i, elf := range elves {
		calories_per_elf[i] = sum(elf)
	}
	fmt.Println(max(calories_per_elf))
}

func readLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Could not open file", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := make([]string, 0, ESTIMATED_INPUTS)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
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
