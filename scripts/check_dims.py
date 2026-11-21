from PIL import Image
import os

folder = r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads"
images = [
    "LH-31_A.jpg", "LH-31_cut.jpg",
    "LH-32_A.jpg", "LH-32_B.jpg", "LH-32_C.jpg", "LH-32_D.jpg",
    "LH-33_A.jpg", "LH-33_B.jpg", "LH-33_C.jpg",
    "LH-34_A.jpg", "LH-34_cut.jpg",
    "LH-35_A.jpg", "LH-35_B.jpg", "LH-35_C.jpg",
]

for img in images:
    path = os.path.join(folder, img)
    if os.path.exists(path):
        im = Image.open(path)
        print(f"{img}: size={im.size} (width={im.width}, height={im.height})")
    else:
        print(f"{img}: NOT FOUND locally")
