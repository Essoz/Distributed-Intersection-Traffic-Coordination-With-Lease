from Quanser.q_misc import BasicStream 
import time 
import numpy as np
import cv2

# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
# Image Parameters
imageWidth = 640
imageHeight = 480
imageChannels = 3
buffer_size = imageHeight*imageWidth*imageChannels # 4 bytes for float32 data, 3 channels for RGB
image_data = np.zeros((imageHeight, imageWidth, imageChannels), dtype=np.uint8)

# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
# Create a BasicStream object configured as a Server agent, with buffer sizes enough to send/receive the image above.
myServer = BasicStream('tcpip://192.168.2.12:18001', agent='s', send_buffer_size=buffer_size, recv_buffer_size=buffer_size)
prev_con = False

# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
## Timing Parameters and methods 
startTime = time.time()
def elapsed_time():
    return time.time() - startTime

sampleRate = 30.0
sampleTime = 1/sampleRate
simulationTime = 30.0
print('Sample Time: ', sampleTime)

# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
# Main Loop
try:
    while elapsed_time() < simulationTime:

        if not myServer.connected:
            myServer.checkConnection()

        if myServer.connected and not prev_con:
            print('Connection to Client was successful.')
        prev_con = myServer.connected

        if myServer.connected:

            # Start timing this iteration
            start = time.time()    

            # Receive data from client 
            image_data, bytes_received = myServer.receive(image_data)
            
            if bytes_received < len(image_data.tobytes()):
                print('Client stopped sending data over.')
                break
            # End timing this iteration
            end = time.time()

            # Calculate the computation time, and the time that the thread should pause/sleep for
            computationTime = end - start
            sleepTime = sampleTime - ( computationTime % sampleTime )

            # Pause/sleep for sleepTime in milliseconds
            msSleepTime = int(1000*sleepTime)
            if msSleepTime <= 0:
                msSleepTime = 1 # this check prevents an indefinite sleep as cv2.waitKey waits indefinitely if input is 0
                            
            cv2.imshow('Server Image Received', image_data)
            cv2.waitKey(msSleepTime)

except KeyboardInterrupt:
    print("User interrupted!")

finally:
    # Terminate Server
    myServer.terminate()
    print('All the right turns in all the right places.')




