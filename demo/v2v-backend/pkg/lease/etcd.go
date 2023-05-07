package lease

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
	yaml "gopkg.in/yaml.v3"
)

const BLOCK_ETCD_PREFIX = "block/"
const INTERSECTION_ETCD_PREFIX = "intersection/"

func GetBlockEtcd(cli *clientv3.Client, ctx context.Context, blockName string) Block {
	resp, err := cli.Get(ctx, BLOCK_ETCD_PREFIX+blockName)
	if err != nil {
		panic(err)
	}

	var block Block
	err = yaml.Unmarshal(resp.Kvs[0].Value, &block)
	if err != nil {
		panic(err)
	}
	return block
}

func GetAllBlocksEtcd(cli *clientv3.Client, ctx context.Context) []Block {
	resp, err := cli.Get(ctx, BLOCK_ETCD_PREFIX, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}

	var blocks []Block
	for _, ev := range resp.Kvs {
		var block Block
		err = yaml.Unmarshal(ev.Value, &block)
		if err != nil {
			panic(err)
		}
		blocks = append(blocks, block)
	}
	return blocks
}

func (b *Block) PutBlockEtcd(cli *clientv3.Client, ctx context.Context) {
	block_bytes, err := yaml.Marshal(*b)
	if err != nil {
		panic(err)
	}

	_, err = cli.Put(ctx, BLOCK_ETCD_PREFIX+b.GetName(), string(block_bytes))
	if err != nil {
		panic(err)
	}
}

func GetLeasesRelatedToCarEtcd(cli *clientv3.Client, ctx context.Context, carName string) []Lease {
	blocks := GetAllBlocksEtcd(cli, ctx)

	var leases []Lease
	for _, block := range blocks {
		for _, lease := range block.GetLeases() {
			if lease.GetCarName() == carName {
				leases = append(leases, lease)
			}
		}
	}

	return leases
}

func GetIntersectionEtcd(cli *clientv3.Client, ctx context.Context, intersectionName string) Intersection {
	resp, err := cli.Get(ctx, INTERSECTION_ETCD_PREFIX+intersectionName)
	if err != nil {
		panic(err)
	}

	var intersection Intersection
	err = yaml.Unmarshal(resp.Kvs[0].Value, &intersection)
	if err != nil {
		panic(err)
	}

	return intersection
}
