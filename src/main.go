package main

import (
	"context"
	"flag"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var global_car_name string

func main() {
	// parse the command line arguments
	intersectionPath := flag.String("i", "data/intersection.yaml", "path to the intersection yaml file")
	carPath := flag.String("c", "data/car.yaml", "path to the car yaml file")

	flag.Parse()

	// Create a client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	initializeIntersectionAndBlock(cli, context.Background(), *intersectionPath)
	initializeCar(cli, context.Background(), *carPath)

	// start the main loop
	run(cli, context.Background())
	log.Println("Done")
}
