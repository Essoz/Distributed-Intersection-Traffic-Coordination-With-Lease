package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/essoz/v2v-backend/pkg/car"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func GetAllCars(cli *clientv3.Client, ctx context.Context) ([]byte, error) {
	cars := car.GetAllCarsWithKeyEtcd(cli, ctx, "")

	// // for test purposes only
	// car_1 := car.Car{
	// 	Metadata: car.CarMetadata{
	// 		Name: "car_1",
	// 	},
	// 	Dynamics: car.CarDynamics{
	// 		Location: []float64{10.0, 20.0},
	// 	},
	// }
	// car_2 := car.Car{
	// 	Metadata: car.CarMetadata{
	// 		Name: "car_2",
	// 	},
	// 	Dynamics: car.CarDynamics{
	// 		Location: []float64{30.0, 40.0},
	// 	},
	// }
	// cars := []car.Car{car_1, car_2}

	return json.Marshal(cars)
}

func run(cli *clientv3.Client, ctx context.Context) {
	// start the http server
	http.HandleFunc("/car/getAll", func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// getting request for all cars
		cars, err := GetAllCars(cli, ctx)
		if err != nil {
			log.Fatalf("Failed to get all cars: %v", err)
		}
		// cars is a json document
		w.Write(cars)
	})
	http.HandleFunc("/utils/clearAllCarDestination", func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		// getting request for all cars
		print("Getting request for clearing all cars' destination\n")

		ClearAllCarDestinationEtcd(cli, ctx)
	})
	http.HandleFunc("/car/setLocation", func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		// If it's an OPTIONS request (preflight), return immediately
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// get car name and location from request
		var requestData struct {
			CarName  string    `json:"carName"`
			Location []float64 `json:"location"`
		}

		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			log.Fatalf("Failed to decode request JSON: %v", err)
		}

		log.Printf("carName: %s, destination: %v\n", requestData.CarName, requestData.Location)
	})
	http.HandleFunc("/car/setDestination", func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		// If it's an OPTIONS request (preflight), return immediately
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		var requestData struct {
			CarName     string    `json:"carName"`
			Destination []float64 `json:"destination"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			log.Fatalf("Failed to decode request JSON: %v", err)
		}

		log.Printf("carName: %s, destination: %v\n", requestData.CarName, requestData.Destination)

		// set car destination
		SetCarDestination(cli, ctx, requestData.CarName, requestData.Destination)
	})

	print("Starting server on port 11002")
	log.Fatal(http.ListenAndServe(":11002", nil))
}

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx := context.Background()
	run(cli, ctx)
}
