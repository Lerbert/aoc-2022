package aoc

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Coord struct {
	X int
	Y int
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

func Contains[T comparable](s *[]T, e T) bool {
	for _, v := range *s {
		if v == e {
			return true
		}
	}
	return false
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
