import subprocess
import json
import uuid

def run_sql(query):
    cmd = [
        "ssh", "-o", "StrictHostKeyChecking=no", "root@187.77.112.197",
        "docker exec -i shared_mysql mysql -u atlasuser -p'AtlasUserPassword2026!' db_atlas_food"
    ]
    res = subprocess.run(cmd, input=query, text=True, capture_output=True)
    if res.returncode != 0:
        raise Exception(f"SQL error: {res.stderr}")
    return res.stdout

# Get all unlinked as_served_images
select_query = """
SELECT asi.id, asi.label, asi.weight_gram, asi.image_url, COALESCE(asi.thumbnail_url, asi.image_url) as thumb, s.food_id, f.name
FROM as_served_images asi
JOIN as_served_sets s ON asi.set_id = s.id
JOIN foods f ON s.food_id = f.id
LEFT JOIN food_images fi ON asi.description = CONCAT('food_image:', fi.id)
WHERE fi.id IS NULL;
"""

output = run_sql(select_query)
lines = [l.strip() for l in output.strip().split('\n') if l.strip()]
if len(lines) <= 1:
    print("No unlinked as_served_images found!")
    exit(0)

header = lines[0].split('\t')
rows = []
for l in lines[1:]:
    parts = l.split('\t')
    if len(parts) >= 7:
        rows.append(dict(zip(header, parts)))

print(f"Found {len(rows)} unlinked as_served_images. Generating SQL insert batch...")

batch_size = 100
for i in range(0, len(rows), batch_size):
    batch = rows[i:i+batch_size]
    fi_values = []
    asi_cases = []
    asi_ids = []
    
    for r in batch:
        fi_id = str(uuid.uuid4())
        asi_id = r['id']
        food_id = r['food_id']
        label = r['label'].replace("'", "\\'")
        food_name = r['name'].replace("'", "\\'")
        title = f"{food_name} - Porsi {label} ({float(r['weight_gram']):.1f}g)"
        img_url = r['image_url'].replace("'", "\\'")
        thumb_url = r['thumb'].replace("'", "\\'")
        
        fi_values.append(
            f"('{fi_id}', '{title}', '{img_url}', '{thumb_url}', 1064, 598, 'published', '{food_id}', '7d376cd7-d20c-11f1-b499-e2699381f9eb', NOW(), NOW(), NOW())"
        )
        asi_cases.append(f"WHEN id = '{asi_id}' THEN 'food_image:{fi_id}'")
        asi_ids.append(f"'{asi_id}'")
        
    sql_fi = f"""
    INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
    VALUES {', '.join(fi_values)};
    """
    
    sql_asi = f"""
    UPDATE as_served_images
    SET description = CASE {' '.join(asi_cases)} END
    WHERE id IN ({', '.join(asi_ids)});
    """
    
    combined = sql_fi + "\n" + sql_asi
    run_sql(combined)
    print(f"Synced batch {i} - {i + len(batch)} / {len(rows)}")

print("Done syncing as_served_images to food_images successfully!")
