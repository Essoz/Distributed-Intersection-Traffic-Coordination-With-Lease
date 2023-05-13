package lease

import (
	"errors"
	"log"
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

func (b *Block) GetRecentPastLease(currentTime int) Lease {
	max_time := 0
	var recent_lease Lease
	for _, lease := range b.GetLeases() {
		if lease.EndTime < currentTime {
			if lease.EndTime > max_time {
				max_time = lease.EndTime
				recent_lease = lease
			}
		}
	}
	return recent_lease
}

func (b *Block) GetMatchedLeaseIndex(input_lease Lease) int {
	for i, lease := range b.GetLeases() {
		if lease.CarName == input_lease.CarName && lease.BlockName == input_lease.BlockName {
			return i
		}
	}
	log.Println("No matched lease found")
	return -1
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

func (b *Block) GetCarLease(carName string) *Lease {
	for _, lease := range b.Spec.Leases {
		if lease.GetCarName() == carName {
			return &lease
		}
	}
	return nil
}
