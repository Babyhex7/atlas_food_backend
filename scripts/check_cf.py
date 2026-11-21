import urllib.request

for host in ["https://atlasmakananku.id", "https://api.atlasmakananku.id"]:
    req = urllib.request.Request(host, headers={"User-Agent": "Mozilla/5.0"})
    try:
        with urllib.request.urlopen(req) as resp:
            headers = dict(resp.headers)
            print(f"=== {host} ===")
            print(f"  Server: {headers.get('server')}")
            print(f"  CF-RAY: {headers.get('cf-ray')}")
            print(f"  CF-Cache-Status: {headers.get('cf-cache-status')}")
    except Exception as e:
        print(f"{host}: {e}")
