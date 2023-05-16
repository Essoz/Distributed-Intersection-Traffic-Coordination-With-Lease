CAR_URL_1=192.168.2.12
CAR_URL_2=192.168.2.13
CONTROL_URL=192.168.2.16

CAR_1_HOSTNAME="qcar-50577"
CAR_2_HOSTNAME="qcar-50578"
CONTROL_HOSTNAME="MattJiangs-MacBook-Air.local"

ETCD_VER=v3.5.8
REPO_URL=quay.io/coreos/etcd
IMG_URL=$REPO_URL:$ETCD_VER
CONTAINER_NAME=etcd
HOSTNAME=$(hostname)

if [[ "$OSTYPE" == "darwin"* ]]; then 
    CURRENT_URL=$(ipconfig getifaddr en0)
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
	CURRENT_URL=$(hostname -I | awk '{print $1}')
fi

INITIAL_CLUSTER="$CAR_1_HOSTNAME=http://$CAR_URL_1:2380,$CAR_2_HOSTNAME=http://$CAR_URL_2:2380,$CONTROL_HOSTNAME=http://$CONTROL_URL:2380"

# echo all the variables
echo "CAR_URL_1: $CAR_URL_1"
echo "CAR_URL_2: $CAR_URL_2"
echo "CONTROL_URL: $CONTROL_URL"

echo "CURRENT_URL: $CURRENT_URL"
echo "CURRENT_HOSTNAME: $HOSTNAME"
echo "INITIAL_CLUSTER: $INITIAL_CLUSTER"


# docker rmi $IMG_URL || true && \
# 	docker run \
rm -rf $(pwd)/etcd-data
etcd \
		--name $HOSTNAME \
		--data-dir $(pwd)/etcd-data \
		--listen-client-urls http://$CURRENT_URL:2379,http://127.0.0.1:2379 \
		--advertise-client-urls http://$CURRENT_URL:2379 \
		--listen-peer-urls http://$CURRENT_URL:2380 \
		--initial-advertise-peer-urls http://$CURRENT_URL:2380 \
		--initial-cluster $INITIAL_CLUSTER \
		--initial-cluster-token tkn \
		--initial-cluster-state new \
		--log-level info \
		--logger zap \
		--log-outputs stderr

