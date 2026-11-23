-- Migration 024: Samakan judul foto Ikan Lele Goreng (LH-35) dengan berat porsi di buku Atlas
-- (A = 26 g, B = 51 g, C = 68 g). Sebelumnya judul A/B masih placeholder 10.0g / 20.0g,
-- padahal polygon (food_areas) dan as_served_images sudah memakai 26 / 51 / 68 g.

START TRANSACTION;

UPDATE food_images SET title = 'Ikan Lele Goreng - Porsi A (26g - Kecil)', updated_at = NOW(3)
WHERE id = '31581dd0-5a99-473d-add6-0f67b9615fc1';

UPDATE food_images SET title = 'Ikan Lele Goreng - Porsi B (51g - Sedang)', updated_at = NOW(3)
WHERE id = 'e2cbb90e-ce43-441d-8d54-eb576be985f3';

UPDATE food_images SET title = 'Ikan Lele Goreng - Porsi C (68g - Besar/Utuh)', updated_at = NOW(3)
WHERE id = '88e2fec8-a77b-4bf6-baae-f694201afd85';

COMMIT;
