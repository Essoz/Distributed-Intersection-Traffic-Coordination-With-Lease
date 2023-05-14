#!/usr/bin/env bash

# this script runs service on the car

session="ece445-demo"

# create a new tmux session called "ece445-demo", delete the session if it already exists
bash ./stop-service-local.sh
tmux new-session -d -s ece445-demo

# create a new window called "etcd" in the session
tmux new-window -t $session -n 'etcd'
tmux send-keys -t $session:etcd 'bash ./etcd-up.sh' C-m

tmux new-window -t $session -n 'tracking-service'
tmux send-keys -t $session:tracking-service 'cd ./object_detection/yolov5/; bash object_tracking_kf.sh' C-m

tmux new-window -t $session -n 'service'
tmux send-keys -t $session:service './service_arm64' C-m
