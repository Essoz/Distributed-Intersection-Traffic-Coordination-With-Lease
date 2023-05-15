import numpy as np
import socket

hostname = socket.gethostname()

manual_drive = False
speed = 0
imgWidth = 1300
imgHeight = 980
transformed_imgWidth = 1024
transformed_imgHeight = 800
pix2angle = 0.183988

# parameters for car 192.168.2.12
if hostname.strip() == 'qcar-50577':
    camera_offset = [-1.931, 1.548, 1.931, 4.235]
    wheel_offset = -0.0605
else: 
    camera_offset = [-2.144, 0, -3.215, 0]
    wheel_offset = -0.055

dist_to_head = 0.23
dist_to_tail = 0.19
car_width = 0.24
brake_acc = 3
head_angle = np.arctan2(car_width/2, dist_to_head)
tail_angle = np.arctan2(car_width/2, dist_to_tail)
collision_range_head = dist_to_head + 0.2
collision_range_tail = dist_to_tail + 0.2
