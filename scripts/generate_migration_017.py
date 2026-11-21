import json

with open(r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads\all_fish_polygons.json") as f:
    polygons = json.load(f)

sql = []
sql.append("-- Migration 017: Fix all fish annotations with organic smooth polygons, correct Kakap thumbnails, and precise weights")
sql.append("START TRANSACTION;")
sql.append("")

# 1. LH-34 Gurame
sql.append("-- 1. LH-34 Ikan Gurame Goreng")
for porsi, area_id, name, weight in [
    ("A", "fa340000-c362-492a-92c2-6153f7245a00", "Ikan Gurame Porsi A (254g - Utuh)", 254.0),
    ("B", "fa340001-c362-492a-92c2-6153f7245b01", "Ikan Gurame Porsi B (52g - Kepala)", 52.0),
    ("C", "fa340002-c362-492a-92c2-6153f7245c02", "Ikan Gurame Porsi C (95g - Badan Tengah)", 95.0),
    ("D", "fa340003-c362-492a-92c2-6153f7245d03", "Ikan Gurame Porsi D (69g - Ekor)", 69.0),
]:
    poly_str = json.dumps(polygons["LH-34"][porsi])
    sql.append(f"UPDATE food_areas SET polygon = '{poly_str}', name = '{name}', weight_gram = {weight}, updated_at = NOW(3) WHERE id = '{area_id}';")

sql.append("UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-34_A.jpg', thumbnail_url = '/uploads/atlas/LH/LH-34_A.jpg', description = 'food_image:c4e340a1-7341-4567-89ab-lh34cut00001', weight_gram = 254.0 WHERE id = '8ed122cf-ce26-11f1-b499-e2699381f9eb';")
sql.append("UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-34_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-34_cut.jpg', description = 'food_image:c4e340a1-7341-4567-89ab-lh34cut00001', weight_gram = 52.0 WHERE id = '8ed1275e-ce26-11f1-b499-e2699381f9eb';")
sql.append("UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-34_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-34_cut.jpg', description = 'food_image:c4e340a1-7341-4567-89ab-lh34cut00001', weight_gram = 95.0 WHERE id = '8ed127e3-ce26-11f1-b499-e2699381f9eb';")
sql.append("UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-34_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-34_cut.jpg', description = 'food_image:c4e340a1-7341-4567-89ab-lh34cut00001', weight_gram = 69.0 WHERE id = '8ed1284b-ce26-11f1-b499-e2699381f9eb';")
sql.append("")

# 2. LH-31 Patin
sql.append("-- 2. LH-31 Ikan Patin Goreng")
for porsi, area_id, name, weight in [
    ("A", "fa310000-c362-492a-92c2-6153f7245a00", "Ikan Patin Porsi A (219g - Utuh)", 219.0),
    ("B", "fa310001-c362-492a-92c2-6153f7245b01", "Ikan Patin Porsi B (59g - Kepala)", 59.0),
    ("C", "fa310002-c362-492a-92c2-6153f7245c02", "Ikan Patin Porsi C (84g - Badan Tengah)", 84.0),
    ("D", "fa310003-c362-492a-92c2-6153f7245d03", "Ikan Patin Porsi D (59g - Ekor)", 59.0),
]:
    poly_str = json.dumps(polygons["LH-31"][porsi])
    sql.append(f"UPDATE food_areas SET polygon = '{poly_str}', name = '{name}', weight_gram = {weight}, updated_at = NOW(3) WHERE id = '{area_id}';")

sql.append("UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-31_A.jpg', thumbnail_url = '/uploads/atlas/LH/LH-31_A.jpg', description = 'food_image:c31000a1-7341-4567-89ab-lh31cut00001', weight_gram = 219.0 WHERE id = '8ecec51f-ce26-11f1-b499-e2699381f9eb';")
sql.append("UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-31_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-31_cut.jpg', description = 'food_image:c31000a1-7341-4567-89ab-lh31cut00001', weight_gram = 59.0 WHERE id = '8ecec6fa-ce26-11f1-b499-e2699381f9eb';")
sql.append("UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-31_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-31_cut.jpg', description = 'food_image:c31000a1-7341-4567-89ab-lh31cut00001', weight_gram = 84.0 WHERE id = '8ecec7fe-ce26-11f1-b499-e2699381f9eb';")
sql.append("UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-31_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-31_cut.jpg', description = 'food_image:c31000a1-7341-4567-89ab-lh31cut00001', weight_gram = 59.0 WHERE id = '8ecf04ab-ce26-11f1-b499-e2699381f9eb';")
sql.append("")

# 3. LH-32 Kakap
sql.append("-- 3. LH-32 Ikan Kakap Goreng")
for porsi, area_id, name, weight, fi_id, asi_id in [
    ("A", "03d50718-62ed-40af-acb2-a6a46220133a", "Ikan Kakap Porsi A (348g - Utuh)", 348.0, "6b0c0734-c7e3-49da-9e9b-1836ac8816ca", "8ecfa900-ce26-11f1-b499-e2699381f9eb"),
    ("B", "faf4ab4d-9206-46df-84d8-3f40f8536671", "Ikan Kakap Porsi B (101g)", 101.0, "ed6583aa-c597-468c-8553-71671d540199", "8ecfaaf9-ce26-11f1-b499-e2699381f9eb"),
    ("C", "000e0d72-a850-469e-8550-ff1bf2496ac8", "Ikan Kakap Porsi C (103g)", 103.0, "0832d7d9-72ae-4c0f-a003-986096673a50", "8ecfabbc-ce26-11f1-b499-e2699381f9eb"),
    ("D", "584ac28a-eea6-4e21-a4db-d80bd88576e5", "Ikan Kakap Porsi D (111g)", 111.0, "5d1a435d-44d5-4a29-8685-0427076621d0", "8ecfac97-ce26-11f1-b499-e2699381f9eb"),
]:
    poly_str = json.dumps(polygons["LH-32"][porsi])
    sql.append(f"UPDATE food_images SET width = 508, height = 286, image_url = '/uploads/atlas/LH/LH-32_{porsi}.jpg', thumbnail_url = '/uploads/atlas/LH/LH-32_{porsi}.jpg', updated_at = NOW(3) WHERE id = '{fi_id}';")
    sql.append(f"UPDATE food_areas SET polygon = '{poly_str}', name = '{name}', weight_gram = {weight}, updated_at = NOW(3) WHERE id = '{area_id}';")
    sql.append(f"UPDATE as_served_images SET label = '{porsi}', weight_gram = {weight}, image_url = '/uploads/atlas/LH/LH-32_{porsi}.jpg', thumbnail_url = '/uploads/atlas/LH/LH-32_{porsi}.jpg', description = 'food_image:{fi_id}' WHERE id = '{asi_id}';")

sql.append("")

# 4. LH-33 Nila
sql.append("-- 4. LH-33 Ikan Nila/Mujair Goreng")
for porsi, area_id, name, weight, fi_id, asi_id in [
    ("A", "0d9b10ad-f6f6-4011-a5c9-e9d4a141b3ff", "Ikan Nila Porsi A (105g)", 105.0, "af898243-9abf-474e-84dc-8991a516493b", "8ed054f6-ce26-11f1-b499-e2699381f9eb"),
    ("B", "e291f2f3-b375-4ea3-ad89-241179c7e4f5", "Ikan Nila Porsi B (124g)", 124.0, "074d1ee5-a3d4-4e33-aef9-cb85c174c13c", "8ed0569c-ce26-11f1-b499-e2699381f9eb"),
    ("C", "3eca5c8b-68d1-48df-8665-08791409f3dd", "Ikan Nila Porsi C (351g - Utuh)", 351.0, "a3d1e35b-b5e2-490b-882a-94a9a3fad2df", "8ed0575c-ce26-11f1-b499-e2699381f9eb"),
]:
    poly_str = json.dumps(polygons["LH-33"][porsi])
    sql.append(f"UPDATE food_images SET width = 770, height = 434, image_url = '/uploads/atlas/LH/LH-33_{porsi}.jpg', thumbnail_url = '/uploads/atlas/LH/LH-33_{porsi}.jpg', updated_at = NOW(3) WHERE id = '{fi_id}';")
    sql.append(f"UPDATE food_areas SET polygon = '{poly_str}', name = '{name}', weight_gram = {weight}, updated_at = NOW(3) WHERE id = '{area_id}';")
    sql.append(f"UPDATE as_served_images SET label = '{porsi}', weight_gram = {weight}, image_url = '/uploads/atlas/LH/LH-33_{porsi}.jpg', thumbnail_url = '/uploads/atlas/LH/LH-33_{porsi}.jpg', description = 'food_image:{fi_id}' WHERE id = '{asi_id}';")

sql.append("")

# 5. LH-35 Lele
sql.append("-- 5. LH-35 Ikan Lele Goreng")
for porsi, area_id, name, weight, fi_id, asi_id in [
    ("A", "c8b4b751-eb75-4527-8ddf-72d9fd677e87", "Ikan Lele Porsi A (26g)", 26.0, "31581dd0-5a99-473d-add6-0f67b9615fc1", "8ed1ce9f-ce26-11f1-b499-e2699381f9eb"),
    ("B", "59668e99-8018-4273-97f8-49180c9461e2", "Ikan Lele Porsi B (51g)", 51.0, "e2cbb90e-ce43-441d-8d54-eb576be985f3", "8ed1d041-ce26-11f1-b499-e2699381f9eb"),
    ("C", "6cc85c82-3998-48ba-90e0-ff18603625d5", "Ikan Lele Porsi C (68g - Utuh)", 68.0, "88e2fec8-a77b-4bf6-baae-f694201afd85", "8ed1d0dd-ce26-11f1-b499-e2699381f9eb"),
]:
    poly_str = json.dumps(polygons["LH-35"][porsi])
    sql.append(f"UPDATE food_images SET width = 774, height = 437, image_url = '/uploads/atlas/LH/LH-35_{porsi}.jpg', thumbnail_url = '/uploads/atlas/LH/LH-35_{porsi}.jpg', updated_at = NOW(3) WHERE id = '{fi_id}';")
    sql.append(f"UPDATE food_areas SET polygon = '{poly_str}', name = '{name}', weight_gram = {weight}, updated_at = NOW(3) WHERE id = '{area_id}';")
    sql.append(f"UPDATE as_served_images SET label = '{porsi}', weight_gram = {weight}, image_url = '/uploads/atlas/LH/LH-35_{porsi}.jpg', thumbnail_url = '/uploads/atlas/LH/LH-35_{porsi}.jpg', description = 'food_image:{fi_id}' WHERE id = '{asi_id}';")

sql.append("")
sql.append("COMMIT;")

out_path = r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\migrations\017_fix_all_fish_annotations.sql"
with open(out_path, "w", encoding="utf-8") as f:
    f.write("\n".join(sql) + "\n")

print(f"Generated {out_path} ({len(sql)} lines)")
