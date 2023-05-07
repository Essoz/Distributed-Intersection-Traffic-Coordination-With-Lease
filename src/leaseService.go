package main

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func RunLeaseService(cli *clientv3.Client, ctx context.Context, carName string, intersectionName string) {
	// THIS FUNCTION WILL INTERFACE THE LEASE SERVICE

	// TODO: Observe all cars passing through the intersection

	// TODO: Make leases (and necessary, override leases) for non-v2v cars passing through the intersection

	// TODO:
	
}
