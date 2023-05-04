from Quanser.q_essential import Camera2D
import time
import struct
import numpy as np 
import cv2
import sys

# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
photo_path = sys.argv[1]

## Timing Parameters and methods 
startTime = time.time()
def elapsed_time():
    return time.time() - startTime

sampleRate = 30.0
sampleTime = 1/sampleRate
simulationTime = 6
print('Sample Time: ', sampleTime)

# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
# Additional parameters
counter = 0
imageWidth = 640
imageHeight = 480
imageBuffer360 = np.zeros((imageHeight, 4*imageWidth + 60, 3), dtype=np.uint8) # 20 px padding between pieces

# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
## Initialize the CSI cameras
myCam1 = Camera2D(camera_id="0", frame_width=imageWidth, frame_height=imageHeight, frame_rate=sampleRate)
myCam2 = Camera2D(camera_id="1", frame_width=imageWidth, frame_height=imageHeight, frame_rate=sampleRate)
myCam3 = Camera2D(camera_id="2", frame_width=imageWidth, frame_height=imageHeight, frame_rate=sampleRate)
myCam4 = Camera2D(camera_id="3", frame_width=imageWidth, frame_height=imageHeight, frame_rate=sampleRate)
# myCam = Camera2D(camera_id="3", frame_width=imageWidth, frame_height=imageHeight, frame_rate=sampleRate)

# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
## Main Loop
try:
    while elapsed_time() < simulationTime:
        
        # Start timing this iteration
        start = time.time()

        # Capture RGB Image from CSI
        myCam1.read()
        myCam2.read()
        myCam3.read()
        myCam4.read()
        # myCam.read()

        counter += 1

        # End timing this iteration
        end = time.time()

        # Calculate the computation time, and the time that the thread should pause/sleep for
        computationTime = end - start
        sleepTime = sampleTime - ( computationTime % sampleTime )
        
        # Stitch images together with black padding
        # horizontalBlank     = np.zeros((20, 4*imageWidth+100, 3), dtype=np.uint8)
        # verticalBlank       = np.zeros((imageHeight, 20, 3), dtype=np.uint8)
        horizontalBlank     = np.zeros((20, 2*imageWidth+20, 3), dtype=np.uint8)
        verticalBlank       = np.zeros((imageHeight, 20, 3), dtype=np.uint8)

        # imageBuffer360 = np.concatenate((myCam4.image_data, 
        #                                 verticalBlank, 
        #                                 myCam2.image_data, 
        #                                 verticalBlank, 
        #                                 myCam3.image_data, 
        #                                 verticalBlank, 
        #                                 myCam1.image_data), 
        #                                 axis = 1)
        imageBuffer360 = np.concatenate((np.concatenate((myCam4.image_data, verticalBlank, myCam2.image_data), axis = 1), horizontalBlank, np.concatenate((myCam3.image_data, verticalBlank, myCam1.image_data), axis = 1)), axis = 0)


        # imageBuffer = myCam.image_data
        if counter == 60:
            # cv2.imwrite(photo_path, imageBuffer)
            # cv2.imwrite(photo_path, imageBuffer360)
            cv2.imwrite(photo_path, imageBuffer360)
            # cv2.imwrite(photo_path, cv2.resize(imageBuffer360, (int(2*imageWidth), int(imageHeight/2))))
            

        # Display the stitched image at half the resolution
        cv2.imshow('Combined View', cv2.resize(imageBuffer360, (int(imageWidth), int(imageHeight))))
        # cv2.imshow('Combined View', cv2.resize(image, (int(imageWidth), int(imageHeight))))
        # cv2.imwrite("img/test.jpg", cv2.resize(imageBuffer360, (int(2*imageWidth), int(imageHeight/2))))
        # Pause/sleep for sleepTime in milliseconds
        msSleepTime = int(1000*sleepTime)
        if msSleepTime <= 0:
            msSleepTime = 1 # this check prevents an indefinite sleep as cv2.waitKey waits indefinitely if input is 0
        cv2.waitKey(msSleepTime)

except KeyboardInterrupt:
    print("User interrupted!")

finally:
    # Terminate all webcam objects    
    myCam1.terminate()
    myCam2.terminate()
    myCam3.terminate()
    myCam4.terminate()
    # myCam.terminate()
# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 