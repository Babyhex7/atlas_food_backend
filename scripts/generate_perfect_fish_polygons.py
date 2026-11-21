import cv2
import numpy as np
import json
import os

downloads_dir = r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads"

polygons = {}

# -------------------------------------------------------------
# 1. LH-34 Gurame (cut image)
# -------------------------------------------------------------
print("Processing LH-34 Gurame...")
img34 = cv2.imread(os.path.join(downloads_dir, "LH-34_cut.jpg"))
h, w, _ = img34.shape
gray = cv2.cvtColor(img34, cv2.COLOR_BGR2GRAY)
hsv = cv2.cvtColor(img34, cv2.COLOR_BGR2HSV)
fish_mask = ((gray < 225) & ((hsv[:,:,1] > 20) | (gray < 140))).astype(np.uint8)
fish_mask[:, int(w*0.72):] = 0
fish_mask[:, :int(w*0.14)] = 0
fish_mask[:int(h*0.20), :] = 0
fish_mask[int(h*0.83):, :] = 0
kernel = cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (5, 5))
clean_fish = cv2.morphologyEx(fish_mask, cv2.MORPH_CLOSE, kernel)

# Piece D (Ekor)
mask_D = clean_fish.copy()
for y in range(h):
    cut1_x = int(355 + (y - 150) * 0.03)
    mask_D[y, cut1_x:] = 0
cnts_D, _ = cv2.findContours(mask_D, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
cnt_D = max(cnts_D, key=cv2.contourArea)
approx_D = cv2.approxPolyDP(cnt_D, 2.5, True)

# Piece B (Kepala)
mask_B = clean_fish.copy()
for y in range(h):
    cut2_x = int(512 + (y - 135) * 0.05)
    mask_B[y, :cut2_x] = 0
cnts_B, _ = cv2.findContours(mask_B, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
cnt_B = max(cnts_B, key=cv2.contourArea)
approx_B = cv2.approxPolyDP(cnt_B, 2.5, True)

# Piece C (Badan Tengah)
mask_C = clean_fish.copy()
for y in range(h):
    cut1_x = int(355 + (y - 150) * 0.03)
    cut2_x = int(512 + (y - 135) * 0.05)
    mask_C[y, :cut1_x] = 0
    mask_C[y, cut2_x:] = 0
cnts_C, _ = cv2.findContours(mask_C, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
cnt_C = max(cnts_C, key=cv2.contourArea)
approx_C = cv2.approxPolyDP(cnt_C, 2.5, True)

# Piece A (Utuh - entire clean fish)
cnts_A, _ = cv2.findContours(clean_fish, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
cnt_A = max(cnts_A, key=cv2.contourArea)
approx_A = cv2.approxPolyDP(cnt_A, 3.0, True)

polygons["LH-34"] = {
    "A": [[int(p[0][0]), int(p[0][1])] for p in approx_A],
    "B": [[int(p[0][0]), int(p[0][1])] for p in approx_B],
    "C": [[int(p[0][0]), int(p[0][1])] for p in approx_C],
    "D": [[int(p[0][0]), int(p[0][1])] for p in approx_D],
}

# -------------------------------------------------------------
# 2. LH-31 Patin (cut image)
# -------------------------------------------------------------
print("Processing LH-31 Patin...")
img31 = cv2.imread(os.path.join(downloads_dir, "LH-31_cut.jpg"))
h, w, _ = img31.shape
gray = cv2.cvtColor(img31, cv2.COLOR_BGR2GRAY)
hsv = cv2.cvtColor(img31, cv2.COLOR_BGR2HSV)
fish_mask = ((gray < 225) & ((hsv[:,:,1] > 20) | (gray < 140))).astype(np.uint8)
fish_mask[:, int(w*0.78):] = 0
fish_mask[:, :int(w*0.24)] = 0
fish_mask[:int(h*0.10), :] = 0
fish_mask[int(h*0.93):, :] = 0
clean_fish = cv2.morphologyEx(fish_mask, cv2.MORPH_CLOSE, kernel)

# Piece D (Ekor)
mask_D = clean_fish.copy()
mask_D[:, 435:] = 0
cnts_D, _ = cv2.findContours(mask_D, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
cnt_D = max(cnts_D, key=cv2.contourArea)
approx_D = cv2.approxPolyDP(cnt_D, 2.5, True)

# Piece B (Kepala)
mask_B = clean_fish.copy()
mask_B[:, :598] = 0
cnts_B, _ = cv2.findContours(mask_B, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
cnt_B = max(cnts_B, key=cv2.contourArea)
approx_B = cv2.approxPolyDP(cnt_B, 2.5, True)

# Piece C (Badan)
mask_C = clean_fish.copy()
mask_C[:, :435] = 0
mask_C[:, 598:] = 0
cnts_C, _ = cv2.findContours(mask_C, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
cnt_C = max(cnts_C, key=cv2.contourArea)
approx_C = cv2.approxPolyDP(cnt_C, 2.5, True)

# Piece A (Utuh)
cnts_A, _ = cv2.findContours(clean_fish, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
cnt_A = max(cnts_A, key=cv2.contourArea)
approx_A = cv2.approxPolyDP(cnt_A, 3.0, True)

polygons["LH-31"] = {
    "A": [[int(p[0][0]), int(p[0][1])] for p in approx_A],
    "B": [[int(p[0][0]), int(p[0][1])] for p in approx_B],
    "C": [[int(p[0][0]), int(p[0][1])] for p in approx_C],
    "D": [[int(p[0][0]), int(p[0][1])] for p in approx_D],
}

# -------------------------------------------------------------
# 3. LH-32 Kakap (4 separate images: A, B, C, D)
# -------------------------------------------------------------
print("Processing LH-32 Kakap...")
polygons["LH-32"] = {}
for porsi in ["A", "B", "C", "D"]:
    img = cv2.imread(os.path.join(downloads_dir, f"LH-32_{porsi}.jpg"))
    h, w, _ = img.shape
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)
    mask = ((gray < 225) & ((hsv[:,:,1] > 20) | (gray < 140))).astype(np.uint8)
    mask[:, int(w*0.80):] = 0
    mask[:, :int(w*0.10)] = 0
    mask[:int(h*0.15), :] = 0
    mask[int(h*0.85):, :] = 0
    cnts, _ = cv2.findContours(mask, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    cnt = max(cnts, key=cv2.contourArea)
    approx = cv2.approxPolyDP(cnt, 2.0, True)
    polygons["LH-32"][porsi] = [[int(p[0][0]), int(p[0][1])] for p in approx]

# -------------------------------------------------------------
# 4. LH-33 Nila (3 separate images: A, B, C)
# -------------------------------------------------------------
print("Processing LH-33 Nila...")
polygons["LH-33"] = {}
for porsi in ["A", "B", "C"]:
    img = cv2.imread(os.path.join(downloads_dir, f"LH-33_{porsi}.jpg"))
    h, w, _ = img.shape
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)
    mask = ((gray < 225) & ((hsv[:,:,1] > 20) | (gray < 140))).astype(np.uint8)
    mask[:, int(w*0.80):] = 0
    mask[:, :int(w*0.10)] = 0
    mask[:int(h*0.05), :] = 0
    mask[int(h*0.95):, :] = 0
    cnts, _ = cv2.findContours(mask, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    cnt = max(cnts, key=cv2.contourArea)
    approx = cv2.approxPolyDP(cnt, 2.0, True)
    polygons["LH-33"][porsi] = [[int(p[0][0]), int(p[0][1])] for p in approx]

# -------------------------------------------------------------
# 5. LH-35 Lele (3 separate images: A, B, C)
# -------------------------------------------------------------
print("Processing LH-35 Lele...")
polygons["LH-35"] = {}
for porsi in ["A", "B", "C"]:
    img = cv2.imread(os.path.join(downloads_dir, f"LH-35_{porsi}.jpg"))
    h, w, _ = img.shape
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)
    mask = ((gray < 225) & ((hsv[:,:,1] > 20) | (gray < 140))).astype(np.uint8)
    mask[:, int(w*0.80):] = 0
    mask[:, :int(w*0.10)] = 0
    mask[:int(h*0.05), :] = 0
    mask[int(h*0.95):, :] = 0
    cnts, _ = cv2.findContours(mask, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    cnt = max(cnts, key=cv2.contourArea)
    approx = cv2.approxPolyDP(cnt, 2.0, True)
    polygons["LH-35"][porsi] = [[int(p[0][0]), int(p[0][1])] for p in approx]

# Save JSON
out_json = os.path.join(downloads_dir, "perfect_fish_polygons.json")
with open(out_json, "w") as f:
    json.dump(polygons, f, indent=2)

print(f"Saved {out_json}")
for code, portions in polygons.items():
    print(f"{code}: {list(portions.keys())}")
    for p, pts in portions.items():
        print(f"  {p}: {len(pts)} pts")
