package lease

import (
	"errors"
)

func NewBlock(id int) Block {
	return Block{
		Id:     id,
		Leases: []Lease{},
	}
}

func (block *Block) CleanPastLeases(currentTime int) {
	for i, lease := range block.Leases {
		if lease.EndTime < currentTime {
			block.Leases = append(block.Leases[:i], block.Leases[i+1:]...)
		}
	}
}

func (block *Block) GetCurrentLease(currentTime int) Lease {
	for _, lease := range block.Leases {
		if lease.StartTime <= currentTime && lease.EndTime >= currentTime {
			return lease
		}
	}
	// No current lease
	return Lease{}
}

func (block *Block) addNewLease(lease Lease) {
	block.Leases = append(block.Leases, lease)
}

func (block *Block) ApplyNewLease(lease Lease, currentTime int) error {
	// If there is a current lease, then we need to end it
	currentLease := block.GetCurrentLease(currentTime)
	if currentLease != (Lease{}) {
		// TODO: if the applied lease cannot fit, add constraint solving logic. For now, just return failure
		return errors.New("cannot apply new lease as it overlaps with the current lease")
	}

	// Add the new lease
	block.addNewLease(lease)
	return nil
}
