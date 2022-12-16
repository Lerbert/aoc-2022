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

func Between[T constraints.Ordered](x T, b1 T, b2 T) bool {
	if b1 <= b2 {
		return b1 <= x && x <= b2
	} else {
		return b2 <= x && x <= b1
	}
}

func Map[T interface{}, S interface{}](s *[]T, f func(T) S) []S {
	res := make([]S, len(*s))
	for i, v := range *s {
		res[i] = f(v)
	}
	return res
}

func Filter[T interface{}](s *[]T, f func(T) bool) []T {
	res := make([]T, 0, len(*s))
	for _, v := range *s {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}

func Reduce[T interface{}, S interface{}](s *[]T, init S, f func(S, T) S) S {
	res := init
	for _, v := range *s {
		res = f(res, v)
	}
	return res
}

func Any(s *[]bool) bool {
	for _, v := range *s {
		if v {
			return true
		}
	}
	return false
}

func All(s *[]bool) bool {
	for _, v := range *s {
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

func SetEqual[T comparable](s1 []T, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	setS1 := make(map[T]struct{})
	for _, v := range s1 {
		setS1[v] = struct{}{}
	}
	for _, v := range s2 {
		if _, ok := setS1[v]; !ok {
			return false
		}
	}
	return true
}

func Sum(s []int) int {
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

func Max(s []int) int {
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
