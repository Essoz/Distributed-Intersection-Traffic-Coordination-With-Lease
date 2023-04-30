package car

import (
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
func (c *Car) SetSpeed(speed float64) {
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

func (c *Car) GetSpeed() float64 {
	return c.Dynamics.Speed
}
