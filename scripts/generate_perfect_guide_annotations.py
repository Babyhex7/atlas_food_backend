import cv2
import numpy as np
import os
import json
import math

base_dir = r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\guide_downloads"

def contour_to_poly(cnt, eps=0.003):
    peri = cv2.arcLength(cnt, True)
    approx = cv2.approxPolyDP(cnt, eps * peri, True)
    return [[int(p[0][0]), int(p[0][1])] for p in approx]

guide_data = {}

# =========================================================================
# 1. LH-48: TELUR AYAM RAS CEPLOK (1 area: Porsi A 48.0g)
# =========================================================================
cx, cy = 482.5, 330
rx, ry = 87.5, 70
poly_telur = []
for i in range(24):
    theta = 2 * math.pi * i / 24
    x = int(cx + rx * math.cos(theta))
    y = int(cy + ry * math.sin(theta))
    poly_telur.append([x, y])

guide_data["LH-48"] = {
    "A": {
        "name": "Telur Ayam Ras Ceplok Porsi A (48g)",
        "weight": 48.0,
        "poly": poly_telur
    }
}

# =========================================================================
# 2. LH-12: AYAM GORENG DADA (3 areas: A=57g bawah, B=96g atas, C=149g utuh)
# =========================================================================
img_ayam = cv2.imread(os.path.join(base_dir, "LH-12_guide.jpg"))
hsv_ayam = cv2.cvtColor(img_ayam, cv2.COLOR_BGR2HSV)
mask_ayam = ((hsv_ayam[:,:,1] > 35) & (hsv_ayam[:,:,2] > 60) & (hsv_ayam[:,:,0] < 30)).astype(np.uint8) * 255
mask_ayam[:, 670:] = 0 # exclude spoon

# Cut line strictly between top-right and bottom-right pieces (x: 435 to 590)
cv2.line(mask_ayam, (435, 290), (590, 305), 0, 5)

cnts_ayam, _ = cv2.findContours(mask_ayam, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
pieces_ayam = [c for c in cnts_ayam if cv2.contourArea(c) > 3000]

# Sort: left piece has min cx. Among other two, top has min cy, bottom has max cy
sorted_ayam = sorted(pieces_ayam, key=lambda c: cv2.boundingRect(c)[0])
c_left = sorted_ayam[0] # Left = Dada Utuh (Porsi C: 149g)
right_two = sorted(sorted_ayam[1:], key=lambda c: cv2.boundingRect(c)[1])
c_top_right = right_two[0] # Top Right = Dada Potong Atas (Porsi B: 96g)
c_bot_right = right_two[1] # Bottom Right = Dada Potong Bawah (Porsi A: 57g)

guide_data["LH-12"] = {
    "A": {
        "name": "Ayam Goreng Dada Porsi A (57g - Potongan Bawah)",
        "weight": 57.0,
        "poly": contour_to_poly(c_bot_right, 0.002)
    },
    "B": {
        "name": "Ayam Goreng Dada Porsi B (96g - Potongan Atas)",
        "weight": 96.0,
        "poly": contour_to_poly(c_top_right, 0.002)
    },
    "C": {
        "name": "Ayam Goreng Dada Porsi C (149g - Potong Utuh)",
        "weight": 149.0,
        "poly": contour_to_poly(c_left, 0.002)
    }
}

# =========================================================================
# 3. MP-12: KETUPAT (3 areas: A=13.5g dadu, B=108g belah, C=215g utuh)
# =========================================================================
img_k = cv2.imread(os.path.join(base_dir, "MP-12_guide.jpg"))
gray_k = cv2.cvtColor(img_k, cv2.COLOR_BGR2GRAY)

# 1. Whole ketupat (Top Left, x: 295-515, y: 180-365)
mask_k_whole = np.zeros_like(gray_k)
for y in range(180, 365):
    for x in range(295, 515):
        b, g, r = img_k[y, x]
        if int(g) - int(b) > 10 or int(r) - int(b) > 12 or gray_k[y, x] < 228:
            mask_k_whole[y, x] = 255
cnts, _ = cv2.findContours(cv2.morphologyEx(mask_k_whole, cv2.MORPH_CLOSE, np.ones((7,7), np.uint8)), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
c_k_whole = max(cnts, key=cv2.contourArea)

# 2. Cut-half ketupat (Right, x: 460-645, y: 235-485)
mask_k_half = np.zeros_like(gray_k)
for y in range(235, 485):
    for x in range(460, 645):
        b, g, r = img_k[y, x]
        if int(g) - int(b) > 8 or int(r) - int(b) > 10 or gray_k[y, x] < 234:
            mask_k_half[y, x] = 255
cnts, _ = cv2.findContours(cv2.morphologyEx(mask_k_half, cv2.MORPH_CLOSE, np.ones((7,7), np.uint8)), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
c_k_half = max(cnts, key=cv2.contourArea)

# 3. Cubed pieces (Bottom Left, x: 310-565, y: 365-510)
mask_k_dadu = np.zeros_like(gray_k)
for y in range(365, 510):
    for x in range(310, 565):
        b, g, r = img_k[y, x]
        if gray_k[y, x] < 236 or int(r) - int(b) > 6:
            mask_k_dadu[y, x] = 255
cnts, _ = cv2.findContours(cv2.morphologyEx(mask_k_dadu, cv2.MORPH_CLOSE, np.ones((11,11), np.uint8)), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
c_k_dadu = max(cnts, key=cv2.contourArea)

guide_data["MP-12"] = {
    "A": {
        "name": "Ketupat Porsi A (13.5g - Potongan Dadu)",
        "weight": 13.5,
        "poly": contour_to_poly(c_k_dadu, 0.0025)
    },
    "B": {
        "name": "Ketupat Porsi B (108g - Setengah Belah)",
        "weight": 108.0,
        "poly": contour_to_poly(c_k_half, 0.0025)
    },
    "C": {
        "name": "Ketupat Porsi C (215g - Buah Utuh)",
        "weight": 215.0,
        "poly": contour_to_poly(c_k_whole, 0.0025)
    }
}

# =========================================================================
# 4. AB-01: PEPAYA SEGAR (10 pieces: A-J with exact contours & weights)
# =========================================================================
img_p = cv2.imread(os.path.join(base_dir, "AB-01_guide.jpg"))
hsv_p = cv2.cvtColor(img_p, cv2.COLOR_BGR2HSV)
mask_p = cv2.inRange(hsv_p, np.array([5, 80, 80]), np.array([25, 255, 255]))
mask_p[:, 750:] = 0 # exclude spoon

# 6 pieces in left column:
# We know exact bounding regions:
# F (Kipas Besar): y: 180-235, x: 190-305
# E (Kipas Kecil): y: 235-290, x: 200-295
# D (Segitiga):    y: 290-355, x: 210-280
# C (Koktail Bsr): y: 355-400, x: 220-270
# B (Kotak):       y: 400-440, x: 225-265
# A (Koktail Kcl): y: 440-475, x: 225-265

def get_sub_cnt(mask, y1, y2, x1, x2):
    sub = np.zeros_like(mask)
    sub[y1:y2, x1:x2] = mask[y1:y2, x1:x2]
    cnts, _ = cv2.findContours(cv2.morphologyEx(sub, cv2.MORPH_CLOSE, np.ones((5,5), np.uint8)), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    return max(cnts, key=cv2.contourArea)

cnt_A = get_sub_cnt(mask_p, 440, 480, 220, 265) # 1.7g
cnt_B = get_sub_cnt(mask_p, 400, 440, 220, 265) # 3.3g
cnt_C = get_sub_cnt(mask_p, 355, 400, 220, 270) # 4.7g
cnt_D = get_sub_cnt(mask_p, 290, 355, 210, 280) # 11.7g
cnt_E = get_sub_cnt(mask_p, 235, 290, 200, 295) # 12.3g
cnt_F = get_sub_cnt(mask_p, 180, 235, 190, 305) # 23.7g

cnt_H = get_sub_cnt(mask_p, 150, 315, 335, 485) # 161.7g (ujung belakang atas)
cnt_G = get_sub_cnt(mask_p, 315, 480, 330, 485) # 143g (ujung depan bawah)
cnt_I = get_sub_cnt(mask_p, 155, 455, 505, 585) # 146g (vertikal)
cnt_J = get_sub_cnt(mask_p, 150, 460, 610, 755) # 621g (setengah buah)

guide_data["AB-01"] = {
    "A": {"name": "Pepaya Segar Porsi A (1.7g - Potongan Koktail Kecil)", "weight": 1.7, "poly": contour_to_poly(cnt_A, 0.003)},
    "B": {"name": "Pepaya Segar Porsi B (3.3g - Potongan Kotak/Dadu)", "weight": 3.3, "poly": contour_to_poly(cnt_B, 0.003)},
    "C": {"name": "Pepaya Segar Porsi C (4.7g - Potongan Koktail Besar)", "weight": 4.7, "poly": contour_to_poly(cnt_C, 0.003)},
    "D": {"name": "Pepaya Segar Porsi D (11.7g - Potongan Segitiga)", "weight": 11.7, "poly": contour_to_poly(cnt_D, 0.003)},
    "E": {"name": "Pepaya Segar Porsi E (12.3g - Potongan Kipas Kecil)", "weight": 12.3, "poly": contour_to_poly(cnt_E, 0.003)},
    "F": {"name": "Pepaya Segar Porsi F (23.7g - Potongan Kipas Besar)", "weight": 23.7, "poly": contour_to_poly(cnt_F, 0.003)},
    "G": {"name": "Pepaya Segar Porsi G (143g - ⅛ Ptg Horizontal Ujung Depan)", "weight": 143.0, "poly": contour_to_poly(cnt_G, 0.002)},
    "H": {"name": "Pepaya Segar Porsi H (161.7g - ⅛ Ptg Horizontal Ujung Belakang)", "weight": 161.7, "poly": contour_to_poly(cnt_H, 0.002)},
    "I": {"name": "Pepaya Segar Porsi I (146g - ⅛ Ptg Vertikal)", "weight": 146.0, "poly": contour_to_poly(cnt_I, 0.002)},
    "J": {"name": "Pepaya Segar Porsi J (621g - ½ Buah Utuh Belah)", "weight": 621.0, "poly": contour_to_poly(cnt_J, 0.002)}
}

# Write visualizations
vis_ayam = img_ayam.copy()
for label, col in [("A", (255, 0, 0)), ("B", (0, 0, 255)), ("C", (0, 255, 0))]:
    poly = np.array(guide_data["LH-12"][label]["poly"])
    cv2.drawContours(vis_ayam, [poly], -1, col, 2)
cv2.imwrite(os.path.join(base_dir, "vis_final_ayam.jpg"), vis_ayam)

vis_k = img_k.copy()
for label, col in [("A", (255, 0, 0)), ("B", (0, 0, 255)), ("C", (0, 255, 0))]:
    poly = np.array(guide_data["MP-12"][label]["poly"])
    cv2.drawContours(vis_k, [poly], -1, col, 2)
cv2.imwrite(os.path.join(base_dir, "vis_final_ketupat.jpg"), vis_k)

vis_p = img_p.copy()
for label, item in guide_data["AB-01"].items():
    poly = np.array(item["poly"])
    cv2.drawContours(vis_p, [poly], -1, (0, 255, 0), 2)
    M = cv2.moments(poly)
    if M["m00"] > 0:
        cx, cy = int(M["m10"]/M["m00"]), int(M["m01"]/M["m00"])
        cv2.putText(vis_p, label, (cx-7, cy+5), cv2.FONT_HERSHEY_SIMPLEX, 0.5, (0, 0, 255), 2)
cv2.imwrite(os.path.join(base_dir, "vis_final_pepaya.jpg"), vis_p)

vis_tel = cv2.imread(os.path.join(base_dir, "LH-48_guide.jpg"))
cv2.drawContours(vis_tel, [np.array(guide_data["LH-48"]["A"]["poly"])], -1, (0, 255, 0), 2)
cv2.imwrite(os.path.join(base_dir, "vis_final_telur.jpg"), vis_tel)

with open(os.path.join(base_dir, "perfect_guide_polygons.json"), "w") as f:
    json.dump(guide_data, f, indent=2)

print("SUCCESS: Generated all 4 guide polygons and visualizations!")
