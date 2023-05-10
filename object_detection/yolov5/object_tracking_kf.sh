#!/bin/bash
sudo -E python3 object_tracking_kf.py --weights ./yolov5n6.pt --half --imgsz 800 1024 --conf-thres 0.25 --classes 0 2 3 5 7 9 11 56 --start_pos 0 --heading 0
