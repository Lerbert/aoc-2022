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

// ============Preprocessing============

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

// ============Common============

type openedScore struct {
	opened map[string]struct{}
	score  int
}

func findOpenedScore(s []openedScore, opened map[string]struct{}) (int, bool) {
	for _, v := range s {
		if util.SetEqual(v.opened, opened) {
			return v.score, true
		}
	}
	return 0, false
}

// ============Part 1============

type valveTracker struct {
	v      valve
	minute int
}

func findMaxPressureRelease(valves map[string]valve, tunnels map[string]map[string]int, minutes int) int {
	cur := valveTracker{v: valves["AA"], minute: minutes}
	opened := make(map[string]struct{})
	memo := make(map[valveTracker][]openedScore)
	return findMaxRec(valves, tunnels, cur, opened, memo)
}

func findMaxRec(valves map[string]valve, tunnels map[string]map[string]int, cur valveTracker, opened map[string]struct{}, memo map[valveTracker][]openedScore) int {
	opScore, ok := memo[cur]
	if ok {
		score, found := findOpenedScore(opScore, opened)
		if found {
			return score
		}
	}

	bestScore := 0
	nextOpened := make(map[string]struct{})
	for v := range opened {
		nextOpened[v] = struct{}{}
	}
	nextOpened[cur.v.id] = struct{}{}
	for _, v := range valves {
		_, alreadyOpen := nextOpened[v.id]
		if !alreadyOpen {
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

// ============Part 2============

type valveTrackerElephant struct {
	trackers [2]valveTracker
}

func (vte valveTrackerElephant) swapped() valveTrackerElephant {
	return valveTrackerElephant{trackers: [2]valveTracker{vte.trackers[1], vte.trackers[0]}}
}

func findMaxPressureReleaseElephant(valves map[string]valve, tunnels map[string]map[string]int, minutes int) int {
	cur := valveTrackerElephant{trackers: [2]valveTracker{{v: valves["AA"], minute: minutes}, {v: valves["AA"], minute: minutes}}}
	opened := make(map[string]struct{}, 0)
	memo := make(map[valveTrackerElephant][]openedScore)
	return findMaxRecElephant(valves, tunnels, cur, opened, memo)
}

func findMaxRecElephant(valves map[string]valve, tunnels map[string]map[string]int, cur valveTrackerElephant, opened map[string]struct{}, memo map[valveTrackerElephant][]openedScore) int {
	opScore, ok := memo[cur]
	if ok {
		score, found := findOpenedScore(opScore, opened)
		if found {
			return score
		}
	} else {
		opScore, ok := memo[cur.swapped()]
		if ok {
			score, found := findOpenedScore(opScore, opened)
			if found {
				return score
			}
		}
	}

	bestScore := 0
	nextOpened := make(map[string]struct{})
	for v := range opened {
		nextOpened[v] = struct{}{}
	}
	nextOpened[cur.trackers[0].v.id] = struct{}{}
	nextOpened[cur.trackers[1].v.id] = struct{}{}
	if len(nextOpened) < len(valves) {
		skipped := true
		for _, v := range valves {
			_, alreadyOpenV := nextOpened[v.id]
			if alreadyOpenV {
				continue
			}
			reachedV := cur.trackers[0].minute - tunnels[cur.trackers[0].v.id][v.id]
			if reachedV < 0 {
				continue
			}
			skipped = false
			skippedInner := true
			for _, w := range valves {
				if w == v {
					continue
				}
				_, alreadyOpenW := nextOpened[w.id]
				if alreadyOpenW {
					continue
				}
				reachedW := cur.trackers[1].minute - tunnels[cur.trackers[1].v.id][w.id]
				if reachedW < 0 {
					continue
				}
				skippedInner = false
				next := valveTrackerElephant{[2]valveTracker{{v: v, minute: reachedV}, {v: w, minute: reachedW}}}
				score := findMaxRecElephant(valves, tunnels, next, nextOpened, memo)
				if score > bestScore {
					bestScore = score
				}
			}
			if skippedInner {
				next := valveTrackerElephant{[2]valveTracker{{v: v, minute: reachedV}, cur.trackers[1]}}
				score := findMaxRecElephant(valves, tunnels, next, nextOpened, memo)
				if score > bestScore {
					bestScore = score
				}
			}
		}
		if skipped {
			for _, w := range valves {
				_, alreadyOpenW := nextOpened[w.id]
				if alreadyOpenW {
					continue
				}
				reachedW := cur.trackers[1].minute - tunnels[cur.trackers[1].v.id][w.id]
				if reachedW < 0 {
					continue
				}
				next := valveTrackerElephant{[2]valveTracker{cur.trackers[0], {v: w, minute: reachedW}}}
				score := findMaxRecElephant(valves, tunnels, next, nextOpened, memo)
				if score > bestScore {
					bestScore = score
				}
			}
		}
	}
	// We will never be at the same valve except at the start, but there flow rate is 0
	score := bestScore
	for i := range cur.trackers {
		// If we skipped a loop above we must not count the same valve twice
		if _, ok := opened[cur.trackers[i].v.id]; !ok {
			score += cur.trackers[i].minute * cur.trackers[i].v.flowRate
		}
	}

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

	pressureReleaseElephant := findMaxPressureReleaseElephant(nodes, edges, 26)
	fmt.Printf("Part 2: %d\n", pressureReleaseElephant)
}
