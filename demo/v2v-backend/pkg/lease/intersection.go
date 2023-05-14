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
