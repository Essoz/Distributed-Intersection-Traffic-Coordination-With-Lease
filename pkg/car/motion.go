package car

import "math"

// TimeToReach returns the time it takes to reach a destination from a current location
// with a given speed and acceleration
// Estimation Method:
// \int_{0}^{t} (v + at) dt = | destination - currentLocation |
// Input:
//   - currentLocation: the current location of the car in meters
//   - destination: the destination of the car in meters
//   - speed: the current speed of the car in meters per second
//   - acceleration: the current acceleration of the car in meters per second squared
//
// Output:
//   - the time it takes to reach the destination in seconds
func TimeToReach(currentLocation float64, destination float64, speed float64, acceleration float64) float64 {
	distanceToReach := destination - currentLocation

	if distanceToReach == 0 {
		// already at the destination
		return 0
	}

	// not at the destination
	if speed == 0 && acceleration == 0 {
		// no speed, no acceleration
		return math.Inf(1)
	}

	if acceleration == 0 {
		// no acceleration
		return distanceToReach / speed
	}

	// acceleration != 0
	// solve the quadratic equation
	if speed*speed+2*acceleration*distanceToReach < 0 {
		// no solution
		return math.Inf(1)
	}

	// there is a solution
	t1 := (-speed + math.Sqrt(speed*speed+2*acceleration*distanceToReach)) / acceleration
	t2 := (-speed - math.Sqrt(speed*speed+2*acceleration*distanceToReach)) / acceleration
	if t1 > 0 && t2 > 0 {
		// both solutions are positive
		return math.Min(t1, t2)
	}

	return math.Max(t1, t2)
}
