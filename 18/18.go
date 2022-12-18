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

func openSides(cubes []cube, getKeyValue func(cube) (pair, int)) int {
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

func innerCubes(cubes []cube, getKeyValue func(cube) (pair, int), buildCube func(pair, int) cube) (map[cube]struct{}, []map[cube]struct{}) {
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

	innerCubes := make(map[cube]struct{})
	couplings := make([]map[cube]struct{}, 0)
	for k, ints := range cubeMap {
		sort.Ints(ints)
		for i := 0; i < len(ints); i++ {
			if i < len(ints)-1 {
				coupling := make(map[cube]struct{})
				for j := ints[i] + 1; j < ints[i+1]; j++ {
					innerCubes[buildCube(k, j)] = struct{}{}
					coupling[buildCube(k, j)] = struct{}{}
				}
				couplings = append(couplings, coupling)
			}
		}
	}

	return innerCubes, couplings
}

func surfaceArea(cubes []cube) int {
	surface := openSides(cubes, func(c cube) (pair, int) { return pair{c.x, c.y}, c.z })
	surface += openSides(cubes, func(c cube) (pair, int) { return pair{c.x, c.z}, c.y })
	surface += openSides(cubes, func(c cube) (pair, int) { return pair{c.y, c.z}, c.x })
	return surface
}

func enclosedCubes(cubes []cube) []cube {
	enclosedZ, couplingsZ := innerCubes(cubes, func(c cube) (pair, int) { return pair{c.x, c.y}, c.z }, func(p pair, i int) cube { return cube{x: p.fst, y: p.snd, z: i} })
	enclosedY, couplingsY := innerCubes(cubes, func(c cube) (pair, int) { return pair{c.x, c.z}, c.y }, func(p pair, i int) cube { return cube{x: p.fst, y: i, z: p.snd} })
	enclosedX, couplingsX := innerCubes(cubes, func(c cube) (pair, int) { return pair{c.y, c.z}, c.x }, func(p pair, i int) cube { return cube{x: i, y: p.fst, z: p.snd} })

	enclosed := util.SetIntersect(enclosedX, util.SetIntersect(enclosedY, enclosedZ))
	couplings := append(append(couplingsX, couplingsY...), couplingsZ...)

couplings:
	for _, coupling := range couplings {
		for cube := range coupling {
			if _, ok := enclosed[cube]; !ok {
				for toDelete := range coupling {
					delete(enclosed, toDelete)
				}
				continue couplings
			}
		}
	}

	cubeList := make([]cube, len(enclosed))
	i := 0
	for c := range enclosed {
		cubeList[i] = c
		i++
	}
	return cubeList
}

func main() {
	lines := inp.ReadLines("input")
	cubes := util.Map(lines, cubeFromLine)

	surface := surfaceArea(cubes)
	fmt.Printf("Part 1: %d\n", surface)

	enclosed := enclosedCubes(cubes)
	enclosedSurface := surfaceArea(enclosed)
	fmt.Printf("Part 2: %d\n", surface-enclosedSurface)
}
