package aoc

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Coord struct {
	X int
	Y int
}

func (c Coord) ManhattanDistance(other Coord) int {
	return Abs(c.X-other.X) + Abs(c.Y-other.Y)
}

func (c Coord) Neighbors() []Coord {
	neighbors := make([]Coord, 8)
	neighbors[0] = Coord{X: c.X - 1, Y: c.Y - 1}
	neighbors[1] = Coord{X: c.X, Y: c.Y - 1}
	neighbors[2] = Coord{X: c.X + 1, Y: c.Y - 1}
	neighbors[3] = Coord{X: c.X + 1, Y: c.Y}
	neighbors[4] = Coord{X: c.X + 1, Y: c.Y + 1}
	neighbors[5] = Coord{X: c.X, Y: c.Y + 1}
	neighbors[6] = Coord{X: c.X - 1, Y: c.Y + 1}
	neighbors[7] = Coord{X: c.X - 1, Y: c.Y}
	return neighbors
}

func (c Coord) OrthogonalNeighbors() []Coord {
	neighbors := make([]Coord, 4)
	neighbors[0] = Coord{X: c.X, Y: c.Y - 1}
	neighbors[1] = Coord{X: c.X + 1, Y: c.Y}
	neighbors[2] = Coord{X: c.X, Y: c.Y + 1}
	neighbors[3] = Coord{X: c.X - 1, Y: c.Y}
	return neighbors
}

type Range struct {
	Lower int
	Upper int
}

func (r Range) Contains(number int) bool {
	return r.Lower <= number && r.Upper >= number
}

func (r Range) Includes(other Range) bool {
	return other.Lower >= r.Lower && other.Upper <= r.Upper
}

func (r Range) Overlaps(other Range) bool {
	return r.Contains(other.Lower) || r.Contains(other.Upper) || other.Includes(r)
}

func WrapIndex(index int, length int) int {
	if mod := index % length; mod >= 0 {
		return mod
	} else {
		return length + mod
	}
}

func Between[T constraints.Ordered](x T, b1 T, b2 T) bool {
	if b1 <= b2 {
		return b1 <= x && x <= b2
	} else {
		return b2 <= x && x <= b1
	}
}

func Map[T interface{}, S interface{}](s []T, f func(T) S) []S {
	res := make([]S, len(s))
	for i, v := range s {
		res[i] = f(v)
	}
	return res
}

func Filter[T interface{}](s []T, f func(T) bool) []T {
	res := make([]T, 0, len(s))
	for _, v := range s {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}

func Reduce[T interface{}, S interface{}](s []T, init S, f func(S, T) S) S {
	res := init
	for _, v := range s {
		res = f(res, v)
	}
	return res
}

func Any(s []bool) bool {
	for _, v := range s {
		if v {
			return true
		}
	}
	return false
}

func All(s []bool) bool {
	for _, v := range s {
		if !v {
			return false
		}
	}
	return true
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func SetEqual[T comparable](s1 map[T]struct{}, s2 map[T]struct{}) bool {
	if len(s1) != len(s2) {
		return false
	}

	for v := range s2 {
		if _, ok := s1[v]; !ok {
			return false
		}
	}
	return true
}

func SetSubset[T comparable](sub map[T]struct{}, super map[T]struct{}) bool {
	if len(sub) > len(super) {
		return false
	}

	for v := range sub {
		if _, ok := super[v]; !ok {
			return false
		}
	}
	return true
}

func SetIntersect[T comparable](s1 map[T]struct{}, s2 map[T]struct{}) map[T]struct{} {
	intersection := make(map[T]struct{})
	for v := range s2 {
		if _, ok := s1[v]; ok {
			intersection[v] = struct{}{}
		}
	}
	return intersection
}

func SetIn[T comparable](s map[T]struct{}, e T) bool {
	_, ok := s[e]
	return ok
}

func Sum(s []int) int {
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

func Max(s ...int) int {
	max := math.MinInt
	for _, v := range s {
		if v > max {
			max = v
		}
	}
	return max
}

func MaxN(s []int, n int) []int {
	// index i is the top i+1-th element
	maxs := make([]int, n)
	for i := range maxs {
		maxs[i] = math.MinInt
	}
	for _, v := range s {
		comp := v
		swapAlways := false
		for i := range maxs {
			if swapAlways || comp > maxs[i] {
				comp, maxs[i] = maxs[i], comp
				// when we swapped once, we can alwas swap because maxs is sorted --> find correct spot for v and then shift the remaining elements to the right
				swapAlways = true
			}
		}
	}
	return maxs
}

func Abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func Sgn(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	} else {
		return 0
	}
}

func Pow(base int, exp int) int {
	res := 1
	for x := 0; x < exp; x++ {
		res *= base
	}
	return res
}
