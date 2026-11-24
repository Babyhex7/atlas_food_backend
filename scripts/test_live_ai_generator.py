import urllib.request
import json

login_payload = json.dumps({'email': 'admin@atlasmakananku.id', 'password': 'Password123!'}).encode('utf-8')
req = urllib.request.Request('https://atlasmakananku.id/api/v1/auth/login', data=login_payload, headers={'Content-Type': 'application/json'})
try:
    with urllib.request.urlopen(req) as resp:
        res = json.loads(resp.read().decode())
        print('Login success! User:', res.get('data', {}).get('user', {}).get('name'))
        token = res.get('data', {}).get('access_token')
except Exception as e:
    print('Trying password 12345678...')
    login_payload = json.dumps({'email': 'admin@atlasmakananku.id', 'password': '12345678'}).encode('utf-8')
    req = urllib.request.Request('https://atlasmakananku.id/api/v1/auth/login', data=login_payload, headers={'Content-Type': 'application/json'})
    with urllib.request.urlopen(req) as resp:
        res = json.loads(resp.read().decode())
        print('Login success! User:', res.get('data', {}).get('user', {}).get('name'))
        token = res.get('data', {}).get('access_token')

# 1. Test missing-foods
print('\n--- 1. Testing GET /admin/ai/missing-foods ---')
req_mf = urllib.request.Request('https://atlasmakananku.id/api/v1/admin/ai/missing-foods', headers={'Authorization': 'Bearer ' + token})
with urllib.request.urlopen(req_mf) as resp:
    mf_res = json.loads(resp.read().decode())
    items = mf_res.get('data', [])
    print(f'Total missing foods reported: {len(items)}')
    for it in items[:5]:
        print(f"  - {it.get('name')} (Laporan: {it.get('report_count')}, In Catalog: {it.get('is_in_catalog')})")

# 2. Test generate-food
print('\n--- 2. Testing POST /admin/ai/generate-food ---')
gen_payload = json.dumps({
    'food_name': 'Seblak Ceker Pedas',
    'description': 'Seblak khas Sunda dengan kerupuk kenyal basah dan ceker ayam'
}).encode('utf-8')
req_gen = urllib.request.Request(
    'https://atlasmakananku.id/api/v1/admin/ai/generate-food',
    data=gen_payload,
    headers={'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token}
)

with urllib.request.urlopen(req_gen) as resp:
    gen_res = json.loads(resp.read().decode())
    data = gen_res.get('data', {})
    print(f"Status Code: {resp.status}")
    print(f"Food Name: {data.get('food_name')}")
    print(f"Category: {data.get('category_name')} (ID: {data.get('category_id')})")
    print(f"Kalori per 100g: {data.get('calories_per_100g')} kcal")
    print(f"Protein: {data.get('protein_per_100g')}g, Lemak: {data.get('fat_per_100g')}g, Karbo: {data.get('carbs_per_100g')}g")
    portions = data.get('portions', [])
    print(f"Total Porsi Generated: {len(portions)}")
    for p in portions:
        img_preview = p.get('image_url', '')[:50] + '...' if len(p.get('image_url', '')) > 50 else p.get('image_url', '')
        print(f"  * Porsi {p.get('code')} ({p.get('label')}): {p.get('weight_gram')}g | Img: {img_preview}")
