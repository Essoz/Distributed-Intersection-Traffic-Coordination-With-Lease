package lease

import (
	"encoding/json"
	"fmt"
)

type BlockMetadata struct {
	// Name corresponds to the JSON schema field "name".
	Name string `json:"name" yaml:"name"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *BlockMetadata) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name in BlockMetadata: required")
	}
	type Plain BlockMetadata
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = BlockMetadata(plain)
	return nil
}

type Lease struct {
	// BlockName corresponds to the JSON schema field "blockName".
	BlockName string `json:"blockName" yaml:"blockName"`

	// CarName corresponds to the JSON schema field "carName".
	CarName string `json:"carName" yaml:"carName"`

	// EndTime corresponds to the JSON schema field "endTime".
	EndTime int `json:"endTime" yaml:"endTime"`

	// StartTime corresponds to the JSON schema field "startTime".
	StartTime int `json:"startTime" yaml:"startTime"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Lease) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["blockName"]; !ok || v == nil {
		return fmt.Errorf("field blockName in Lease: required")
	}
	if v, ok := raw["carName"]; !ok || v == nil {
		return fmt.Errorf("field carName in Lease: required")
	}
	if v, ok := raw["endTime"]; !ok || v == nil {
		return fmt.Errorf("field endTime in Lease: required")
	}
	if v, ok := raw["startTime"]; !ok || v == nil {
		return fmt.Errorf("field startTime in Lease: required")
	}
	type Plain Lease
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Lease(plain)
	return nil
}

type BlockSpec struct {
	// Leases corresponds to the JSON schema field "leases".
	Leases []Lease `json:"leases" yaml:"leases"`

	// Location corresponds to the JSON schema field "location".
	Location []float64 `json:"location" yaml:"location"`

	// Size corresponds to the JSON schema field "size".
	Size []float64 `json:"size" yaml:"size"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *BlockSpec) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["leases"]; !ok || v == nil {
		return fmt.Errorf("field leases in BlockSpec: required")
	}
	if v, ok := raw["location"]; !ok || v == nil {
		return fmt.Errorf("field location in BlockSpec: required")
	}
	if v, ok := raw["size"]; !ok || v == nil {
		return fmt.Errorf("field size in BlockSpec: required")
	}
	type Plain BlockSpec
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
	if len(plain.Size) < 2 {
		return fmt.Errorf("field %s length: must be >= %d", "size", 2)
	}
	if len(plain.Size) > 2 {
		return fmt.Errorf("field %s length: must be <= %d", "size", 2)
	}
	*j = BlockSpec(plain)
	return nil
}

type Block struct {
	// Metadata corresponds to the JSON schema field "metadata".
	Metadata BlockMetadata `json:"metadata" yaml:"metadata"`

	// Spec corresponds to the JSON schema field "spec".
	Spec BlockSpec `json:"spec" yaml:"spec"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Block) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["metadata"]; !ok || v == nil {
		return fmt.Errorf("field metadata in Block: required")
	}
	if v, ok := raw["spec"]; !ok || v == nil {
		return fmt.Errorf("field spec in Block: required")
	}
	type Plain Block
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Block(plain)
	return nil
}

type IntersectionMetadata struct {
	// name of the intersection, this name can be used to retrieve the intersection
	// from the database
	Name string `json:"name" yaml:"name"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *IntersectionMetadata) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name in IntersectionMetadata: required")
	}
	type Plain IntersectionMetadata
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = IntersectionMetadata(plain)
	return nil
}

type IntersectionSpec struct {
	// position of the intersection in centimeters
	Position []float64 `json:"position" yaml:"position"`

	// size of the intersection in centimeters
	Size []float64 `json:"size" yaml:"size"`

	// the intersection will be equally split into 2^splitIndex blocks
	SplitIndex float64 `json:"splitIndex" yaml:"splitIndex"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *IntersectionSpec) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["position"]; !ok || v == nil {
		return fmt.Errorf("field position in IntersectionSpec: required")
	}
	if v, ok := raw["size"]; !ok || v == nil {
		return fmt.Errorf("field size in IntersectionSpec: required")
	}
	if v, ok := raw["splitIndex"]; !ok || v == nil {
		return fmt.Errorf("field splitIndex in IntersectionSpec: required")
	}
	type Plain IntersectionSpec
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if len(plain.Position) < 2 {
		return fmt.Errorf("field %s length: must be >= %d", "position", 2)
	}
	if len(plain.Position) > 2 {
		return fmt.Errorf("field %s length: must be <= %d", "position", 2)
	}
	if len(plain.Size) < 2 {
		return fmt.Errorf("field %s length: must be >= %d", "size", 2)
	}
	if len(plain.Size) > 2 {
		return fmt.Errorf("field %s length: must be <= %d", "size", 2)
	}
	*j = IntersectionSpec(plain)
	return nil
}

type Intersection struct {
	// Metadata corresponds to the JSON schema field "metadata".
	Metadata IntersectionMetadata `json:"metadata" yaml:"metadata"`

	// Spec corresponds to the JSON schema field "spec".
	Spec IntersectionSpec `json:"spec" yaml:"spec"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Intersection) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["metadata"]; !ok || v == nil {
		return fmt.Errorf("field metadata in Intersection: required")
	}
	if v, ok := raw["spec"]; !ok || v == nil {
		return fmt.Errorf("field spec in Intersection: required")
	}
	type Plain Intersection
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Intersection(plain)
	return nil
}
