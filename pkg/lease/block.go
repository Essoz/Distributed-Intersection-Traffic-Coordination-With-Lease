package lease

import (
	"errors"
)

// get block name
func (b *Block) GetName() string {
	return b.Metadata.Name
}

func (b *Block) GetLeases() []Lease {
	return b.Spec.Leases
}

func (b *Block) CleanPastLeases(currentTime int) {
	for i, lease := range b.GetLeases() {
		if lease.EndTime < currentTime {
			b.Spec.Leases = append(b.GetLeases()[:i], b.GetLeases()[i+1:]...)
		}
	}
}

func (b *Block) GetCurrentLease(currentTime int) Lease {
	for _, lease := range b.GetLeases() {
		if lease.StartTime <= currentTime && lease.EndTime >= currentTime {
			return lease
		}
	}
	// No current lease
	return Lease{}
}

func (b *Block) addNewLease(lease Lease) {
	b.Spec.Leases = append(b.GetLeases(), lease)
}

func (b *Block) ApplyNewLease(lease Lease, currentTime int) error {
	// If there is a current lease, then we need to end it
	currentLease := b.GetCurrentLease(currentTime)
	if currentLease != (Lease{}) {
		// TODO: if the applied lease cannot fit, add constraint solving logic. For now, just return failure
		return errors.New("cannot apply new lease as it overlaps with the current lease")
	}

	// Add the new lease
	b.addNewLease(lease)
	return nil
}
