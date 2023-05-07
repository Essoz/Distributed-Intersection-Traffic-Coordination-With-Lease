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
	cars := car.GetAllCarsWithKeyEtcd(cli, ctx, car.CAR_ETCD_PREFIX)
	return json.Marshal(cars)
}

func run(cli *clientv3.Client, ctx context.Context) {
	// start the http server
	http.HandleFunc("/car/getAll", func(w http.ResponseWriter, r *http.Request) {
		// getting request for all cars
		print("Getting request for all cars")

		cars, err := GetAllCars(cli, ctx)
		if err != nil {
			log.Fatalf("Failed to get all cars: %v", err)
		}
		// cars is a json document
		w.Write(cars)
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
