import time
import control

t0 = time.time()
dist = 0
sampleRate = 60
sampleTime = 1/sampleRate

try:
    while True:
        start = time.time()
        print(dist, control.current_speed)
        dist += control.current_speed/60
        if (dist) < 1:
            control.speed_ctrl(0.5)
        else:
            control.speed_ctrl(0.1)
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
    # Terminate the LIDAR object
    control.myCar.terminate()