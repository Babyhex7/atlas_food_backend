import urllib.request
import re
import json

folder_id = '1XlxkGAQrq64G3_94wSiMdDkW0ITSHRdY'
url = f'https://drive.google.com/drive/folders/{folder_id}'

req = urllib.request.Request(url, headers={'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)'})
html = urllib.request.urlopen(req).read().decode('utf-8', errors='ignore')

# Cari file ID dari file di Google Drive
file_ids = re.findall(r'\[\"([a-zA-Z0-9_-]{25,})\",\[\"(DSC[^\"]+\.jpg)\"', html)
print("Found files with IDs:", file_ids)
if not file_ids:
    # Coba pattern lain untuk file id di drive
    all_pairs = re.findall(r'\"(DSC\d+\.jpg)\".*?\"([a-zA-Z0-9_-]{28,35})\"', html)
    print("All pairs:", all_pairs)
