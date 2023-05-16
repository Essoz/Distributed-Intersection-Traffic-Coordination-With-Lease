package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/essoz/car-backend/pkg/car"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	PERCEPTION_SERVICE_PORT = "11001"
	LOCATION_ERROR          = 0.3
)

func getCarSelfSpeed(ctx context.Context) []float64 {
	url := "http://localhost:" + PERCEPTION_SERVICE_PORT + "/perception/getSelfSpeed"
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

func getCarSelfLocation(ctx context.Context) []float64 {
	url := "http://localhost:" + PERCEPTION_SERVICE_PORT + "/perception/getSelfLocation"
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

func getCarSurrounding(ctx context.Context) []car.Car {
	// THIS FUNCTION WILL GET THE SURROUNDING OBJECTS OF THE CAR
	// the server is running at localhost:11001
	url := "http://localhost:" + PERCEPTION_SERVICE_PORT + "/perception/getSurrounding"
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
	// printout body as json string
	// log.Printf("Body: %s", string(body))
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		log.Fatalf("Failed to cast data to map[string]interface{}")
	}

	// convert data to a list of cars, the key will be the car's name
	var cars []car.Car
	for key, value := range data {
		carMap, ok := value.(map[string]interface{})
		if !ok {
			log.Fatalf("Failed to cast carMap to map[string]interface{}")
		}

		car := car.Car{
			Metadata: car.CarMetadata{
				Name: key,
			},
			Dynamics: car.CarDynamics{
				Location: []float64{
					carMap["location"].([]interface{})[0].(float64),
					carMap["location"].([]interface{})[1].(float64),
				},
				Speed: []float64{
					carMap["speed"].([]interface{})[0].(float64),
					carMap["speed"].([]interface{})[1].(float64),
				},
			},
		}
		cars = append(cars, car)
	}

	return cars
}

func IsSameCar(car1 car.Car, car2 car.Car, allowedError float64) bool {
	if len(car1.GetLocation()) != 2 || len(car2.GetLocation()) != 2 {
		return false
	}
	loc_x_diff := math.Abs(car1.GetLocation()[0] - car2.GetLocation()[0])
	loc_y_diff := math.Abs(car1.GetLocation()[1] - car2.GetLocation()[1])

	if loc_x_diff <= allowedError && loc_y_diff <= allowedError {
		return true
	}
	return false
}

func RunPerceptionService(cli *clientv3.Client, ctx context.Context, carName string) {
	// THIS FUNCTION WILL INTERFACE THE PERCEPTION SERVICE
	log.Printf("Running perception service for car %s", carName)
	// register this car
	currCar := car.Car{
		Metadata: car.CarMetadata{
			Name: carName,
		},
	}
	currCar.PutEtcd(cli, ctx, "")

	for {
		// log.Printf("Updating car %s at %s", carName, time.Now().Format("2006-01-02 15:04:05"))

		// update the car's location and speed
		// read car's past state from etcd (need transaction)
		txn := cli.Txn(ctx)
		currCar := car.GetCarEtcd(cli, ctx, carName)
		currCar.Dynamics.Location = getCarSelfLocation(ctx)
		currCar.Dynamics.Speed = getCarSelfSpeed(ctx)
		currCar.PutEtcd(cli, ctx, "")
		txn.Commit()

		// Get all surrounding objects
		currSurrCars := getCarSurrounding(ctx)
		allCars := car.GetAllCarsEtcd(cli, ctx, "")
		toRemove := make([]int, 0)
		for idx, surrCar := range currSurrCars {
			for _, existingCar := range allCars {
				if IsSameCar(surrCar, *existingCar, LOCATION_ERROR) {
					// remove the surrCar
					if surrCar.GetName() == existingCar.GetName() {
						continue
					}

					log.Printf("Removing Car %s because it is duplicated with %s", surrCar.GetName(), existingCar.GetName())
					toRemove = append(toRemove, idx)
				}
			}
		}

		// remove the surrCar
		for i := len(toRemove) - 1; i >= 0; i-- {
			if len(currSurrCars) == 1 {
				currSurrCars = make([]car.Car, 0)
				break
			}
			if toRemove[i] == len(currSurrCars)-1 {
				currSurrCars = currSurrCars[:toRemove[i]]
				break
			}
			currSurrCars = append(currSurrCars[:toRemove[i]], currSurrCars[toRemove[i]+1:]...)
		}

		currSurrCarPtrs := make([]*car.Car, len(currSurrCars))
		for i := range currSurrCars {
			// log.Printf("Surrounding car %s", currSurrCars[i].Metadata.Name)
			currSurrCarPtrs[i] = &currSurrCars[i]
		}
		SurroundingCarsLeasing(cli, ctx, currSurrCars, currCar)
		currCar.UpdateSurroundingCarsEtcd(cli, ctx, currSurrCarPtrs)

		// sleep for 0.05 second
		time.Sleep(50 * time.Millisecond)

	}
}
