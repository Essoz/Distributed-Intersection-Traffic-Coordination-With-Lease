from Quanser.q_ui import gamepadViaTarget
import time
import numpy as np
import os

# Timing and Initialization
startTime = time.time()
def elapsed_time():
	return time.time() - startTime

simulationTime = 60
sampleRate = 100
sampleTime = 1/sampleRate

# Unless there are more than 1 usb receiver plugged in, now the gamepadViaTarget will always be 1
gpad = gamepadViaTarget(1) 

# Restart starTime just before Main Loop
startTime = time.time()

## Main Loop
try:
	while elapsed_time() < simulationTime:
		# Start timing this iteration
		start = elapsed_time()

		# Basic IO - write motor commands
		new = gpad.read()
		
		if new:
			# Clear the Screen for better readability
			os.system('clear')

			# Print out the gamepad IO read
			print("Left Laterial:\t\t{0:.2f}\nLeft Longitudonal:\t{1:.2f}\nLeft Trigger:\t\t{2:.2f}\nRight Lateral:\t\t{3:.2f}\nRight Longitudonal:\t{4:.2f}\nRight Trigger:\t\t{5:.2f}"
																.format(gpad.LLA, gpad.LLO, gpad.LT, gpad.RLA, gpad.RLO, gpad.RT))
			print("Button A:\t\t{0:.0f}\nButton B:\t\t{1:.0f}\nButton X:\t\t{2:.0f}\nButton Y:\t\t{3:.0f}\nButton LB:\t\t{4:.0f}\nButton RB:\t\t{5:.0f}"
																.format(gpad.A, gpad.B, gpad.X, gpad.Y, gpad.LB, gpad.RB))
			print("Up:\t\t\t{0:.0f}\nRight:\t\t\t{1:.0f}\nDown:\t\t\t{2:.0f}\nLeft:\t\t\t{3:.0f}"
																.format(gpad.up, gpad.right, gpad.down, gpad.left))

    	# End timing this iteration
		end = elapsed_time()

    	# Calculate computation time, and the time that the thread should pause/sleep for
		computation_time = end - start
		sleep_time = sampleTime - computation_time%sampleTime
		
    	# Pause/sleep and print out the current timestamp
		time.sleep(sleep_time)

except KeyboardInterrupt:
	print("User interrupted!")

finally:    
	# Terminate Joystick properly
	gpad.terminate()
