import json
import cv2
import numpy as np

img_path = r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads\LH-34_cut.jpg"
img = cv2.imread(img_path)

with open(r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads\all_fish_polygons.json") as f:
    data = json.load(f)

gurame = data["LH-34"]

# Draw each portion
colors = {
    "A": (255, 0, 0),    # Blue
    "B": (0, 255, 0),    # Green (Kepala)
    "C": (0, 165, 255),  # Orange (Badan)
    "D": (0, 0, 255),    # Red (Ekor)
}

vis = img.copy()
for label in ["D", "C", "B"]:
    pts = np.array(gurame[label], np.int32).reshape((-1, 1, 2))
    cv2.polylines(vis, [pts], isClosed=True, color=colors[label], thickness=2)

out_path = r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads\vis_gurame_m17.jpg"
cv2.imwrite(out_path, vis)
print(f"Saved visualization to {out_path}")
