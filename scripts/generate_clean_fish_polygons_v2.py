import cv2
import numpy as np
import json
import os

downloads_dir = r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads"

polygons = {}

# =========================================================================
# 1. LH-34 Gurame (3 BAGIAN: A=Kepala 52g, B=Badan Tengah 95g, C=Ekor 69g)
# =========================================================================
print("Processing LH-34 Gurame (3 pieces)...")
img34 = cv2.imread(os.path.join(downloads_dir, "LH-34_cut.jpg"))
h, w, _ = img34.shape
hsv = cv2.cvtColor(img34, cv2.COLOR_BGR2HSV)
gray = cv2.cvtColor(img34, cv2.COLOR_BGR2GRAY)

food = ((hsv[:,:,1] > 20) & (gray < 215)) | (gray < 150)
food[:, int(w*0.72):] = 0
food[:, :int(w*0.14)] = 0
food[:int(h*0.20), :] = 0
food[int(h*0.74):, :] = 0

clean34 = cv2.morphologyEx(food.astype(np.uint8), cv2.MORPH_CLOSE, cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (7, 7)))

mask_A = clean34.copy()
mask_B = clean34.copy()
mask_C = clean34.copy()

for y in range(h):
    c1 = int(360 + 0.02 * (y - 200))
    c2 = int(515 + 0.05 * (y - 200))
    mask_A[y, :c2] = 0
    mask_B[y, :c1] = 0
    mask_B[y, c2:] = 0
    mask_C[y, c1:] = 0

cnts_A, _ = cv2.findContours(mask_A, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
approx_A = cv2.approxPolyDP(max(cnts_A, key=cv2.contourArea), 2.5, True)

cnts_B, _ = cv2.findContours(mask_B, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
approx_B = cv2.approxPolyDP(max(cnts_B, key=cv2.contourArea), 2.5, True)

cnts_C, _ = cv2.findContours(mask_C, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
approx_C = cv2.approxPolyDP(max(cnts_C, key=cv2.contourArea), 2.5, True)

polygons["LH-34"] = {
    "A": [[int(p[0][0]), int(p[0][1])] for p in approx_A], # Kepala (52g)
    "B": [[int(p[0][0]), int(p[0][1])] for p in approx_B], # Badan (95g)
    "C": [[int(p[0][0]), int(p[0][1])] for p in approx_C], # Ekor (69g)
}

# =========================================================================
# 2. LH-31 Patin (3 BAGIAN: A=Kepala 59g, B=Badan Tengah 84g, C=Ekor 59g)
# =========================================================================
print("Processing LH-31 Patin (3 pieces)...")
img31 = cv2.imread(os.path.join(downloads_dir, "LH-31_cut.jpg"))
h, w, _ = img31.shape
hsv = cv2.cvtColor(img31, cv2.COLOR_BGR2HSV)
gray = cv2.cvtColor(img31, cv2.COLOR_BGR2GRAY)

food = ((hsv[:,:,1] > 18) & (gray < 220)) | (gray < 150)
food[:, int(w*0.80):] = 0
food[:, :int(w*0.22)] = 0
food[:int(h*0.10), :] = 0
food[int(h*0.93):, :] = 0

clean31 = cv2.morphologyEx(food.astype(np.uint8), cv2.MORPH_CLOSE, cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (5, 5)))

Y, X = np.ogrid[:h, :w]
is_kepala = (X + 0.58*Y < 465)
is_ekor = (X - 0.65*Y > 295)
is_badan = (~is_kepala) & (~is_ekor)

cnts_A, _ = cv2.findContours((clean31 & is_kepala).astype(np.uint8), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
approx_A = cv2.approxPolyDP(max(cnts_A, key=cv2.contourArea), 2.5, True)

cnts_B, _ = cv2.findContours((clean31 & is_badan).astype(np.uint8), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
approx_B = cv2.approxPolyDP(max(cnts_B, key=cv2.contourArea), 2.5, True)

cnts_C, _ = cv2.findContours((clean31 & is_ekor).astype(np.uint8), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
approx_C = cv2.approxPolyDP(max(cnts_C, key=cv2.contourArea), 2.5, True)

polygons["LH-31"] = {
    "A": [[int(p[0][0]), int(p[0][1])] for p in approx_A], # Kepala (59g)
    "B": [[int(p[0][0]), int(p[0][1])] for p in approx_B], # Badan (84g)
    "C": [[int(p[0][0]), int(p[0][1])] for p in approx_C], # Ekor (59g)
}

# =========================================================================
# 3. LH-32 Kakap (4 PORSI ASCENDING: A=101g, B=103g, C=111g, D=348g)
# =========================================================================
print("Processing LH-32 Kakap (full fish contours, no plate leakage)...")
polygons["LH-32"] = {}
kakap_mapping = {
    "A": ("LH-32_B.jpg", 101.0),
    "B": ("LH-32_C.jpg", 103.0),
    "C": ("LH-32_D.jpg", 111.0),
    "D": ("LH-32_A.jpg", 348.0),
}
for porsi, (fn, w_gram) in kakap_mapping.items():
    img = cv2.imread(os.path.join(downloads_dir, fn))
    h, w, _ = img.shape
    hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    fish = (((hsv[:,:,1] > 18) & (gray < 235)) | (gray < 165)).astype(np.uint8)
    fish[:, 315:] = 0
    clean = cv2.morphologyEx(fish, cv2.MORPH_CLOSE, cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (7, 7)))
    cnts, _ = cv2.findContours(clean, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    approx = cv2.approxPolyDP(max(cnts, key=cv2.contourArea), 2.0, True)
    polygons["LH-32"][porsi] = [[int(p[0][0]), int(p[0][1])] for p in approx]

# =========================================================================
# 4. LH-33 Nila (3 PORSI ASCENDING: A=105g, B=124g, C=351g)
# =========================================================================
print("Processing LH-33 Nila...")
polygons["LH-33"] = {}
for porsi in ["A", "B", "C"]:
    img = cv2.imread(os.path.join(downloads_dir, f"LH-33_{porsi}.jpg"))
    h, w, _ = img.shape
    hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    fish = (((hsv[:,:,1] > 18) & (gray < 225)) | (gray < 160)).astype(np.uint8)
    fish[:, int(w*0.80):] = 0
    fish[:, :int(w*0.12)] = 0
    clean = cv2.morphologyEx(fish, cv2.MORPH_CLOSE, cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (7, 7)))
    cnts, _ = cv2.findContours(clean, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    approx = cv2.approxPolyDP(max(cnts, key=cv2.contourArea), 2.5, True)
    polygons["LH-33"][porsi] = [[int(p[0][0]), int(p[0][1])] for p in approx]

# =========================================================================
# 5. LH-35 Lele (3 PORSI ASCENDING: A=26g, B=51g, C=68g)
# =========================================================================
print("Processing LH-35 Lele...")
polygons["LH-35"] = {}
for porsi in ["A", "B", "C"]:
    img = cv2.imread(os.path.join(downloads_dir, f"LH-35_{porsi}.jpg"))
    h, w, _ = img.shape
    hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    fish = (((hsv[:,:,1] > 18) & (gray < 225)) | (gray < 160)).astype(np.uint8)
    fish[:, int(w*0.80):] = 0
    fish[:, :int(w*0.12)] = 0
    clean = cv2.morphologyEx(fish, cv2.MORPH_CLOSE, cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (7, 7)))
    cnts, _ = cv2.findContours(clean, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    approx = cv2.approxPolyDP(max(cnts, key=cv2.contourArea), 2.5, True)
    polygons["LH-35"][porsi] = [[int(p[0][0]), int(p[0][1])] for p in approx]

out_json = os.path.join(downloads_dir, "clean_fish_polygons_v2.json")
with open(out_json, "w") as f:
    json.dump(polygons, f, indent=2)

print("Saved clean_fish_polygons_v2.json successfully.")
