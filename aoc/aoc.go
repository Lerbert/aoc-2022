package aoc

import "math"

func Map[T interface{}, S interface{}](s *[]T, f func(T) S) []S {
	res := make([]S, len(*s))
	for i, v := range *s {
		res[i] = f(v)
	}
	return res
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
