package main

import (
	inp "aoc2022/input"
	"fmt"
)

func getVisibleTreesRow(treeRow []byte) []bool {
	visibility := make([]bool, len(treeRow))

	// forward
	highest := -1
	for i, height := range treeRow {
		intHeight := int(height - '0')
		visibility[i] = intHeight > highest
		if intHeight > highest {
			highest = intHeight
		}
	}

	// backward
	highest = -1
	for i := len(treeRow) - 1; i >= 0; i-- {
		intHeight := int(treeRow[i] - '0')
		visibility[i] = visibility[i] || intHeight > highest
		if intHeight > highest {
			highest = intHeight
		}
	}

	return visibility
}

func getVisibleTrees(forestMap [][]byte) [][]bool {
	visibility := make([][]bool, len(forestMap))
	for i, row := range forestMap {
		visibility[i] = getVisibleTreesRow(row)
	}

	transpose(&visibility)
	transpose(&forestMap)
	for i, row := range forestMap {
		for j, visible := range getVisibleTreesRow(row) {
			visibility[i][j] = visibility[i][j] || visible
		}
	}

	return visibility
}

func countTrue(visibility [][]bool) int {
	visibleCnt := 0
	for _, row := range visibility {
		for _, visible := range row {
			if visible {
				visibleCnt++
			}
		}
	}
	return visibleCnt
}

func transpose[T interface{}](m *[][]T) {
	for row := 0; row < len(*m); row++ {
		for col := row; col < len((*m)[row]); col++ {
			(*m)[row][col], (*m)[col][row] = (*m)[col][row], (*m)[row][col]
		}
	}
}

func main() {
	lines := inp.ReadLines("test")
	forestMap := make([][]byte, len(lines))
	for i, l := range lines {
		forestMap[i] = []byte(l)
	}

	visibleTrees := countTrue(getVisibleTrees(forestMap))
	fmt.Printf("Part 1: %d\n", visibleTrees)
}
