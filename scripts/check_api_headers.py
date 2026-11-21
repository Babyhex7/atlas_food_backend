import urllib.request

urls = [
    "https://api.atlasmakananku.id/api/v1/public/food-images/c4e340a1-7341-4567-89ab-lh34cut00001",
    "https://api.atlasmakananku.id/api/v1/public/foods/8ed0845f-ce26-11f1-b499-e2699381f9eb",
    "https://api.atlasmakananku.id/api/v1/public/foods/search?q=LH-32",
]

for u in urls:
    req = urllib.request.Request(u, headers={"User-Agent": "Mozilla/5.0"})
    with urllib.request.urlopen(req) as resp:
        print(f"=== {u} ===")
        for k, v in resp.headers.items():
            if k.lower() in ['cache-control', 'etag', 'expires', 'last-modified', 'cf-cache-status', 'cf-ray', 'age']:
                print(f"  {k}: {v}")
    print()
