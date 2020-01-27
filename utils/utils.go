package utils

import (
	"math"
)

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

// Itof converts two ints to two floats.
func Itof(a, b int) (float64, float64) {
	return float64(a), float64(b)
}

// Normalize returns 1, -1, or 0 if the number is
// above, below, or equal to zero respectively.
func Normalize(x float64) float64 {
	if x == 0 {
		return 0
	} else if x < 0 {
		return -1
	}
	return 1
}

// Normalize2 - normalize a 2d vector
func Normalize2(x, y float64) (float64, float64) {
	mag := Magnitude2(x, y)
	return x / mag, y / mag
}

// Normalize3 - normalize a 3d vector
func Normalize3(x, y, z float64) (float64, float64, float64) {
	mag := Magnitude3(x, y, z)
	return x / mag, y / mag, z / mag
}

// Magnitude2 - get the length of a 2d vector
func Magnitude2(x, y float64) float64 {
	return math.Sqrt((x * x) + (y * y))
}

// Magnitude3 - get the length of a 3d vector
func Magnitude3(x, y, z float64) float64 {
	return math.Sqrt((x * x) + (y * y) + (z * z))
}

// Pass returns b if and only if both a and b are the same positive or negative.
// if a is zero, the return value will always be b
func Pass(a int, b float64) float64 {
	if (a < 0 && b > 0) || (a > 0 && b < 0) {
		return 0
	}
	return b
}

// Unpass returns b if and only if a is 0
func Unpass(a int, b float64) float64 {
	if a != 0 {
		return 0
	}
	return b
}

// Flint - floors a float and returns it as an int.
func Flint(f float64) int { return int(math.Floor(f)) }
