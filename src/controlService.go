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
	CONTROL_SERVICE_PORT         = "11001"
	ALLOWED_ERROR_LOCATION       = 0.05 // in meters
	RECOMMENDED_FORWARD_SPEED    = 0.3  // in meters per second
	RECOMMENDED_BACKWARD_SPEED   = -0.1 // in meters per second
	STOP_SPEED                   = 0.0  // in meters per second
	CHECK_INTERVAL               = 100 * time.Millisecond
	CAR_DIST_HEAD                = 0.23 // in meters
	CAR_DIST_TAIL                = 0.19 // in meters
	ALLOWED_ERROR_TIME_EXTENDING = 100  // in milliseconds
	LEASE_EXTEND_DURATION        = 500  // in milliseconds
)

func getCurrTimeMilli() int {
	return int(time.Now().UnixMilli())
}

func setCarSelfSpeed(ctx context.Context, speed float64) {
	url := "http://localhost:" + CONTROL_SERVICE_PORT + "/control/setSelfSpeed"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(`{"speed": %f}`, speed))))
	if err != nil {
		log.Fatalf("Failed to make a request: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Setting car speed to %f", speed)
	// no need to check data as long as the request is 200
}

func getCarSelfHeading(ctx context.Context) []float64 {
	url := "http://localhost:" + CONTROL_SERVICE_PORT + "/perception/getSelfHeading"
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
		log.Printf("Body is %s", body)
		panic("Failed to parse JSON" + err.Error())
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
		return currCar.Dynamics.Location[1] - currBlock.Spec.Location[1]
	} else if heading[0] == 1 && heading[1] == 0 {
		// heading postive x
		return currBlock.Spec.Location[0] + currBlock.Spec.Size[0] - currCar.Dynamics.Location[0]
	} else if heading[0] == -1 && heading[1] == 0 {
		// heading negative x
		return currCar.Dynamics.Location[0] - currBlock.Spec.Location[0]
	} else {
		panic("Invalid heading at getDistToBlockStart")
	}
}

func getCarStage(currCar car.Car, currBlock lease.Block, heading []float64) string {
	var closestBoundary float64
	var farthestBoundary float64

	currLoc := currCar.GetLocation()
	if heading[0] == 0 && heading[1] == 1 {
		// heading postive y
		closestBoundary = currBlock.Spec.Location[1]
		farthestBoundary = currBlock.Spec.Location[1] + currBlock.Spec.Size[1]
		if currLoc[1] < closestBoundary {
			return "planning"
		} else if currLoc[1] < farthestBoundary {
			return "crossing"
		} else {
			return "post-crossing"
		}
	} else if heading[0] == 0 && heading[1] == -1 {
		// heading negative y
		closestBoundary = currBlock.Spec.Location[1] + currBlock.Spec.Size[1]
		farthestBoundary = currBlock.Spec.Location[1]
		if currLoc[1] > closestBoundary {
			return "planning"
		} else if currLoc[1] > farthestBoundary {
			return "crossing"
		} else {
			return "post-crossing"
		}
	} else if heading[0] == 1 && heading[1] == 0 {
		// heading postive x
		closestBoundary = currBlock.Spec.Location[0]
		farthestBoundary = currBlock.Spec.Location[0] + currBlock.Spec.Size[0]
		if currLoc[0] < closestBoundary {
			return "planning"
		} else if currLoc[0] < farthestBoundary {
			return "crossing"
		} else {
			return "post-crossing"
		}
	} else if heading[0] == -1 && heading[1] == 0 {
		// heading negative x
		closestBoundary = currBlock.Spec.Location[0] + currBlock.Spec.Size[0]
		farthestBoundary = currBlock.Spec.Location[0]
		if currLoc[0] > closestBoundary {
			return "planning"
		} else if currLoc[0] > farthestBoundary {
			return "crossing"
		} else {
			return "post-crossing"
		}
	}

	log.Printf("Heading is %v", heading)
	panic("Invalid heading at getCarStage")
}

func IsCarNearIntersection(currCar car.Car, currBlock lease.Block, heading []float64) bool {
	dist := getDistToBlockStart(currCar, currBlock, heading)
	return dist < CAR_DIST_HEAD
}

func getCarRecommendedSpeed(currCar car.Car, currBlock lease.Block, heading []float64) float64 {
	distMin := getDistToBlockStart(currCar, currBlock, heading)
	distMax := getDistToBlockEnd(currCar, currBlock, heading)

	currTimeMilli := getCurrTimeMilli()
	currLease := currBlock.GetCarLease(currCar.Metadata.Name)
	if currLease == nil {
		// recommend a random speed
		return RECOMMENDED_FORWARD_SPEED
	}

	if currLease.EndTime < currTimeMilli {
		log.Fatal("The lease should not be expired")
	}

	// speed's lower bound is distMax / (leaseEndTime - currTime)
	// speed's upper bound is distMin / (leaseEndTime - currTime)
	// speed's recommended value is (distMin + distMax) / 2 / (leaseEndTime - currTime)

	currTimeSeconds := float64(currTimeMilli) / 1000
	leaseEndTimeSeconds := float64(currLease.EndTime) / 1000
	leaseStartTimeSeconds := float64(currLease.StartTime) / 1000
	speedLowerBound := distMax / (leaseEndTimeSeconds - currTimeSeconds)
	speedUpperBound := distMin / (leaseStartTimeSeconds - currTimeSeconds)

	if speedUpperBound < 0 {
		// keep the curr speed
		return getCurrSpeedScalar(currCar, heading)
	}

	if speedLowerBound > speedUpperBound {
		log.Printf("Current time is %v", currTimeSeconds)
		log.Printf("Lease's start time is %v", leaseStartTimeSeconds)
		log.Printf("Lease's end time is %v", leaseEndTimeSeconds)
		log.Printf("Dist's min is %v", distMin)
		log.Printf("Dist's max is %v", distMax)
		log.Printf("Speed's lower bound is %v", speedLowerBound)
		log.Printf("Speed's upper bound is %v", speedUpperBound)
		log.Fatal("Speed's lower bound should not be greater than speed's upper bound")
		// TODO: the car cannot make it to the intersection, what should we do?
	}

	if RECOMMENDED_FORWARD_SPEED >= speedLowerBound && RECOMMENDED_FORWARD_SPEED <= speedUpperBound {
		return RECOMMENDED_FORWARD_SPEED
	}

	if RECOMMENDED_FORWARD_SPEED < (speedLowerBound+speedUpperBound)/2 {
		return RECOMMENDED_FORWARD_SPEED
	} else {
		return (speedLowerBound + speedUpperBound) / 2
	}
}

func getCurrSpeedScalar(currCar car.Car, heading []float64) float64 {
	return currCar.Dynamics.Speed[0]*heading[0] + currCar.Dynamics.Speed[1]*heading[1]
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

	time.Sleep(2 * time.Second) // FIXME: remove this hack
	currCar := car.GetCarEtcd(cli, ctx, carName)
	currBlock := getCurrBlock(cli, ctx) // FIXME: we assume there is only one block
	currCarHeading := getCarSelfHeading(ctx)

	log.Printf("Current heading is %v\n", currCarHeading)

	for {
		currCar = car.GetCarEtcd(cli, ctx, carName)
		for !currCar.IsCarAtDestination(ALLOWED_ERROR_LOCATION) {
			stage := getCarStage(currCar, currBlock, currCarHeading)
			log.Printf("CurrentStage: %s\n", stage)

			currTimeMilli := getCurrTimeMilli()
			currBlock.CleanPastLeases(currTimeMilli)

			if !currCar.IsDestinationAhead(currCarHeading) {
				log.Println("destination is behind the car, not running lease control service")
				setCarSelfSpeed(ctx, -1*RECOMMENDED_FORWARD_SPEED)
				time.Sleep(CHECK_INTERVAL)
				currCar = car.GetCarEtcd(cli, ctx, carName)
				continue
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

			if stage == "post-crossing" {
				log.Println("post-crossing state")

				// CONTROL
				setCarSelfSpeed(ctx, RECOMMENDED_FORWARD_SPEED)

				// LEASING
				currBlock.CleanCarLeases(carName)

			} else if stage == "crossing" {
				log.Println("crossing state")

				// CONTROL
				// if there is lease, accelerate to the recommended speed
				setCarSelfSpeed(ctx, RECOMMENDED_FORWARD_SPEED)

				// LEASING
				// get current time in milliseconds
				currBlock := getCurrBlock(cli, ctx)
				recentPastLease := currBlock.GetCarLease(carName)
				if currTimeMilli > recentPastLease.EndTime+ALLOWED_ERROR_TIME_EXTENDING {
					// extend the lease of the car itself (end only)
					currentEndTime := currBlock.GetCarLease(carName).EndTime
					currBlock.GetCarLease(carName).EndTime = currTimeMilli + LEASE_EXTEND_DURATION
					// extend the upcoming leases (start and end)
					currBlock.ExtendUpcomingLease(currentEndTime, LEASE_EXTEND_DURATION)
				}
				currBlock.PutEtcd(cli, ctx, "")

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
						setCarSelfSpeed(ctx, RECOMMENDED_FORWARD_SPEED)
					}
				} else {
					log.Println("lease found, trying to plan")
					// calculate the recommended speed
					recommendedSpeed := getCarRecommendedSpeed(currCar, currBlock, currCarHeading)

					setCarSelfSpeed(ctx, recommendedSpeed)
				}

				// LEASING
				if currLease == nil {
					// if there's no lease, try to make leases
					// given the current speed, predict the time when the car will reach the intersection
					// predStartTime, predEndTime
					// make the lease at the first available time
					firstAvailableTimeMilli := currBlock.GetLastEndTime() + 1 // FIXME: There are actually two available times, the other one is the current time till the start time of the first lease
					reachTimeMilliUnderRecSpeed := int(getDistToBlockStart(currCar, currBlock, currCarHeading)/RECOMMENDED_FORWARD_SPEED*1000) - ALLOWED_ERROR_TIME_EXTENDING*2 + currTimeMilli
					leaveTimeMilliUnderRecSpeed := int(getDistToBlockEnd(currCar, currBlock, currCarHeading)/RECOMMENDED_FORWARD_SPEED*1000) + ALLOWED_ERROR_TIME_EXTENDING*2 + currTimeMilli
					if firstAvailableTimeMilli < currTimeMilli || firstAvailableTimeMilli < reachTimeMilliUnderRecSpeed {
						// reserve the lease according to the recommended speed
						newLease := lease.NewLease(carName, currBlock.GetName(), reachTimeMilliUnderRecSpeed, leaveTimeMilliUnderRecSpeed)
						err := currBlock.ApplyNewLease(*newLease)
						if err != nil {
							log.Println("failed to apply new lease")
						} else {
							log.Println("successfully applied new lease")
							currBlock.PutEtcd(cli, ctx, "")
						}
					} else {
						leaseDurationUnderRecSpeed := leaveTimeMilliUnderRecSpeed - reachTimeMilliUnderRecSpeed
						newLease := lease.NewLease(carName, currBlock.GetName(), firstAvailableTimeMilli, firstAvailableTimeMilli+leaseDurationUnderRecSpeed)
						err := currBlock.ApplyNewLease(*newLease)
						if err != nil {
							log.Println("failed to apply new lease")
						} else {
							log.Println("successfully applied new lease")
							currBlock.PutEtcd(cli, ctx, "")
						}
					}
				} else {
					// there is a lease
					// check if the lease is still valid
					if currLease.EndTime < currTimeMilli {
						// if not, cancel the lease
						log.Println("lease expired, cancelling lease")
						currBlock.CleanPastLeases(currTimeMilli)
						currBlock.PutEtcd(cli, ctx, "")
					} else {
						// if yes, check if we need to bring the lease forward
						if getCurrSpeedScalar(currCar, currCarHeading) > RECOMMENDED_FORWARD_SPEED {
							// if yes, attempt to bring the lease forward
							prevLease := currBlock.GetRecentPastLease(currLease.StartTime)
							if prevLease != (lease.Lease{}) {
								if prevLease.EndTime < currLease.StartTime-ALLOWED_ERROR_TIME_EXTENDING {
									// move the lease forward
									currLease.EndTime = currLease.StartTime - prevLease.EndTime - 1
									currLease.StartTime = prevLease.EndTime + 1
									err := currBlock.UpdateLease(*currLease)
									if err != nil {
										log.Println("failed to update lease")
									} else {
										log.Println("successfully updated lease")
										currBlock.PutEtcd(cli, ctx, "")
									}
								}
							}
						}
					}
				}

			} else {
				log.Println("unknown state")
				setCarSelfSpeed(ctx, STOP_SPEED)
			}

			time.Sleep(CHECK_INTERVAL)
			currCar = car.GetCarEtcd(cli, ctx, carName)

		}

		setCarSelfSpeed(ctx, STOP_SPEED)

		log.Println("no new destination found, checking again in 1 second")
		time.Sleep(10 * CHECK_INTERVAL)
	}
}
