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

/* Remove all expired lease */
func (b *Block) CleanPastLeases(currentTimeMilli int) {
	for i, lease := range b.GetLeases() {
		if lease.EndTime < currentTimeMilli {
			log.Printf("Deleting lease %#v because of expiration", lease)
			b.Spec.Leases = append(b.GetLeases()[:i], b.GetLeases()[i+1:]...)
		}
	}
}

/* Get the most recent lease previous to currentTimeMilli */
func (b *Block) GetRecentPastLease(currentTimeMilli int) Lease {
	max_time := 0
	var recent_lease Lease
	for _, lease := range b.GetLeases() {
		if lease.EndTime < currentTimeMilli {
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

func (b *Block) GetCurrentLease(currentTimeMilli int) Lease {
	for _, lease := range b.GetLeases() {
		if lease.StartTime <= currentTimeMilli && lease.EndTime >= currentTimeMilli {
			return lease
		}
	}
	// No current lease
	return Lease{}
}

func (b *Block) ExtendUpcomingLease(inputTime int, extendTime int) {
	for _, lease := range b.Spec.Leases {
		if lease.StartTime > inputTime {
			lease.DelayStartTime(lease.StartTime + extendTime)
			lease.ExtendEndTime(lease.EndTime + extendTime)
		}
	}
}

func (b *Block) addNewLease(lease Lease) {
	// check if the lease already exists
	b.Spec.Leases = append(b.GetLeases(), lease)
}

func (b *Block) ApplyNewLease(lease Lease) error {
	// If there is a current lease, then we need to end it
	prevLease := b.GetCurrentLease(lease.StartTime)
	pastLease := b.GetCurrentLease(lease.EndTime)
	if prevLease != (Lease{}) || pastLease != (Lease{}) {
		// TODO: if the applied lease cannot fit, add constraint solving logic. For now, just return failure
		return errors.New("cannot apply new lease as it overlaps with existing leases")
	}

	// Add the new lease
	b.addNewLease(lease)
	return nil
}

func (b *Block) UpdateLease(lease Lease) error {
	// find the lease
	existingLease := b.GetCarLease(lease.CarName)
	if existingLease == nil {
		return errors.New("cannot find existing lease")
	}

	// get rid of the old lease
	b.Spec.Leases = append(b.GetLeases()[:b.GetMatchedLeaseIndex(*existingLease)], b.GetLeases()[b.GetMatchedLeaseIndex(*existingLease)+1:]...)
	return b.ApplyNewLease(lease)
}

func (b *Block) GetCarLease(carName string) *Lease {
	for _, lease := range b.Spec.Leases {
		if lease.GetCarName() == carName {
			return &lease
		}
	}
	return nil
}

/* Get the end time of the last lease */
func (b *Block) GetLastEndTime() int {
	maxTime := 0
	for _, lease := range b.Spec.Leases {
		if lease.EndTime > maxTime {
			maxTime = lease.EndTime
		}
	}
	return maxTime
}

func (b *Block) CleanCarLeases(carName string) {
	for i, lease := range b.Spec.Leases {
		if lease.CarName == carName {
			log.Printf("Deleting lease %#v for car %s", lease, carName)
			b.Spec.Leases = append(b.GetLeases()[:i], b.GetLeases()[i+1:]...)
		}
	}
}
