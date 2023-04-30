package main

import (
	"context"

	"github.com/essoz/car-backend/pkg/car"
	"github.com/essoz/car-backend/pkg/lease"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Given a car key, it will run the control service for that car
// Three kinds of actions will happen, depending on the car's state:
//  1. Planning state: The car is about to enter the intersection, and it will plan its path
//     During this state, the control service will compare the current location, current speed, first boundary's location, and the lease starting time. To decide how to accelerate and decelerate.
//  2. Crossing state: The car is waiting for the intersection to be clear
//  3. Post-crossing state:
func RunControlService(cli *clientv3.Client, ctx context.Context, carName string, intersectionName string) {
	currCarKey := CAR_PREFIX + carName
	intersectionKey := INTERSECTION_PREFIX + intersectionName

	intersection := lease.GetIntersectionEtcd(cli, ctx, intersectionKey)
	currCar := car.GetCarEtcd(cli, ctx, currCarKey)
	for {
		currCar = car.GetCarEtcd(cli, ctx, currCarKey) // No need to keep fetching from etcd as we can maintain no other entities will write to the car's etcd
		switch currCar.GetStage() {
		case car.CarDynamicsStagePlanning:
			panic("Not implemented")
			// get the leases related to the current car
			leases := lease.GetLeasesRelatedToCarEtcd(cli, ctx, carName)
			speed := currCar.GetSpeed()
			location := currCar.GetLocation()

		case car.CarDynamicsStageCrossing:
			panic("Not implemented")
		case car.CarDynamicsStageCrossed:
			panic("Not implemented")
		}

		// car stage update
		if currCar.GetStage() == car.CarDynamicsStagePlanning &&
			intersection.IsVehicleInsideIntersection(currCar.GetLocation()) {
			currCar.UpdateStage(car.CarDynamicsStageCrossing).PutEtcd(cli, ctx, CAR_PREFIX)
		} else if currCar.GetStage() == car.CarDynamicsStageCrossing &&
			!intersection.IsVehicleInsideIntersection(currCar.GetLocation()) {
			currCar.UpdateStage(car.CarDynamicsStageCrossed).PutEtcd(cli, ctx, CAR_PREFIX)
		}
	}
}
