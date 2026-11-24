import urllib.request, json

for name, fi in [('Gurame Cut', 'c4e340a1-7341-4567-89ab-lh34cut00001'), ('Patin Cut', 'c31000a1-7341-4567-89ab-lh31cut00001')]:
    url = f'https://atlasmakananku.id/api/v1/public/food-images/{fi}'
    req = urllib.request.Request(url, headers={'User-Agent': 'Mozilla/5.0'})
    with urllib.request.urlopen(req) as resp:
        data = json.loads(resp.read().decode())
        img = data.get('data', {})
        print(f"=== {name} ({img.get('width')}x{img.get('height')}) ===")
        for a in img.get('areas', []):
            poly = json.loads(a.get('polygon', '[]')) if isinstance(a.get('polygon'), str) else a.get('polygon', [])
            print(f"  Area: {a.get('name')} | {a.get('weight_gram')}g | pts: {len(poly)}")
