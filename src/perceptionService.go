package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/essoz/car-backend/pkg/car"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	PERCEPTION_SERVICE_PORT = "11001"
)

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
					carMap["location"].([]float64)[0],
					carMap["location"].([]float64)[1],
				},
				Speed: []float64{
					carMap["speed"].([]float64)[0],
					carMap["speed"].([]float64)[1],
				},
			},
		}
		cars = append(cars, car)
	}

	return cars
}

func RunPerceptionService(cli *clientv3.Client, ctx context.Context, carName string) {
	// THIS FUNCTION WILL INTERFACE THE PERCEPTION SERVICE

	for true {
		currCar := car.GetCarEtcd(cli, ctx, carName)
		// Get all surrounding objects
		currSurrCars := getCarSurrounding(ctx)
		currSurrCarPtrs := make([]*car.Car, len(currSurrCars))
		for i := range currSurrCars {
			currSurrCarPtrs[i] = &currSurrCars[i]
		}
		currCar.UpdateSurroundingCarsEtcd(cli, ctx, currSurrCarPtrs)
	}
}
