const { chromium } = require("@playwright/test");
const fs = require("fs");
const crypto = require("crypto");
const { execSync } = require("child_process");

function convexHull(points) {
  if (points.length <= 3) return points;
  points.sort((a, b) => a[0] === b[0] ? a[1] - b[1] : a[0] - b[0]);
  const cross = (o, a, b) => (a[0] - o[0]) * (b[1] - o[1]) - (a[1] - o[1]) * (b[0] - o[0]);
  
  const lower = [];
  for (let p of points) {
    while (lower.length >= 2 && cross(lower[lower.length - 2], lower[lower.length - 1], p) <= 0) {
      lower.pop();
    }
    lower.push(p);
  }
  
  const upper = [];
  for (let i = points.length - 1; i >= 0; i--) {
    const p = points[i];
    while (upper.length >= 2 && cross(upper[upper.length - 2], upper[upper.length - 1], p) <= 0) {
      upper.pop();
    }
    upper.push(p);
  }
  
  upper.pop();
  lower.pop();
  return lower.concat(upper);
}

function getPlatePolygon(code = "") {
  const points = [];
  const cx = code === "AS-07" ? 620 : 750;
  const cy = code === "AS-07" ? 516 : 492;
  const rx = code === "AS-07" ? 425 : 380;
  const ry = code === "AS-07" ? 314 : 282;
  const numPoints = 32;
  for (let i = 0; i < numPoints; i++) {
    const angle = (i / numPoints) * 2 * Math.PI;
    points.push([Math.round(cx + rx * Math.cos(angle)), Math.round(cy + ry * Math.sin(angle))]);
  }
  return points;
}

(async () => {
  // Query foods and portions via mysql CLI
  const mysqlCmd = (query) => {
    return execSync(`/opt/homebrew/bin/mysql -h 127.0.0.1 -u root -N -e "${query}" db_atlas_food`, { encoding: 'utf8' });
  };

  const adminId = mysqlCmd("SELECT id FROM users WHERE role = 'admin' LIMIT 1").trim();
  console.log("Admin ID:", adminId);

  const targetCodes = ['LN-01', 'AS-02', 'AS-04', 'AS-05', 'AS-07', 'AS-35'];
  const labels = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H'];

  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();

  const sqlStatements = [];

  for (const code of targetCodes) {
    const foodRaw = mysqlCmd(`SELECT id, name FROM foods WHERE code = '${code}'`).trim().split('\t');
    if (!foodRaw || !foodRaw[0]) continue;
    const foodId = foodRaw[0];
    const foodName = foodRaw[1];

    const portionsRaw = mysqlCmd(`
      SELECT i.label, i.weight_gram 
      FROM as_served_sets s 
      JOIN as_served_images i ON i.set_id = s.id 
      WHERE s.food_id = '${foodId}'
    `).trim().split('\n');

    const weightMap = {};
    portionsRaw.forEach(line => {
      if (!line) return;
      const parts = line.split('\t');
      weightMap[parts[0]] = parseFloat(parts[1]);
    });

    const catFolder = code.split('-')[0];
    console.log(`Processing ${code} - ${foodName}...`);

    for (const label of labels) {
      const weight = weightMap[label] || 0;
      const filename = `${code}_${label}.jpg`;
      const localFilePath = `/Users/macbookair/Developer/BRIN/atlas_food_backend/uploads/atlas/${catFolder}/${filename}`;
      const relativeImageUrl = `/uploads/atlas/${catFolder}/${filename}`;

      if (!fs.existsSync(localFilePath)) {
        console.log(`  File missing: ${localFilePath}`);
        continue;
      }

      const imgBase64 = fs.readFileSync(localFilePath).toString("base64");

      const detection = await page.evaluate(async (base64) => {
        return new Promise((resolve) => {
          const img = new Image();
          img.onload = () => {
            const canvas = document.createElement("canvas");
            canvas.width = img.width;
            canvas.height = img.height;
            const ctx = canvas.getContext("2d");
            ctx.drawImage(img, 0, 0);

            const imgData = ctx.getImageData(0, 0, canvas.width, canvas.height);
            const data = imgData.data;
            const foodPixels = [];

            for (let y = 250; y < 800; y += 4) {
              for (let x = 350; x < 1150; x += 4) {
                const idx = (y * canvas.width + x) * 4;
                const r = data[idx];
                const g = data[idx + 1];
                const b = data[idx + 2];

                const max = Math.max(r, g, b);
                const min = Math.min(r, g, b);
                const sat = max === 0 ? 0 : (max - min) / max;
                const brightness = (r + g + b) / 3;

                // Threshold contrast from white/grey plate
                const isFood = (brightness < 210 && sat > 0.12) || (brightness < 160);
                if (isFood) {
                  foodPixels.push([x, y]);
                }
              }
            }

            resolve({
              width: canvas.width,
              height: canvas.height,
              pixels: foodPixels
            });
          };
          img.src = "data:image/jpeg;base64," + base64;
        });
      }, imgBase64);

      let foodPolygon = [];
      if (detection.pixels.length > 10) {
        const sampled = detection.pixels.filter((_, idx) => idx % 2 === 0);
        foodPolygon = convexHull(sampled);
      } else {
        foodPolygon = [[700, 480], [780, 480], [780, 560], [700, 560]];
      }

      const platePolygon = getPlatePolygon(code);
      const foodImageId = crypto.randomUUID();
      const plateAreaId = crypto.randomUUID();

      const title = `${foodName} - Porsi ${label} (${weight}g)`;

      // SQL statements
      sqlStatements.push(`DELETE FROM food_images WHERE image_url = '${relativeImageUrl}';`);
      sqlStatements.push(`
        INSERT INTO food_images (id, title, image_url, thumbnail_url, width, height, status, primary_food_id, created_by, published_at, created_at, updated_at)
        VALUES ('${foodImageId}', '${title.replace(/'/g, "\\'")}', '${relativeImageUrl}', '${relativeImageUrl}', ${detection.width}, ${detection.height}, 'published', '${foodId}', '${adminId}', NOW(), NOW(), NOW());
      `);
      sqlStatements.push(`
        INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
        VALUES ('${plateAreaId}', '${foodImageId}', '${foodName} (Porsi ${label})', '${foodId}', '${JSON.stringify(platePolygon)}', 0, ${weight}, NOW(), NOW());
      `);
      sqlStatements.push(`UPDATE as_served_images SET description = 'food_image:${foodImageId}' WHERE image_url = '${relativeImageUrl}';`);

      console.log(`  ✓ Porsi ${label} (${weight}g): food_polygon=${foodPolygon.length} pts`);
    }
  }

  await browser.close();

  fs.writeFileSync("/tmp/seed_annotations.sql", sqlStatements.join('\n'));
  console.log(`\nGenerated ${sqlStatements.length} SQL statements in /tmp/seed_annotations.sql`);
})();
