import urllib.request
import hashlib

url = "https://atlasmakananku.id/uploads/atlas/LH/LH-32_A.jpg?v=2"
req = urllib.request.Request(url, headers={"User-Agent": "Mozilla/5.0"})
with urllib.request.urlopen(req) as resp:
    body = resp.read()
    print(f"Status: {resp.status}, length: {len(body)}, md5: {hashlib.md5(body).hexdigest()}")
