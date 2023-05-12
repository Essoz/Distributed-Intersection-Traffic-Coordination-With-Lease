#!/usr/bin/env zsh

car_ip_1="192.168.2.12"
car_ip_2="192.168.2.13"
user_name="nvidia"

# 1. build the service go binary
cd src; zsh build_for_car.sh; cd ..

# 2. scp the binary to the car
scp src/service_arm64 $user_name@$car_ip_1:/home/$user_name
scp src/service_arm64 $user_name@$car_ip_2:/home/$user_name

# 3. scp the object_detection model to the car
scp -r object_detection $user_name@$car_ip_1:/home/$user_name
scp -r object_detection $user_name@$car_ip_2:/home/$user_name
