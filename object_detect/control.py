from Quanser.product_QCar import QCar
from Quanser.q_control import *
from Quanser.q_dp import *
from Quanser.q_interpretation import *
from Quanser.q_misc import *
from Quanser.q_ui import *
from Quanser.q_essential import *
import time
import struct
import numpy as np 
import math


sampleRate = 60
sampleTime = 1/sampleRate
print('Sample Time: ', sampleTime)

diff = Calculus().differentiator_variable(sampleTime)
_ = next(diff)
myCar = QCar()
t0 = time.time()
dist = 0

try:
	while True:
		start = time.time()
		# # Find slope and intercept of linear fit from the binary image
		# slope, intercept = find_slope_intercept_from_binary(binary)

		# # steering from slope and intercept
		# raw_steering = 1.5*(slope - 0.3419) + (1/150)*(intercept+5)
		# steering = steering_filter.send((saturate(raw_steering, 0.5, -0.5), dt))

		# steering = 0

		mtr_cmd = np.array([0,-0.066]) # the parameter for straight run
		# mtr_cmd[0] = mtr_cmd[0]*np.cos(steering)
		# myCar.write_mtrs(mtr_cmd)
		LEDs = np.array([0, 0, 0, 0, 0, 0, 1, 1])
		current, batteryVoltage, encoderCounts = myCar.read_write_std(mtr_cmd, LEDs)  
		encoderSpeed = diff.send((encoderCounts, sampleTime))
		
		current_speed = basic_speed_estimation(encoderSpeed)
		print("Current Speed: {}".format(current_speed))
		dist += current_speed*sampleTime
		mtr_cmd[0] = speed_control(0.5,current_speed,2,sampleTime)
		print("mtr_cmd: {}".format(mtr_cmd))
		myCar.write_mtrs(mtr_cmd)
		
		end = time.time()
		dt = end - start
		computationTime = end - start
		sleepTime = sampleTime - ( computationTime % sampleTime )
		time.sleep(sleepTime)
		if (dist) > 2: 
			break

except KeyboardInterrupt:
	print("User interrupted!")

finally:
	# Terminate camera and QCar
	myCar.terminate()

	