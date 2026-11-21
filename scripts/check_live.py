import urllib.request
import json

fi_ids = [
    ('LH-34 Cut', 'c4e340a1-7341-4567-89ab-lh34cut00001'),
    ('LH-31 Cut', 'c31000a1-7341-4567-89ab-lh31cut00001'),
    ('LH-32 A', '6b0c0734-c7e3-49da-9e9b-1836ac8816ca'),
    ('LH-32 B', 'ed6583aa-c597-468c-8553-71671d540199'),
    ('LH-32 C', '0832d7d9-72ae-4c0f-a003-986096673a50'),
    ('LH-32 D', '5d1a435d-44d5-4a29-8685-0427076621d0'),
    ('LH-33 A', 'af898243-9abf-474e-84dc-8991a516493b'),
    ('LH-33 B', '074d1ee5-a3d4-4e33-aef9-cb85c174c13c'),
    ('LH-33 C', 'a3d1e35b-b5e2-490b-882a-94a9a3fad2df'),
    ('LH-35 A', '31581dd0-5a99-473d-add6-0f67b9615fc1'),
    ('LH-35 B', 'e2cbb90e-ce43-441d-8d54-eb576be985f3'),
    ('LH-35 C', '88e2fec8-a77b-4bf6-baae-f694201afd85'),
]

for name, fi_id in fi_ids:
    url = f"https://api.atlasmakananku.id/api/v1/public/food-images/{fi_id}"
    req = urllib.request.Request(url, headers={"User-Agent": "Mozilla/5.0"})
    try:
        with urllib.request.urlopen(req) as resp:
            data = json.loads(resp.read().decode())
            img = data.get('data', {})
            print(f"=== {name} ({fi_id}) ===")
            print(f"  title: {img.get('title')}, status: {img.get('status')}, size: {img.get('width')}x{img.get('height')}")
            areas = img.get('areas', [])
            print(f"  areas count: {len(areas)}")
            for a in areas:
                poly = a.get('polygon') or []
                print(f"    - {a.get('name')}: weight={a.get('weight_gram')}g, pts={len(poly)}")
    except Exception as e:
        print(f"{name} ({fi_id}) ERROR: {e}")
    print()
