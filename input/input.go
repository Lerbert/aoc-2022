package input

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func ReadLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Could not open file", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func MustAtoi(s string) int {
	x, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("Could not parse int ", err)
	}
	return x
}
