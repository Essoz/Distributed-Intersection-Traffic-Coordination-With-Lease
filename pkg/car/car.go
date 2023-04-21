package car

import (
	"os"

	"gopkg.in/yaml.v3"
)

func NewCar(path string) *Car {
	dat, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var car Car
	err = yaml.Unmarshal(dat, &car)
	if err != nil {
		panic(err)
	}

	return &car
}
