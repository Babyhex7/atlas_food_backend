const fs = require("fs");
const crypto = require("crypto");
const { execSync } = require("child_process");

function getPolygon(cx, cy, rx, ry, numPoints = 32) {
  const points = [];
  for (let i = 0; i < numPoints; i++) {
    const angle = (i / numPoints) * 2 * Math.PI;
    points.push([
      Math.round(cx + rx * Math.cos(angle)),
      Math.round(cy + ry * Math.sin(angle))
    ]);
  }
  return points;
}

const polyStandard = getPolygon(750, 492, 380, 282, 32);
const polyAS07 = getPolygon(620, 516, 425, 314, 32);

function main() {
  const query = `
    SELECT fi.id as food_image_id, fi.image_url, asi.weight_gram, asi.label, f.id as food_id, f.name as food_name, f.code as food_code
    FROM food_images fi
    JOIN as_served_images asi ON asi.description = CONCAT('food_image:', fi.id)
    JOIN as_served_sets ass ON ass.id = asi.set_id
    JOIN foods f ON f.id = ass.food_id
    ORDER BY f.code, asi.label;
  `;

  const raw = execSync(`/opt/homebrew/bin/mysql -h 127.0.0.1 -u root -N -e "${query}" db_atlas_food`, { encoding: 'utf8' });
  const lines = raw.trim().split('\n').filter(Boolean);

  console.log(`Found ${lines.length} food images to re-annotate.`);

  const sqlStatements = [];

  for (const line of lines) {
    const [food_image_id, image_url, weight_gram, label, food_id, food_name, food_code] = line.split('\t');
    const polygon = (food_code === "AS-07") ? polyAS07 : polyStandard;
    const areaId = crypto.randomUUID();
    const areaName = `${food_name} (Porsi ${label})`;

    sqlStatements.push(`DELETE FROM food_areas WHERE food_image_id = '${food_image_id}';`);
    sqlStatements.push(`
      INSERT INTO food_areas (id, food_image_id, name, food_id, polygon, z_index, weight_gram, created_at, updated_at)
      VALUES ('${areaId}', '${food_image_id}', '${areaName.replace(/'/g, "\\'")}', '${food_id}', '${JSON.stringify(polygon)}', 0, ${weight_gram}, NOW(), NOW());
    `);
  }

  const sqlFile = "/tmp/reannotate.sql";
  fs.writeFileSync(sqlFile, sqlStatements.join('\n'));
  console.log(`Generated ${sqlStatements.length} SQL statements in ${sqlFile}`);

  // Execute against database
  execSync(`/opt/homebrew/bin/mysql -h 127.0.0.1 -u root db_atlas_food < ${sqlFile}`, { stdio: 'inherit' });
  console.log("Successfully executed SQL updates in db_atlas_food!");
}

main();
