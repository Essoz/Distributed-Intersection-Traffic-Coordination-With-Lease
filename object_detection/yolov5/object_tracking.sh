#!/bin/bash
sudo -E python3 object_tracking.py --weights ./yolov5n.pt --imgsz 800 1024 --conf-thres 0.25 --classes 0 2 3 5 7 9 11 56
