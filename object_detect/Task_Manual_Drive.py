from Quanser.product_QCar import QCar
from Quanser.q_ui import gamepadViaTarget
from Quanser.q_misc import Calculus
from Quanser.q_interpretation import basic_speed_estimation
import os
import time
import struct
import numpy as np 

# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
## Timing Parameters and methods 
startTime = time.time()
def elapsed_time():
    return time.time() - startTime

sampleRate = 50
sampleTime = 1/sampleRate
simulationTime = 60.0
print('Sample Time: ', sampleTime)

# Additional parameters
counter = 0

# Initialize motor command array
mtr_cmd = np.array([0,0])

# Set up a differentiator to get encoderSpeed from encoderCounts
diff = Calculus().differentiator_variable(sampleTime)
_ = next(diff)
timeStep = sampleTime

# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
## QCar and Gamepad Initialization
myCar = QCar()
gpad = gamepadViaTarget(1)

# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
## Driving Configuration: Use 3 toggles or 4 toggles mode as you see fit:
# Common to both 3 or 4 mode
#   Steering                    - Left Lateral axis
#   Arm                         - LB
# In 3 mode: 
#   Throttle (Drive or Reverse) - Right Longitudonal axis
# In 4 mode:
#   Throttle                    - Right Trigger (always positive)
#   Button A                    - Reverse if held, Drive otherwise
configuration = '3' # change to '4' if required

# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
# Reset startTime before Main Loop
startTime = time.time()

## Main Loop
try:
    while elapsed_time() < simulationTime:
        # Start timing this iteration
        start = elapsed_time()

        # Read Gamepad states
        new = gpad.read()

        # Basic IO - write motor commands
        if configuration == '3':
            if new and gpad.LB:
                mtr_cmd = np.array([0.3*gpad.RLO, 0.5*gpad.LLA])            
        elif configuration == '4':
            if new and gpad.LB:
                if gpad.A:
                    mtr_cmd = np.array([-0.3*gpad.RT, 0.5*gpad.LLA])
                else:
                    mtr_cmd = np.array([0.3*gpad.RT, 0.5*gpad.LLA])            
        LEDs = np.array([0, 0, 0, 0, 0, 0, 1, 1])

        # Adjust LED indicators based on steering and reverse indicators based on reverse gear
        if mtr_cmd[1] > 0.3:
            LEDs[0] = 1
            LEDs[2] = 1
        elif mtr_cmd[1] < -0.3:
            LEDs[1] = 1
            LEDs[3] = 1
        if mtr_cmd[0] < 0:
            LEDs[5] = 1

        # Perform I/O
        current, batteryVoltage, encoderCounts = myCar.read_write_std(mtr_cmd, LEDs)        
        
        # Differentiate encoder counts and then estimate linear speed in m/s
        encoderSpeed = diff.send((encoderCounts, timeStep))
        linearSpeed = basic_speed_estimation(encoderSpeed)

        # End timing this iteration
        end = elapsed_time()

        # Calculate computation time, and the time that the thread should pause/sleep for
        computation_time = end - start
        sleep_time = sampleTime - computation_time%sampleTime

        # Pause/sleep and print out the current timestamp
        time.sleep(sleep_time)

        if new:
            os.system('clear')
            print("Car Speed:\t\t\t{0:1.2f}\tm/s\nRemaining battery capacity:\t{1:4.2f}\t%\nMotor throttle:\t\t\t{2:4.2f}\t% PWM\nSteering:\t\t\t{3:3.2f}\trad"
                                                            .format(linearSpeed, 100 - (batteryVoltage - 10.5)*100/(12.6 - 10.5), mtr_cmd[0], mtr_cmd[1]))
        timeAfterSleep = elapsed_time()
        timeStep = timeAfterSleep - start
        counter += 1

except KeyboardInterrupt:
    print("User interrupted!")

finally:    
    myCar.terminate()
# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 