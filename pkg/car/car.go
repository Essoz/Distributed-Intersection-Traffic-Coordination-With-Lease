package car

import (
	"math"
	"os"

	"gopkg.in/yaml.v3"
)

func NewCar(path string) Car {
	dat, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var car Car
	err = yaml.Unmarshal(dat, &car)
	if err != nil {
		panic(err)
	}

	return car
}

func (c *Car) GetDynamics() *CarDynamics {
	return &c.Dynamics
}

// set the speed of the car to a given speed
// Input:
//   - speed: the speed to set the car to in centimeters per second
func (c *Car) SetSpeed(speed []float64) {
	c.Dynamics.Speed = speed
}

func (c *Car) SetDynamics(carDynamics *CarDynamics) {
	c.Dynamics = *carDynamics
}

func (c *Car) UpdateStage(stage CarDynamicsStage) *Car {
	c.Dynamics.Stage = stage
	return c
}

func (c *Car) GetName() string {
	return c.Metadata.Name
}

func (c *Car) GetStage() CarDynamicsStage {
	return c.Dynamics.Stage
}

func (c *Car) GetLocation() []float64 {
	return c.Dynamics.Location
}

func (c *Car) GetSpeed() []float64 {
	return c.Dynamics.Speed
}

func (c *Car) GetDestination() []float64 {
	return c.Dynamics.Destination
}

/*
 * Returns true if the car is at the destination, false otherwise
 * Input:
 *   - allowedError: the allowed error in meters
 */
func (c *Car) IsCarAtDestination(allowedError float64) bool {
	currLocation := c.GetLocation()
	destination := c.GetDestination()

	// if there's no destination, then the car is at the destination
	if len(destination) == 0 {
		return true
	}

	// check the distance between the current location and the destination
	distance := math.Sqrt(math.Pow(currLocation[0]-destination[0], 2) + math.Pow(currLocation[1]-destination[1], 2))
	return distance <= allowedError
}
