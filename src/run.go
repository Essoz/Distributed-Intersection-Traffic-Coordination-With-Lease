package main

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// func initializeCar(cli *clientv3.Client, ctx context.Context, pathToCar string) {
// 	txn := cli.Txn(ctx)
// 	defer txn.Commit()

// 	car := car.NewCar(pathToCar)

// 	car.PutEtcd(cli, ctx, "")

// 	global_car_name = car.Metadata.Name

// 	// commit transaction
// 	_, err := txn.Commit()
// 	if err != nil {
// 		panic(err)
// 	}
// }

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

func Run(cli *clientv3.Client, ctx context.Context) {
	RunPerceptionService(cli, ctx, global_car_name)
}
