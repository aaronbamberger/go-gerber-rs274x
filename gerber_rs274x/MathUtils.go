package gerber_rs274x

import "math"

type Quadrant int

const (
	QUADRANT_1 Quadrant = iota
	QUADRANT_2
	QUADRANT_3
	QUADRANT_4
)

func inQuadrant(angle float64, quadrant Quadrant) bool {
	switch quadrant {
		case QUADRANT_1:
			return (angle >= 0.0) && (angle <= ONE_HALF_PI)
			
		case QUADRANT_2:
			return (angle >= ONE_HALF_PI) && (angle <= math.Pi)
			
		case QUADRANT_3:
			return (angle >= -math.Pi) && (angle <= -ONE_HALF_PI)
			
		case QUADRANT_4:
			return (angle >= -ONE_HALF_PI) && (angle <= 0.0)
		
		default:
			return false
	}
}

func epsilonEquals(x float64, y float64, drawPrecision float64) bool {
	epsilon := drawPrecision / math.Pow10(3) // arbitrarily making epsilon 3 orders of magnitude smaller than the drawing precision
	return math.Abs(x - y) < epsilon
}

func convertAngleBetweenCoordinateFrames(angle float64) (convertedAngle float64) {
	// Convert an angle calculated in the gerber coordinate frame into the corresponding
	// angle in the cairo coordinate frame, or an angle calculated in the cairo coordinate frame
	// into the gerber coordinate frame
	
	// First, we subtract the given angle from 360 to swap the sign on the y axis
	angle = TWO_PI - angle
	
	// We then normalize the angle to between 0 and 360
	for angle > TWO_PI {
		angle -= TWO_PI
	}
	
	for angle < 0 {
		angle += TWO_PI
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