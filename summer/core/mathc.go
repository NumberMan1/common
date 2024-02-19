package core

import "math"

func Equal(a, b float64) bool {
	return math.Abs(a-b) < 10e-6
}
