package main

import (
	util "aoc2022/aoc"
	inp "aoc2022/input"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

type pair struct {
	fst int
	snd int
}

type cube struct {
	x int
	y int
	z int
}

func cubeFromLine(l string) cube {
	split := strings.Split(l, ",")
	x, errX := strconv.Atoi(split[0])
	y, errY := strconv.Atoi(split[1])
	z, errZ := strconv.Atoi(split[2])

	if errX != nil || errY != nil || errZ != nil {
		log.Fatal("Could not parse coordinates")
	}

	return cube{x: x, y: y, z: z}
}

func countOpenSides(cubes []cube, getKeyValue func(cube) (pair, int)) int {
	cubeMap := make(map[pair][]int)
	for _, c := range cubes {
		k, v := getKeyValue(c)
		l, ok := cubeMap[k]
		if ok {
			cubeMap[k] = append(l, v)
		} else {
			cubeMap[k] = []int{v}
		}
	}

	openSides := 0
	for _, ints := range cubeMap {
		sort.Ints(ints)
		openSides += 1 // first cube always has front side open
		for i := 0; i < len(ints); i++ {
			if i < len(ints)-1 && ints[i]+1 < ints[i+1] {
				openSides += 2 // back side of this cube and front side of next cube
			} else if i == len(ints)-1 {
				openSides += 1 // last cube always has back side open
			}
		}
	}

	return openSides
}

func surfaceArea(cubes []cube) int {
	surface := countOpenSides(cubes, func(c cube) (pair, int) { return pair{c.x, c.y}, c.z })
	surface += countOpenSides(cubes, func(c cube) (pair, int) { return pair{c.x, c.z}, c.y })
	surface += countOpenSides(cubes, func(c cube) (pair, int) { return pair{c.y, c.z}, c.x })
	return surface
}

func main() {
	lines := inp.ReadLines("input")
	cubes := util.Map(lines, cubeFromLine)

	surface := surfaceArea(cubes)
	fmt.Printf("Part 1: %d\n", surface)
}
