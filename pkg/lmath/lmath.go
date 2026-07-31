package lmath

import (
	"math"
)

// This file holds common functions which are used through-out the package

const (
	epsilon = 0.000000001
)

// Signed matches all built-in signed integer + float types.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

// Ordered matches all built-in ordered numeric types (no complex).
type Ordered interface {
	Signed |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Float matches all built-in float types.
type Float interface {
	~float32 | ~float64
}

func Ceil[T Float](v T) int {
	return int(math.Ceil(float64(v)))
}

func Floor[T Float](v T) int {
	return int(math.Floor(float64(v)))
}

func Round[T Float](v T) int {
	return int(math.Round(float64(v)))
}

func Exp[T Ordered](v T) T {
	return T(math.Exp(float64(v)))
}

func Sqrt[T Ordered](v T) T {
	return T(math.Sqrt(float64(v)))
}

// RoundTo rounds to N decimal places.
func RoundTo[T Float](v T, places int) T {
	pow := math.Pow(10, float64(places))
	return T(math.Round(float64(v)*pow) / pow)
}

func Mod[T Float](a, b T) T {
	return T(math.Mod(float64(a), float64(b)))
}

func Filter[T any](a []T, keep func(i int, b T) bool) []T {
	n := 0
	for i, x := range a {
		if keep(i, x) {
			a[n] = x
			n++
		}
	}
	return a[:n]
}

// Checks if two floats are equal. Doing a comparision using a small epsilon value
func closeEq(a, b, eps float64) bool {
	if a > b {
		return ((a - b) < eps)
	} else {
		return ((b - a) < eps)
	}
}

// Calculate the determinant of a 2x2 matrix.
// Values are givein in Row-Major order
func det2x2(x, y, z, w float64) float64 {
	return x*w - y*z
}

// Calculate the determinant of a 3x3 matrix.
// Values are given in Row-Major order
func det3x3(a1, a2, a3, b1, b2, b3, c1, c2, c3 float64) float64 {
	// a1 a2 a3
	// b1 b2 b3
	// c1 c2 c3
	return (a1*det2x2(b2, b3, c2, c3) -
		b1*det2x2(a2, a3, c2, c3) +
		c1*det2x2(a2, a3, b2, b3))
}

// Convert from a degree to a radian
func Radians(a float64) float64 {
	return a * math.Pi / 180.0
}

// Convert radian to degree
func Degrees(a float64) float64 {
	return a * 180.0 / math.Pi
}

func Min[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Clamp[T Ordered](v, lo, hi T) T {
	if lo > hi {
		lo, hi = hi, lo
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func Abs[T Signed](v T) T {
	if v < 0 {
		return -v
	}
	return v
}

// Linearly interpolates between the start and end values.
//
//	inc is specified between the range 0 -1
//	Lerp(0,2,0) ==> 0
//	Lerp(0,2,0.5) ==> 1
//	Lerp(0,2,1) ==> 2
func Lerp(start, end, inc float64) float64 {
	return (1-inc)*start + inc*end
}
