import urllib.request
import hashlib
import os

def check_url(url):
    req = urllib.request.Request(url, headers={"User-Agent": "Mozilla/5.0"})
    try:
        with urllib.request.urlopen(req) as resp:
            headers = dict(resp.headers)
            body = resp.read()
            md5 = hashlib.md5(body).hexdigest()
            print(f"{url}:")
            print(f"  Status: {resp.status}, length: {len(body)}, md5: {md5}")
            print(f"  cf-cache-status: {headers.get('cf-cache-status')}")
            print(f"  cache-control: {headers.get('cache-control')}")
            print(f"  age: {headers.get('age')}")
            print(f"  etag: {headers.get('etag')}")
            print(f"  last-modified: {headers.get('last-modified')}")
            return body, md5
    except Exception as e:
        print(f"{url} ERROR: {e}")
        return None, None

print("=== CHECKING LH-32_A ===")
b1, m1 = check_url("https://api.atlasmakananku.id/uploads/atlas/LH/LH-32_A.jpg")
b2, m2 = check_url("https://atlasmakananku.id/uploads/atlas/LH/LH-32_A.jpg")

local_x264 = r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads\LH-32_pdf_3_X264.jpg.jpg"
local_x261 = r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\atlas_food_backend\scripts\drive_downloads\LH-32_pdf_0_X261.jpg.jpg"
if os.path.exists(local_x264):
    with open(local_x264, "rb") as f:
        print(f"Local X264 (Full Kakap): length={os.path.getsize(local_x264)}, md5={hashlib.md5(f.read()).hexdigest()}")
if os.path.exists(local_x261):
    with open(local_x261, "rb") as f:
        print(f"Local X261 (Potongan Kecil Kakap): length={os.path.getsize(local_x261)}, md5={hashlib.md5(f.read()).hexdigest()}")
