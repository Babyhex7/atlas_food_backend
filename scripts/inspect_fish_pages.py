import pypdf

reader = pypdf.PdfReader(r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\Atlas Makananku Full + Pembatas_compressed.pdf")

pages = [
    (71, "LH-31 Patin"),
    (72, "LH-32 Kakap"),
    (73, "LH-34 Gurame"),
    (74, "LH-36 Bandeng"),
    (75, "LH-38 Kembung"),
    (76, "LH-39 Tenggiri"),
    (77, "LH-40 Tongkol"),
    (78, "LH-33 Nila"),
    (79, "LH-35 Lele"),
    (80, "LH-41 Bawal"),
]

for pnum, name in pages:
    page = reader.pages[pnum - 1]
    print("=" * 40)
    print(f"PAGE {pnum}: {name}")
    txt = page.extract_text().strip().replace("\n", " | ")
    print(f"Text: {txt}")
    content = page.get_contents().get_data().decode("latin-1")
    lines = content.split("\n")
    for i, line in enumerate(lines):
        if line.endswith(" Do") and line.startswith("/"):
            img_name = line.split()[0]
            cm_line = ""
            for j in range(max(0, i-6), i):
                if lines[j].endswith(" cm"):
                    cm_line = lines[j]
            print(f"  {img_name} -> cm: {cm_line}")
