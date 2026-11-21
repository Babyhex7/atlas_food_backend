import json
import cv2
import os
import numpy as np

downloads_dir = r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads"

# =========================================================================
# 1. Patin (LH-31)
# =========================================================================
img31 = cv2.imread(os.path.join(downloads_dir, "LH-31_cut.jpg"))
h31, w31, _ = img31.shape
gray31 = cv2.cvtColor(img31, cv2.COLOR_BGR2GRAY)
mask31 = (gray31 < 228).astype(np.uint8)
mask31[:, 800:] = 0
mask31[:, :80] = 0
mask31[:30, :] = 0
mask31[570:, :] = 0
kernel31 = cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (15, 15))
closed31 = cv2.morphologyEx(mask31, cv2.MORPH_CLOSE, kernel31)

Y31, X31 = np.ogrid[:h31, :w31]
is_kepala31 = (6*X31 + 7*Y31 < 3290)
is_ekor31 = (2*X31 - Y31 > 650)
is_badan31 = (~is_kepala31) & (~is_ekor31)

approx_patin_A = cv2.approxPolyDP(max(cv2.findContours((closed31 & is_kepala31).astype(np.uint8), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)[0], key=cv2.contourArea), 2.5, True)
approx_patin_B = cv2.approxPolyDP(max(cv2.findContours((closed31 & is_badan31).astype(np.uint8), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)[0], key=cv2.contourArea), 2.5, True)
approx_patin_C = cv2.approxPolyDP(max(cv2.findContours((closed31 & is_ekor31).astype(np.uint8), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)[0], key=cv2.contourArea), 2.5, True)

poly_patin = {
    "A": [[int(p[0][0]), int(p[0][1])] for p in approx_patin_A],
    "B": [[int(p[0][0]), int(p[0][1])] for p in approx_patin_B],
    "C": [[int(p[0][0]), int(p[0][1])] for p in approx_patin_C],
}

# =========================================================================
# 2. Gurame (LH-34)
# =========================================================================
img34 = cv2.imread(os.path.join(downloads_dir, "LH-34_cut.jpg"))
h34, w34, _ = img34.shape
gray34 = cv2.cvtColor(img34, cv2.COLOR_BGR2GRAY)
hsv34 = cv2.cvtColor(img34, cv2.COLOR_BGR2HSV)
food34 = ((hsv34[:,:,1] > 20) & (gray34 < 215)) | (gray34 < 150)
food34[:, int(w34*0.72):] = 0
food34[:, :int(w34*0.14)] = 0
food34[:int(h34*0.20), :] = 0
food34[int(h34*0.74):, :] = 0
clean34 = cv2.morphologyEx(food34.astype(np.uint8), cv2.MORPH_CLOSE, cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (7, 7)))

mask_A34 = clean34.copy()
mask_B34 = clean34.copy()
mask_C34 = clean34.copy()
for y in range(h34):
    c1 = int(360 + 0.02 * (y - 200))
    c2 = int(515 + 0.05 * (y - 200))
    mask_A34[y, :c2] = 0
    mask_B34[y, :c1] = 0
    mask_B34[y, c2:] = 0
    mask_C34[y, c1:] = 0

approx_gurame_A = cv2.approxPolyDP(max(cv2.findContours(mask_A34, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)[0], key=cv2.contourArea), 2.5, True)
approx_gurame_B = cv2.approxPolyDP(max(cv2.findContours(mask_B34, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)[0], key=cv2.contourArea), 2.5, True)
approx_gurame_C = cv2.approxPolyDP(max(cv2.findContours(mask_C34, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)[0], key=cv2.contourArea), 2.5, True)

poly_gurame = {
    "A": [[int(p[0][0]), int(p[0][1])] for p in approx_gurame_A],
    "B": [[int(p[0][0]), int(p[0][1])] for p in approx_gurame_B],
    "C": [[int(p[0][0]), int(p[0][1])] for p in approx_gurame_C],
}

# =========================================================================
# 3. Kakap (LH-32) - 4 images (A: Utuh 348g, B: Kepala 101g, C: Tengah 103g, D: Ekor 111g)
# =========================================================================
kakap_meta = {
    "A": ("LH-32_A.jpg", 348.0, "Ikan Kakap Porsi A (348g - Utuh)", "03d50718-62ed-40af-acb2-a6a46220133a", "6b0c0734-c7e3-49da-9e9b-1836ac8816ca", "8ecfa900-ce26-11f1-b499-e2699381f9eb", 1),
    "B": ("LH-32_B.jpg", 101.0, "Ikan Kakap Porsi B (101g - Bagian Kepala)", "faf4ab4d-9206-46df-84d8-3f40f8536671", "ed6583aa-c597-468c-8553-71671d540199", "8ecfaaf9-ce26-11f1-b499-e2699381f9eb", 2),
    "C": ("LH-32_C.jpg", 103.0, "Ikan Kakap Porsi C (103g - Bagian Tengah)", "000e0d72-a850-469e-8550-ff1bf2496ac8", "0832d7d9-72ae-4c0f-a003-986096673a50", "8ecfabbc-ce26-11f1-b499-e2699381f9eb", 3),
    "D": ("LH-32_D.jpg", 111.0, "Ikan Kakap Porsi D (111g - Bagian Ekor)", "584ac28a-eea6-4e21-a4db-d80bd88576e5", "5d1a435d-44d5-4a29-8685-0427076621d0", "8ecfac97-ce26-11f1-b499-e2699381f9eb", 4),
}

poly_kakap = {}
for porsi, (fn, w, name, area_id, fi_id, asi_id, d_order) in kakap_meta.items():
    img = cv2.imread(os.path.join(downloads_dir, fn))
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    hsv = cv2.cvtColor(img, cv2.COLOR_BGR2HSV)
    fish = (((hsv[:,:,1] > 18) & (gray < 235)) | (gray < 165)).astype(np.uint8)
    fish[:, 315:] = 0
    clean = cv2.morphologyEx(fish, cv2.MORPH_CLOSE, cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (7, 7)))
    approx = cv2.approxPolyDP(max(cv2.findContours(clean, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)[0], key=cv2.contourArea), 2.0, True)
    poly_kakap[porsi] = [[int(p[0][0]), int(p[0][1])] for p in approx]

# =========================================================================
# Build SQL Migration 020
# =========================================================================
sql = []
sql.append("-- Migration 020: Accurate contours & correct book ordering for Patin, Gurame, and Kakap")
sql.append("START TRANSACTION;")
sql.append("")

# 1. Patin
sql.append("-- 1. LH-31 Patin (3 Bagian Potongan)")
sql.append(f"UPDATE food_areas SET polygon = '{json.dumps(poly_patin['A'])}', name = 'Ikan Patin Porsi A (59g - Bagian Kepala)', weight_gram = 59.0, updated_at = NOW(3) WHERE id = 'fa310000-c362-492a-92c2-6153f7245a00';")
sql.append(f"UPDATE food_areas SET polygon = '{json.dumps(poly_patin['B'])}', name = 'Ikan Patin Porsi B (84g - Bagian Badan Tengah)', weight_gram = 84.0, updated_at = NOW(3) WHERE id = 'fa310001-c362-492a-92c2-6153f7245b01';")
sql.append(f"UPDATE food_areas SET polygon = '{json.dumps(poly_patin['C'])}', name = 'Ikan Patin Porsi C (59g - Bagian Ekor)', weight_gram = 59.0, updated_at = NOW(3) WHERE id = 'fa310002-c362-492a-92c2-6153f7245c02';")
sql.append("DELETE FROM food_areas WHERE id = 'fa310003-c362-492a-92c2-6153f7245d03';")
sql.append("UPDATE as_served_images SET label = 'A', weight_gram = 59.0, image_url = '/uploads/atlas/LH/LH-31_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-31_cut.jpg', description = 'food_image:c31000a1-7341-4567-89ab-lh31cut00001', display_order = 1 WHERE id = '8ecec51f-ce26-11f1-b499-e2699381f9eb';")
sql.append("UPDATE as_served_images SET label = 'B', weight_gram = 84.0, image_url = '/uploads/atlas/LH/LH-31_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-31_cut.jpg', description = 'food_image:c31000a1-7341-4567-89ab-lh31cut00001', display_order = 2 WHERE id = '8ecec6fa-ce26-11f1-b499-e2699381f9eb';")
sql.append("UPDATE as_served_images SET label = 'C', weight_gram = 59.0, image_url = '/uploads/atlas/LH/LH-31_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-31_cut.jpg', description = 'food_image:c31000a1-7341-4567-89ab-lh31cut00001', display_order = 3 WHERE id = '8ecec7fe-ce26-11f1-b499-e2699381f9eb';")
sql.append("DELETE FROM as_served_images WHERE id = '8ecf04ab-ce26-11f1-b499-e2699381f9eb';")
sql.append("")

# 2. Gurame
sql.append("-- 2. LH-34 Gurame (3 Bagian Potongan)")
sql.append(f"UPDATE food_areas SET polygon = '{json.dumps(poly_gurame['A'])}', name = 'Ikan Gurame Porsi A (52g - Bagian Kepala)', weight_gram = 52.0, updated_at = NOW(3) WHERE id = 'fa340000-c362-492a-92c2-6153f7245a00';")
sql.append(f"UPDATE food_areas SET polygon = '{json.dumps(poly_gurame['B'])}', name = 'Ikan Gurame Porsi B (95g - Bagian Badan Tengah)', weight_gram = 95.0, updated_at = NOW(3) WHERE id = 'fa340001-c362-492a-92c2-6153f7245b01';")
sql.append(f"UPDATE food_areas SET polygon = '{json.dumps(poly_gurame['C'])}', name = 'Ikan Gurame Porsi C (69g - Bagian Ekor)', weight_gram = 69.0, updated_at = NOW(3) WHERE id = 'fa340002-c362-492a-92c2-6153f7245c02';")
sql.append("DELETE FROM food_areas WHERE id = 'fa340003-c362-492a-92c2-6153f7245d03';")
sql.append("UPDATE as_served_images SET label = 'A', weight_gram = 52.0, image_url = '/uploads/atlas/LH/LH-34_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-34_cut.jpg', description = 'food_image:c4e340a1-7341-4567-89ab-lh34cut00001', display_order = 1 WHERE id = '8ed122cf-ce26-11f1-b499-e2699381f9eb';")
sql.append("UPDATE as_served_images SET label = 'B', weight_gram = 95.0, image_url = '/uploads/atlas/LH/LH-34_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-34_cut.jpg', description = 'food_image:c4e340a1-7341-4567-89ab-lh34cut00001', display_order = 2 WHERE id = '8ed1275e-ce26-11f1-b499-e2699381f9eb';")
sql.append("UPDATE as_served_images SET label = 'C', weight_gram = 69.0, image_url = '/uploads/atlas/LH/LH-34_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-34_cut.jpg', description = 'food_image:c4e340a1-7341-4567-89ab-lh34cut00001', display_order = 3 WHERE id = '8ed127e3-ce26-11f1-b499-e2699381f9eb';")
sql.append("DELETE FROM as_served_images WHERE id = '8ed1284b-ce26-11f1-b499-e2699381f9eb';")
sql.append("")

# 3. Kakap (Book order: A=348g Utuh, B=101g Kepala, C=103g Tengah, D=111g Ekor)
sql.append("-- 3. LH-32 Kakap (A=348g Utuh, B=101g Kepala, C=103g Tengah, D=111g Ekor)")
for porsi, (fn, w, name, area_id, fi_id, asi_id, d_order) in kakap_meta.items():
    poly_str = json.dumps(poly_kakap[porsi])
    sql.append(f"UPDATE food_images SET width = 508, height = 286, image_url = '/uploads/atlas/LH/{fn}', thumbnail_url = '/uploads/atlas/LH/{fn}', updated_at = NOW(3) WHERE id = '{fi_id}';")
    sql.append(f"UPDATE food_areas SET food_image_id = '{fi_id}', polygon = '{poly_str}', name = '{name}', weight_gram = {w}, updated_at = NOW(3) WHERE id = '{area_id}';")
    sql.append(f"UPDATE as_served_images SET label = '{porsi}', weight_gram = {w}, image_url = '/uploads/atlas/LH/{fn}', thumbnail_url = '/uploads/atlas/LH/{fn}', description = 'food_image:{fi_id}', display_order = {d_order} WHERE id = '{asi_id}';")

sql.append("")
sql.append("COMMIT;")

out_sql = r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\migrations\020_accurate_fish_portions.sql"
with open(out_sql, "w", encoding="utf-8") as f:
    f.write("\n".join(sql) + "\n")

print(f"Generated {out_sql} ({len(sql)} lines)")
