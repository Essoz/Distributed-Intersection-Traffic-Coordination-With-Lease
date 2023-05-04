# YOLOv5 by Ultralytics, GPL-3.0 license
"""
Run inference on images, videos, directories, streams, etc.

Usage:
    $ python path/to/detect.py --source path/to/img.jpg --weights yolov5s.pt --img 640
"""

import argparse
import os
import sys
from pathlib import Path

import cv2
import numpy as np
import torch
import torch.backends.cudnn as cudnn

FILE = Path(__file__).resolve()
ROOT = FILE.parents[0]  # YOLOv5 root directory
if str(ROOT) not in sys.path:
    sys.path.append(str(ROOT))  # add ROOT to PATH
ROOT = Path(os.path.relpath(ROOT, Path.cwd()))  # relative

from models.experimental import attempt_load
from utils.datasets import LoadImages, LoadStreams
from utils.general import apply_classifier, check_img_size, check_imshow, check_requirements, check_suffix, colorstr, \
    increment_path, non_max_suppression, print_args, save_one_box, scale_coords, set_logging, \
    strip_optimizer, xyxy2xywh
from utils.plots import Annotator, colors
from utils.torch_utils import load_classifier, select_device, time_sync
from utils.augmentations import letterbox

# multi thread
import threading

# import Qcar packages
from Quanser.q_essential import Camera2D, LIDAR
from Quanser.q_misc import Utilities
from Quanser.q_interpretation import *
import time


# Timing
# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
## Timing Parameters and methods 
startTime = time.time()
def elapsed_time():
    return time.time() - startTime

simulationTime = 60.0
frameRate = 30.0
thread_terminate = False
# print('Sample Time: ', sampleTime)


# Qcar 360-camera module
# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
# Additional parameters
counter = 0
imageWidth = 640
imageHeight = 480
horizontalBlank     = np.zeros((20, 2*imageWidth+20, 3), dtype=np.uint8)
verticalBlank       = np.zeros((imageHeight, 20, 3), dtype=np.uint8)
# imageBuffer360 = np.zeros((2*imageHeight + 20, 2*imageWidth + 20, 3), dtype=np.uint8) # 20 px padding between pieces

# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
## Initialize the CSI cameras
myCam1 = Camera2D(camera_id="0", frame_width=imageWidth, frame_height=imageHeight, frame_rate=frameRate)
myCam2 = Camera2D(camera_id="1", frame_width=imageWidth, frame_height=imageHeight, frame_rate=frameRate)
myCam3 = Camera2D(camera_id="2", frame_width=imageWidth, frame_height=imageHeight, frame_rate=frameRate)
myCam4 = Camera2D(camera_id="3", frame_width=imageWidth, frame_height=imageHeight, frame_rate=frameRate)


# Qcar LIDAR module
# -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
## Additional parameters and buffers
# gain = 50 # pixels per meter
# dim = 8 * gain # 8 meters width, or 400 pixels side length
# decay = 0.9 # 90% decay rate on old map data
# map = np.zeros((dim, dim), dtype=np.float32) # map object
max_distance = 4
new_map = []
map = []
mutex_map = threading.Lock()

img_transform = None
mutex_img = threading.Lock()

# LIDAR initialization and measurement buffers
myLidar = LIDAR(num_measurements=1000)

img_w = 1300
img_h = 980
compressed_w = 1024
compressed_h = 800

camera_offset = [-1.931, 1.548, 1.931, 4.235]
pix2angle_const = 0.183988
resize_rate_x = img_w / compressed_w
resize_rate_y = img_h / compressed_h

# original png size: x=1300 y=980, after resize: x=640, y=512
# process prediction of cars
# partially predict the position of vanished car
# map = zip(angles_in_body_frame[idx], myLidar.distances[idx])
# return [(x_position, y_position, radius), ...]
def process_pred(pred, map):
    pred = pred[0]
    cars_angle = []
    cars_position = []
    for obj_pred in pred:
        obj_pred = obj_pred.detach().cpu().numpy()
        # print(obj_pred)
        if int(obj_pred[5]) in [2, 3, 5, 7]:
            # Cam4
            if obj_pred[2] * resize_rate_x < 640 and obj_pred[3] * resize_rate_y < 480:
                real_angle_l = ((-(obj_pred[0]*resize_rate_x - 320)*pix2angle_const + camera_offset[3])*np.pi/180)
                real_angle_r = ((-(obj_pred[2]*resize_rate_x - 320)*pix2angle_const + camera_offset[3])*np.pi/180)
                cars_angle.append((real_angle_l, real_angle_r))
            # Cam2
            elif obj_pred[2] * resize_rate_x > 660 and obj_pred[3] * resize_rate_y < 480:
                real_angle_l = ((-(obj_pred[0]*resize_rate_x - 980)*pix2angle_const + camera_offset[3] + 180)*np.pi/180)
                real_angle_r = ((-(obj_pred[2]*resize_rate_x - 980)*pix2angle_const + camera_offset[3] + 180)*np.pi/180)
                cars_angle.append((real_angle_l, real_angle_r))
            # Cam3
            elif obj_pred[2] * resize_rate_x < 640 and obj_pred[3] * resize_rate_y > 500:
                real_angle_l = ((-(obj_pred[0]*resize_rate_x - 320)*pix2angle_const + camera_offset[3] + 90)*np.pi/180)
                real_angle_r = ((-(obj_pred[2]*resize_rate_x - 320)*pix2angle_const + camera_offset[3] + 90)*np.pi/180)
                cars_angle.append((real_angle_l, real_angle_r))
            # Cam1
            elif obj_pred[2] * resize_rate_x > 660 and obj_pred[3] * resize_rate_y > 500:
                real_angle_l = ((-(obj_pred[0]*resize_rate_x - 980)*pix2angle_const + camera_offset[3] + 270)*np.pi/180)
                real_angle_r = ((-(obj_pred[2]*resize_rate_x - 980)*pix2angle_const + camera_offset[3] + 270)*np.pi/180)
                cars_angle.append((real_angle_l, real_angle_r))

    for car in cars_angle:
        dist_list = []

        for angle, dist in map:
            if (angle > 0.15*car[0] + 0.85*car[1] and angle < 0.85*car[0] + 0.15*car[1]) or\
               (angle+2*np.pi > 0.15*car[0] + 0.85*car[1] and angle+2*np.pi < 0.85*car[0] + 0.15*car[1]):
                dist_list.append(dist)

        print("car:", car)
        print("dist_list:", dist_list)

        if dist_list:
            median = np.median(dist_list)
            dist_list = [i for i in dist_list if np.abs(i - median) < 0.1]
            if len(dist_list) != 0:
                real_dist = sum(dist_list) / len(dist_list)
                real_angle = (real_angle_l + real_angle_r)/2
                cars_position.append((np.cos(real_angle)*real_dist, np.sin(real_angle)*real_dist))
        else:
            print("cnm, empty dist_list")
            print(map)

    # remove redundant cars
    cars = []
    for car_1 in cars_position:
        if cars:
            for car_2 in cars:
                if ((car_1[0] - car_2[0])**2 + (car_1[1] - car_2[1])**2)**0.5 < 0.1:
                    car_1 = None
                    break
            if car_1:
                cars.append(car_1)
        else:
            cars.append(car_1)

    return cars


# Kalman filter and predictor
def KalmanFilter(z,  n_iter = 20):
    # suppose A=1，H=1
    # intial parameters  
    sz = (n_iter,) # size of array   

    #Q = 1e-5 # process variance  
    Q = 1e-6 # process variance   
    # allocate space for arrays  
    xhat=np.zeros(sz)      # a posteri estimate of x  
    P=np.zeros(sz)         # a posteri error estimate  
    xhatminus=np.zeros(sz) # a priori estimate of x  
    Pminus=np.zeros(sz)    # a priori error estimate  
    K=np.zeros(sz)         # gain or blending factor  

    R = 0.1**2 # estimate of measurement variance, change to see effect  

    # intial guesses  
    xhat[0] = 0.0  
    P[0] = 1.0  
    A = 1
    H = 1

    for k in range(1,n_iter):  
        # time update  
        xhatminus[k] = A * xhat[k-1]  #X(k|k-1) = AX(k-1|k-1) + BU(k) + W(k),A=1,BU(k) = 0  
        Pminus[k] = A * P[k-1]+Q      #P(k|k-1) = AP(k-1|k-1)A' + Q(k) ,A=1  

        # measurement update  
        K[k] = Pminus[k]/( Pminus[k]+R ) #Kg(k)=P(k|k-1)H'/[HP(k|k-1)H' + R],H=1  
        xhat[k] = xhatminus[k]+K[k]*(z[k]-H * xhatminus[k]) #X(k|k) = X(k|k-1) + Kg(k)[Z(k) - HX(k|k-1)], H=1  
        P[k] = (1-K[k] * H) * Pminus[k] #P(k|k) = (1 - Kg(k)H)P(k|k-1), H=1  
    return xhat


def predict_location(last_cars_position, cars_position):
    pass


def camera_receiver(imgsz=640,  # inference size (pixels)
                    stride=64,
                    pt=True,
                    onnx=False,
                    device='',
                    half=False,  # use FP16 half-precision inference
                    ):
    global img_transform, mutex_img, thread_terminate
    sampleRate = 10.0
    sampleTime = 1/sampleRate
    ## Main Loop
    try:
        while elapsed_time() < simulationTime and not thread_terminate:
            start = time.time()

            # 360-camera
            myCam1.read()
            myCam2.read()
            myCam3.read()
            myCam4.read()

            imageBuffer360 = np.concatenate((np.concatenate((myCam4.image_data, verticalBlank, myCam2.image_data), axis = 1), horizontalBlank, np.concatenate((myCam3.image_data, verticalBlank, myCam1.image_data), axis = 1)), axis = 0)
            img = letterbox(imageBuffer360, imgsz, stride, pt)[0]
            img = img.transpose((2, 0, 1))[::-1]  # HWC to CHW, BGR to RGB
            img = np.ascontiguousarray(img)

            mutex_img.acquire()
            img_transform = img
            mutex_img.release()

            end = time.time()
            computationTime = end - start
            sleepTime = sampleTime - (computationTime % sampleTime)
            time.sleep(sleepTime)

    except KeyboardInterrupt:
        print("Camera: User interrupted!")

    finally:
        print("terminate camera")
        # Terminate all webcam objects    
        myCam1.terminate()
        myCam2.terminate()
        myCam3.terminate()
        myCam4.terminate()

    return


def lidar_receiver(*args, **kwargs):
    global new_map, map, mutex_map, thread_terminate
    sampleRate = 10.0
    sampleTime = 1/sampleRate
    ## Main Loop
    try:
        while elapsed_time() < simulationTime and not thread_terminate:
            start = time.time()

            # LIDAR
            # Decay existing map
            # map = decay*map

            # Capture LIDAR data
            myLidar.read()

            # convert angles from lidar frame to body frame
            angles_in_body_frame = lidar_frame_2_body_frame(myLidar.angles)

            # Find the points where it exceed the max distance and drop them off
            idx = [i for i, v in enumerate(myLidar.distances) if v < max_distance and v > 0.1]
            # print(list(zip(angles_in_body_frame[idx], myLidar.distances[idx])))

            # convert distances and angles to XY contour
            # x = myLidar.distances[idx]*np.cos(angles_in_body_frame[idx])
            # y = myLidar.distances[idx]*np.sin(angles_in_body_frame[idx])

            # convert XY contour to pixels contour and update those pixels in the map
            # pX = (dim/2 - x*gain).astype(np.uint16)
            # pY = (dim/2 - y*gain).astype(np.uint16)
            # map[pX, pY] = 1

            if map:
                for angle, dist in map:
                    if angle > angles_in_body_frame[idx[0]]:
                        new_map.append((angle, dist))
                    else:
                        break

            old_angle = 0.5 * np.pi
            for angle, dist in zip(angles_in_body_frame[idx], myLidar.distances[idx]):
                if angle > old_angle:
                    break
                old_angle = angle
                new_map.append((angle, dist))

            # update new map
            mutex_map.acquire()
            map = new_map
            mutex_map.release()
            new_map = []
            # print(map)

            end = time.time()
            computationTime = end - start
            sleepTime = sampleTime - (computationTime % sampleTime)
            time.sleep(sleepTime)

    except KeyboardInterrupt:
        print("LIDAR: User interrupted!")

    finally:
        print("terminate LIDAR")
        # Terminate the LIDAR object
        myLidar.terminate()

    return


@torch.no_grad()
def run(weights=ROOT / 'yolov5s.pt',  # model.pt path(s)
        # source=ROOT / 'data/images',  # file/dir/URL/glob, 0 for webcam
        imgsz=640,  # inference size (pixels)
        conf_thres=0.25,  # confidence threshold
        iou_thres=0.45,  # NMS IOU threshold
        max_det=1000,  # maximum detections per image
        device='',  # cuda device, i.e. 0 or 0,1,2,3 or cpu
        view_img=False,  # show results
        save_txt=False,  # save results to *.txt
        save_conf=False,  # save confidences in --save-txt labels
        save_crop=False,  # save cropped prediction boxes
        nosave=False,  # do not save images/videos
        classes=None,  # filter by class: --class 0, or --class 0 2 3
        agnostic_nms=False,  # class-agnostic NMS
        augment=False,  # augmented inference
        visualize=False,  # visualize features
        update=False,  # update all models
        project=ROOT / 'runs/detect',  # save results to project/name
        name='exp',  # save results to project/name
        exist_ok=False,  # existing project/name ok, do not increment
        line_thickness=3,  # bounding box thickness (pixels)
        hide_labels=False,  # hide labels
        hide_conf=False,  # hide confidences
        half=False,  # use FP16 half-precision inference
        dnn=False,  # use OpenCV DNN for ONNX inference
        ):
    global mutex_img, mutex_map, img_transform, map
    sampleRate = 10.0
    sampleTime = 1/sampleRate
    # source = str(source)
    # save_img = not nosave and not source.endswith('.txt')  # save inference images
    # webcam = source.isnumeric() or source.endswith('.txt') or source.lower().startswith(
    #     ('rtsp://', 'rtmp://', 'http://', 'https://'))

    # Directories
    save_dir = increment_path(Path(project) / name, exist_ok=exist_ok)  # increment run
    (save_dir / 'labels' if save_txt else save_dir).mkdir(parents=True, exist_ok=True)  # make dir

    # Initialize
    set_logging()
    device = select_device(device)
    half &= device.type != 'cpu'  # half precision only supported on CUDA

    # Load model
    w = str(weights[0] if isinstance(weights, list) else weights)
    classify, suffix, suffixes = False, Path(w).suffix.lower(), ['.pt', '.onnx', '.tflite', '.pb', '']
    check_suffix(w, suffixes)  # check weights have acceptable suffix
    pt, onnx, tflite, pb, saved_model = (suffix == x for x in suffixes)  # backend booleans
    stride, names = 64, [f'class{i}' for i in range(1000)]  # assign defaults
    if pt:
        model = torch.jit.load(w) if 'torchscript' in w else attempt_load(weights, map_location=device)
        stride = int(model.stride.max())  # model stride
        names = model.module.names if hasattr(model, 'module') else model.names  # get class names
        if half:
            model.half()  # to FP16
        if classify:  # second-stage classifier
            modelc = load_classifier(name='resnet50', n=2)  # initialize
            modelc.load_state_dict(torch.load('resnet50.pt', map_location=device)['model']).to(device).eval()
    elif onnx:
        if dnn:
            # check_requirements(('opencv-python>=4.5.4',))
            net = cv2.dnn.readNetFromONNX(w)
        else:
            check_requirements(('onnx', 'onnxruntime'))
            import onnxruntime
            session = onnxruntime.InferenceSession(w, None)
    else:  # TensorFlow models
        check_requirements(('tensorflow>=2.4.1',))
        import tensorflow as tf
        if pb:  # https://www.tensorflow.org/guide/migrate#a_graphpb_or_graphpbtxt
            def wrap_frozen_graph(gd, inputs, outputs):
                x = tf.compat.v1.wrap_function(lambda: tf.compat.v1.import_graph_def(gd, name=""), [])  # wrapped import
                return x.prune(tf.nest.map_structure(x.graph.as_graph_element, inputs),
                               tf.nest.map_structure(x.graph.as_graph_element, outputs))

            graph_def = tf.Graph().as_graph_def()
            graph_def.ParseFromString(open(w, 'rb').read())
            frozen_func = wrap_frozen_graph(gd=graph_def, inputs="x:0", outputs="Identity:0")
        elif saved_model:
            model = tf.keras.models.load_model(w)
        elif tflite:
            interpreter = tf.lite.Interpreter(model_path=w)  # load TFLite model
            interpreter.allocate_tensors()  # allocate
            input_details = interpreter.get_input_details()  # inputs
            output_details = interpreter.get_output_details()  # outputs
            int8 = input_details[0]['dtype'] == np.uint8  # is TFLite quantized uint8 model
    imgsz = check_img_size(imgsz, s=stride)  # check image size

    # Dataloader
    # if webcam:
    #     view_img = check_imshow()
    #     cudnn.benchmark = True  # set True to speed up constant image size inference
    #     dataset = LoadStreams(source, img_size=imgsz, stride=stride, auto=pt)
    #     bs = len(dataset)  # batch_size
    # else:
    #     dataset = LoadImages(source, img_size=imgsz, stride=stride, auto=pt)
    #     bs = 1  # batch_size
    # vid_path, vid_writer = [None] * bs, [None] * bs

    lidarThread = threading.Thread(target=lidar_receiver, args=tuple())
    cameraThread = threading.Thread(target=camera_receiver, args=(imgsz, stride, pt, onnx, device, half))
    lidarThread.start()
    cameraThread.start()

    # Run inference
    if pt and device.type != 'cpu':
        model(torch.zeros(1, 3, *imgsz).to(device).type_as(next(model.parameters())))  # run once
    dt, seen = [0.0, 0.0, 0.0], 0

    while img_transform is None:
        continue

    try:
        while elapsed_time() < simulationTime:
            start = time.time()

            # t1 = time_sync()
            mutex_map.acquire()
            map_copy = map.copy()
            mutex_map.release()

            if onnx:
                mutex_img.acquire()
                img = img_transform.astype('float32')
                mutex_img.release()
            else:
                mutex_img.acquire()
                img = torch.from_numpy(img_transform).to(device)
                mutex_img.release()
                img = img.half() if half else img.float()  # uint8 to fp16/32
            img = img / 255.0  # 0 - 255 to 0.0 - 1.0
            if len(img.shape) == 3:
                img = img[None]  # expand for batch dim
            # t2 = time_sync()
            # dt[0] += t2 - t1

            # Inference
            if pt:
                # visualize = increment_path(save_dir / Path(path).stem, mkdir=True) if visualize else False
                visualize = False
                pred = model(img, augment=augment, visualize=visualize)[0]
            elif onnx:
                if dnn:
                    net.setInput(img)
                    pred = torch.tensor(net.forward())
                else:
                    pred = torch.tensor(session.run([session.get_outputs()[0].name], {session.get_inputs()[0].name: img}))
            else:  # tensorflow model (tflite, pb, saved_model)
                imn = img.permute(0, 2, 3, 1).cpu().numpy()  # image in numpy
                if pb:
                    pred = frozen_func(x=tf.constant(imn)).numpy()
                elif saved_model:
                    pred = model(imn, training=False).numpy()
                elif tflite:
                    if int8:
                        scale, zero_point = input_details[0]['quantization']
                        imn = (imn / scale + zero_point).astype(np.uint8)  # de-scale
                    interpreter.set_tensor(input_details[0]['index'], imn)
                    interpreter.invoke()
                    pred = interpreter.get_tensor(output_details[0]['index'])
                    if int8:
                        scale, zero_point = output_details[0]['quantization']
                        pred = (pred.astype(np.float32) - zero_point) * scale  # re-scale
                pred[..., 0] *= imgsz[1]  # x
                pred[..., 1] *= imgsz[0]  # y
                pred[..., 2] *= imgsz[1]  # w
                pred[..., 3] *= imgsz[0]  # h
                pred = torch.tensor(pred)
            # t3 = time_sync()
            # dt[1] += t3 - t2
            # print("pred1: ", pred.shape)

            # NMS
            # pred = non_max_suppression(pred, conf_thres, iou_thres, classes, agnostic_nms, max_det=max_det)
            pred = non_max_suppression(pred, conf_thres, iou_thres, classes, agnostic_nms, max_det=max_det)
            # print("pred2: ", pred)
            # dt[2] += time_sync() - t3

            # information of detected cars
            cars_position = process_pred(pred, map_copy)
            print(f"position of cars: {cars_position}")

            det = pred[0]
            s = '%gx%g ' % img.shape[2:]
            for c in det[:, -1].unique():
                n = (det[:, -1] == c).sum()  # detections per class
                s += f"{n} {names[int(c)]}{'s' * (n > 1)}, "  # add to string

            end = time.time()
            computationTime = end - start
            print(f'{s}Done. ({computationTime:.3f}s) (end time: {time.time()})')
            # sleepTime = sampleTime - (computationTime % sampleTime)
            # time.sleep(sleepTime)

            # Second-stage classifier (optional)
            # var classify is False (by default)
            # if classify:
            #     pred = apply_classifier(pred, modelc, img, im0s)

            # Process predictions
            # print(names)
            # for i, det in enumerate(pred):  # per image
            #     seen += 1
            #     if webcam:  # batch_size >= 1
            #         p, s, im0, frame = path[i], f'{i}: ', im0s[i].copy(), dataset.count
            #     else:
            #         p, s, im0, frame = path, '', im0s.copy(), getattr(dataset, 'frame', 0)

            #     p = Path(p)  # to Path
            #     save_path = str(save_dir / p.name)  # img.jpg
            #     txt_path = str(save_dir / 'labels' / p.stem) + ('' if dataset.mode == 'image' else f'_{frame}')  # img.txt
            #     s += '%gx%g ' % img.shape[2:]  # print string
            #     gn = torch.tensor(im0.shape)[[1, 0, 1, 0]]  # normalization gain whwh
            #     imc = im0.copy() if save_crop else im0  # for save_crop
            #     annotator = Annotator(im0, line_width=line_thickness, example=str(names))
            #     # print(det.shape)
            #     if len(det):
            #         # Rescale boxes from img_size to im0 size
            #         det[:, :4] = scale_coords(img.shape[2:], det[:, :4], im0.shape).round()

            #         # Print results
            #         for c in det[:, -1].unique():
            #             n = (det[:, -1] == c).sum()  # detections per class
            #             s += f"{n} {names[int(c)]}{'s' * (n > 1)}, "  # add to string

            #         # Write results
            #         for *xyxy, conf, cls in reversed(det):
            #             if save_txt:  # Write to file
            #                 xywh = (xyxy2xywh(torch.tensor(xyxy).view(1, 4)) / gn).view(-1).tolist()  # normalized xywh
            #                 line = (cls, *xywh, conf) if save_conf else (cls, *xywh)  # label format
            #                 with open(txt_path + '.txt', 'a') as f:
            #                     f.write(('%g ' * len(line)).rstrip() % line + '\n')

            #             if save_img or save_crop or view_img:  # Add bbox to image
            #                 c = int(cls)  # integer class
            #                 label = None if hide_labels else (names[c] if hide_conf else f'{names[c]} {conf:.2f}')
            #                 annotator.box_label(xyxy, label, color=colors(c, True))
            #                 if save_crop:
            #                     save_one_box(xyxy, imc, file=save_dir / 'crops' / names[c] / f'{p.stem}.jpg', BGR=True)

            #     # Print time (inference-only)
            #     print(f'{s}Done. ({t3 - t2:.3f}s)')

            #     # Stream results
            #     im0 = annotator.result()
            #     if view_img:
            #         cv2.imshow(str(p), im0)
            #         cv2.waitKey(1)  # 1 millisecond

            #     # Save results (image with detections)
            #     if save_img:
            #         if dataset.mode == 'image':
            #             cv2.imwrite(save_path, im0)
            #         else:  # 'video' or 'stream'
            #             if vid_path[i] != save_path:  # new video
            #                 vid_path[i] = save_path
            #                 if isinstance(vid_writer[i], cv2.VideoWriter):
            #                     vid_writer[i].release()  # release previous video writer
            #                 if vid_cap:  # video
            #                     fps = vid_cap.get(cv2.CAP_PROP_FPS)
            #                     w = int(vid_cap.get(cv2.CAP_PROP_FRAME_WIDTH))
            #                     h = int(vid_cap.get(cv2.CAP_PROP_FRAME_HEIGHT))
            #                 else:  # stream
            #                     fps, w, h = 30, im0.shape[1], im0.shape[0]
            #                     save_path += '.mp4'
            #                 vid_writer[i] = cv2.VideoWriter(save_path, cv2.VideoWriter_fourcc(*'mp4v'), fps, (w, h))
            #             vid_writer[i].write(im0)

    except KeyboardInterrupt:
        global thread_terminate
        thread_terminate = True
        print("User interrupted!")

    # finally:
    #     Terminate all webcam objects
    #     myCam1.terminate()
    #     myCam2.terminate()
    #     myCam3.terminate()
    #     myCam4.terminate()
    #     Terminate the LIDAR object
    #     myLidar.terminate()

    lidarThread.join()
    cameraThread.join()

    # Print results
    # t = tuple(x / seen * 1E3 for x in dt)  # speeds per image
    # print(f'Speed: %.1fms pre-process, %.1fms inference, %.1fms NMS per image at shape {(1, 3, *imgsz)}' % t)
    # if save_txt or save_img:
    #     s = f"\n{len(list(save_dir.glob('labels/*.txt')))} labels saved to {save_dir / 'labels'}" if save_txt else ''
    #     print(f"Results saved to {colorstr('bold', save_dir)}{s}")
    if update:
        strip_optimizer(weights)  # update model (to fix SourceChangeWarning)


def parse_opt():
    parser = argparse.ArgumentParser()
    parser.add_argument('--weights', nargs='+', type=str, default=ROOT / 'yolov5s.pt', help='model path(s)')
    # parser.add_argument('--source', type=str, default=ROOT / 'data/images', help='file/dir/URL/glob, 0 for webcam')
    parser.add_argument('--imgsz', '--img', '--img-size', nargs='+', type=int, default=[640], help='inference size h,w')
    parser.add_argument('--conf-thres', type=float, default=0.25, help='confidence threshold')
    parser.add_argument('--iou-thres', type=float, default=0.45, help='NMS IoU threshold')
    parser.add_argument('--max-det', type=int, default=1000, help='maximum detections per image')
    parser.add_argument('--device', default='', help='cuda device, i.e. 0 or 0,1,2,3 or cpu')
    parser.add_argument('--view-img', action='store_true', help='show results')
    parser.add_argument('--save-txt', action='store_true', help='save results to *.txt')
    parser.add_argument('--save-conf', action='store_true', help='save confidences in --save-txt labels')
    parser.add_argument('--save-crop', action='store_true', help='save cropped prediction boxes')
    parser.add_argument('--nosave', action='store_true', help='do not save images/videos')
    parser.add_argument('--classes', nargs='+', type=int, help='filter by class: --classes 0, or --classes 0 2 3')
    parser.add_argument('--agnostic-nms', action='store_true', help='class-agnostic NMS')
    parser.add_argument('--augment', action='store_true', help='augmented inference')
    parser.add_argument('--visualize', action='store_true', help='visualize features')
    parser.add_argument('--update', action='store_true', help='update all models')
    parser.add_argument('--project', default=ROOT / 'runs/detect', help='save results to project/name')
    parser.add_argument('--name', default='exp', help='save results to project/name')
    parser.add_argument('--exist-ok', action='store_true', help='existing project/name ok, do not increment')
    parser.add_argument('--line-thickness', default=3, type=int, help='bounding box thickness (pixels)')
    parser.add_argument('--hide-labels', default=False, action='store_true', help='hide labels')
    parser.add_argument('--hide-conf', default=False, action='store_true', help='hide confidences')
    parser.add_argument('--half', action='store_true', help='use FP16 half-precision inference')
    parser.add_argument('--dnn', action='store_true', help='use OpenCV DNN for ONNX inference')
    opt = parser.parse_args()
    opt.imgsz *= 2 if len(opt.imgsz) == 1 else 1  # expand
    print_args(FILE.stem, opt)
    return opt


def main(opt):
    check_requirements(exclude=('tensorboard', 'thop'))
    run(**vars(opt))


if __name__ == "__main__":
    opt = parse_opt()
    main(opt)
