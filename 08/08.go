package main

import (
	inp "aoc2022/input"
	"fmt"
)

func getVisibleTreesRow(treeRow []byte) ([]bool, []int) {
	visibility := make([]bool, len(treeRow))
	scenicScore := make([]int, len(treeRow))

	// forward
	highest := -1
	blockedAt := make([]int, 10)
	for i, height := range treeRow {
		intHeight := int(height - '0')

		visibility[i] = intHeight > highest
		if intHeight > highest {
			highest = intHeight
		}

		if i == 0 {
			scenicScore[i] = 0
			// no update, because it does not matter if we are blocked at edge or not
		} else {
			scenicScore[i] = i - blockedAt[intHeight]
			for j := 0; j <= intHeight; j++ {
				blockedAt[j] = i
			}
		}
	}

	// backward
	highest = -1
	blockedAt = make([]int, 10)
	for i := len(treeRow) - 1; i >= 0; i-- {
		intHeight := int(treeRow[i] - '0')

		visibility[i] = visibility[i] || intHeight > highest
		if intHeight > highest {
			highest = intHeight
		}

		if i == len(treeRow)-1 {
			scenicScore[i] = 0
			// no update, because it does not matter if we are blocked at edge or not
		} else {
			scenicScore[i] *= len(treeRow) - 1 - i - blockedAt[intHeight]
			for j := 0; j <= intHeight; j++ {
				blockedAt[j] = len(treeRow) - 1 - i
			}
		}
	}

	return visibility, scenicScore
}

func getVisibleTrees(forestMap [][]byte) ([][]bool, [][]int) {
	visibility := make([][]bool, len(forestMap))
	scenicScore := make([][]int, len(forestMap))
	for i, row := range forestMap {
		visibility[i], scenicScore[i] = getVisibleTreesRow(row)
	}

	transpose(&forestMap)
	transpose(&visibility)
	transpose(&scenicScore)
	for i, row := range forestMap {
		v, s := getVisibleTreesRow(row)
		for j := range row {
			visibility[i][j] = visibility[i][j] || v[j]
			scenicScore[i][j] *= s[j]
		}
	}

	return visibility, scenicScore
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

func max(score [][]int) int {
	highest := -1
	for _, row := range score {
		for _, s := range row {
			if s > highest {
				highest = s
			}
		}
	}
	return highest
}

func transpose[T interface{}](m *[][]T) {
	for row := 0; row < len(*m); row++ {
		for col := row; col < len((*m)[row]); col++ {
			(*m)[row][col], (*m)[col][row] = (*m)[col][row], (*m)[row][col]
		}
	}
}

func main() {
	lines := inp.ReadLines("input")
	forestMap := make([][]byte, len(lines))
	for i, l := range lines {
		forestMap[i] = []byte(l)
	}

	visibility, scenicScore := getVisibleTrees(forestMap)

	visibleTrees := countTrue(visibility)
	fmt.Printf("Part 1: %d\n", visibleTrees)

	bestScore := max(scenicScore)
	fmt.Printf("Part 2: %d\n", bestScore)
}
