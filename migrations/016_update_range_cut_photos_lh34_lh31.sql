-- Migration 016: Update LH-34 (Gurame) and LH-31 (Patin) to composite cut interactive photos
-- Keeps full plate / whole fish for thumbnail (A), while viewer uses cut photo for smooth piece-level annotations

START TRANSACTION;

-- 1. LH-34: Ikan Gurame Goreng
-- Hapus food_areas lama untuk LH-34 jika ada
DELETE FROM food_areas WHERE food_image_id IN (
    '73b4010a-3567-42d9-ac26-ff024d7e15b7',
    '9b49dc14-abfb-4bff-86e2-714cf1671136',
    '5558f0a4-8901-4c56-acc9-78734b340e0a',
    'cbd94653-f8c2-46cf-be3a-a69484ca3942',
    'c4e340a1-7341-4567-89ab-lh34cut00001'
);

-- Masukkan food_images tunggal untuk foto potongan LH-34
INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
VALUES (
    'c4e340a1-7341-4567-89ab-lh34cut00001',
    'Ikan Gurame Goreng - Porsi Potongan',
    '/uploads/atlas/LH/LH-34_cut.jpg',
    '/uploads/atlas/LH/LH-34_cut.jpg',
    1064,
    598,
    'published',
    '8ed0845f-ce26-11f1-b499-e2699381f9eb',
    '7d376cd7-d20c-11f1-b499-e2699381f9eb',
    NOW(3),
    NOW(3),
    NOW(3)
) ON DUPLICATE KEY UPDATE
    image_url = VALUES(image_url),
    thumbnail_url = VALUES(thumbnail_url),
    width = VALUES(width),
    height = VALUES(height),
    updated_at = NOW(3);

-- Masukkan area potongan B (Kepala), C (Badan), D (Ekor), dan A (Utuh)
INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at) VALUES
('fa340001-c362-492a-92c2-6153f7245b01', 'c4e340a1-7341-4567-89ab-lh34cut00001', 'Ikan Gurame Porsi B (52g - Kepala)', '8ed0845f-ce26-11f1-b499-e2699381f9eb', '[[634, 135], [573, 134], [480, 163], [480, 412], [493, 432], [594, 341], [658, 317], [679, 209]]', 1, 52.0, NOW(3), NOW(3)),
('fa340002-c362-492a-92c2-6153f7245c02', 'c4e340a1-7341-4567-89ab-lh34cut00001', 'Ikan Gurame Porsi C (95g - Badan Tengah)', '8ed0845f-ce26-11f1-b499-e2699381f9eb', '[[479, 164], [442, 187], [451, 215], [416, 192], [330, 228], [356, 262], [330, 253], [330, 468], [429, 444], [463, 401], [479, 410]]', 2, 95.0, NOW(3), NOW(3)),
('fa340003-c362-492a-92c2-6153f7245d03', 'c4e340a1-7341-4567-89ab-lh34cut00001', 'Ikan Gurame Porsi D (69g - Ekor)', '8ed0845f-ce26-11f1-b499-e2699381f9eb', '[[311, 247], [260, 287], [197, 312], [200, 340], [260, 334], [207, 372], [264, 470], [329, 469], [329, 253]]', 3, 69.0, NOW(3), NOW(3)),
('fa340000-c362-492a-92c2-6153f7245a00', 'c4e340a1-7341-4567-89ab-lh34cut00001', 'Ikan Gurame Porsi A (254g - Utuh)', '8ed0845f-ce26-11f1-b499-e2699381f9eb', '[[638, 139], [454, 169], [451, 215], [385, 199], [325, 230], [355, 263], [311, 247], [200, 309], [200, 340], [260, 334], [207, 371], [254, 459], [311, 473], [429, 444], [467, 400], [502, 432], [594, 341], [650, 327], [679, 209]]', 0, 254.0, NOW(3), NOW(3));

-- Update as_served_images LH-34: Porsi A tetap foto utuh (thumbnail), B, C, D menggunakan foto potongan
UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-34_A.jpg', thumbnail_url = '/uploads/atlas/LH/LH-34_A.jpg', description = 'food_image:c4e340a1-7341-4567-89ab-lh34cut00001', weight_gram = 254.0 WHERE id = '8ed122cf-ce26-11f1-b499-e2699381f9eb';
UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-34_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-34_cut.jpg', description = 'food_image:c4e340a1-7341-4567-89ab-lh34cut00001', weight_gram = 52.0 WHERE id = '8ed1275e-ce26-11f1-b499-e2699381f9eb';
UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-34_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-34_cut.jpg', description = 'food_image:c4e340a1-7341-4567-89ab-lh34cut00001', weight_gram = 95.0 WHERE id = '8ed127e3-ce26-11f1-b499-e2699381f9eb';
UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-34_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-34_cut.jpg', description = 'food_image:c4e340a1-7341-4567-89ab-lh34cut00001', weight_gram = 69.0 WHERE id = '8ed1284b-ce26-11f1-b499-e2699381f9eb';


-- 2. LH-31: Ikan Patin Goreng
-- Hapus food_areas lama untuk LH-31 jika ada
DELETE FROM food_areas WHERE food_image_id IN (
    '9b1240e8-7bff-41dd-a8de-e8958370bb60',
    '79ab7ba8-ba02-4932-88e2-4fad1f621523',
    'd6be433c-cf5f-4991-869b-8bab002fc11d',
    'daed2b1f-3c15-4c17-9530-6a23683343f4',
    'c31000a1-7341-4567-89ab-lh31cut00001'
);

-- Masukkan food_images tunggal untuk foto potongan LH-31
INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
VALUES (
    'c31000a1-7341-4567-89ab-lh31cut00001',
    'Ikan Patin Goreng - Porsi Potongan',
    '/uploads/atlas/LH/LH-31_cut.jpg',
    '/uploads/atlas/LH/LH-31_cut.jpg',
    1063,
    598,
    'published',
    '8ece4ae1-ce26-11f1-b499-e2699381f9eb',
    '7d376cd7-d20c-11f1-b499-e2699381f9eb',
    NOW(3),
    NOW(3),
    NOW(3)
) ON DUPLICATE KEY UPDATE
    image_url = VALUES(image_url),
    thumbnail_url = VALUES(thumbnail_url),
    width = VALUES(width),
    height = VALUES(height),
    updated_at = NOW(3);

-- Masukkan area potongan B (Kepala), C (Badan), D (Ekor), dan A (Utuh)
INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at) VALUES
('fa310001-c362-492a-92c2-6153f7245b01', 'c31000a1-7341-4567-89ab-lh31cut00001', 'Ikan Patin Porsi B (59g - Kepala)', '8ece4ae1-ce26-11f1-b499-e2699381f9eb', '[[275, 110], [315, 209], [381, 209], [388, 197], [385, 209], [512, 209], [478, 208], [327, 74]]', 1, 59.0, NOW(3), NOW(3)),
('fa310002-c362-492a-92c2-6153f7245c02', 'c31000a1-7341-4567-89ab-lh31cut00001', 'Ikan Patin Porsi C (84g - Badan Tengah)', '8ece4ae1-ce26-11f1-b499-e2699381f9eb', '[[316, 210], [296, 279], [339, 243], [375, 291], [383, 344], [602, 344], [540, 219]]', 2, 84.0, NOW(3), NOW(3)),
('fa310003-c362-492a-92c2-6153f7245d03', 'c31000a1-7341-4567-89ab-lh31cut00001', 'Ikan Patin Porsi D (59g - Ekor)', '8ece4ae1-ce26-11f1-b499-e2699381f9eb', '[[383, 345], [394, 380], [456, 376], [442, 401], [452, 442], [481, 456], [599, 442], [610, 493], [670, 420], [607, 345]]', 3, 59.0, NOW(3), NOW(3)),
('fa310000-c362-492a-92c2-6153f7245a00', 'c31000a1-7341-4567-89ab-lh31cut00001', 'Ikan Patin Porsi A (219g - Utuh)', '8ece4ae1-ce26-11f1-b499-e2699381f9eb', '[[283, 103], [324, 220], [296, 279], [339, 243], [388, 375], [456, 376], [452, 442], [599, 442], [616, 489], [670, 418], [582, 318], [560, 240], [478, 208], [327, 74]]', 0, 219.0, NOW(3), NOW(3));

-- Update as_served_images LH-31: Porsi A tetap foto utuh (thumbnail), B, C, D menggunakan foto potongan
UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-31_A.jpg', thumbnail_url = '/uploads/atlas/LH/LH-31_A.jpg', description = 'food_image:c31000a1-7341-4567-89ab-lh31cut00001', weight_gram = 219.0 WHERE id = '8ecec51f-ce26-11f1-b499-e2699381f9eb';
UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-31_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-31_cut.jpg', description = 'food_image:c31000a1-7341-4567-89ab-lh31cut00001', weight_gram = 59.0 WHERE id = '8ecec6fa-ce26-11f1-b499-e2699381f9eb';
UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-31_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-31_cut.jpg', description = 'food_image:c31000a1-7341-4567-89ab-lh31cut00001', weight_gram = 84.0 WHERE id = '8ecec7fe-ce26-11f1-b499-e2699381f9eb';
UPDATE as_served_images SET image_url = '/uploads/atlas/LH/LH-31_cut.jpg', thumbnail_url = '/uploads/atlas/LH/LH-31_cut.jpg', description = 'food_image:c31000a1-7341-4567-89ab-lh31cut00001', weight_gram = 59.0 WHERE id = '8ecf04ab-ce26-11f1-b499-e2699381f9eb';

COMMIT;
