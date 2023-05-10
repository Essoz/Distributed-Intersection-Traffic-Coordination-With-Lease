import numpy as np

manual_drive = True
speed = 0
imgWidth = 1300
imgHeight = 980
transformed_imgWidth = 1024
transformed_imgHeight = 800
pix2angle = 0.183988

# parameters for car 192.168.2.12
# camera_offset_12 = [-1.931, 1.548, 1.931, 4.235]
camera_offset = [-2.144, 0, -3.215, 0]
# wheel_offset_12 = -0.066
wheel_offset = -0.042

dist_to_head = 0.23
dist_to_tail = 0.19
car_width = 0.24
brake_acc = 3.2
head_angle = np.arctan2(car_width/2, dist_to_head)
tail_angle = np.arctan2(car_width/2, dist_to_tail)
collision_range_head = dist_to_head + 0.2
collision_range_tail = dist_to_tail + 0.2
