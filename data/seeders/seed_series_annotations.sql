DELETE FROM food_images WHERE image_url = '/uploads/atlas/LN/LN-01_A.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('7ac26b60-6c19-4651-8383-996e7c58b83d', 'Kacang Tanah Goreng - Porsi A (5g)', '/uploads/atlas/LN/LN-01_A.jpg', '/uploads/atlas/LN/LN-01_A.jpg', 1600, 900, 'published', '47bc9580-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('c45a3f0f-c66e-474a-92b4-b91fdc2d8ac1', '7ac26b60-6c19-4651-8383-996e7c58b83d', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('6cbf8e0f-0024-44e5-a6f4-98a3f281ebe2', '7ac26b60-6c19-4651-8383-996e7c58b83d', 'Kacang Tanah Goreng Porsi A', '47bc9580-cfaf-11f1-ab66-664527f556d1', '[[934,422],[942,410],[970,402],[978,402],[994,414],[1002,426],[1006,462],[998,534],[994,542],[974,566],[962,578],[938,578],[934,574]]', 1, 5, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/LN/LN-01_B.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('667185a2-cd9d-4877-9f0e-010178da1f60', 'Kacang Tanah Goreng - Porsi B (10g)', '/uploads/atlas/LN/LN-01_B.jpg', '/uploads/atlas/LN/LN-01_B.jpg', 1600, 900, 'published', '47bc9580-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('04be3bd6-0cf1-40be-93c5-f9f36078dd0d', '667185a2-cd9d-4877-9f0e-010178da1f60', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('eed301ce-494f-48d1-9efb-a03769efba94', '667185a2-cd9d-4877-9f0e-010178da1f60', 'Kacang Tanah Goreng Porsi B', '47bc9580-cfaf-11f1-ab66-664527f556d1', '[[902,574],[906,414],[910,406],[922,386],[934,378],[966,378],[978,386],[982,394],[1006,458],[1010,498],[994,554],[974,578],[970,582],[942,602],[934,606],[910,606],[902,602]]', 1, 10, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/LN/LN-01_C.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('5a90100a-2663-4702-ad90-2aba2e838056', 'Kacang Tanah Goreng - Porsi C (30g)', '/uploads/atlas/LN/LN-01_C.jpg', '/uploads/atlas/LN/LN-01_C.jpg', 1600, 900, 'published', '47bc9580-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('4392e94e-0399-46f1-af44-edffb62c8f87', '5a90100a-2663-4702-ad90-2aba2e838056', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('1c13ea2e-cee6-448c-aa5c-89b8134cf985', '5a90100a-2663-4702-ad90-2aba2e838056', 'Kacang Tanah Goreng Porsi C', '47bc9580-cfaf-11f1-ab66-664527f556d1', '[[810,498],[814,430],[830,338],[834,334],[842,330],[854,326],[902,334],[934,342],[954,354],[970,378],[1006,458],[1002,506],[998,526],[974,578],[970,582],[926,622],[914,630],[874,650],[846,658],[830,658],[818,654],[810,618]]', 1, 30, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/LN/LN-01_D.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('40f2b699-e868-42a4-803e-d662c3165df3', 'Kacang Tanah Goreng - Porsi D (50g)', '/uploads/atlas/LN/LN-01_D.jpg', '/uploads/atlas/LN/LN-01_D.jpg', 1600, 900, 'published', '47bc9580-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('eac9db75-9824-4701-87ed-109be5db7773', '40f2b699-e868-42a4-803e-d662c3165df3', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('5126e83d-0125-497f-9dea-2f454f12682d', '40f2b699-e868-42a4-803e-d662c3165df3', 'Kacang Tanah Goreng Porsi D', '47bc9580-cfaf-11f1-ab66-664527f556d1', '[[730,558],[738,402],[742,326],[746,314],[754,306],[778,302],[794,302],[846,306],[882,322],[898,330],[946,366],[966,382],[994,414],[1002,426],[1010,482],[1010,490],[1006,518],[986,566],[942,610],[918,626],[874,650],[858,658],[834,666],[778,678],[754,678],[746,674],[742,662],[730,562]]', 1, 50, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/LN/LN-01_E.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('bf37b94d-534f-4cb6-b661-237f363cb350', 'Kacang Tanah Goreng - Porsi E (80g)', '/uploads/atlas/LN/LN-01_E.jpg', '/uploads/atlas/LN/LN-01_E.jpg', 1600, 900, 'published', '47bc9580-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('cc95aa62-dd04-4001-a8dd-f2eeb91303d9', 'bf37b94d-534f-4cb6-b661-237f363cb350', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('ce6895ab-0d6b-4e2e-8655-a4d17132a5f5', 'bf37b94d-534f-4cb6-b661-237f363cb350', 'Kacang Tanah Goreng Porsi E', '47bc9580-cfaf-11f1-ab66-664527f556d1', '[[634,654],[638,466],[646,370],[650,342],[670,306],[682,302],[762,286],[806,294],[894,314],[902,318],[954,350],[970,362],[994,398],[998,406],[1018,458],[1018,478],[1010,518],[1006,534],[974,586],[938,614],[894,646],[846,666],[834,670],[790,678],[750,682],[698,682],[670,678],[638,666]]', 1, 80, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/LN/LN-01_F.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('a3e13176-bb11-4b3d-a277-77435a1c0b05', 'Kacang Tanah Goreng - Porsi F (100g)', '/uploads/atlas/LN/LN-01_F.jpg', '/uploads/atlas/LN/LN-01_F.jpg', 1600, 900, 'published', '47bc9580-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('f79d2e5b-544c-4cd4-a640-ff61bb00bbb9', 'a3e13176-bb11-4b3d-a277-77435a1c0b05', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('70f2716b-0251-4b5c-b841-1cb65685cc0a', 'a3e13176-bb11-4b3d-a277-77435a1c0b05', 'Kacang Tanah Goreng Porsi F', '47bc9580-cfaf-11f1-ab66-664527f556d1', '[[538,542],[562,422],[566,406],[594,338],[598,330],[630,310],[734,290],[790,290],[834,302],[882,318],[934,342],[942,346],[970,370],[974,374],[994,410],[1006,446],[1010,470],[1010,498],[1002,534],[986,562],[962,598],[958,602],[930,626],[894,646],[858,662],[846,666],[818,674],[794,678],[738,686],[722,686],[670,682],[642,678],[594,662],[586,658],[566,646],[558,638],[538,546]]', 1, 100, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/LN/LN-01_G.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('3eab4072-05ea-4608-a444-c8d21dad348b', 'Kacang Tanah Goreng - Porsi G (110g)', '/uploads/atlas/LN/LN-01_G.jpg', '/uploads/atlas/LN/LN-01_G.jpg', 1600, 900, 'published', '47bc9580-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('b493c0fb-2064-426e-b508-afda3faf58f4', '3eab4072-05ea-4608-a444-c8d21dad348b', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('db5c4b4f-1f12-4b63-ad2b-38116653cb0d', '3eab4072-05ea-4608-a444-c8d21dad348b', 'Kacang Tanah Goreng Porsi G', '47bc9580-cfaf-11f1-ab66-664527f556d1', '[[522,594],[546,402],[570,346],[574,342],[590,330],[610,318],[634,306],[646,302],[726,290],[766,290],[914,326],[926,334],[982,390],[1010,442],[1018,462],[1018,474],[1014,506],[1010,518],[990,558],[974,582],[954,610],[934,626],[914,638],[882,654],[838,670],[802,678],[750,686],[710,686],[674,682],[606,666],[546,630],[538,622],[522,598]]', 1, 110, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/LN/LN-01_H.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('2218355d-223a-488f-9eb1-02584455128d', 'Kacang Tanah Goreng - Porsi H (140g)', '/uploads/atlas/LN/LN-01_H.jpg', '/uploads/atlas/LN/LN-01_H.jpg', 1600, 900, 'published', '47bc9580-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('71ec9ba1-93e4-4c5b-858c-d614528167ee', '2218355d-223a-488f-9eb1-02584455128d', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('da631a4c-9eb9-4fd8-b803-43533d6a3f80', '2218355d-223a-488f-9eb1-02584455128d', 'Kacang Tanah Goreng Porsi H', '47bc9580-cfaf-11f1-ab66-664527f556d1', '[[466,486],[474,438],[482,426],[514,382],[530,366],[554,350],[586,334],[642,310],[710,298],[806,290],[818,294],[862,314],[894,330],[958,378],[970,390],[994,426],[1006,450],[1010,462],[998,542],[978,582],[962,602],[942,622],[902,646],[894,650],[846,670],[834,674],[774,686],[678,686],[614,670],[582,658],[558,646],[546,638],[526,622],[506,602],[490,578],[470,534],[466,502]]', 1, 140, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-02_A.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('64130d83-dff8-473a-b202-27638fc14bdb', 'Buncis Rebus - Porsi A (10g)', '/uploads/atlas/AS/AS-02_A.jpg', '/uploads/atlas/AS/AS-02_A.jpg', 1600, 900, 'published', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('08561e85-2555-4af5-a495-ad385102bbc4', '64130d83-dff8-473a-b202-27638fc14bdb', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('2d202ea8-8773-45be-a616-92a53b304ba1', '64130d83-dff8-473a-b202-27638fc14bdb', 'Buncis Rebus Porsi A', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '[[890,442],[918,382],[934,382],[958,386],[962,390],[974,486],[974,518],[966,582],[958,590],[950,594],[930,594],[918,590],[890,546]]', 1, 10, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-02_B.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('7278647b-f53f-4a20-a093-6c613991cc13', 'Buncis Rebus - Porsi B (30g)', '/uploads/atlas/AS/AS-02_B.jpg', '/uploads/atlas/AS/AS-02_B.jpg', 1600, 900, 'published', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('1977b9b9-f4c7-4119-a04f-24457b2d71c0', '7278647b-f53f-4a20-a093-6c613991cc13', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('2cf27b20-0b75-4b7c-8961-8c0a623d3013', '7278647b-f53f-4a20-a093-6c613991cc13', 'Buncis Rebus Porsi B', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '[[814,634],[826,382],[842,330],[846,322],[866,310],[938,334],[950,350],[966,374],[986,406],[994,426],[998,438],[998,462],[994,486],[986,530],[974,562],[942,602],[918,626],[910,630],[850,650],[842,650],[826,646],[814,642]]', 1, 30, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-02_C.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('78fc8003-7b57-4cd0-a307-1543cd6ff275', 'Buncis Rebus - Porsi C (40g)', '/uploads/atlas/AS/AS-02_C.jpg', '/uploads/atlas/AS/AS-02_C.jpg', 1600, 900, 'published', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('f45bedc2-c8dd-4d0f-a220-6336654bce33', '78fc8003-7b57-4cd0-a307-1543cd6ff275', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('313ac499-eb77-41ee-8c29-5e1cd875c788', '78fc8003-7b57-4cd0-a307-1543cd6ff275', 'Buncis Rebus Porsi C', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '[[790,626],[794,350],[806,318],[818,314],[842,310],[898,322],[906,326],[926,338],[946,354],[982,386],[990,394],[1010,486],[1010,494],[1006,522],[1002,538],[994,554],[974,574],[922,618],[910,626],[882,638],[850,650],[830,654],[822,654],[810,650],[794,638],[790,634]]', 1, 40, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-02_D.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('a9e3973c-ca83-42ae-8b5a-950ffef74636', 'Buncis Rebus - Porsi D (50g)', '/uploads/atlas/AS/AS-02_D.jpg', '/uploads/atlas/AS/AS-02_D.jpg', 1600, 900, 'published', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('20816932-cca2-4e31-902c-28f85e8f3bfa', 'a9e3973c-ca83-42ae-8b5a-950ffef74636', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('c380625d-75cb-47ed-b010-31f98f554112', 'a9e3973c-ca83-42ae-8b5a-950ffef74636', 'Buncis Rebus Porsi D', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '[[786,622],[790,370],[798,326],[802,314],[806,306],[866,274],[874,278],[970,338],[978,350],[994,386],[1006,422],[1014,482],[1010,518],[974,574],[966,586],[950,606],[946,610],[894,646],[838,650],[794,650],[786,646]]', 1, 50, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-02_E.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('c5e1d333-4f83-4755-9b42-db074da1c83b', 'Buncis Rebus - Porsi E (60g)', '/uploads/atlas/AS/AS-02_E.jpg', '/uploads/atlas/AS/AS-02_E.jpg', 1600, 900, 'published', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('c7e3b634-5d91-48d6-8143-f707b8e1a25d', 'c5e1d333-4f83-4755-9b42-db074da1c83b', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('b761f2cc-56f6-4e2e-a3ee-2478b0057b02', 'c5e1d333-4f83-4755-9b42-db074da1c83b', 'Buncis Rebus Porsi E', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '[[730,614],[750,334],[786,294],[790,290],[818,282],[846,282],[858,286],[954,322],[958,326],[1018,438],[1018,466],[1002,514],[990,542],[986,546],[894,634],[882,638],[830,654],[814,658],[734,630],[730,626]]', 1, 60, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-02_F.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('5348acdc-2dc4-497d-adc2-38bf554a97cb', 'Buncis Rebus - Porsi F (80g)', '/uploads/atlas/AS/AS-02_F.jpg', '/uploads/atlas/AS/AS-02_F.jpg', 1600, 900, 'published', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('bdacc901-f6e1-426e-9a87-5bddc1855afd', '5348acdc-2dc4-497d-adc2-38bf554a97cb', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('893f6496-e2b5-4d23-b355-aa6cde418cd3', '5348acdc-2dc4-497d-adc2-38bf554a97cb', 'Buncis Rebus Porsi F', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '[[658,554],[666,370],[670,346],[750,274],[758,274],[810,282],[874,306],[978,378],[990,390],[994,402],[1014,482],[1002,514],[982,562],[978,570],[962,586],[902,630],[874,646],[846,658],[766,666],[690,666],[658,654]]', 1, 80, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-02_G.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('7495880f-865c-4463-b4f1-3a650e232cf4', 'Buncis Rebus - Porsi G (100g)', '/uploads/atlas/AS/AS-02_G.jpg', '/uploads/atlas/AS/AS-02_G.jpg', 1600, 900, 'published', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('c27e405b-e93c-419b-bd94-a489b67e87dc', '7495880f-865c-4463-b4f1-3a650e232cf4', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('9a41a5d0-dd7e-48e6-abae-8d87905908af', '7495880f-865c-4463-b4f1-3a650e232cf4', 'Buncis Rebus Porsi G', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '[[586,442],[594,398],[598,382],[610,350],[634,302],[722,270],[838,278],[902,294],[918,310],[930,326],[950,358],[990,434],[994,442],[1006,478],[1002,522],[990,554],[986,562],[978,574],[966,586],[838,650],[762,666],[698,666],[666,662],[638,658],[630,654],[626,650],[594,602],[590,594],[586,470]]', 1, 100, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-02_H.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('2b78bfdc-0bc8-40e8-bf7e-679468ecb603', 'Buncis Rebus - Porsi H (130g)', '/uploads/atlas/AS/AS-02_H.jpg', '/uploads/atlas/AS/AS-02_H.jpg', 1600, 900, 'published', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('ce83b36e-f006-47eb-9c14-7baae20d0b96', '2b78bfdc-0bc8-40e8-bf7e-679468ecb603', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('4c4ff8d9-4e4c-4020-800f-0b50b3ffab4e', '2b78bfdc-0bc8-40e8-bf7e-679468ecb603', 'Buncis Rebus Porsi H', '47ad2a14-cfaf-11f1-ab66-664527f556d1', '[[474,458],[534,366],[554,342],[558,338],[598,310],[634,298],[686,290],[730,290],[870,302],[910,322],[946,346],[974,366],[978,370],[994,398],[998,406],[1010,442],[1014,462],[1018,502],[1006,526],[962,574],[938,598],[910,622],[886,638],[874,642],[798,662],[686,662],[662,658],[646,654],[594,634],[546,606],[526,582],[486,522],[474,462]]', 1, 130, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-04_A.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('3d4bccb4-f565-4b31-b644-5aacd7da6b99', 'Kacang Panjang Rebus - Porsi A (10g)', '/uploads/atlas/AS/AS-04_A.jpg', '/uploads/atlas/AS/AS-04_A.jpg', 1600, 900, 'published', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('48e888b8-57e4-426e-a068-175b6fa230d1', '3d4bccb4-f565-4b31-b644-5aacd7da6b99', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('c770d612-6cc2-4d42-a662-866a246d4c7c', '3d4bccb4-f565-4b31-b644-5aacd7da6b99', 'Kacang Panjang Rebus Porsi A', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '[[886,362],[890,350],[906,346],[914,346],[950,362],[958,374],[970,482],[962,594],[950,598],[898,610],[890,610],[886,602]]', 1, 10, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-04_B.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('d17959f3-7403-400d-84d6-7444526abe77', 'Kacang Panjang Rebus - Porsi B (30g)', '/uploads/atlas/AS/AS-04_B.jpg', '/uploads/atlas/AS/AS-04_B.jpg', 1600, 900, 'published', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('8793d07d-66c2-4b54-afca-c80ef29715ea', 'd17959f3-7403-400d-84d6-7444526abe77', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('eafc4cf5-d089-4825-b428-d2bb6d9cb4ae', 'd17959f3-7403-400d-84d6-7444526abe77', 'Kacang Panjang Rebus Porsi B', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '[[782,654],[790,342],[794,338],[898,334],[942,334],[962,370],[970,390],[990,470],[970,602],[966,606],[946,618],[918,630],[890,638],[858,646],[794,658],[786,658]]', 1, 30, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-04_C.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('2fec0e0f-394b-4718-94e7-703d408f14c9', 'Kacang Panjang Rebus - Porsi C (50g)', '/uploads/atlas/AS/AS-04_C.jpg', '/uploads/atlas/AS/AS-04_C.jpg', 1600, 900, 'published', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('516d0e18-7065-4ca5-b217-7d780c04e060', '2fec0e0f-394b-4718-94e7-703d408f14c9', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('f05ad40c-b142-4c66-ac27-1689e7f06103', '2fec0e0f-394b-4718-94e7-703d408f14c9', 'Kacang Panjang Rebus Porsi C', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '[[742,506],[746,390],[754,326],[774,302],[830,302],[902,322],[934,334],[946,342],[986,378],[990,390],[990,454],[974,586],[970,598],[966,606],[926,634],[918,638],[858,658],[834,662],[798,666],[774,666],[746,630],[742,622]]', 1, 50, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-04_D.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('104b7476-f479-4213-bd79-b4dfe661e52f', 'Kacang Panjang Rebus - Porsi D (70g)', '/uploads/atlas/AS/AS-04_D.jpg', '/uploads/atlas/AS/AS-04_D.jpg', 1600, 900, 'published', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('d17107ee-d65f-431e-aea3-cfad2fc5c07f', '104b7476-f479-4213-bd79-b4dfe661e52f', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('a55e3b5c-5366-44ab-a6b6-e45b168f7ba2', '104b7476-f479-4213-bd79-b4dfe661e52f', 'Kacang Panjang Rebus Porsi D', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '[[626,638],[662,322],[666,310],[682,310],[798,314],[938,334],[946,338],[950,342],[966,362],[982,398],[990,418],[994,438],[998,474],[998,514],[994,530],[958,602],[922,634],[910,642],[902,646],[858,662],[846,666],[790,682],[774,682],[746,678],[630,642]]', 1, 70, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-04_E.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('91d72b92-3e37-410e-8527-661fc36d1c6c', 'Kacang Panjang Rebus - Porsi E (90g)', '/uploads/atlas/AS/AS-04_E.jpg', '/uploads/atlas/AS/AS-04_E.jpg', 1600, 900, 'published', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('d84cd5d8-5e37-4f74-8fcf-29df32a7e512', '91d72b92-3e37-410e-8527-661fc36d1c6c', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('07cf6c15-4dd7-4e9b-92e1-fde7afb3fbf3', '91d72b92-3e37-410e-8527-661fc36d1c6c', 'Kacang Panjang Rebus Porsi E', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '[[610,642],[622,318],[626,298],[630,294],[654,290],[718,290],[902,310],[930,322],[946,330],[998,378],[1002,382],[1014,506],[1014,522],[1010,538],[994,570],[982,590],[974,602],[958,618],[934,634],[906,650],[898,654],[834,678],[798,686],[766,690],[758,690],[726,686],[674,678],[650,674],[630,670],[618,666],[610,662]]', 1, 90, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-04_F.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('95d11baa-ae8c-42e7-97b6-8068731a8f8c', 'Kacang Panjang Rebus - Porsi F (110g)', '/uploads/atlas/AS/AS-04_F.jpg', '/uploads/atlas/AS/AS-04_F.jpg', 1600, 900, 'published', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('4adf0a8b-d9f1-4092-8d70-dc5ab5274fac', '95d11baa-ae8c-42e7-97b6-8068731a8f8c', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('9f6b64e9-6bfa-4cf0-9e8d-408f17d9b4ce', '95d11baa-ae8c-42e7-97b6-8068731a8f8c', 'Kacang Panjang Rebus Porsi F', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '[[458,438],[502,358],[522,338],[578,302],[590,298],[718,290],[874,306],[994,374],[998,378],[1002,402],[1006,454],[1006,458],[1002,562],[998,570],[978,606],[966,622],[878,662],[866,666],[846,670],[814,674],[766,678],[694,682],[670,682],[590,674],[582,670],[566,658],[502,598],[498,594],[482,562],[458,450]]', 1, 110, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-04_G.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('ef1b58d5-4b2a-4ab3-8f7a-c766c020fa6b', 'Kacang Panjang Rebus - Porsi G (130g)', '/uploads/atlas/AS/AS-04_G.jpg', '/uploads/atlas/AS/AS-04_G.jpg', 1600, 900, 'published', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('dd8f3500-6918-4e70-818f-23c5f1dd77a0', 'ef1b58d5-4b2a-4ab3-8f7a-c766c020fa6b', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('d5143f83-137c-471f-833f-694f09b168d5', 'ef1b58d5-4b2a-4ab3-8f7a-c766c020fa6b', 'Kacang Panjang Rebus Porsi G', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '[[474,498],[486,454],[506,418],[546,358],[550,354],[558,350],[650,326],[678,322],[798,318],[822,318],[918,330],[1010,390],[1014,394],[1018,402],[1022,498],[1022,502],[1018,518],[1002,554],[990,578],[978,598],[962,610],[914,642],[906,646],[878,658],[834,670],[774,678],[690,678],[654,674],[634,670],[606,662],[598,658],[494,590],[490,586],[478,562],[474,550]]', 1, 130, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-04_H.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('b7699961-885a-4dc6-a4a2-38cac821e538', 'Kacang Panjang Rebus - Porsi H (190g)', '/uploads/atlas/AS/AS-04_H.jpg', '/uploads/atlas/AS/AS-04_H.jpg', 1600, 900, 'published', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('2c9f4151-3c33-45f0-81ff-acdf7b48145f', 'b7699961-885a-4dc6-a4a2-38cac821e538', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('0167084c-75b0-4544-b0a3-6c7deb92dffd', 'b7699961-885a-4dc6-a4a2-38cac821e538', 'Kacang Panjang Rebus Porsi H', '47ad7abe-cfaf-11f1-ab66-664527f556d1', '[[458,506],[462,490],[478,442],[494,398],[502,382],[542,310],[738,294],[886,294],[934,322],[942,330],[982,386],[990,402],[1002,434],[1010,514],[1010,530],[1002,550],[986,578],[970,598],[898,650],[870,658],[826,666],[770,674],[730,674],[630,658],[594,650],[582,646],[566,638],[462,578],[458,574]]', 1, 190, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-05_A.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('8cc7cde5-4af9-4ca1-88da-8d939639b71d', 'Daun Singkong Rebus - Porsi A (20g)', '/uploads/atlas/AS/AS-05_A.jpg', '/uploads/atlas/AS/AS-05_A.jpg', 1600, 900, 'published', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('1e049752-988a-410c-a643-8143032d22e0', '8cc7cde5-4af9-4ca1-88da-8d939639b71d', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('fbfcbcc9-11b0-4fe4-a438-892ed7c99831', '8cc7cde5-4af9-4ca1-88da-8d939639b71d', 'Daun Singkong Rebus Porsi A', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '[[770,486],[774,390],[786,378],[826,358],[886,346],[902,346],[954,358],[986,430],[986,490],[946,602],[914,618],[838,650],[830,650],[822,646],[790,622],[786,618],[774,582],[770,490]]', 1, 20, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-05_B.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('7e1cbafa-8717-485c-94b6-0a9a4bbf1762', 'Daun Singkong Rebus - Porsi B (30g)', '/uploads/atlas/AS/AS-05_B.jpg', '/uploads/atlas/AS/AS-05_B.jpg', 1600, 900, 'published', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('68f2f2ed-ab90-41f2-ad4f-4b15fcd5a3f5', '7e1cbafa-8717-485c-94b6-0a9a4bbf1762', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('69510507-ac20-4c58-8c4c-83324c01ada4', '7e1cbafa-8717-485c-94b6-0a9a4bbf1762', 'Daun Singkong Rebus Porsi B', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '[[734,530],[738,458],[758,370],[766,350],[778,338],[870,318],[986,406],[990,474],[966,554],[950,586],[914,622],[854,658],[774,666],[766,666],[758,662],[754,658],[746,634],[738,602],[734,538]]', 1, 30, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-05_C.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('6d813152-0a38-47a0-9ac9-8c84db44063f', 'Daun Singkong Rebus - Porsi C (40g)', '/uploads/atlas/AS/AS-05_C.jpg', '/uploads/atlas/AS/AS-05_C.jpg', 1600, 900, 'published', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('5e646148-37f4-4add-bba8-738113cf9f00', '6d813152-0a38-47a0-9ac9-8c84db44063f', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('e89c7463-ae45-44a5-8df1-958c80efecea', '6d813152-0a38-47a0-9ac9-8c84db44063f', 'Daun Singkong Rebus Porsi C', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '[[690,634],[694,382],[702,342],[710,326],[718,318],[770,306],[846,314],[854,318],[934,370],[938,374],[978,418],[974,554],[970,562],[934,602],[922,614],[902,626],[854,650],[834,658],[822,662],[806,666],[782,670],[722,670],[710,658]]', 1, 40, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-05_D.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('f30a60b9-6dda-47ca-a64a-76df8c86524c', 'Daun Singkong Rebus - Porsi D (50g)', '/uploads/atlas/AS/AS-05_D.jpg', '/uploads/atlas/AS/AS-05_D.jpg', 1600, 900, 'published', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('a8d2e7e8-28af-45e3-8939-fb46e8884895', 'f30a60b9-6dda-47ca-a64a-76df8c86524c', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('61648a1b-f592-451a-8dcc-26917d16f272', 'f30a60b9-6dda-47ca-a64a-76df8c86524c', 'Daun Singkong Rebus Porsi D', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '[[686,438],[702,366],[714,350],[738,326],[778,290],[786,290],[858,302],[894,330],[946,378],[978,430],[994,474],[994,534],[970,574],[946,598],[922,618],[890,642],[806,682],[726,670],[718,666],[698,650],[686,574]]', 1, 50, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-05_E.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('c47b8766-d170-488a-92c9-f346c7be8b50', 'Daun Singkong Rebus - Porsi E (60g)', '/uploads/atlas/AS/AS-05_E.jpg', '/uploads/atlas/AS/AS-05_E.jpg', 1600, 900, 'published', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('2e301290-7735-4a09-b0d9-6c06c4c81a96', 'c47b8766-d170-488a-92c9-f346c7be8b50', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('1f54f592-99af-453b-a69a-2f309c0d678c', 'c47b8766-d170-488a-92c9-f346c7be8b50', 'Daun Singkong Rebus Porsi E', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '[[626,486],[654,358],[658,346],[674,310],[714,294],[730,294],[834,310],[858,314],[906,338],[942,358],[954,370],[994,438],[994,494],[982,554],[966,594],[954,602],[854,658],[818,670],[798,674],[742,682],[718,682],[654,626],[650,618],[634,566]]', 1, 60, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-05_F.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('fc1a2996-979f-4545-9942-8b96638b386e', 'Daun Singkong Rebus - Porsi F (70g)', '/uploads/atlas/AS/AS-05_F.jpg', '/uploads/atlas/AS/AS-05_F.jpg', 1600, 900, 'published', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('6ab710b0-1660-407d-81c1-795339293ecd', 'fc1a2996-979f-4545-9942-8b96638b386e', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('302ec486-1f86-4ccb-8bb5-659f96e3f638', 'fc1a2996-979f-4545-9942-8b96638b386e', 'Daun Singkong Rebus Porsi F', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '[[634,462],[638,418],[674,350],[686,330],[742,282],[882,310],[1010,438],[1002,502],[986,538],[974,562],[902,642],[842,666],[766,678],[702,678],[654,654],[650,642],[642,602],[634,474]]', 1, 70, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-05_G.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('c0b6b471-2f80-4c9b-85c6-5b61184c1f2a', 'Daun Singkong Rebus - Porsi G (80g)', '/uploads/atlas/AS/AS-05_G.jpg', '/uploads/atlas/AS/AS-05_G.jpg', 1600, 900, 'published', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('189cd463-cc7b-4e77-9a7b-7560657f82e4', 'c0b6b471-2f80-4c9b-85c6-5b61184c1f2a', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('214d5efe-87ee-40be-afce-5b75aadb91af', 'c0b6b471-2f80-4c9b-85c6-5b61184c1f2a', 'Daun Singkong Rebus Porsi G', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '[[578,538],[582,482],[586,442],[622,338],[626,330],[634,318],[642,310],[690,294],[750,286],[786,286],[818,298],[906,342],[910,346],[994,434],[994,506],[990,530],[978,554],[918,618],[910,626],[874,646],[842,658],[742,670],[714,670],[662,666],[618,654],[586,570]]', 1, 80, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-05_H.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('6e27a1de-7cc7-4cfc-a036-c72497bb0581', 'Daun Singkong Rebus - Porsi H (90g)', '/uploads/atlas/AS/AS-05_H.jpg', '/uploads/atlas/AS/AS-05_H.jpg', 1600, 900, 'published', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('b2ba5aeb-696b-42e5-afbd-1f76b9fbc203', '6e27a1de-7cc7-4cfc-a036-c72497bb0581', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('7911421d-dae5-4aac-b610-c2176fe5746d', '6e27a1de-7cc7-4cfc-a036-c72497bb0581', 'Daun Singkong Rebus Porsi H', '47ad9b5c-cfaf-11f1-ab66-664527f556d1', '[[458,530],[506,378],[546,338],[586,310],[598,306],[678,302],[830,298],[846,302],[930,370],[962,406],[966,418],[970,434],[974,462],[982,546],[938,594],[806,666],[794,670],[758,674],[710,678],[686,678],[606,658],[582,650],[574,646],[562,638],[474,558],[462,538]]', 1, 90, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-07_A.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('47eb70bc-be9b-4de8-8930-59ab0c3f556a', 'Kangkung Rebus - Porsi A (10g)', '/uploads/atlas/AS/AS-07_A.jpg', '/uploads/atlas/AS/AS-07_A.jpg', 1600, 900, 'published', '47addb3a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('88fa368d-4ad4-4382-a3ec-4ab596ad6ab2', '47eb70bc-be9b-4de8-8930-59ab0c3f556a', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('e85eab2e-42c5-452c-8553-33dd21fd344c', '47eb70bc-be9b-4de8-8930-59ab0c3f556a', 'Kangkung Rebus Porsi A', '47addb3a-cfaf-11f1-ab66-664527f556d1', '[[738,482],[766,410],[778,382],[782,378],[846,354],[1146,410],[1146,434],[838,638],[810,650],[738,586]]', 1, 10, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-07_B.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('e2455ad1-2a16-4033-bbc0-f60d94120419', 'Kangkung Rebus - Porsi B (30g)', '/uploads/atlas/AS/AS-07_B.jpg', '/uploads/atlas/AS/AS-07_B.jpg', 1600, 900, 'published', '47addb3a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('15d04523-78d5-461a-9444-46023a52fcec', 'e2455ad1-2a16-4033-bbc0-f60d94120419', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('b3704024-cb43-4abb-a2a6-89c2cc204d07', 'e2455ad1-2a16-4033-bbc0-f60d94120419', 'Kangkung Rebus Porsi B', '47addb3a-cfaf-11f1-ab66-664527f556d1', '[[626,458],[658,406],[690,358],[706,342],[742,322],[750,318],[1134,410],[1146,414],[1146,694],[730,714],[698,710],[670,634],[666,622],[626,470]]', 1, 30, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-07_C.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('71836d4a-a6bd-48e8-9f79-0317fa724644', 'Kangkung Rebus - Porsi C (40g)', '/uploads/atlas/AS/AS-07_C.jpg', '/uploads/atlas/AS/AS-07_C.jpg', 1600, 900, 'published', '47addb3a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('0173ca24-12b6-4f57-9a63-ee19ca043f3f', '71836d4a-a6bd-48e8-9f79-0317fa724644', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('4f6d9962-5f06-4800-8d01-38a9d5a61c13', '71836d4a-a6bd-48e8-9f79-0317fa724644', 'Kangkung Rebus Porsi C', '47addb3a-cfaf-11f1-ab66-664527f556d1', '[[602,626],[610,462],[654,330],[658,322],[694,282],[1146,414],[1146,726],[650,718],[614,702],[606,682],[602,654]]', 1, 40, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-07_D.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('0eff86bd-2a8a-4294-aa02-b9c4b14028ba', 'Kangkung Rebus - Porsi D (60g)', '/uploads/atlas/AS/AS-07_D.jpg', '/uploads/atlas/AS/AS-07_D.jpg', 1600, 900, 'published', '47addb3a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('11451d74-fdad-4f65-9289-721bf23c3d4c', '0eff86bd-2a8a-4294-aa02-b9c4b14028ba', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('584184b2-bdd4-4e27-b5c0-778e37d0738c', '0eff86bd-2a8a-4294-aa02-b9c4b14028ba', 'Kangkung Rebus Porsi D', '47addb3a-cfaf-11f1-ab66-664527f556d1', '[[554,586],[562,458],[566,434],[586,350],[618,290],[622,286],[750,266],[1146,414],[1146,690],[590,722],[570,702]]', 1, 60, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-07_E.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('a3560a7b-0e8f-4ebb-b9a6-20a2b221cc58', 'Kangkung Rebus - Porsi E (80g)', '/uploads/atlas/AS/AS-07_E.jpg', '/uploads/atlas/AS/AS-07_E.jpg', 1600, 900, 'published', '47addb3a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('f38909f6-65a1-4702-be8d-d37937dea7af', 'a3560a7b-0e8f-4ebb-b9a6-20a2b221cc58', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('fea921c9-eb89-450a-b0c5-cd4d83a4c74a', 'a3560a7b-0e8f-4ebb-b9a6-20a2b221cc58', 'Kangkung Rebus Porsi E', '47addb3a-cfaf-11f1-ab66-664527f556d1', '[[474,510],[482,410],[546,310],[550,306],[562,302],[618,290],[706,286],[1146,418],[1146,682],[674,726],[598,726],[538,706],[510,662],[502,646],[474,514]]', 1, 80, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-07_F.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('3c745fae-9bd5-48d5-8abf-8c80d69986d1', 'Kangkung Rebus - Porsi F (90g)', '/uploads/atlas/AS/AS-07_F.jpg', '/uploads/atlas/AS/AS-07_F.jpg', 1600, 900, 'published', '47addb3a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('10bffd1d-0819-4fe2-b783-da598ebb2c14', '3c745fae-9bd5-48d5-8abf-8c80d69986d1', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('df59c9ab-61fb-4b27-81d9-3eefe34824d2', '3c745fae-9bd5-48d5-8abf-8c80d69986d1', 'Kangkung Rebus Porsi F', '47addb3a-cfaf-11f1-ab66-664527f556d1', '[[434,438],[438,422],[462,350],[538,274],[666,270],[1142,270],[1146,430],[1146,722],[530,718],[450,662],[446,658],[434,450]]', 1, 90, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-07_G.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('5bac74b8-8ba2-448a-87a5-1714ae44a12c', 'Kangkung Rebus - Porsi G (100g)', '/uploads/atlas/AS/AS-07_G.jpg', '/uploads/atlas/AS/AS-07_G.jpg', 1600, 900, 'published', '47addb3a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('4bf9a9ce-3cfb-4afd-b2b3-41176ebccddf', '5bac74b8-8ba2-448a-87a5-1714ae44a12c', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('2da4b95b-411a-44b0-9b0d-13e780ec125c', '5bac74b8-8ba2-448a-87a5-1714ae44a12c', 'Kangkung Rebus Porsi G', '47addb3a-cfaf-11f1-ab66-664527f556d1', '[[394,530],[434,366],[438,354],[446,342],[462,330],[558,282],[582,282],[698,290],[1142,410],[1146,414],[1146,726],[610,726],[534,718],[502,714],[486,710],[450,694],[442,690],[430,682],[402,626],[394,554]]', 1, 100, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-07_H.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('fe76ecfe-5f75-4e5a-9826-2e262c75d84e', 'Kangkung Rebus - Porsi H (140g)', '/uploads/atlas/AS/AS-07_H.jpg', '/uploads/atlas/AS/AS-07_H.jpg', 1600, 900, 'published', '47addb3a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('27cb645f-798f-413e-987a-b8ba6530770d', 'fe76ecfe-5f75-4e5a-9826-2e262c75d84e', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('ca87594a-e3e2-46b4-8263-f8b83b269390', 'fe76ecfe-5f75-4e5a-9826-2e262c75d84e', 'Kangkung Rebus Porsi H', '47addb3a-cfaf-11f1-ab66-664527f556d1', '[[350,406],[358,386],[466,310],[498,306],[586,298],[1146,270],[1146,714],[1142,722],[538,730],[514,726],[494,722],[458,710],[442,702],[374,666],[366,658],[350,634]]', 1, 140, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-35_A.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('7aa97984-46e6-4723-ab34-604b4a6250a7', 'Wortel Rebus - Porsi A (10g)', '/uploads/atlas/AS/AS-35_A.jpg', '/uploads/atlas/AS/AS-35_A.jpg', 1600, 900, 'published', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('107d8ed5-5cf0-4420-9a3d-4df191b195e5', '7aa97984-46e6-4723-ab34-604b4a6250a7', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('26ff0c8d-4897-4785-a1ce-6f0502b6f2f3', '7aa97984-46e6-4723-ab34-604b4a6250a7', 'Wortel Rebus Porsi A', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '[[918,490],[922,450],[926,426],[930,422],[942,414],[974,414],[1006,434],[1010,450],[1010,502],[1002,518],[994,526],[974,534],[962,538],[946,538],[934,534],[922,518]]', 1, 10, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-35_B.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('b84327d1-6e77-4d21-bbb3-31c46f3b9383', 'Wortel Rebus - Porsi B (20g)', '/uploads/atlas/AS/AS-35_B.jpg', '/uploads/atlas/AS/AS-35_B.jpg', 1600, 900, 'published', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('7aca7e37-99e2-444f-995a-87d16b7cf721', 'b84327d1-6e77-4d21-bbb3-31c46f3b9383', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('03f2c909-778e-436f-aaf0-10c9d9d54393', 'b84327d1-6e77-4d21-bbb3-31c46f3b9383', 'Wortel Rebus Porsi B', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '[[870,550],[878,410],[882,390],[894,374],[898,370],[906,366],[918,362],[934,362],[946,366],[966,374],[974,378],[990,394],[1014,466],[1014,486],[978,558],[958,574],[950,578],[938,582],[898,582],[886,574],[878,566]]', 1, 20, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-35_C.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('c6fae6b5-1627-4ba5-b8ad-7f9721128170', 'Wortel Rebus - Porsi C (40g)', '/uploads/atlas/AS/AS-35_C.jpg', '/uploads/atlas/AS/AS-35_C.jpg', 1600, 900, 'published', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('acc21633-5462-4fe4-b174-e1365b608853', 'c6fae6b5-1627-4ba5-b8ad-7f9721128170', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('572ad914-2a0c-442d-b8b5-3375dedbf835', 'c6fae6b5-1627-4ba5-b8ad-7f9721128170', 'Wortel Rebus Porsi C', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '[[770,494],[774,474],[778,466],[842,358],[858,342],[874,334],[886,330],[906,334],[930,350],[942,362],[1006,446],[1010,458],[1010,466],[926,606],[922,610],[910,618],[898,622],[878,626],[862,626],[850,622],[834,610],[830,606],[774,514]]', 1, 40, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-35_D.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('3ce625b2-2bf5-46fa-9191-9e5a1562c5c8', 'Wortel Rebus - Porsi D (50g)', '/uploads/atlas/AS/AS-35_D.jpg', '/uploads/atlas/AS/AS-35_D.jpg', 1600, 900, 'published', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('3506d707-38e5-4726-adb3-9fc4bdeb280c', '3ce625b2-2bf5-46fa-9191-9e5a1562c5c8', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('ef4fe6d8-56a7-4e58-b934-acf65b69e6e0', '3ce625b2-2bf5-46fa-9191-9e5a1562c5c8', 'Wortel Rebus Porsi D', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '[[798,538],[802,466],[810,326],[850,306],[866,306],[974,358],[994,378],[1002,398],[1010,454],[1010,474],[1002,530],[994,546],[986,554],[902,630],[870,646],[846,654],[838,654],[814,638],[806,630]]', 1, 50, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-35_E.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('6b9e2b35-6204-442f-acc1-61d39a7da942', 'Wortel Rebus - Porsi E (70g)', '/uploads/atlas/AS/AS-35_E.jpg', '/uploads/atlas/AS/AS-35_E.jpg', 1600, 900, 'published', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('90558695-1346-4e68-ba92-1032a5be64d5', '6b9e2b35-6204-442f-acc1-61d39a7da942', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('d7669516-510d-4612-9ca8-78b020c77a6c', '6b9e2b35-6204-442f-acc1-61d39a7da942', 'Wortel Rebus Porsi E', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '[[722,494],[734,382],[742,370],[750,362],[814,318],[830,310],[854,310],[934,334],[958,346],[970,354],[982,370],[1010,438],[1014,450],[1010,514],[1006,526],[998,538],[954,586],[902,634],[886,642],[866,650],[838,658],[822,658],[758,606],[750,594],[746,586],[722,522]]', 1, 70, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-35_F.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('f76a304d-aec9-413b-84d7-ae30a7df0ccf', 'Wortel Rebus - Porsi F (80g)', '/uploads/atlas/AS/AS-35_F.jpg', '/uploads/atlas/AS/AS-35_F.jpg', 1600, 900, 'published', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('ddbf5e19-3649-4458-8222-aefe1abd6a93', 'f76a304d-aec9-413b-84d7-ae30a7df0ccf', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('032f25e8-e04d-4a10-a0b6-4a8f6b2d90b6', 'f76a304d-aec9-413b-84d7-ae30a7df0ccf', 'Wortel Rebus Porsi F', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '[[742,622],[754,382],[758,374],[802,326],[810,318],[822,310],[834,306],[842,306],[854,310],[970,354],[978,358],[986,366],[1026,442],[1026,454],[1014,506],[1010,518],[970,570],[938,606],[930,614],[906,634],[890,642],[854,654],[826,662],[810,666],[790,670],[766,670],[750,650],[746,642]]', 1, 80, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-35_G.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('45a2dc8c-6b95-464b-bb75-99b41f590102', 'Wortel Rebus - Porsi G (100g)', '/uploads/atlas/AS/AS-35_G.jpg', '/uploads/atlas/AS/AS-35_G.jpg', 1600, 900, 'published', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('887490d9-8376-49e9-8b3e-f0e45448b3aa', '45a2dc8c-6b95-464b-bb75-99b41f590102', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('492e7792-aa03-4225-a88c-415c2e97dfcc', '45a2dc8c-6b95-464b-bb75-99b41f590102', 'Wortel Rebus Porsi G', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '[[702,534],[706,466],[722,406],[726,394],[730,386],[794,318],[810,302],[818,298],[834,294],[850,298],[922,334],[946,350],[1010,434],[1014,446],[1014,518],[1010,526],[926,622],[922,626],[906,634],[882,642],[782,670],[766,666],[754,662],[734,654],[730,650],[726,642],[702,538]]', 1, 100, NOW(), NOW());
      
DELETE FROM food_images WHERE image_url = '/uploads/atlas/AS/AS-35_H.jpg';

        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('a0e6294c-6fa4-458f-85a6-653f8d77f0fe', 'Wortel Rebus - Porsi H (120g)', '/uploads/atlas/AS/AS-35_H.jpg', '/uploads/atlas/AS/AS-35_H.jpg', 1600, 900, 'published', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '7ff4eb2c-1a20-4e8c-9542-43c1b204f91d', NOW(), NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('2ee23abf-4ba1-4c73-b267-7503b38e209d', 'a0e6294c-6fa4-458f-85a6-653f8d77f0fe', 'Piring Saji Baku (AH-01)', NULL, '[[1120,540],[1101,623],[1047,699],[963,758],[857,797],[740,810],[623,797],[517,758],[433,699],[379,623],[360,540],[379,457],[433,381],[517,322],[623,283],[740,270],[857,283],[963,322],[1047,381],[1101,457]]', 0, NULL, NOW(), NOW());
      

        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('8e5d082b-1b93-4b09-b8ca-2d672409b9e9', 'a0e6294c-6fa4-458f-85a6-653f8d77f0fe', 'Wortel Rebus Porsi H', '47b1f44a-cfaf-11f1-ab66-664527f556d1', '[[670,562],[674,374],[678,354],[714,302],[718,298],[726,294],[794,290],[802,290],[918,322],[930,330],[946,346],[978,394],[1018,466],[1018,494],[990,550],[982,562],[918,626],[898,638],[882,646],[826,666],[794,674],[730,674],[702,658],[690,646],[674,586],[670,570]]', 1, 120, NOW(), NOW());
      