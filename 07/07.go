package main

import (
	inp "aoc2022/input"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func isCommand(l string) bool {
	return strings.HasPrefix(l, "$ ")
}

func dirStackToName(stack []string) string {
	name := "/"
	for _, s := range stack {
		name += s + "/"
	}
	return name
}

func parseCd(s string) (string, error) {
	if !strings.HasPrefix(s, "cd ") {
		return "", errors.New("not a cd command")
	}
	return s[3:], nil
}

func cd(dirStack *[]string, target string) {
	if target == ".." {
		*dirStack = (*dirStack)[:len(*dirStack)-1]
	} else if target == "/" {
		*dirStack = make([]string, 0)
	} else {
		*dirStack = append(*dirStack, target)
	}
}

func fileSizeFromLsOutput(s string) int {
	if strings.HasPrefix(s, "dir ") {
		return 0
	}
	split := strings.Split(s, " ")
	size, err := strconv.Atoi(split[0])
	if err != nil {
		log.Fatal("Could not convert size to int", err)
	}
	return size
}

func addSizeToAll(dirStack []string, dirSizes *map[string]int, size int) {
	if size == 0 {
		return
	}
	for i := 0; i <= len(dirStack); i++ {
		currentDir := dirStackToName(dirStack[:i])
		_, ok := (*dirSizes)[currentDir]
		if !ok {
			(*dirSizes)[currentDir] = 0
		}
		(*dirSizes)[currentDir] += size
	}
}

func main() {
	lines := inp.ReadLines("input")
	dirStack := make([]string, 0)
	dirSizes := make(map[string]int)

	for _, l := range lines {
		if isCommand(l) {
			target, err := parseCd(l[2:])
			if err == nil {
				cd(&dirStack, target)
			}
		} else {
			addSizeToAll(dirStack, &dirSizes, fileSizeFromLsOutput(l))
		}
	}

	sizeSum := 0
	for _, size := range dirSizes {
		if size <= 100000 {
			sizeSum += size
		}
	}
	fmt.Printf("Part 1: %d\n", sizeSum)
}
