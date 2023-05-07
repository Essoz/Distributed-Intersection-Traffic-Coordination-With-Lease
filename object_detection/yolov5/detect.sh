#!/bin/bash
python3 detect.py --source ../photo360test1.png --nosave --weights ./yolov5n6.pt --imgsz 800 1024 --classes 0 2 3 5 7 9 11 56
