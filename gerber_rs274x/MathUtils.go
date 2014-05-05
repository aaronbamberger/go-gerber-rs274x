package gerber_rs274x

import "math"

func epsilonEquals(x float64, y float64, drawPrecision float64) bool {
	epsilon := drawPrecision / math.Pow10(3) // arbitrarily making epsilon 3 orders of magnitude smaller than the drawing precision
	return math.Abs(x - y) < epsilon
}