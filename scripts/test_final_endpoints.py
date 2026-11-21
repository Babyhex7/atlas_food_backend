import urllib.request
import json

urls = [
    "https://atlasmakananku.id/api/v1/public/foods/search?q=LH-32",
    "https://atlasmakananku.id/api/v1/public/food-images/c4e340a1-7341-4567-89ab-lh34cut00001",
]

for u in urls:
    req = urllib.request.Request(u, headers={"User-Agent": "Mozilla/5.0"})
    with urllib.request.urlopen(req) as resp:
        data = json.loads(resp.read().decode())
        print(f"=== {u} ===")
        print(f"Status: {resp.status}")
        if 'food-images' in u:
            print("Areas count:", len(data['data']['areas']))
            for a in data['data']['areas']:
                print(f"  {a['name']}: {a['weight_gram']}g, pts={len(a['polygon'])}")
        else:
            for item in data['data']:
                print(f"  {item['name']} ({item['code']}): image={item['image_url']}")
    print()
