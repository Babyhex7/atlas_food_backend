import urllib.request
import json

for code in ['LH-31', 'LH-32', 'LH-33', 'LH-34', 'LH-35']:
    url = f"https://api.atlasmakananku.id/api/v1/public/foods/search?q={code}"
    req = urllib.request.Request(url, headers={"User-Agent": "Mozilla/5.0"})
    with urllib.request.urlopen(req) as resp:
        f = json.loads(resp.read().decode())['data'][0]
        fid = f['id']
        detail_url = f"https://api.atlasmakananku.id/api/v1/public/foods/{fid}"
        req_d = urllib.request.Request(detail_url, headers={"User-Agent": "Mozilla/5.0"})
        with urllib.request.urlopen(req_d) as resp_d:
            d = json.loads(resp_d.read().decode())['data']
            print(f"=== {code} ({d['name']}) ===")
            for p in d.get('portion_photos', []):
                print(f"  Porsi {p.get('label')}: {p.get('weight_gram')}g | img: {p.get('image_url')}")
