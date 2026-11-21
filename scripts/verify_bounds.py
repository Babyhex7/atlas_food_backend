import json

with open(r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads\all_fish_polygons.json") as f:
    data = json.load(f)

dims = {
    "LH-31": (1063, 598),
    "LH-32": (508, 286),
    "LH-33": (770, 434),
    "LH-34": (1064, 598),
    "LH-35": (774, 437)
}

for fish, portions in data.items():
    w, h = dims[fish]
    print(f"=== {fish} (image: {w}x{h}) ===")
    for label, poly in portions.items():
        xs = [p[0] for p in poly]
        ys = [p[1] for p in poly]
        min_x, max_x = min(xs), max(xs)
        min_y, max_y = min(ys), max(ys)
        valid = (0 <= min_x and max_x <= w and 0 <= min_y and max_y <= h)
        print(f"  Porsi {label}: pts={len(poly)}, x=[{min_x}, {max_x}], y=[{min_y}, {max_y}] -> Valid in bounds: {valid}")
