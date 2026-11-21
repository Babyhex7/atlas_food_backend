import pypdf

reader = pypdf.PdfReader(r"c:\Users\mybook_bagas\Projek_Agency\atlas_food\Atlas Makananku Full + Pembatas_compressed.pdf")
page = reader.pages[71] # page 72 (0-indexed 71)
print(page.extract_text())
