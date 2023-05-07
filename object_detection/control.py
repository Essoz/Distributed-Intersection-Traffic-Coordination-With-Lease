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


myCar = QCar()
t0 = time.time()
dist = 0
sampleRate = 60
sampleTime = 1/sampleRate
diff = Calculus().differentiator_variable(sampleTime)
next(diff)
current_speed = 0.0

def speed_ctrl(target_speed):
    global current_speed
    sampleRate = 60
    sampleTime = 1/sampleRate

    LEDs = np.array([0, 0, 0, 0, 0, 0, 1, 1])
    mtr_cmd = np.array([0,-0.066]) 
    _,_, encoderCounts = myCar.read_std()  
    encoderSpeed = diff.send((encoderCounts, sampleTime))
    current_speed = basic_speed_estimation(encoderSpeed)

    # the parameter for straight run
    mtr_cmd[0] = speed_control(target_speed,current_speed,1,sampleTime)
    myCar.write_mtrs(mtr_cmd)
    myCar.write_LEDs(LEDs)

    return

if __name__ == "__main__":
    try:
        while True:
            start = time.time()
            print(dist, current_speed)
            dist += current_speed/60
            if (dist) < 1:
                speed_ctrl(0.5)
            else:
                speed_ctrl(0.1)
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
        myCar.terminate()
