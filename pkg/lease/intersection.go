package lease

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

func NewIntersection(path string) *Intersection {
	// read the yaml file and marshal it into the Intersection struct
	dat, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var intersection Intersection
	err = yaml.Unmarshal(dat, &intersection)
	if err != nil {
		panic(err)
	}

	return &intersection
}

// func (i *Intersection) IsVehicleGoingToPassIntersection(speed float64, location float64, acceleration float64) bool {
// 	var closestIntersectionBoundary float64

// 	if location < i.Spec.Position[0] {
// 		closestIntersectionBoundary = i.Spec.Position[0]
// 	} else if location > i.Spec.Position[1] {
// 		closestIntersectionBoundary = i.Spec.Position[1]
// 	} else {
// 		// already in the intersection TODO: behavior to be specified
// 		panic("Unimplemented")
// 	}

// 	timeToReach := car.TimeToReach(location, speed, closestIntersectionBoundary, acceleration)
// 	// if timeToReach is inf, then the vehicle is not going to pass the intersection
// 	if timeToReach == math.Inf(1) || timeToReach == math.Inf(-1) {
// 		return false
// 	}

// 	return true
// }

func (i *Intersection) IsVehicleInsideIntersection(location []float64) bool {
	locationX := location[0]
	locationY := location[1]

	if locationX >= i.Spec.Position[0] &&
		locationY <= i.Spec.Position[1] &&
		locationX <= i.Spec.Position[0]+i.Spec.Size[0] &&
		locationY <= i.Spec.Position[1]+i.Spec.Size[1] {
		return true
	}

	return false
}
