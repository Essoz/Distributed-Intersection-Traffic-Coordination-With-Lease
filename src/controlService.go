package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/essoz/car-backend/pkg/car"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	CONTROL_SERVICE_PORT   = "11001"
	ALLOWED_ERROR_LOCATION = 0.1 // in meters
	RECOMMENDED_SPEED      = 0.3 // in meters per second
	STOP_SPEED             = 0.0 // in meters per second
	CHECK_INTERVAL         = 100 * time.Millisecond
)

func setCarSelfSpeed(ctx context.Context, speed float64) {
	url := "http://localhost:" + CONTROL_SERVICE_PORT + "/control/setSelfSpeed"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(`{"speed": %f}`, speed))))
	if err != nil {
		log.Fatalf("Failed to make a request: %v", err)
	}
	defer resp.Body.Close()

	// no need to check data as long as the request is 200
}

// Given a car key, it will run the control service for that car
// Three kinds of actions will happen, depending on the car's state:
//  1. Planning state: The car is about to enter the intersection, and it will plan its path
//     During this state, the control service will compare the current location, current speed, first boundary's location, and the lease starting time. To decide how to accelerate and decelerate.
//  2. Crossing state: The car is waiting for the intersection to be clear
//  3. Post-crossing state:
func RunControlService(cli *clientv3.Client, ctx context.Context, carName string, intersectionName string) {
	// intersection := lease.GetIntersectionEtcd(cli, ctx, intersectionName)
	log.Print("Control service is running for car ", carName)

	currCar := car.GetCarEtcd(cli, ctx, carName)
	// TODO: wait for MASTER's command to start
	for {
		currCar = car.GetCarEtcd(cli, ctx, carName)
		// TODO: Check if the car is at its current destination, if so, loop back to the beginning

		for !currCar.IsCarAtDestination(ALLOWED_ERROR_LOCATION) {
			log.Println("Car is not at destination, running")
			// NOTE: this is a simplified version of the control service, which only checks the current location and speed
			// The destination does not determine the orientation of the car. In other words, we assume the destination to be always at the front of the car
			setCarSelfSpeed(ctx, RECOMMENDED_SPEED)

			time.Sleep(CHECK_INTERVAL)
			currCar = car.GetCarEtcd(cli, ctx, carName) // No need to keep fetching from etcd as we can maintain no other entities will write to the car's etcd
		}

		setCarSelfSpeed(ctx, STOP_SPEED)

		log.Println("no new destination found, checking again in 1 second")
		time.Sleep(10 * CHECK_INTERVAL)
	}

}
