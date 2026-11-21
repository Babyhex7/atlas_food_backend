import urllib.request
import json

for code in ['LH-31', 'LH-32', 'LH-33', 'LH-34', 'LH-35']:
    url = f"https://api.atlasmakananku.id/api/v1/public/foods/search?q={code}"
    req = urllib.request.Request(url, headers={"User-Agent": "Mozilla/5.0"})
    with urllib.request.urlopen(req) as resp:
        d = json.loads(resp.read().decode())['data'][0]
        print(f"Code: {code} | Name: {d['name']} | photo_type: '{d.get('photo_type')}' | thumbnail: {d.get('image_url')}")
