package lease_test

import (
	"reflect"
	"testing"

	"github.com/essoz/v2v-backend/pkg/lease"
	"github.com/google/go-cmp/cmp"
)

func TestMarshalIntersection(t *testing.T) {
	intersection := lease.NewIntersection("../../data/intersection_test.yaml")
	intersection_expected := lease.Intersection{
		Metadata: lease.IntersectionMetadata{
			Name: "test",
		},
		Spec: lease.IntersectionSpec{
			Position:   []float64{0, 0},
			Size:       []float64{1, 1},
			SplitIndex: 1,
		},
	}
	// compare the intersection with the expected intersection structs
	if !reflect.DeepEqual(intersection, &intersection_expected) {
		t.Errorf("Intersection is not equal to the expected intersection: %s", cmp.Diff(intersection, &intersection_expected))
	}
}
