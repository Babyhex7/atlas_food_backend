import cv2
import numpy as np
import json
import os

downloads_dir = r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads"

def extract_contour(img, mask_func, approx_eps=2.0):
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)
    h, w = gray.shape
    mask = mask_func(gray, hsv, h, w)
    
    # morphological close to fill internal gaps
    kernel = cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (5, 5))
    mask = cv2.morphologyEx(mask, cv2.MORPH_CLOSE, kernel)
    mask = cv2.morphologyEx(mask, cv2.MORPH_OPEN, cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (3, 3)))
    
    contours, _ = cv2.findContours(mask, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    if not contours:
        return []
    # Largest contour
    cnt = max(contours, key=cv2.contourArea)
    # Approx polygon
    epsilon = approx_eps
    approx = cv2.approxPolyDP(cnt, epsilon, True)
    pts = [[int(pt[0][0]), int(pt[0][1])] for pt in approx]
    return pts

print("Script template ready")
