package main

import (
	"context"

	"github.com/essoz/v2v-backend/pkg/car"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func ClearAllCarDestinationEtcd(cli *clientv3.Client, ctx context.Context) {
	txn := cli.Txn(ctx)
	defer txn.Commit()

	// get all cars from etcd
	cars := car.GetAllCarsEtcd(cli, ctx, "")

	// clear all cars' destination
	for _, c := range cars {
		c.Dynamics.Destination = []float64{}
		c.PutEtcd(cli, ctx, "")
	}

	// done
}

func SetCarLocation(cli *clientv3.Client, ctx context.Context, carKey string, location []float64) {
	// FIXME: TO BE INTEGRATED WITH XINWEN'S PERCEPTION CODE
	car := car.GetCarEtcd(cli, ctx, carKey)
	car.Dynamics.Location = location
	car.PutEtcd(cli, ctx, "")
}

func SetCarDestination(cli *clientv3.Client, ctx context.Context, carKey string, destination []float64) {
	txn := cli.Txn(ctx)
	car := car.GetCarEtcd(cli, ctx, carKey)
	car.Dynamics.Destination = destination
	car.PutEtcd(cli, ctx, "")
	txn.Commit()
}
