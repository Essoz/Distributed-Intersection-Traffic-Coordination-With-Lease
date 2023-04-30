package car

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v3"
)

const CAR_ETCD_PREFIX = "car/"

func GetCarEtcd(cli *clientv3.Client, ctx context.Context, carName string) Car {
	resp, err := cli.Get(ctx, CAR_ETCD_PREFIX+carName)
	if err != nil {
		panic(err)
	}

	var car Car
	err = yaml.Unmarshal(resp.Kvs[0].Value, &car)
	if err != nil {
		panic(err)
	}

	return car
}

func GetAllCarsEtcd(cli *clientv3.Client, ctx context.Context) []*Car {
	resp, err := cli.Get(ctx, CAR_ETCD_PREFIX, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}

	var cars []*Car
	for _, ev := range resp.Kvs {
		var car Car
		err = yaml.Unmarshal(ev.Value, &car)
		if err != nil {
			panic(err)
		}
		cars = append(cars, &car)
	}

	return cars
}

func (c *Car) PutEtcd(cli *clientv3.Client, ctx context.Context) {
	car_bytes, err := yaml.Marshal(*c)
	if err != nil {
		panic(err)
	}

	_, err = cli.Put(ctx, CAR_ETCD_PREFIX+c.GetName(), string(car_bytes))
	if err != nil {
		panic(err)
	}
}
