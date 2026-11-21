import urllib.request
import json

url = "https://api.atlasmakananku.id/api/v1/public/categories/LH/foods?limit=50"
req = urllib.request.Request(url, headers={"User-Agent": "Mozilla/5.0"})
try:
    with urllib.request.urlopen(req) as resp:
        data = json.loads(resp.read().decode())
        foods = data.get("data", [])
        if isinstance(foods, dict):
            foods = foods.get("foods", []) or foods.get("items", [])
        print(f"Total foods in category LH: {len(foods)}")
        for f in foods:
            code = f.get("code")
            if code in ['LH-31', 'LH-32', 'LH-33', 'LH-34', 'LH-35']:
                print(f"{code}: name={f.get('name')}, image_url={f.get('image_url')}")
except Exception as e:
    print(f"Error: {e}")
