package main

import (
	"context"

	"github.com/essoz/car-backend/pkg/car"
	"github.com/essoz/car-backend/pkg/lease"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/appengine/log"
)

// Given a car key, it will run the control service for that car
// Three kinds of actions will happen, depending on the car's state:
//  1. Planning state: The car is about to enter the intersection, and it will plan its path
//     During this state, the control service will compare the current location, current speed, first boundary's location, and the lease starting time. To decide how to accelerate and decelerate.
//  2. Crossing state: The car is waiting for the intersection to be clear
//  3. Post-crossing state:
func RunControlService(cli *clientv3.Client, ctx context.Context, carName string, intersectionName string) {

	intersection := lease.GetIntersectionEtcd(cli, ctx, intersectionName)
	currCar := car.GetCarEtcd(cli, ctx, carName)

	// TODO: wait for MASTER's command to start

	for {
		currCar = car.GetCarEtcd(cli, ctx, carName) // No need to keep fetching from etcd as we can maintain no other entities will write to the car's etcd
		switch currCar.GetStage() {
		case car.CarDynamicsStagePlanning:
			panic("Not implemented")
			// get the leases related to the current car
			leases := lease.GetLeasesRelatedToCarEtcd(cli, ctx, carName)

			// if there are no leases, then the car is not going to enter the intersection
			if len(leases) == 0 {
				continue
			}

			// get the first lease that the car will enter
			FirstLease := leases[0]
			if FirstLease.IsUpcoming() == false {
				log.Errorf("Car %s is not going to enter the intersection", carName)
				continue
			}

			// estimate the time it takes to reach the first lease
			timeToFirstLease := currCar.TimeToEnter(intersection.Spec.Position[0]) // FIXME: Coordinate system

			// calculate the speed the car should be at in order to meet the leasing requirements

			recommenedSpeed := 123
			// make control movements
			currCar.UpdateSpeed(recommenedSpeed).PutEtcd(cli, ctx)

		case car.CarDynamicsStageCrossing:
			continue
			panic("Not implemented")
		case car.CarDynamicsStageCrossed:
			continue
			panic("Not implemented")
		}

		// car stage update
		if currCar.GetStage() == car.CarDynamicsStagePlanning &&
			intersection.IsVehicleInsideIntersection(currCar.GetLocation()) {
			currCar.UpdateStage(car.CarDynamicsStageCrossing).PutEtcd(cli, ctx)
		} else if currCar.GetStage() == car.CarDynamicsStageCrossing &&
			!intersection.IsVehicleInsideIntersection(currCar.GetLocation()) {
			currCar.UpdateStage(car.CarDynamicsStageCrossed).PutEtcd(cli, ctx)
		}
	}
}
