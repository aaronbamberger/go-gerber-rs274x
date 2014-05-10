package gerber_rs274x

import "math"

func epsilonEquals(x float64, y float64, drawPrecision float64) bool {
	epsilon := drawPrecision / math.Pow10(3) // arbitrarily making epsilon 3 orders of magnitude smaller than the drawing precision
	return math.Abs(x - y) < epsilon
}

func convertAngleBetweenCoordinateFrames(angle float64) (convertedAngle float64) {
	// Convert an angle calculated in the gerber coordinate frame into the corresponding
	// angle in the cairo coordinate frame, or an angle calculated in the cairo coordinate frame
	// into the gerber coordinate frame
	
	// First, we subtract the given angle from 360 to swap the sign on the y axis
	angle = (2.0 * math.Pi) - angle
	
	// We then normalize the angle to between 0 and 360
	for angle > (2.0 * math.Pi) {
		angle -= (2.0 * math.Pi)
	}
	
	for angle < 0 {
		angle += (2.0 * math.Pi)
	}
	
	return angle
}