package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/essoz/car-backend/pkg/car"
	"github.com/essoz/car-backend/pkg/lease"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	CONTROL_SERVICE_PORT   = "11001"
	ALLOWED_ERROR_LOCATION = 0.05 // in meters
	RECOMMENDED_SPEED      = 0.3  // in meters per second
	STOP_SPEED             = 0.0  // in meters per second
	CHECK_INTERVAL         = 100 * time.Millisecond
	CAR_DIST_HEAD          = 0.23 // in meters
	CAR_DIST_TAIL          = 0.19 // in meters
	ALLOWED_ERROR_TIME_EXTENDING = 100 // in milliseconds
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

func getCarSelfHeading(ctx context.Context) []float64 {
	url := "http://localhost:" + CONTROL_SERVICE_PORT + "/control/getSelfHeading"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to make a request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read the response body: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		log.Fatalf("Failed to cast data to []interface{}")
	}

	return []float64{data[0].(float64), data[1].(float64)}
}

func getCurrBlock(cli *clientv3.Client, ctx context.Context) lease.Block {
	blocks := lease.GetAllBlocksEtcd(cli, ctx)
	if len(blocks) != 1 {
		log.Fatal("There should be only one block, othercases are not supported yet")
	}
	return blocks[0]
}

func getDistToBlockStart(currCar car.Car, currBlock lease.Block, heading []float64) float64 {
	// find the closest boundary to the car
	if heading[0] == 0 && heading[1] == 1 {
		// heading postive y
		return currBlock.Spec.Location[1] - currCar.Dynamics.Location[1]
	} else if heading[0] == 0 && heading[1] == -1 {
		// heading negative y
		return currCar.Dynamics.Location[1] - currBlock.Spec.Location[1] - currBlock.Spec.Size[1]
	} else if heading[0] == 1 && heading[1] == 0 {
		// heading postive x
		return currBlock.Spec.Location[0] - currCar.Dynamics.Location[0]
	} else if heading[0] == -1 && heading[1] == 0 {
		// heading negative x
		return currCar.Dynamics.Location[0] - currBlock.Spec.Location[0] - currBlock.Spec.Size[0]
	} else {
		panic("Invalid heading at getDistToBlockStart")
	}
}

func getDistToBlockEnd(currCar car.Car, currBlock lease.Block, heading []float64) float64 {
	// find the closest boundary to the car
	if heading[0] == 0 && heading[1] == 1 {
		// heading postive y
		return currBlock.Spec.Location[1] + currBlock.Spec.Size[1] - currCar.Dynamics.Location[1]
	} else if heading[0] == 0 && heading[1] == -1 {
		// heading negative y
		return currCar.Dynamics.Location[1] - currBlock.Spec.Location[1] - currBlock.Spec.Size[1]
	} else if heading[0] == 1 && heading[1] == 0 {
		// heading postive x
		return currBlock.Spec.Location[0] - currCar.Dynamics.Location[0]
	} else if heading[0] == -1 && heading[1] == 0 {
		// heading negative x
		return currCar.Dynamics.Location[0] - currBlock.Spec.Location[0] - currBlock.Spec.Size[0]
	} else {
		panic("Invalid heading at getDistToBlockStart")
	}
}

func IsCarNearIntersection(currCar car.Car, currBlock lease.Block, heading []float64) bool {
	dist := getDistToBlockStart(currCar, currBlock, heading)
	return dist < CAR_DIST_HEAD
}

// func getCarRecommendedSpeed(currCar car.Car, currBlock lease.Block, heading []float64) float64 {
// 	dist := getDistToBlockStart(currCar, currBlock, heading)

// }

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
	currBlock := getCurrBlock(cli, ctx) // FIXME: we assume there is only one block
	currCarHeading := getCarSelfHeading(ctx)

	for {
		currCar = car.GetCarEtcd(cli, ctx, carName)

		for !currCar.IsCarAtDestination(ALLOWED_ERROR_LOCATION) {
			if !currCar.IsDestinationAhead(currCarHeading) {
				log.Println("destination is behind the car, not running lease control service")
				setCarSelfSpeed(ctx, -1*RECOMMENDED_SPEED)
				time.Sleep(CHECK_INTERVAL)
				currCar = car.GetCarEtcd(cli, ctx, carName)
			}

			/* main control logic - a finite state machine
			1. Planning state
			  Description: The car is about to enter the intersection, and it will plan its speed.
			  Condition: The car is not at the destination, and the block is ahead of the car
			  Action:
				- Reserve the block (need to keep monitoring the block's lease status as your lease may be overriden by others)
				- Plan the speed
				- Adjust the speed accordingly

			2. Crossing state
			  Description: The car is waiting for the intersection to be clear
			  Condition: The car is not at the destination, and the block is not ahead of the car
			  Action: Decelerate to the stop speed

			3. Post-crossing state
			  Description: The car has crossed the intersection, and it will accelerate to the recommended speed
			  Condition: The car has crossed the intersection
			  Action: Accelerate to the recommended speed
			*/
			var stage string
			if stage == "post-crossing" {
				log.Println("post-crossing state")

				// CONTROL
				setCarSelfSpeed(ctx, RECOMMENDED_SPEED)

				// LEASING
				// TODO: CLEANUP all leases related to this car

			} else if stage == "crossing" {
				log.Println("crossing state")

				// CONTROL
				// if there is lease, accelerate to the recommended speed
				setCarSelfSpeed(ctx, RECOMMENDED_SPEED)

				// LEASING
				// TODO: check if there's a need to extend the lease
				/*
					time := getCurrTime()
					currBlock = getBlock()
					if time > currBlock.LeaseEndTime + allowedErrorRange {
						// extend the lease
						getBlock
						findTheLease
						modifyLeaseEndTime
					    if overlapping {
							clear all the following leases
						}
						putBlock
					}
				*/

				// get current time in seconds
				currTime := int(time.Now().UnixNano() / int64(time.Millisecond))

				currBlock := getCurrBlock(cli, ctx)
				recentPastLease := currBlock.GetRecentPastLease(currTime)
				if currTime > currBlock.Spec.Leases[currBlock.GetMatchedLeaseIndex(recentPastLease)].EndTime + ALLOWED_ERROR_TIME_EXTENDING {

				}

			} else if stage == "planning" {
				log.Println("planning state")

				currLease := currBlock.GetCarLease(carName)

				// CONTROL FIRST
				if currLease == nil {
					log.Println("no lease found, trying to lease")

					// if not near the intersection, accelerate to the recommended speed
					if IsCarNearIntersection(currCar, currBlock, currCarHeading) {
						setCarSelfSpeed(ctx, STOP_SPEED)
					} else {
						setCarSelfSpeed(ctx, RECOMMENDED_SPEED)
					}
				} else {
					log.Println("lease found, trying to plan")

					// calculate the recommended speed

				}

				// LEASING

				// if there is no lease, try to lease

			} else {
				log.Println("unknown state")
				setCarSelfSpeed(ctx, STOP_SPEED)
			}

			// state transition

			time.Sleep(CHECK_INTERVAL)
			currCar = car.GetCarEtcd(cli, ctx, carName)
		}

		setCarSelfSpeed(ctx, STOP_SPEED)

		log.Println("no new destination found, checking again in 1 second")
		time.Sleep(10 * CHECK_INTERVAL)
	}
}
