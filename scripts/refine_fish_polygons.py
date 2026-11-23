# Penyempurnaan polygon anotasi ikan (Patin LH-31, Kakap LH-32, Gurame LH-34).
# Tujuan: polygon mengikuti tepi daging ikan saja (tanpa piring/sendok/bayangan) berdasarkan
# segmentasi warna + GrabCut, lalu tiap piksel ikan dibagi ke potongan terdekat dari polygon lama
# (batas potongan kepala/badan/ekor tetap mengikuti garis potong yang sudah ada).
# Pemakaian: python refine_fish_polygons.py <live.json> <out_dir>
import cv2, numpy as np, json, os, sys

D = os.path.join(os.path.dirname(os.path.abspath(__file__)), "drive_downloads")
live = json.load(open(sys.argv[1]))
OUT = sys.argv[2]
os.makedirs(OUT, exist_ok=True)


def fish_mask(img, old_polys, spoon_x):
    h, w = img.shape[:2]
    lab = cv2.cvtColor(img, cv2.COLOR_BGR2LAB).astype(np.int32)
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    # Piksel "pasti ikan": berwarna hangat (kanal b LAB tinggi) atau cukup gelap.
    # Piring netral (b ~128) sehingga tidak ikut terambil.
    sure = ((lab[..., 2] - 128 > 9) | (gray < 150)).astype(np.uint8)
    # Batasi ke sekitar polygon lama (maks 14 px) agar bayangan/sendok jauh tidak terbawa
    union = np.zeros((h, w), np.uint8)
    for p in old_polys:
        cv2.fillPoly(union, [np.array(p, np.int32)], 1)
    near = cv2.dilate(union, cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (29, 29)))
    sure &= near
    sure[:, spoon_x:] = 0
    sure = cv2.morphologyEx(sure, cv2.MORPH_CLOSE, cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (7, 7)))
    # GrabCut: inti (erode) = pasti FG, sisanya di area 'near' = mungkin
    gc = np.full((h, w), cv2.GC_BGD, np.uint8)
    gc[near > 0] = cv2.GC_PR_BGD
    gc[sure > 0] = cv2.GC_PR_FGD
    core = cv2.erode(sure, np.ones((5, 5), np.uint8))
    gc[core > 0] = cv2.GC_FGD
    gc[:, spoon_x:] = cv2.GC_BGD
    bgd = np.zeros((1, 65), np.float64)
    fgd = np.zeros((1, 65), np.float64)
    cv2.grabCut(img, gc, None, bgd, fgd, 4, cv2.GC_INIT_WITH_MASK)
    m = ((gc == cv2.GC_FGD) | (gc == cv2.GC_PR_FGD)).astype(np.uint8)
    m = cv2.morphologyEx(m, cv2.MORPH_OPEN, np.ones((3, 3), np.uint8))
    m = cv2.morphologyEx(m, cv2.MORPH_CLOSE, cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (5, 5)))
    return m


def largest(m):
    n, lab, st, _ = cv2.connectedComponentsWithStats(m)
    if n <= 1:
        return m
    return (lab == 1 + np.argmax(st[1:, 4])).astype(np.uint8)


def to_poly(m, eps=0.9):
    cs, _ = cv2.findContours(m, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_NONE)
    c = max(cs, key=cv2.contourArea)
    a = cv2.approxPolyDP(c, eps, True)
    return [[int(p[0][0]), int(p[0][1])] for p in a]


result = {}
for fname, items in live.items():
    img = cv2.imread(os.path.join(D, fname))
    h, w = img.shape[:2]
    spoon_x = int(w * 0.70) if w > 800 else int(w * 0.72)
    # Kakap A: ekor menumpuk di leher sendok (sampai x~385). Batas sendok digeser ke kanan,
    # tetapi gagang/leher sendok di bawah ekor (y>157, x>=360) tetap diblok agar tidak ikut.
    block = None
    if fname == "LH-32_A.jpg":
        spoon_x = 388
        block = (360, 157)
    old = [p for _, p in items]
    M = fish_mask(img, old, spoon_x)
    if block:
        M[block[1]:, block[0]:] = 0
        # Leher sendok yang menonjol tepat di atas ekor (y<133, x>=370) juga dibuang
        M[:133, 370:] = 0
    # Jarak tiap piksel ke tiap polygon lama -> pembagian potongan
    dists = []
    for p in old:
        f = np.ones((h, w), np.uint8)
        cv2.fillPoly(f, [np.array(p, np.int32)], 0)
        dists.append(cv2.distanceTransform(f, cv2.DIST_L2, 3))
    dists = np.stack(dists)
    lab = np.argmin(dists, axis=0)
    out = []
    for k, (aid, p) in enumerate(items):
        pm = ((lab == k) & (M > 0)).astype(np.uint8) if len(old) > 1 else M
        pm = cv2.morphologyEx(pm, cv2.MORPH_OPEN, np.ones((3, 3), np.uint8))
        pm = largest(pm)
        poly = to_poly(pm)
        out.append((aid, poly))
        old_m = np.zeros((h, w), np.uint8)
        cv2.fillPoly(old_m, [np.array(p, np.int32)], 1)
        new_m = np.zeros((h, w), np.uint8)
        cv2.fillPoly(new_m, [np.array(poly, np.int32)], 1)
        print(fname, aid[:8], "old px", int(old_m.sum()), "new px", int(new_m.sum()), "pts", len(poly))
    result[fname] = out
    vis = img.copy()
    cols = [(0, 0, 255), (0, 160, 0), (255, 0, 0)]
    for k, (aid, poly) in enumerate(out):
        cv2.polylines(vis, [np.array(poly, np.int32)], True, cols[k % 3], 1)
    cv2.imwrite(os.path.join(OUT, "new_" + fname), vis)
json.dump(result, open(os.path.join(OUT, "new.json"), "w"))
