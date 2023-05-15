package main

import (
	"context"
	"errors"
	"log"
	"math"

	"github.com/essoz/car-backend/pkg/car"
	"github.com/essoz/car-backend/pkg/lease"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	ZERO_ERROR_RANGE = 0.1
)

func speedVectorToScalar(currSpeed []float64) float64 {
	return math.Sqrt(currSpeed[0]*currSpeed[0] + currSpeed[1]*currSpeed[1])
}

func speedToDirection(currSpeed []float64) []int {
	x := 0
	y := 0
	if currSpeed[0] <= ZERO_ERROR_RANGE && currSpeed[0] >= -1*ZERO_ERROR_RANGE {
		x = 0
	} else if currSpeed[0] > ZERO_ERROR_RANGE {
		x = 1
	} else {
		x = -1
	}

	if currSpeed[1] <= ZERO_ERROR_RANGE && currSpeed[1] >= -1*ZERO_ERROR_RANGE {
		y = 0
	} else if currSpeed[1] > ZERO_ERROR_RANGE {
		y = 1
	} else {
		y = -1
	}

	return []int{x, y}
}

func isCarTowardIntersection(surrCar car.Car, currBlock lease.Block) bool {
	direction := speedToDirection(surrCar.Dynamics.Speed)
	if direction[0] == 0 && direction[1] == 0 {
		return false
	}

	if direction[0] == 0 {
		if direction[1] == 1 {
			return surrCar.Dynamics.Location[1] > currBlock.Spec.Location[1]
		} else {
			return surrCar.Dynamics.Location[1] <= currBlock.Spec.Location[1]+currBlock.Spec.Size[1]
		}
	} else {
		if direction[0] == 1 {
			return surrCar.Dynamics.Location[0] > currBlock.Spec.Location[0]
		} else {
			return surrCar.Dynamics.Location[0] <= currBlock.Spec.Location[0]+currBlock.Spec.Size[0]
		}
	}
}

func isCarInsideBlock(surrCar car.Car, currBlock lease.Block) bool {
	return surrCar.Dynamics.Location[0] >= currBlock.Spec.Location[0] &&
		surrCar.Dynamics.Location[0] <= currBlock.Spec.Location[0]+currBlock.Spec.Size[0] &&
		surrCar.Dynamics.Location[1] >= currBlock.Spec.Location[1] &&
		surrCar.Dynamics.Location[1] <= currBlock.Spec.Location[1]+currBlock.Spec.Size[1]
}

func nonV2VTimeToEnterBlock(surrCar car.Car, currBlock lease.Block) (int, error) {
	// IF CAR IS INSIDE THE BLOCK, THEN RETURN 0
	if isCarInsideBlock(surrCar, currBlock) {
		return 0, errors.New("car is already inside the block")
	}

	direction := speedToDirection(surrCar.Dynamics.Speed)
	if direction[0] == 0 && direction[1] == 0 || direction[0]*direction[1] != 0 {
		return 0, errors.New("car's speed is not valid")
	}
	log.Printf("direction: %v", direction)
	nearDist := getDistToBlockStart(surrCar, currBlock, []float64{float64(direction[0]), float64(direction[1])})
	return int(nearDist / speedVectorToScalar(surrCar.Dynamics.Speed)), nil
}

func nonV2VTimeToLeaveBlock(surrCar car.Car, currBlock lease.Block) (int, error) {
	// IF CAR IS OUTSIDE THE BLOCK, THEN RETURN 0
	if !isCarInsideBlock(surrCar, currBlock) {
		return 0, errors.New("car is already inside the block")
	}

	// TODO:
	direction := speedToDirection(surrCar.Dynamics.Speed)
	if direction[0] == 0 && direction[1] == 0 || direction[0]*direction[1] != 0 {
		return 0, errors.New("car's speed is not valid")
	}
	log.Printf("direction: %v", direction)
	farDist := getDistToBlockEnd(surrCar, currBlock, []float64{float64(direction[0]), float64(direction[1])})
	return int(farDist / speedVectorToScalar(surrCar.Dynamics.Speed)), nil
}

/*
 * Description:
 * Leasing service for non-V2V cars
 * Leasing service is responsible for leasing the car to other cars
 * when the car is going to cross the intersection.
 * Args:
 * surroundingCars: a list of surrounding cars
 * carHeading: the qcar's heading
 * Note:
 * This function assumes that the cars from perception are already de-duplicated, i.e.,
 * one physical car should be only percepted once.
 *
 */
func SurroundingCarsLeasing(cli *clientv3.Client, ctx context.Context, surroundingCars []car.Car, currCar car.Car) {
	/*
		Description:
		1. For each surrounding car, determine whether the car is going to cross the intersection
		2. If the car is going to cross the intersection, then lease the car
		3. If the car is not going to cross the intersection, then do nothing

		Criteria for determining whether the car is going to cross the intersection:
			1. The car is going to cross the intersection if the car's speed is greater
				than 0 and the car's location is near a range of the intersection.
			2. The car is not going to cross the intersection if the car's speed is 0.
	*/
	surrCarNames := []string{}

	for _, surrCar := range surroundingCars {
		currBlock := getCurrBlock(cli, ctx)
		currTimeMilli := getCurrTimeMilli()
		surrCarName := "car/" + currCar.Metadata.Name + "/surrounding/" + surrCar.Metadata.Name
		surrCarNames = append(surrCarNames, surrCarName)
		if !isCarTowardIntersection(surrCar, currBlock) {
			// the car is not going to cross the intersection, delete all related leases
			currBlock.CleanCarLeases(surrCarName)
			currBlock.PutEtcd(cli, ctx, "")
			continue
		}
		// the car is going to cross the intersection

		currLease := currBlock.GetCarLease(surrCarName)
		// there is no lease for the car, create a new lease
		timeToEnterBlock, err := nonV2VTimeToEnterBlock(surrCar, currBlock)
		if err != nil {
			// cannot determine the time to enter the block, continue
			continue
		}
		leaseStart := currTimeMilli + timeToEnterBlock - ALLOWED_ERROR_TIME_EXTENDING*5

		timeToLeaveBlock, err := nonV2VTimeToLeaveBlock(surrCar, currBlock)
		if err != nil {
			// cannot determine the time to leave the block, continue
			continue
		}
		leaseEnd := currTimeMilli + timeToLeaveBlock + ALLOWED_ERROR_TIME_EXTENDING*5
		if currLease == nil {
			currLease = lease.NewLease(surrCarName, currBlock.GetName(), leaseStart, leaseEnd)
			currBlock.ApplyNewLeaseNonV2V(*currLease)
		} else {
			// there is already a lease for the car, update the lease
			currLease.StartTime = leaseStart
			currLease.EndTime = leaseEnd
			currBlock.ApplyNewLeaseNonV2V(*currLease)
		}

		currBlock.PutEtcd(cli, ctx, "")
	}

	// delete all leases that are not in the surrounding cars
	currBlock := getCurrBlock(cli, ctx)
	currBlock.CleanNonExistSurrCarLeases(surrCarNames)
	currBlock.PutEtcd(cli, ctx, "")
}
