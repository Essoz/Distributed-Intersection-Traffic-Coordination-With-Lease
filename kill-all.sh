#!/usr/bin/env zsh

car_ip_1="192.168.2.12"
car_ip_2="192.168.2.13"
user_name="nvidia"

session="ece445-demo"

# create a new tmux session called "ece445-demo", delete the session if it already exists
ssh $user_name@$car_ip_1 "tmux kill-session -t $session 2>/dev/null || true"
ssh $user_name@$car_ip_2 "tmux kill-session -t $session 2>/dev/null || true"
ssh $user_name@$car_ip_1 "cd object_detection; echo nvidia | sudo -S python3 HardwareStop.py; cd .."
ssh $user_name@$car_ip_2 "cd object_detection; echo nvidia | sudo -S python3 HardwareStop.py; cd .."
