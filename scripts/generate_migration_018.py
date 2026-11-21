import json

with open(r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads\perfect_fish_polygons.json") as f:
    polygons = json.load(f)

sql = []
sql.append("-- Migration 018: Exact contour-matched annotations for 5 fish")
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

# 3. LH-32 Kakap
sql.append("-- 3. LH-32 Ikan Kakap Goreng")
for porsi, area_id, name, weight in [
    ("A", "03d50718-62ed-40af-acb2-a6a46220133a", "Ikan Kakap Porsi A (348g - Utuh)", 348.0),
    ("B", "faf4ab4d-9206-46df-84d8-3f40f8536671", "Ikan Kakap Porsi B (101g)", 101.0),
    ("C", "000e0d72-a850-469e-8550-ff1bf2496ac8", "Ikan Kakap Porsi C (103g)", 103.0),
    ("D", "584ac28a-eea6-4e21-a4db-d80bd88576e5", "Ikan Kakap Porsi D (111g)", 111.0),
]:
    poly_str = json.dumps(polygons["LH-32"][porsi])
    sql.append(f"UPDATE food_areas SET polygon = '{poly_str}', name = '{name}', weight_gram = {weight}, updated_at = NOW(3) WHERE id = '{area_id}';")

# 4. LH-33 Nila
sql.append("-- 4. LH-33 Ikan Nila/Mujair Goreng")
for porsi, area_id, name, weight in [
    ("A", "0d9b10ad-f6f6-4011-a5c9-e9d4a141b3ff", "Ikan Nila Porsi A (105g)", 105.0),
    ("B", "e291f2f3-b375-4ea3-ad89-241179c7e4f5", "Ikan Nila Porsi B (124g)", 124.0),
    ("C", "3eca5c8b-68d1-48df-8665-08791409f3dd", "Ikan Nila Porsi C (351g - Utuh)", 351.0),
]:
    poly_str = json.dumps(polygons["LH-33"][porsi])
    sql.append(f"UPDATE food_areas SET polygon = '{poly_str}', name = '{name}', weight_gram = {weight}, updated_at = NOW(3) WHERE id = '{area_id}';")

# 5. LH-35 Lele
sql.append("-- 5. LH-35 Ikan Lele Goreng")
for porsi, area_id, name, weight in [
    ("A", "c8b4b751-eb75-4527-8ddf-72d9fd677e87", "Ikan Lele Porsi A (26g)", 26.0),
    ("B", "59668e99-8018-4273-97f8-49180c9461e2", "Ikan Lele Porsi B (51g)", 51.0),
    ("C", "6cc85c82-3998-48ba-90e0-ff18603625d5", "Ikan Lele Porsi C (68g - Utuh)", 68.0),
]:
    poly_str = json.dumps(polygons["LH-35"][porsi])
    sql.append(f"UPDATE food_areas SET polygon = '{poly_str}', name = '{name}', weight_gram = {weight}, updated_at = NOW(3) WHERE id = '{area_id}';")

sql.append("")
sql.append("COMMIT;")

out_path = r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\migrations\018_perfect_fish_annotations.sql"
with open(out_path, "w", encoding="utf-8") as f:
    f.write("\n".join(sql) + "\n")

print(f"Generated {out_path} ({len(sql)} lines)")
