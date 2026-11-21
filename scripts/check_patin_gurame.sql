SELECT f.code, f.name, i.id as img_id, i.label, i.image_url, i.weight_gram, i.description
FROM foods f
JOIN as_served_sets s ON s.food_id = f.id
JOIN as_served_images i ON i.set_id = s.id
WHERE f.code IN ('LH-31', 'LH-34')
ORDER BY f.code, i.label;
