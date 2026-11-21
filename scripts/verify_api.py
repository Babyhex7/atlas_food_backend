import urllib.request
import json

def check_food(code):
    url = f"https://api.atlasmakananku.id/api/v1/public/foods/search?q={code}"
    req = urllib.request.Request(url, headers={"User-Agent": "Mozilla/5.0"})
    with urllib.request.urlopen(req) as resp:
        data = json.loads(resp.read().decode())
        foods = data.get("data", [])
        if not foods:
            print(f"Food {code} not found")
            return
        f = foods[0]
        fid = f["id"]
        print(f"Food: {f['name']} (ID: {fid})")
        print(f"  Primary Thumbnail: {f.get('image_url')}")
        
        detail_url = f"https://api.atlasmakananku.id/api/v1/public/foods/{fid}"
        req_d = urllib.request.Request(detail_url, headers={"User-Agent": "Mozilla/5.0"})
        with urllib.request.urlopen(req_d) as resp_d:
            d_data = json.loads(resp_d.read().decode())
            food_detail = d_data.get("data", {})
            photos = food_detail.get("portion_photos", [])
            print(f"  Portion photos ({len(photos)}):")
            for p in photos:
                print(f"    - Porsi {p.get('label')}: {p.get('weight_gram')}g -> Image: {p.get('image_url')} (FoodImage: {p.get('food_image_id')})")

check_food("LH-34")
print()
check_food("LH-31")
