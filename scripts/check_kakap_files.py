from PIL import Image
import os

folder = r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads"
for f in ["LH-32_A.jpg", "LH-32_B.jpg", "LH-32_C.jpg", "LH-32_D.jpg", "LH-32_pdf_3_X264.jpg.jpg", "LH-32_pdf_0_X261.jpg.jpg"]:
    p = os.path.join(folder, f)
    if os.path.exists(p):
        print(f"{f}: size={os.path.getsize(p)} bytes")
