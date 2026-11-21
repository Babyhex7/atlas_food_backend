import json

with open(r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads\clean_fish_polygons_v2.json") as f:
    polygons = json.load(f)

sql = []
sql.append("-- Migration 019: Fix fish portions to 3 cut pieces (Gurame, Patin), ascending weights for Kakap, and full fish contours")
sql.append("START TRANSACTION;")
sql.append("")

# 1. LH-34 Gurame (3 BAGIAN: A=Kepala 52g, B=Badan Tengah 95g, C=Ekor 69g)
sql.append("-- 1. LH-34 Ikan Gurame Goreng (3 Bagian Potongan)")
# Re-use areas: fa340000 -> A (Kepala), fa340001 -> B (Badan Tengah), fa340002 -> C (Ekor). Delete fa340003!
sql.append(f"UPDATE food_areas SET polygon = '{json.dumps(polygons['LH-34']['A'])}', name = 'Ikan Gurame Porsi A (52g - Bagian Kepala)', weight_gram = 52.0, updated_at = NOW(3) WHERE id = 'fa340000-c362-492a-92c2-6153f7245a00';")
sql.append(f"UPDATE food_areas SET polygon = '{json.dumps(polygons['LH-34']['B'])}', name = 'Ikan Gurame Porsi B (95g - Bagian Badan Tengah)', weight_gram = 95.0, updated_at = NOW(3) WHERE id = 'fa340001-c362-492a-92c2-6153f7245b01';")
sql.append(f"UPDATE food_areas SET polygon = '{json.dumps(polygons['LH-34']['C'])}', name = 'Ikan Gurame Porsi C (69g - Bagian Ekor)', weight_gram = 69.0, updated_at = NOW(3) WHERE id = 'fa340002-c362-492a-92c2-6153f7245c02';")
sql.append("DELETE FROM food_areas WHERE id = 'fa340003-c362-492a-92c2-6153f7245d03';")

# Re-use as_served_images:
# 8ed122cf -> A (52g), 8ed1275e -> B (95g), 8ed127e3 -> C (69g). Delete 8ed1284b!
sql.append("UPDATE as_served_images SET label = 'A', weight_gram = 52.0, image_url = '/uploads/atlas/LH/LH-34_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-34_cut.jpg', description = 'food_image:c4e340a1-7341-4567-89ab-lh34cut00001', display_order = 1 WHERE id = '8ed122cf-ce26-11f1-b499-e2699381f9eb';")
sql.append("UPDATE as_served_images SET label = 'B', weight_gram = 95.0, image_url = '/uploads/atlas/LH/LH-34_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-34_cut.jpg', description = 'food_image:c4e340a1-7341-4567-89ab-lh34cut00001', display_order = 2 WHERE id = '8ed1275e-ce26-11f1-b499-e2699381f9eb';")
sql.append("UPDATE as_served_images SET label = 'C', weight_gram = 69.0, image_url = '/uploads/atlas/LH/LH-34_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-34_cut.jpg', description = 'food_image:c4e340a1-7341-4567-89ab-lh34cut00001', display_order = 3 WHERE id = '8ed127e3-ce26-11f1-b499-e2699381f9eb';")
sql.append("DELETE FROM as_served_images WHERE id = '8ed1284b-ce26-11f1-b499-e2699381f9eb';")
sql.append("")

# 2. LH-31 Patin (3 BAGIAN: A=Kepala 59g, B=Badan Tengah 84g, C=Ekor 59g)
sql.append("-- 2. LH-31 Ikan Patin Goreng (3 Bagian Potongan)")
sql.append(f"UPDATE food_areas SET polygon = '{json.dumps(polygons['LH-31']['A'])}', name = 'Ikan Patin Porsi A (59g - Bagian Kepala)', weight_gram = 59.0, updated_at = NOW(3) WHERE id = 'fa310000-c362-492a-92c2-6153f7245a00';")
sql.append(f"UPDATE food_areas SET polygon = '{json.dumps(polygons['LH-31']['B'])}', name = 'Ikan Patin Porsi B (84g - Bagian Badan Tengah)', weight_gram = 84.0, updated_at = NOW(3) WHERE id = 'fa310001-c362-492a-92c2-6153f7245b01';")
sql.append(f"UPDATE food_areas SET polygon = '{json.dumps(polygons['LH-31']['C'])}', name = 'Ikan Patin Porsi C (59g - Bagian Ekor)', weight_gram = 59.0, updated_at = NOW(3) WHERE id = 'fa310002-c362-492a-92c2-6153f7245c02';")
sql.append("DELETE FROM food_areas WHERE id = 'fa310003-c362-492a-92c2-6153f7245d03';")

# Re-use as_served_images:
# 8ecec51f -> A (59g), 8ecec6fa -> B (84g), 8ecec7fe -> C (59g). Delete 8ecf04ab!
sql.append("UPDATE as_served_images SET label = 'A', weight_gram = 59.0, image_url = '/uploads/atlas/LH/LH-31_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-31_cut.jpg', description = 'food_image:c31000a1-7341-4567-89ab-lh31cut00001', display_order = 1 WHERE id = '8ecec51f-ce26-11f1-b499-e2699381f9eb';")
sql.append("UPDATE as_served_images SET label = 'B', weight_gram = 84.0, image_url = '/uploads/atlas/LH/LH-31_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-31_cut.jpg', description = 'food_image:c31000a1-7341-4567-89ab-lh31cut00001', display_order = 2 WHERE id = '8ecec6fa-ce26-11f1-b499-e2699381f9eb';")
sql.append("UPDATE as_served_images SET label = 'C', weight_gram = 59.0, image_url = '/uploads/atlas/LH/LH-31_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-31_cut.jpg', description = 'food_image:c31000a1-7341-4567-89ab-lh31cut00001', display_order = 3 WHERE id = '8ecec7fe-ce26-11f1-b499-e2699381f9eb';")
sql.append("DELETE FROM as_served_images WHERE id = '8ecf04ab-ce26-11f1-b499-e2699381f9eb';")
sql.append("")

# 3. LH-32 Kakap (4 PORSI ASCENDING: A=101g, B=103g, C=111g, D=348g)
sql.append("-- 3. LH-32 Ikan Kakap Goreng (Ascending Kecil -> Besar)")
# Food images & areas:
# Porsi A (101g): img LH-32_B.jpg, fi: ed6583aa, area: faf4ab4d, asi: 8ecfaaf9
# Porsi B (103g): img LH-32_C.jpg, fi: 0832d7d9, area: 000e0d72, asi: 8ecfabbc
# Porsi C (111g): img LH-32_D.jpg, fi: 5d1a435d, area: 584ac28a, asi: 8ecfac97
# Porsi D (348g): img LH-32_A.jpg, fi: 6b0c0734, area: 03d50718, asi: 8ecfa900
for porsi, area_id, name, weight, fi_id, asi_id, img_file, d_order in [
    ("A", "faf4ab4d-9206-46df-84d8-3f40f8536671", "Ikan Kakap Porsi A (101g - Kecil)", 101.0, "ed6583aa-c597-468c-8553-71671d540199", "8ecfaaf9-ce26-11f1-b499-e2699381f9eb", "LH-32_B.jpg", 1),
    ("B", "000e0d72-a850-469e-8550-ff1bf2496ac8", "Ikan Kakap Porsi B (103g - Sedang 1)", 103.0, "0832d7d9-72ae-4c0f-a003-986096673a50", "8ecfabbc-ce26-11f1-b499-e2699381f9eb", "LH-32_C.jpg", 2),
    ("C", "584ac28a-eea6-4e21-a4db-d80bd88576e5", "Ikan Kakap Porsi C (111g - Sedang 2)", 111.0, "5d1a435d-44d5-4a29-8685-0427076621d0", "8ecfac97-ce26-11f1-b499-e2699381f9eb", "LH-32_D.jpg", 3),
    ("D", "03d50718-62ed-40af-acb2-a6a46220133a", "Ikan Kakap Porsi D (348g - Besar/Utuh)", 348.0, "6b0c0734-c7e3-49da-9e9b-1836ac8816ca", "8ecfa900-ce26-11f1-b499-e2699381f9eb", "LH-32_A.jpg", 4),
]:
    poly_str = json.dumps(polygons["LH-32"][porsi])
    sql.append(f"UPDATE food_images SET width = 508, height = 286, image_url = '/uploads/atlas/LH/{img_file}', thumbnail_url = '/uploads/atlas/LH/{img_file}', updated_at = NOW(3) WHERE id = '{fi_id}';")
    sql.append(f"UPDATE food_areas SET polygon = '{poly_str}', name = '{name}', weight_gram = {weight}, updated_at = NOW(3) WHERE id = '{area_id}';")
    sql.append(f"UPDATE as_served_images SET label = '{porsi}', weight_gram = {weight}, image_url = '/uploads/atlas/LH/{img_file}', thumbnail_url = '/uploads/atlas/LH/{img_file}', description = 'food_image:{fi_id}', display_order = {d_order} WHERE id = '{asi_id}';")

sql.append("")

# 4. LH-33 Nila (3 PORSI ASCENDING: A=105g, B=124g, C=351g)
sql.append("-- 4. LH-33 Ikan Nila/Mujair Goreng")
for porsi, area_id, name, weight in [
    ("A", "0d9b10ad-f6f6-4011-a5c9-e9d4a141b3ff", "Ikan Nila Porsi A (105g - Kecil)", 105.0),
    ("B", "e291f2f3-b375-4ea3-ad89-241179c7e4f5", "Ikan Nila Porsi B (124g - Sedang)", 124.0),
    ("C", "3eca5c8b-68d1-48df-8665-08791409f3dd", "Ikan Nila Porsi C (351g - Besar/Utuh)", 351.0),
]:
    poly_str = json.dumps(polygons["LH-33"][porsi])
    sql.append(f"UPDATE food_areas SET polygon = '{poly_str}', name = '{name}', weight_gram = {weight}, updated_at = NOW(3) WHERE id = '{area_id}';")

sql.append("")

# 5. LH-35 Lele (3 PORSI ASCENDING: A=26g, B=51g, C=68g)
sql.append("-- 5. LH-35 Ikan Lele Goreng")
for porsi, area_id, name, weight in [
    ("A", "c8b4b751-eb75-4527-8ddf-72d9fd677e87", "Ikan Lele Porsi A (26g - Kecil)", 26.0),
    ("B", "59668e99-8018-4273-97f8-49180c9461e2", "Ikan Lele Porsi B (51g - Sedang)", 51.0),
    ("C", "6cc85c82-3998-48ba-90e0-ff18603625d5", "Ikan Lele Porsi C (68g - Besar/Utuh)", 68.0),
]:
    poly_str = json.dumps(polygons["LH-35"][porsi])
    sql.append(f"UPDATE food_areas SET polygon = '{poly_str}', name = '{name}', weight_gram = {weight}, updated_at = NOW(3) WHERE id = '{area_id}';")

sql.append("")
sql.append("COMMIT;")

out_path = r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\migrations\019_fix_fish_portions_and_order.sql"
with open(out_path, "w", encoding="utf-8") as f:
    f.write("\n".join(sql) + "\n")

print(f"Generated {out_path} ({len(sql)} lines)")
