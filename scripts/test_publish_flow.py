import urllib.request
import json

# Login
login_payload = json.dumps({'email': 'admin@atlasmakananku.id', 'password': 'Password123!'}).encode('utf-8')
req = urllib.request.Request('https://atlasmakananku.id/api/v1/auth/login', data=login_payload, headers={'Content-Type': 'application/json'})
try:
    with urllib.request.urlopen(req) as resp:
        token = json.loads(resp.read().decode())['data']['access_token']
except Exception:
    login_payload = json.dumps({'email': 'admin@atlasmakananku.id', 'password': '12345678'}).encode('utf-8')
    req = urllib.request.Request('https://atlasmakananku.id/api/v1/auth/login', data=login_payload, headers={'Content-Type': 'application/json'})
    with urllib.request.urlopen(req) as resp:
        token = json.loads(resp.read().decode())['data']['access_token']

# Generate
print("1. Testing AI Generator...")
gen_payload = json.dumps({'food_name': 'Karedok Sayur Segar'}).encode('utf-8')
req_gen = urllib.request.Request('https://atlasmakananku.id/api/v1/admin/ai/generate-food', data=gen_payload, headers={'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token})
with urllib.request.urlopen(req_gen) as resp:
    draft = json.loads(resp.read().decode())['data']
    print('Generated draft:', draft['food_name'], '| Category:', draft['category_name'], '| Kalori:', draft['calories_per_100g'])

# Publish test
print("2. Testing Publish to Catalog...")
pub_payload = json.dumps({
    'food_name': draft['food_name'],
    'local_name': draft.get('local_name', ''),
    'description': draft.get('description', ''),
    'category_id': draft.get('category_id'),
    'calories_per_100g': draft.get('calories_per_100g', 0),
    'protein_per_100g': draft.get('protein_per_100g', 0),
    'fat_per_100g': draft.get('fat_per_100g', 0),
    'carbs_per_100g': draft.get('carbs_per_100g', 0),
    'portions': draft.get('portions', [])
}).encode('utf-8')

req_pub = urllib.request.Request('https://atlasmakananku.id/api/v1/admin/ai/publish-food', data=pub_payload, headers={'Content-Type': 'application/json', 'Authorization': 'Bearer ' + token})
with urllib.request.urlopen(req_pub) as resp:
    pub_res = json.loads(resp.read().decode())['data']
    print('Publish success!')
    print('Food ID:', pub_res['food_id'], '| Code:', pub_res['food_code'])
    print('Message:', pub_res['message'])
    food_id = pub_res['food_id']

# Verify search in public foods
print("3. Verifying Public Food Search...")
req_search = urllib.request.Request('https://atlasmakananku.id/api/v1/public/foods/search?q=Karedok')
with urllib.request.urlopen(req_search) as resp:
    search_res = json.loads(resp.read().decode())['data']
    foods = search_res if isinstance(search_res, list) else search_res.get('foods', [])
    print(f'Search result count: {len(foods)}')
    for f in foods:
        print(f"  * {f.get('name')} | Code: {f.get('code')} | AI Generated: {f.get('is_ai_generated')}")

# Verify detail in public food detail endpoint
print("4. Verifying Public Food Detail...")
req_detail = urllib.request.Request(f'https://atlasmakananku.id/api/v1/public/foods/{food_id}')
with urllib.request.urlopen(req_detail) as resp:
    food_detail = json.loads(resp.read().decode())['data']
    print(f"Detail Food Name: {food_detail.get('name')} | Photo Type: {food_detail.get('photo_type')}")
    portions = food_detail.get('portions', [])
    print(f"Portions count: {len(portions)}")
    for p in portions:
        print(f"  * Porsi {p.get('label')}: {p.get('weight_gram')}g")
