package lease

func NewLease(carId int, blockId int, startTime int, endTime int) *Lease {
	return &Lease{
		CarId:     carId,
		BlockId:   blockId,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

func (lease *Lease) OverlapsWith(other Lease) bool {
	return !(lease.EndTime < other.StartTime || lease.StartTime > other.EndTime)
}

func (lease *Lease) ExtendEndTime(newEndTime int) {
	lease.EndTime = newEndTime
}

func (lease *Lease) GetDuration() int {
	return lease.EndTime - lease.StartTime
}

func (lease *Lease) GetStartTime() int {
	return lease.StartTime
}

func (lease *Lease) GetEndTime() int {
	return lease.EndTime
}
