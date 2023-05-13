#!/bin/bash

curr_host=$(hostname -I | awk '{print $1}')

# echo nvidia | sudo -S -E python3 object_tracking_kf.py --weights ./yolov5n6.pt --half --imgsz 800 1024 --conf-thres 0.25 --classes 0 2 3 5 7 9 11 56 --start_pos 0 --heading 0
# if curr_host is 192.168.2.12, set start_pos to 0 and heading to 0
if [ "$curr_host" = "192.168.2.12" ] ; then
    echo nvidia | sudo -S -E python3 object_tracking_kf.py --weights ./yolov5n6.pt --half --imgsz 800 1024 --conf-thres 0.25 --classes 0 2 3 5 7 9 11 56 --start_pos 0 --heading 0
else
    echo nvidia | sudo -S -E python3 object_tracking_kf.py --weights ./yolov5n6.pt --half --imgsz 800 1024 --conf-thres 0.25 --classes 0 2 3 5 7 9 11 56 --start_pos 1 --heading 270
fi    

