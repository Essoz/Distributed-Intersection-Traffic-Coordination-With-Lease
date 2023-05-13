package main

import (
	"context"
	"log"
	"net"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var global_car_name string

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func main() {
	// parse the command line arguments
	// intersectionPath := flag.String("i", "data/intersection.yaml", "path to the intersection yaml file")
	// carPath := flag.String("c", "data/car.yaml", "path to the car yaml file")
	// flag.Parse()

	// Create a client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// wait for etcd connection before proceeding
	for {
		log.Println("Waiting for etcd connection...")
		_, err := cli.Get(context.Background(), "foo")
		if err == nil {
			log.Println("Connected to etcd")
			break
		}
		time.Sleep(1 * time.Second)
	}

	// wait for 11001 port to be open
	for {
		log.Println("Waiting for perceptionService (port 11001) to be running...")
		conn, err := net.Dial("tcp", "localhost:11001")
		if err == nil {
			log.Println("perceptionService is running")
			conn.Close()
			break
		}
		time.Sleep(1 * time.Second)
	}

	// finally, wait for 10 seconds to make sure everything is ready
	time.Sleep(10 * time.Second)

	// let global_car_name be the hostname network address
	global_car_name = GetLocalIP()

	// initializeIntersectionAndBlock(cli, context.Background(), *intersectionPath)
	// initializeCar(cli, context.Background(), *carPath)

	// start the main loop
	Run(cli, context.Background())
	log.Println("Done")
}
