#!/usr/bin/env zsh

car_ip_1="192.168.2.12"
car_ip_2="192.168.2.13"
user_name="nvidia"

ssh $user_name@$car_ip_1 "bash start-service-local.sh" &
ssh $user_name@$car_ip_2 "bash start-service-local.sh" &
