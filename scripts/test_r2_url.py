import urllib.request

images = [
    "LH-31_A.jpg", "LH-31_cut.jpg",
    "LH-32_A.jpg", "LH-32_B.jpg", "LH-32_C.jpg", "LH-32_D.jpg",
    "LH-33_A.jpg", "LH-33_B.jpg", "LH-33_C.jpg",
    "LH-34_A.jpg", "LH-34_cut.jpg",
    "LH-35_A.jpg", "LH-35_B.jpg", "LH-35_C.jpg",
]

for img in images:
    u = f"https://api.atlasmakananku.id/uploads/atlas/LH/{img}"
    try:
        req = urllib.request.Request(u, headers={"User-Agent": "Mozilla/5.0"}, method="HEAD")
        with urllib.request.urlopen(req) as resp:
            print(f"{img} -> {resp.status}")
    except Exception as e:
        print(f"{img} -> ERROR: {e}")
