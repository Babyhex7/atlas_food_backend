import urllib.request

req = urllib.request.Request("https://atlasmakananku.id/find-food", headers={"User-Agent": "Mozilla/5.0"})
with urllib.request.urlopen(req) as resp:
    html = resp.read().decode('utf-8')
    print(f"Status: {resp.status}, length: {len(html)}")
    print("Has NEXT_DATA:", "__NEXT_DATA__" in html)
