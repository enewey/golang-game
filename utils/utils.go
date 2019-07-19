package utils

import "math"

// Max returns the max int from an arbitrary number of ints
func Max(args ...int) int {
	max := args[0]
	for _, v := range args {
		if v > max {
			max = v
		}
	}
	return max
}

// Min returns the min int from an arbitrary number of ints
func Min(args ...int) int {
	min := args[0]
	for _, v := range args {
		if v < min {
			min = v
		}
	}
	return min
}

// Abs is like math.Abs but for ints
func Abs(v int) int {
	if v < 0 {
		return v * -1
	}
	return v
}

// Flint - floors a float and returns it as an int.
func Flint(f float64) int { return int(math.Floor(f)) }
