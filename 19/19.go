package main

import (
	util "aoc2022/aoc"
	inp "aoc2022/input"
	"fmt"
	"log"
	"regexp"
)

type blueprint struct {
	id                int
	oreCostOre        int
	oreCostClay       int
	oreCostObsidian   int
	clayCostObsidian  int
	oreCostGeode      int
	obsidianCostGeode int
}

var blueprintRegex = regexp.MustCompile(`Blueprint (\d*): Each ore robot costs (\d*) ore. Each clay robot costs (\d*) ore. Each obsidian robot costs (\d*) ore and (\d*) clay. Each geode robot costs (\d*) ore and (\d*) obsidian.`)

func blueprintFromLine(l string) blueprint {
	matches := blueprintRegex.FindStringSubmatch(l)
	numbers := util.Map(matches[1:], inp.MustAtoi)
	return blueprint{
		id:                numbers[0],
		oreCostOre:        numbers[1],
		oreCostClay:       numbers[2],
		oreCostObsidian:   numbers[3],
		clayCostObsidian:  numbers[4],
		oreCostGeode:      numbers[5],
		obsidianCostGeode: numbers[6],
	}
}

const (
	ORE      = 1 << 0
	CLAY     = 1 << 2
	OBSIDIAN = 1 << 3
	GEODE    = 1 << 4
	NONE     = 1 << 5
)

type state struct {
	oreRobots      int
	clayRobots     int
	obsidianRobots int
	geodeRobots    int

	ore      int
	clay     int
	obsidian int
	geodes   int

	couldBuildPrev int

	remainingMinutes int
}

func initState(minutes int) state {
	return state{
		oreRobots:      1,
		clayRobots:     0,
		obsidianRobots: 0,
		geodeRobots:    0,

		ore:      0,
		clay:     0,
		obsidian: 0,
		geodes:   0,

		couldBuildPrev: 0,

		remainingMinutes: minutes,
	}
}

func (s state) track() state {
	// Don't track blocked robots as they are not relevant to the best possible outcome from this state
	s.couldBuildPrev = 0
	return s
}

func (s state) canBuildNow(bp blueprint) int {
	canBuild := 0
	if s.ore >= bp.oreCostOre {
		canBuild |= ORE
	}
	if s.ore >= bp.oreCostClay {
		canBuild |= CLAY
	}
	if s.ore >= bp.oreCostObsidian && s.clay >= bp.clayCostObsidian {
		canBuild |= OBSIDIAN
	}
	if s.ore >= bp.oreCostGeode && s.obsidian >= bp.obsidianCostGeode {
		canBuild |= GEODE
	}
	return canBuild
}

func (s state) nextMinute(build int, bp blueprint) state {
	// Block all robots we could build now until we build another robot
	// There is no point in delaying to build a robot and doing nothing instead
	if build == NONE {
		s.couldBuildPrev = s.canBuildNow(bp)
	} else {
		s.couldBuildPrev = 0
	}

	s.ore += s.oreRobots
	s.clay += s.clayRobots
	s.obsidian += s.obsidianRobots
	s.geodes += s.geodeRobots

	switch build {
	case ORE:
		s.oreRobots += 1
		s.ore -= bp.oreCostOre
	case CLAY:
		s.clayRobots += 1
		s.ore -= bp.oreCostClay
	case OBSIDIAN:
		s.obsidianRobots += 1
		s.ore -= bp.oreCostObsidian
		s.clay -= bp.clayCostObsidian
	case GEODE:
		s.geodeRobots += 1
		s.ore -= bp.oreCostGeode
		s.obsidian -= bp.obsidianCostGeode
	case NONE:
	default:
		log.Fatal("Unknown build ", build)
	}

	s.remainingMinutes -= 1

	return s
}

func (s state) geodesUpperBound() int {
	// Assume we can build a geode robot in every remaining minute
	return s.geodes + s.geodeRobots*s.remainingMinutes + ((s.remainingMinutes-1)*s.remainingMinutes)/2
}

func findMaxGeodes(bp blueprint, minutes int) int {
	currentState := initState(minutes)
	memo := make(map[state]int)
	return findMaxGeodesRec(bp, currentState, 0, memo)
}

func findMaxGeodesRec(bp blueprint, currentState state, currentBest int, memo map[state]int) int {
	if currentState.remainingMinutes <= 0 {
		return currentState.geodes
	}
	if best, ok := memo[currentState.track()]; ok {
		return best
	}
	if maxGeodes := currentState.geodesUpperBound(); maxGeodes < currentBest {
		return maxGeodes
	}

	canBuild := currentState.canBuildNow(bp)
	best := 0
	if canBuild&GEODE == GEODE && currentState.couldBuildPrev&GEODE == 0 {
		next := findMaxGeodesRec(bp, currentState.nextMinute(GEODE, bp), currentBest, memo)
		best = util.Max(best, next)
		currentBest = util.Max(currentBest, best)
	}
	if canBuild&OBSIDIAN == OBSIDIAN && currentState.couldBuildPrev&OBSIDIAN == 0 && currentState.obsidianRobots < bp.obsidianCostGeode {
		next := findMaxGeodesRec(bp, currentState.nextMinute(OBSIDIAN, bp), currentBest, memo)
		best = util.Max(best, next)
		currentBest = util.Max(currentBest, best)
	}
	if canBuild&CLAY == CLAY && currentState.couldBuildPrev&CLAY == 0 && currentState.clayRobots < bp.clayCostObsidian {
		next := findMaxGeodesRec(bp, currentState.nextMinute(CLAY, bp), currentBest, memo)
		best = util.Max(best, next)
		currentBest = util.Max(currentBest, best)
	}
	if canBuild&ORE == ORE && currentState.couldBuildPrev&ORE == 0 && currentState.oreRobots < util.Max(bp.oreCostOre, bp.oreCostClay, bp.oreCostObsidian, bp.oreCostGeode) {
		next := findMaxGeodesRec(bp, currentState.nextMinute(ORE, bp), currentBest, memo)
		best = util.Max(best, next)
		currentBest = util.Max(currentBest, best)
	}
	next := findMaxGeodesRec(bp, currentState.nextMinute(NONE, bp), currentBest, memo)
	best = util.Max(best, next)

	memo[currentState.track()] = best

	return best
}

type job struct {
	bp      blueprint
	minutes int
}

type result struct {
	j      job
	geodes int
}

func worker(jobChan chan job, resultChan chan result) {
	for job := range jobChan {
		geodes := findMaxGeodes(job.bp, job.minutes)
		resultChan <- result{j: job, geodes: geodes}
	}
}

func sendJobs(jobChan chan job, blueprints []blueprint, minutes int) {
	for _, bp := range blueprints {
		jobChan <- job{bp: bp, minutes: minutes}
	}
	close(jobChan)
}

func findMaxGeodesPar(blueprints []blueprint, minutes int, workers int) []result {
	jobChan := make(chan job)
	resultChan := make(chan result, workers)

	for i := 0; i < workers; i++ {
		go worker(jobChan, resultChan)
	}

	go sendJobs(jobChan, blueprints, minutes)

	results := make([]result, len(blueprints))
	for i := range results {
		results[i] = <-resultChan
	}
	close(resultChan)

	return results
}

const NUM_WORKERS = 1

func main() {
	lines := inp.ReadLines("input")
	blueprints := util.Map(lines, blueprintFromLine)

	maxGeodes := findMaxGeodesPar(blueprints, 24, NUM_WORKERS)
	qualityLevelSum := 0
	for _, res := range maxGeodes {
		qualityLevelSum += res.geodes * res.j.bp.id
	}
	fmt.Printf("Part 1: %d\n", qualityLevelSum)

	// Elephants munch the blueprints
	if len(blueprints) > 3 {
		blueprints = blueprints[:3]
	}

	maxGeodes2 := findMaxGeodesPar(blueprints, 32, NUM_WORKERS)
	geodeProduct := 1
	for _, res := range maxGeodes2 {
		geodeProduct *= res.geodes
	}
	fmt.Printf("Part 2: %d\n", geodeProduct)
}
