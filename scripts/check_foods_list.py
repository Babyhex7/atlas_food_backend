import urllib.request
import json

url = "https://api.atlasmakananku.id/api/v1/public/foods?limit=100"
req = urllib.request.Request(url, headers={"User-Agent": "Mozilla/5.0"})
with urllib.request.urlopen(req) as resp:
    data = json.loads(resp.read().decode())
    items = data.get("data", [])
    print(f"Total items from /api/v1/public/foods: {len(items)}")
    for item in items:
        code = item.get("code")
        if code in ['LH-31', 'LH-32', 'LH-33', 'LH-34', 'LH-35']:
            print(f"Code: {code} | Name: {item.get('name')} | ImageURL: {item.get('image_url')}")
