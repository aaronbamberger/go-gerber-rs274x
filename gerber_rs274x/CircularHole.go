package gerber_rs274x

import "fmt"

type CircularHole struct {
	holeDiameter float64
}

func (circle* CircularHole) HolePlaceholder() {

}

func (circle* CircularHole) String() string {
	return fmt.Sprintf("{CH, Diameter: %f}", circle.holeDiameter)
}