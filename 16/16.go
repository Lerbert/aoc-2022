package main

import (
	util "aoc2022/aoc"
	inp "aoc2022/input"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type valve struct {
	id       string
	flowRate int
}

var valveRe = regexp.MustCompile(`Valve (..) has flow rate=(\d*); tunnels? leads? to valves? ((?:(?:, )?..)*)`)

func valveFromline(s string) (valve, []string) {
	matches := valveRe.FindStringSubmatch(s)
	flowRate, err := strconv.Atoi(matches[2])
	if err != nil {
		log.Fatal("Could not parse flow rate", err)
	}
	return valve{id: matches[1], flowRate: flowRate}, strings.Split(matches[3], ", ")
}

func apsp(valves map[string]valve, neighbors map[string][]string) map[string]map[string]int {
	// Floyd-Warshall
	dist := make(map[string]map[string]int)
	for v := range valves {
		dist[v] = make(map[string]int)
		for w := range valves {
			if util.Contains(neighbors[v], w) {
				dist[v][w] = 1
			} else {
				dist[v][w] = math.MaxInt
			}
		}
		dist[v][v] = 0
	}

	for u := range valves {
		for v := range valves {
			for w := range valves {
				if newDist := dist[v][u] + dist[u][w]; newDist > 0 && dist[v][w] > newDist {
					dist[v][w] = newDist
				}
			}
		}
	}

	return dist
}

func reduceGraph(valves map[string]valve, neighbors map[string][]string) (map[string]valve, map[string]map[string]int) {
	nonZeroFlowOrStart := make(map[string]valve)
	for _, v := range valves {
		if v.id == "AA" || v.flowRate > 0 {
			nonZeroFlowOrStart[v.id] = v
		}
	}
	dist := apsp(valves, neighbors)
	paths := make(map[string]map[string]int)
	for v := range nonZeroFlowOrStart {
		paths[v] = make(map[string]int)
		for w := range nonZeroFlowOrStart {
			if v == w || w == "AA" {
				continue
			}
			paths[v][w] = dist[v][w] + 1 // + 1 to account for opening the valve at the destination
		}
	}
	return nonZeroFlowOrStart, paths
}

type valveTracker struct {
	v      valve
	minute int
}

type openedScore struct {
	opened []string
	score  int
}

func findOpenedScore(s []openedScore, opened []string) (int, bool) {
	for _, v := range s {
		if util.SetEqual(v.opened, opened) {
			return v.score, true
		}
	}
	return 0, false
}

func findMaxPressureRelease(valves map[string]valve, tunnels map[string]map[string]int, minutes int) int {
	cur := valveTracker{v: valves["AA"], minute: minutes}
	opened := make([]string, 0)
	memo := make(map[valveTracker][]openedScore)
	return findMaxRec(valves, tunnels, cur, opened, memo)
}

func findMaxRec(valves map[string]valve, tunnels map[string]map[string]int, cur valveTracker, opened []string, memo map[valveTracker][]openedScore) int {
	opScore, ok := memo[cur]
	if ok {
		score, found := findOpenedScore(opScore, opened)
		if found {
			return score
		}
	}

	bestScore := 0
	nextOpened := append(append(make([]string, 0, len(opened)+1), opened...), cur.v.id)
	for _, v := range valves {
		if !util.Contains(nextOpened, v.id) {
			reached := cur.minute - tunnels[cur.v.id][v.id]
			if reached < 0 {
				continue
			}
			next := valveTracker{v: v, minute: reached}
			score := findMaxRec(valves, tunnels, next, nextOpened, memo)
			if score > bestScore {
				bestScore = score
			}
		}
	}
	score := bestScore + cur.minute*cur.v.flowRate
	if !ok {
		memo[cur] = make([]openedScore, 0)
	}
	memo[cur] = append(memo[cur], openedScore{opened: opened, score: score})
	return score
}

func main() {
	lines := inp.ReadLines("input")
	valves := make(map[string]valve)
	neighbors := make(map[string][]string)
	for _, l := range lines {
		valve, next := valveFromline(l)
		valves[valve.id] = valve
		neighbors[valve.id] = next
	}
	nodes, edges := reduceGraph(valves, neighbors)
	pressureRelease := findMaxPressureRelease(nodes, edges, 30)
	fmt.Printf("Part 1: %d\n", pressureRelease)
}
