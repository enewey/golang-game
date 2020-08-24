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

// Sub3 subtracts vector a from vector b
func Sub3(ax, ay, az, bx, by, bz float64) (float64, float64, float64) {
	return ax - bx, ay - by, az - bz
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

// Cast - find the normalized difference between two 3d vectors (from -> to)
func Cast(x, y, z, x2, y2, z2 float64) (float64, float64, float64) {
	return Normalize3(x-x2, y-y2, z-z2)
}

// DominantAxis - returns only the highest magnitude vector as a 1 or -1 with the other two as zeroes
func DominantAxis(x, y, z float64) (float64, float64, float64) {
	highest := 0.0
	if math.Abs(x) > math.Abs(y) {
		y = 0
		highest = x
	} else {
		x = 0
		highest = y
	}
	if math.Abs(z) > math.Abs(highest) {
		x = 0
		y = 0
	} else {
		z = 0
	}
	return Normalize(x), Normalize(y), Normalize(z)
}

// Carry returns the value of the float after the decimal point.
func Carry(x, y, z float64) (float64, float64, float64) {
	return x - math.Trunc(x), y - math.Trunc(y), z - math.Trunc(z)

}

// Flint - floors a float and returns it as an int.
func Flint(f float64) int { return int(math.Floor(f)) }

// Flint3 - converts a thruple of floats into ints
func Flint3(a, b, c float64) (int, int, int) { return int(a), int(b), int(c) }

// Itoa3 - conversts a thruple of ints into floats
func Itoa3(a, b, c int) (float64, float64, float64) {
	return float64(a), float64(b), float64(c)
}
