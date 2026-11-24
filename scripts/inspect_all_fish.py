import urllib.request
import json

req = urllib.request.urlopen('https://atlasmakananku.id/api/v1/public/foods/search?q=ikan')
foods = json.loads(req.read().decode())['data']

for f in foods:
    code = f['code']
    fid = f['id']
    name = f['name']
    ptype = f.get('photo_type')
    
    try:
        dreq = urllib.request.urlopen(f'https://atlasmakananku.id/api/v1/public/foods/{fid}')
        detail = json.loads(dreq.read().decode())['data']
        photos = detail.get('portion_photos', [])
        print(f"=== {code}: {name} (type: {ptype}) ===")
        print(f"  Food image_url: {detail.get('image_url')}")
        print(f"  Portion photos count: {len(photos)}")
        seen_images = set()
        for p in photos:
            print(f"    Photo {p.get('label')}: {p.get('weight_gram')}g, url={p.get('image_url')}, food_img_id={p.get('food_image_id')}")
            fimg_id = p.get('food_image_id')
            if fimg_id and fimg_id not in seen_images:
                seen_images.add(fimg_id)
                try:
                    areq = urllib.request.urlopen(f'https://atlasmakananku.id/api/v1/public/food-images/{fimg_id}')
                    adata = json.loads(areq.read().decode())['data']
                    areas = adata.get('areas', [])
                    print(f"      -> FoodImage {fimg_id}: {len(areas)} areas, img_url={adata.get('image_url')}")
                    for a in areas:
                        coords = a.get('coordinates', [])
                        pts_summary = f"{len(coords)} pts"
                        if coords:
                            pts_summary += f" [x0={coords[0].get('x')}, y0={coords[0].get('y')} ... x_last={coords[-1].get('x')}, y_last={coords[-1].get('y')}]"
                        print(f"         Area '{a.get('name')}': {a.get('weight_gram')}g, {pts_summary}")
                except Exception as e:
                    print(f"      -> Error fetching food-image {fimg_id}: {e}")
    except Exception as e:
        print(f"Error getting detail for {code}: {e}")
