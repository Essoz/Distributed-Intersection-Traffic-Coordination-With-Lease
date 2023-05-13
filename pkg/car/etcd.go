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

func GetAllCarsEtcd(cli *clientv3.Client, ctx context.Context, prefix string) []*Car {
	carPrefix := CAR_ETCD_PREFIX + prefix

	resp, err := cli.Get(ctx, carPrefix, clientv3.WithPrefix())
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

func GetAllCarsWithKeyEtcd(cli *clientv3.Client, ctx context.Context, prefix string) map[string]*Car {
	carPrefix := CAR_ETCD_PREFIX + prefix

	resp, err := cli.Get(ctx, carPrefix, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}

	cars := make(map[string]*Car)
	for _, ev := range resp.Kvs {
		var car Car
		err = yaml.Unmarshal(ev.Value, &car)
		if err != nil {
			panic(err)
		}
		cars[string(ev.Key)] = &car
	}

	return cars
}

func (c *Car) PutEtcd(cli *clientv3.Client, ctx context.Context, prefix string) {
	carPrefix := CAR_ETCD_PREFIX + prefix

	carBytes, err := yaml.Marshal(*c)
	if err != nil {
		panic(err)
	}

	_, err = cli.Put(ctx, carPrefix+c.GetName(), string(carBytes))
	if err != nil {
		panic(err)
	}
}

func (c *Car) GetSurroundingCarsEtcd(cli *clientv3.Client, ctx context.Context) []*Car {
	// TODO: if car is not v2v, return nothing
	prefix := c.Metadata.Name + "/surrounding/"
	return GetAllCarsEtcd(cli, ctx, prefix)
}

func (c *Car) UpdateSurroundingCarsEtcd(cli *clientv3.Client, ctx context.Context, surroundingCars []*Car) {
	prefix := c.Metadata.Name + "/surrounding/"

	prevSurrCars := c.GetSurroundingCarsEtcd(cli, ctx)
	// delete the cars that are not in the new list but in the old list
	for _, prevSurrCar := range prevSurrCars {
		// if the car's name is not in the new list, delete it
		name := prevSurrCar.Metadata.Name
		found := false
		for _, surrCar := range surroundingCars {
			if surrCar.Metadata.Name == name {
				found = true
				break
			}
		}
		if !found {
			DeleteCarEtcd(cli, ctx, name, prefix)
		}
	}

	// put the cars in the new list
	for _, surrCar := range surroundingCars {
		surrCar.PutEtcd(cli, ctx, prefix)
	}
}

func DeleteCarEtcd(cli *clientv3.Client, ctx context.Context, carName string, prefix string) {
	carPrefix := CAR_ETCD_PREFIX + prefix
	_, err := cli.Delete(ctx, carPrefix+carName)
	if err != nil {
		panic(err)
	}
}
