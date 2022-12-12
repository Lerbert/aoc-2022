package main

import (
	inp "aoc2022/input"
	"fmt"
	"log"
	"strings"
)

type coord struct {
	x int
	y int
}

func (c *coord) neighbors(lenX int, lenY int) []coord {
	neighbors := make([]coord, 0, 4)
	if c.y > 0 {
		neighbors = append(neighbors, coord{c.x, c.y - 1})
	}
	if c.y < lenY-1 {
		neighbors = append(neighbors, coord{c.x, c.y + 1})
	}
	if c.x > 0 {
		neighbors = append(neighbors, coord{c.x - 1, c.y})
	}
	if c.x < lenX-1 {
		neighbors = append(neighbors, coord{c.x + 1, c.y})
	}
	return neighbors
}

type node struct {
	c       coord
	height  int
	visited bool
	end     bool
}

func makeNode(height byte, x int, y int, end bool) node {
	var intHeight int
	if height == 'S' {
		intHeight = 0
	} else if height == 'E' {
		intHeight = 25
	} else {
		intHeight = int(height - 'a')
	}
	return node{coord{x, y}, intHeight, false, end}
}

func bfs(heightMap *[][]node, init coord, isEdge func(from *node, to *node) bool) int {
	toVisit := map[coord]struct{}{init: {}}
	return bfsRec(heightMap, toVisit, isEdge)
}

func bfsRec(heightMap *[][]node, toVisit map[coord]struct{}, isEdge func(from *node, to *node) bool) int {
	for steps := 0; len(toVisit) > 0 && steps < 1000; steps++ {
		visitNext := make(map[coord]struct{}) // Use map to avoid duplicate insertions
		for n := range toVisit {
			curNode := &(*heightMap)[n.y][n.x]
			curNode.visited = true
			if curNode.end {
				return steps
			}
			// Look at all neighbors
			for _, c := range n.neighbors(len((*heightMap)[0]), len(*heightMap)) {
				neighborNode := &(*heightMap)[c.y][c.x]
				if !neighborNode.visited && isEdge(curNode, neighborNode) {
					visitNext[c] = struct{}{}
				}
			}
		}
		toVisit = visitNext
	}
	log.Fatal("no reachable end node")
	panic("")
}

func isEdge1(from *node, to *node) bool {
	return to.height <= from.height+1
}

func isEdge2(from *node, to *node) bool {
	return from.height <= to.height+1
}

func main() {
	lines := inp.ReadLines("input")
	heightMap := make([][]node, len(lines))

	var init coord
	for i, l := range lines {
		sIndex := strings.Index(l, "S")
		if sIndex != -1 {
			init = coord{sIndex, i}
		}
		heightMap[i] = make([]node, len(l))
		for j, b := range []byte(l) {
			heightMap[i][j] = makeNode(b, j, i, b == 'E')
		}
	}
	stepsToE := bfs(&heightMap, init, isEdge1)
	fmt.Printf("Part 1: %d\n", stepsToE)

	for i, l := range lines {
		sIndex := strings.Index(l, "E")
		if sIndex != -1 {
			init = coord{sIndex, i}
		}
		heightMap[i] = make([]node, len(l))
		for j, b := range []byte(l) {
			heightMap[i][j] = makeNode(b, j, i, b == 'S' || b == 'a')
		}
	}
	stepsToA := bfs(&heightMap, init, isEdge2)
	fmt.Printf("Part 2: %d\n", stepsToA)
}
