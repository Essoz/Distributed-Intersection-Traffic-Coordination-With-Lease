package main

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/essoz/car-backend/pkg/lease"
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

// store the intersection struct and block structs into etcd
func initializeIntersectionAndBlock(cli *clientv3.Client, ctx context.Context, pathToIntersection string) string {
	intersection := lease.NewIntersection(pathToIntersection)

	intersection.PutEtcd(cli, ctx, "")

	// put blocks into etcd
	// the size and location of each block is controlled by the split index
	numBlocks := math.Pow(2, intersection.Spec.SplitIndex)
	numBlocksRow := int(math.Sqrt(numBlocks))

	blockSize := make([]float64, 2)
	blockSize[0] = intersection.Spec.Size[0] / float64(numBlocksRow)
	blockSize[1] = intersection.Spec.Size[1] / float64(numBlocksRow)

	for i := 0; i < numBlocksRow; i++ {
		for j := 0; j < numBlocksRow; j++ {
			block := lease.Block{
				Metadata: lease.BlockMetadata{
					Name: "block" + strconv.FormatInt(int64(i), 10) + strconv.FormatInt(int64(j), 10),
				},
				Spec: lease.BlockSpec{
					Location: []float64{
						intersection.Spec.Position[0] + float64(i)*blockSize[0],
						intersection.Spec.Position[1] + float64(j)*blockSize[1],
					},
					Size: blockSize,
				},
			}
			block.PutEtcd(cli, ctx, "")
		}
	}

	return intersection.Metadata.Name
}

// 	// commit transaction
// 	_, err = txn.Commit()
// 	if err != nil {
// 		panic(err)
// 	}
// }

func Run(cli *clientv3.Client, ctx context.Context) {
	// we assume that the car already know which intersection it is at, now it tries to register the intersection into etcd
	intersectionName := initializeIntersectionAndBlock(cli, ctx, "intersection.yaml")

	go RunPerceptionService(cli, ctx, global_car_name)
	go RunControlService(cli, ctx, global_car_name, intersectionName)

	for {
		// fmt.Println("Running")
		time.Sleep(10 * time.Second)
	}
}
