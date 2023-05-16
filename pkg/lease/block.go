package lease

import (
	"errors"
	"log"
	"sort"
	"strings"
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
			log.Printf("[CLEAN PAST LEASES] CurrentTime: %d, Deleting lease %#v because of expiration", currentTimeMilli, lease)
			if i+1 < len(b.Spec.Leases) {
				b.Spec.Leases = append(b.GetLeases()[:i], b.GetLeases()[i+1:]...)
			} else {
				b.Spec.Leases = b.GetLeases()[:i]
			}
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

func (b *Block) DelayUpcomingLease(inputTime int) {
	// first sort the leases according to start time
	sort.Slice(b.Spec.Leases, func(i, j int) bool {
		return b.Spec.Leases[i].StartTime < b.Spec.Leases[j].StartTime
	})

	for idx := range b.Spec.Leases {
		if b.Spec.Leases[idx].StartTime > inputTime {
			if idx > 1 && b.Spec.Leases[idx].StartTime < b.Spec.Leases[idx-1].EndTime {
				// if one of the lease at idx is nonV2V, ignore this one
				if strings.Contains(b.Spec.Leases[idx].CarName, "/surrounding/") {
					continue
				}

				// if the lease is conflicting with the previous lease, then we need to delay it
				newStartTime := b.Spec.Leases[idx-1].EndTime + 1
				newEndTime := newStartTime + b.Spec.Leases[idx].GetDuration()
				b.Spec.Leases[idx].DelayStartTime(newStartTime)
				b.Spec.Leases[idx].ExtendEndTime(newEndTime)
			} else {
				// no more conflict, so we can stop
				break
			}
		}
	}
}

func (b *Block) addNewLease(lease Lease) {
	// check if the lease already exists
	log.Printf("Added new lease %#v", lease)
	b.Spec.Leases = append(b.GetLeases(), lease)
}

func (b *Block) ApplyNewLease(lease Lease) error {
	// If there is a current lease, then we need to end it
	prevLease := b.GetCurrentLease(lease.StartTime)
	pastLease := b.GetCurrentLease(lease.EndTime)
	if prevLease != (Lease{}) || pastLease != (Lease{}) {
		// TODO: if the applied lease cannot fit, add constraint solving logic. For now, just return failure
		log.Printf("Cannot apply new lease %#v as it overlaps with the current lease %#v or past lease %#v", lease, prevLease, pastLease)
		return errors.New("cannot apply new lease as it overlaps with existing leases")
	}

	// Add the new lease
	b.addNewLease(lease)
	return nil
}

func (b *Block) ApplyNewLeaseNonV2V(lease Lease) error {
	// if there are conflicting leases, we want to delay all the leases after the new lease
	log.Println("Applying new lease for non-V2V", lease)
	b.addNewLease(lease)
	b.DelayUpcomingLease(lease.EndTime)
	return nil
}

func (b *Block) UpdateLeaseNonV2V(newLease Lease) error {
	// find the lease
	log.Println("Updating lease for non-V2V", newLease)
	for i, lease := range b.Spec.Leases {
		if lease.CarName == newLease.CarName {
			log.Printf("Updating lease %#v", lease)
			b.Spec.Leases = append(b.GetLeases()[:i], b.GetLeases()[i+1:]...)
			return b.ApplyNewLeaseNonV2V(lease)
		}
	}
	return errors.New("cannot find existing lease")
}

func (b *Block) UpdateLease(lease Lease) error {
	log.Printf("Updating lease %#v", lease)
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
			if i+1 < len(b.Spec.Leases) {
				b.Spec.Leases = append(b.GetLeases()[:i], b.GetLeases()[i+1:]...)
			} else {
				b.Spec.Leases = b.GetLeases()[:i]
			}
		}
	}
}

func (b *Block) CleanNonExistSurrCarLeases(surrCarNames []string) {
	toClean := []string{}
	for _, lease := range b.Spec.Leases {
		if strings.Contains(lease.CarName, "/surrounding/") {
			found := false
			for _, surrCarName := range surrCarNames {
				if lease.CarName == surrCarName {
					found = true
					break
				}
			}
			if !found {
				log.Printf("Deleting lease %#v for car %s because it's not being observed anymore", lease, lease.CarName)
				toClean = append(toClean, lease.CarName)
			}
		}
	}

	for _, carName := range toClean {
		b.CleanCarLeases(carName)
	}
}
