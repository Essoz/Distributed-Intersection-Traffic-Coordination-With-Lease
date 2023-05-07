// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package car

import "fmt"
import "encoding/json"
import "reflect"

type Car struct {
	// These are dynamic properties that are updated by the simulator
	Dynamics CarDynamics `json:"dynamics" yaml:"dynamics"`

	// Metadata corresponds to the JSON schema field "metadata".
	Metadata CarMetadata `json:"metadata" yaml:"metadata"`
}

// These are dynamic properties that are updated by the simulator
type CarDynamics struct {
	// Acceleration corresponds to the JSON schema field "acceleration".
	Acceleration float64 `json:"acceleration" yaml:"acceleration"`

	// Destination corresponds to the JSON schema field "destination".
	Destination []float64 `json:"destination" yaml:"destination"`

	// Heading corresponds to the JSON schema field "heading".
	Heading float64 `json:"heading" yaml:"heading"`

	// Location corresponds to the JSON schema field "location".
	Location []float64 `json:"location" yaml:"location"`

	// PassingBlocks corresponds to the JSON schema field "passingBlocks".
	PassingBlocks []CarDynamicsPassingBlocksElem `json:"passingBlocks" yaml:"passingBlocks"`

	// Speed corresponds to the JSON schema field "speed".
	Speed []float64 `json:"speed" yaml:"speed"`

	// Stage corresponds to the JSON schema field "stage".
	Stage CarDynamicsStage `json:"stage" yaml:"stage"`
}

type CarDynamicsPassingBlocksElem struct {
	// Name corresponds to the JSON schema field "name".
	Name string `json:"name" yaml:"name"`

	// Position corresponds to the JSON schema field "position".
	Position []float64 `json:"position" yaml:"position"`

	// Size corresponds to the JSON schema field "size".
	Size []float64 `json:"size" yaml:"size"`
}

type CarDynamicsStage string

const CarDynamicsStageCrossed CarDynamicsStage = "crossed"
const CarDynamicsStageCrossing CarDynamicsStage = "crossing"
const CarDynamicsStagePlanning CarDynamicsStage = "planning"

type CarMetadata struct {
	// IsV2V corresponds to the JSON schema field "isV2v".
	IsV2V bool `json:"isV2v" yaml:"isV2v"`

	// Name corresponds to the JSON schema field "name".
	Name string `json:"name" yaml:"name"`

	// Type corresponds to the JSON schema field "type".
	Type string `json:"type" yaml:"type"`
}

var enumValues_CarDynamicsStage = []interface{}{
	"planning",
	"crossing",
	"crossed",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CarDynamics) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["acceleration"]; !ok || v == nil {
		return fmt.Errorf("field acceleration in CarDynamics: required")
	}
	if v, ok := raw["destination"]; !ok || v == nil {
		return fmt.Errorf("field destination in CarDynamics: required")
	}
	if v, ok := raw["heading"]; !ok || v == nil {
		return fmt.Errorf("field heading in CarDynamics: required")
	}
	if v, ok := raw["location"]; !ok || v == nil {
		return fmt.Errorf("field location in CarDynamics: required")
	}
	if v, ok := raw["passingBlocks"]; !ok || v == nil {
		return fmt.Errorf("field passingBlocks in CarDynamics: required")
	}
	if v, ok := raw["speed"]; !ok || v == nil {
		return fmt.Errorf("field speed in CarDynamics: required")
	}
	if v, ok := raw["stage"]; !ok || v == nil {
		return fmt.Errorf("field stage in CarDynamics: required")
	}
	type Plain CarDynamics
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if len(plain.Location) < 2 {
		return fmt.Errorf("field %s length: must be >= %d", "location", 2)
	}
	if len(plain.Location) > 2 {
		return fmt.Errorf("field %s length: must be <= %d", "location", 2)
	}
	if len(plain.Speed) < 2 {
		return fmt.Errorf("field %s length: must be >= %d", "speed", 2)
	}
	if len(plain.Speed) > 2 {
		return fmt.Errorf("field %s length: must be <= %d", "speed", 2)
	}
	*j = CarDynamics(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CarDynamicsStage) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_CarDynamicsStage {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_CarDynamicsStage, v)
	}
	*j = CarDynamicsStage(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CarMetadata) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["isV2v"]; !ok || v == nil {
		return fmt.Errorf("field isV2v in CarMetadata: required")
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name in CarMetadata: required")
	}
	if v, ok := raw["type"]; !ok || v == nil {
		return fmt.Errorf("field type in CarMetadata: required")
	}
	type Plain CarMetadata
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CarMetadata(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CarDynamicsPassingBlocksElem) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name in CarDynamicsPassingBlocksElem: required")
	}
	if v, ok := raw["position"]; !ok || v == nil {
		return fmt.Errorf("field position in CarDynamicsPassingBlocksElem: required")
	}
	if v, ok := raw["size"]; !ok || v == nil {
		return fmt.Errorf("field size in CarDynamicsPassingBlocksElem: required")
	}
	type Plain CarDynamicsPassingBlocksElem
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = CarDynamicsPassingBlocksElem(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Car) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["dynamics"]; !ok || v == nil {
		return fmt.Errorf("field dynamics in Car: required")
	}
	if v, ok := raw["metadata"]; !ok || v == nil {
		return fmt.Errorf("field metadata in Car: required")
	}
	type Plain Car
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Car(plain)
	return nil
}
