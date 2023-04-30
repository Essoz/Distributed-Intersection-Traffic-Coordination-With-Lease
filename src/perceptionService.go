package main

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func RunPerceptionService(cli *clientv3.Client, ctx context.Context) {
	// THIS FUNCTION WILL INTERFACE THE PERCEPTION SERVICE

	// TODO: Get all surrounding objects

	// TODO: Register / Update all surrounding objects into etcd

	// TODO: Update the car's state in etcd

	panic("Not implemented")
}
