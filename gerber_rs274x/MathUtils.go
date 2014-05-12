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

func snapCoordinate(coord float64) float64 {
	if (math.Ceil(coord) - coord) > 0.5 {
		return math.Floor(coord) + 0.5
	} else {
		return math.Ceil(coord) + 0.5
	}
}

func lawOfCosines(aX float64, aY float64, bX float64, bY float64, cX float64, cY float64) (angle float64) {
	// Use the law of cosines to compute an interior angle of a triangle, given all 3 points
	sideA := math.Hypot(bX - cX, bY - cY)
	sideB := math.Hypot(aX - cX, aY - cY)
	sideC := math.Hypot(aX - bX, aY - bY)
	
	return math.Acos((math.Pow(sideA, 2) + math.Pow(sideB, 2) - math.Pow(sideC, 2)) / (2 * sideA * sideB))
}