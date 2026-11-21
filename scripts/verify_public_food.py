import urllib.request
import json

codes = ['LH-31', 'LH-32', 'LH-33', 'LH-34', 'LH-35']

for code in codes:
    url = f"https://api.atlasmakananku.id/api/v1/public/foods/search?q={code}"
    req = urllib.request.Request(url, headers={"User-Agent": "Mozilla/5.0"})
    with urllib.request.urlopen(req) as resp:
        data = json.loads(resp.read().decode())
        foods = data.get("data", [])
        if not foods:
            print(f"{code}: NOT FOUND")
            continue
        f = foods[0]
        print(f"=== {code}: {f['name']} ===")
        print(f"  Primary Thumbnail: {f.get('image_url')}")
        
        detail_url = f"https://api.atlasmakananku.id/api/v1/public/foods/{f['id']}"
        req_d = urllib.request.Request(detail_url, headers={"User-Agent": "Mozilla/5.0"})
        with urllib.request.urlopen(req_d) as resp_d:
            d = json.loads(resp_d.read().decode())['data']
            portions = d.get('portion_photos', [])
            print(f"  Portions ({len(portions)}):")
            for p in portions:
                print(f"    - Porsi {p.get('label')}: {p.get('weight_gram')}g -> img: {p.get('image_url')} (fi_id: {p.get('food_image_id')})")
    print()
