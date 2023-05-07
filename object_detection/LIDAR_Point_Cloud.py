from Quanser.q_essential import LIDAR
from Quanser.q_misc import Utilities
from Quanser.q_interpretation import *
import time
import struct
import numpy as np 
import cv2

saturate = Utilities.saturate
# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
## Timing Parameters and methods 
startTime = time.time()
def elapsed_time():
    return time.time() - startTime

sampleRate = 30
sampleTime = 1/sampleRate
simulationTime = 60.0
print('Sample Time: ', sampleTime)

# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
## Additional parameters and buffers
counter = 0
gain = 50 # pixels per meter
dim = 8 * gain # 8 meters width, or 400 pixels side length
decay = 0.2 # 90% decay rate on old map data
map = np.zeros((dim, dim), dtype=np.float32) # map object
max_distance = 3.9

# LIDAR initialization and measurement buffers
myLidar = LIDAR(num_measurements=720)

# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
## Main Loop
try:
    while elapsed_time() < simulationTime:
        # Decay existing map
        map = decay*map

        # Start timing this iteration
        start = time.time()

        # Capture LIDAR data
        myLidar.read()

        # convert angles from lidar frame to body frame
        angles_in_body_frame = lidar_frame_2_body_frame(myLidar.angles)

        # Find the points where it exceed the max distance and drop them off
        idx = [i for i, v in enumerate(myLidar.distances) if v < max_distance and v > 0.15]
        
        # convert distances and angles to XY contour
        x = myLidar.distances[idx]*np.cos(angles_in_body_frame[idx])
        y = myLidar.distances[idx]*np.sin(angles_in_body_frame[idx])

        # print(idx)
        # print(list(zip(angles_in_body_frame[idx], myLidar.distances[idx])))
        # print('')
        # print(list(zip(angles_in_body_frame, myLidar.distances)))
        # print('')

        # convert XY contour to pixels contour and update those pixels in the map
        pX = (dim/2 - x*gain).astype(np.uint16)
        pY = (dim/2 - y*gain).astype(np.uint16)

        map[pX, pY] = 1

        # End timing this iteration
        end = time.time()

        # Calculate the computation time, and the time that the thread should pause/sleep for
        computationTime = end - start
        sleepTime = sampleTime - ( computationTime % sampleTime )
        # time.sleep(sleepTime)
        
        # Display the map at full resolution
        cv2.imshow('Map', map)
        
        # Pause/sleep for sleepTime in milliseconds
        msSleepTime = int(1000*sleepTime)
        if msSleepTime <= 0:
        	msSleepTime = 1 # this check prevents an indefinite sleep as cv2.waitKey waits indefinitely if input is 0
        cv2.waitKey(msSleepTime)


except KeyboardInterrupt:
    print("User interrupted!")

finally:
    # Terminate the LIDAR object
    myLidar.terminate()
# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 