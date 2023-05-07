package main

import (
	"context"
	"log"

	"github.com/essoz/car-backend/pkg/car"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func initializeCar(cli *clientv3.Client, ctx context.Context, pathToCar string) {
	txn := cli.Txn(ctx)
	defer txn.Commit()

	car := car.NewCar(pathToCar)

	car.PutEtcd(cli, ctx, "")

	global_car_name = car.Metadata.Name

	// commit transaction
	_, err := txn.Commit()
	if err != nil {
		panic(err)
	}
}

// // store the intersection struct and block structs into etcd
// func initializeIntersectionAndBlock(cli *clientv3.Client, ctx context.Context, pathToIntersection string) {
// 	// start transaction
// 	txn := cli.Txn(ctx)
// 	intersection := lease.NewIntersection(pathToIntersection)
// 	intersection_bytes, err := yaml.Marshal(intersection)
// 	if err != nil {
// 		panic(err)
// 	}

// 	intersectionKey := INTERSECTION_PREFIX + intersection.Metadata.Name
// 	txn.If(clientv3.Compare(clientv3.Version(intersectionKey), "=", 0)).Then(
// 		clientv3.OpPut(intersectionKey, string(intersection_bytes)),
// 	)

// 	for _, block := range intersection.Spec.Blocks {
// 		block_bytes, err := yaml.Marshal(block)
// 		if err != nil {
// 			panic(err)
// 		}

// 		blockKey := BLOCK_PREFIX + block.Name
// 		txn.If(clientv3.Compare(clientv3.Version(blockKey), "=", 0)).Then(
// 			clientv3.OpPut(blockKey, string(block_bytes)),
// 		)
// 	}

// 	// commit transaction
// 	_, err = txn.Commit()
// 	if err != nil {
// 		panic(err)
// 	}
// }

func run(cli *clientv3.Client, ctx context.Context) {
	var carSelf car.Car

	// put a key
	_, err := cli.Put(context.Background(), "abc", "bar")
	if err != nil {
		log.Fatal(err)
	}

	// read the key
	resp, err := cli.Get(context.Background(), "abc")
	if err != nil {
		log.Fatal(err)
	}

	for _, ev := range resp.Kvs {
		log.Printf("%s : %s", ev.Key, ev.Value)
	}

	for {

		// TODO: get the current time, position, speed, acceleration, direction of the car itself
		// TODO: get_info(car_id, time, position, speed, acceleration, direction)
		carSelf = car.GetCarEtcd(cli, ctx, global_car_name)
		if carSelf.GetDynamics() == nil {
			log.Println("car spec is nil")
			continue
		}

		// TODO: get the current time, position, speed, acceleration, direction of the surrounding cars
		// TODO: get_info_surroundings(car_id, time, position, speed, acceleration, direction)

		// REGION: lease management for non-V2V cars
		// if there are no surrounding cars trying to pass the intersection, continue
		// ENDREGION

		// REGION: lease management for V2V cars

		// if the car is not near intersections, continue
		// TODO: crossing_intersection(car_id, intersection_id, time, speed, acceleration, direction)
		if false {
			continue
		}
		// if the car is not trying to pass the intersection, continue
		// ENDREGION
	}
}
