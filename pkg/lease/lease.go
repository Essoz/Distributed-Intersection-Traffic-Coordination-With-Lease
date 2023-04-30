package lease

func NewLease(carName string, blockName string, startTime int, endTime int) *Lease {
	return &Lease{
		CarName:   carName,
		BlockName: blockName,
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

func (lease *Lease) GetCarName() string {
	return lease.CarName
}

func (lease *Lease) GetBlockName() string {
	return lease.BlockName
}
