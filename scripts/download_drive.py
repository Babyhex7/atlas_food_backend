import urllib.request
import os

out_dir = r"C:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads"
os.makedirs(out_dir, exist_ok=True)

files = [
    ('DSC02496.jpg', '1lyZpp5tMjwXEAAJwF6EOa-qdNlFwOm91'),
    ('DSC02507.jpg', '1ehCTHWfogk-QUXeC6g72W2PG-fID1iK5')
]

for name, fid in files:
    download_url = f"https://drive.google.com/uc?export=download&id={fid}"
    dest = os.path.join(out_dir, name)
    print(f"Downloading {name} from {download_url}...")
    try:
        urllib.request.urlretrieve(download_url, dest)
        size = os.path.getsize(dest)
        print(f"  Downloaded {name}: {size} bytes")
    except Exception as e:
        print(f"  Error downloading {name}: {e}")
