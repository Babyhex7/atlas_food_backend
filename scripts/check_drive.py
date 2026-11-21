import urllib.request
import re
import json

url = 'https://drive.google.com/drive/folders/1XlxkGAQrq64G3_94wSiMdDkW0ITSHRdY'
req = urllib.request.Request(url, headers={'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)'})
try:
    with urllib.request.urlopen(req) as resp:
        html = resp.read().decode('utf-8', errors='ignore')
        title = re.search(r'<title>(.*?)</title>', html)
        print('Title:', title.group(1) if title else 'None')
        # Temukan strings yang mirip file foto
        matches = re.findall(r'[\w\-_\s\(\)]+\.(?:jpg|jpeg|png|webp)', html, re.IGNORECASE)
        print('Matches count:', len(matches))
        print('Unique matches:', sorted(list(set(matches))))
except Exception as e:
    print('Error:', e)
