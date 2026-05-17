# 👤 User Stories

## Admin Stories

### Auth Management

**Story 1: Admin Login**

- **Sebagai** admin
- **Saya ingin** login ke dashboard
- **Sehingga** saya bisa mengelola survey dan data makanan

**Acceptance Criteria:**

- Admin dapat login dengan email dan password
- Setelah login, admin mendapatkan JWT token
- Token valid selama 24 jam

---

### Survey Management

**Story 2: Create Survey**

- **Sebagai** admin
- **Saya ingin** membuat survey baru
- **Sehingga** saya bisa mengumpulkan data recall makanan dari responden

**Acceptance Criteria:**

- Admin dapat input nama survey, deskripsi, dan periode
- Admin dapat konfigurasi meals (Sarapan, Makan Siang, dll)
- Admin dapat set prompts/welcome message
- System generate access token unik untuk survey
- Survey dibuat dengan status "draft"

---

**Story 3: Share Survey Link**

- **Sebagai** admin
- **Saya ingin** membagikan link survey ke responden
- **Sehingga** responden bisa mengakses survey setelah login

**Acceptance Criteria:**

- Admin dapat copy access token/link survey
- Link bisa di-share via WhatsApp, email, dll
- Responden login terlebih dahulu sebelum akses link
- Survey bisa di-set "active" untuk dibagikan

---

**Story 4: View Submissions**

- **Sebagai** admin
- **Saya ingin** melihat hasil submission dari responden
- **Sehingga** saya bisa menganalisis data gizi

**Acceptance Criteria:**

- Admin dapat lihat list submissions per survey
- Tampilkan jumlah responden dan completion rate
- Bisa filter berdasarkan tanggal
- Bisa export data ke CSV/JSON

---

### Food Database Management

**Story 5: Add Food Item**

- **Sebagai** admin
- **Saya ingin** menambahkan makanan baru ke database
- **Sehingga** responden bisa memilih makanan tersebut di survey

**Acceptance Criteria:**

- Admin dapat input kode, nama, nama lokal, deskripsi
- Admin dapat pilih kategori makanan
- Admin dapat input nilai gizi per 100g
- Makanan langsung searchable setelah ditambahkan

---

**Story 6: Add Portion Images**

- **Sebagai** admin
- **Saya ingin** upload gambar porsi untuk makanan
- **Sehingga** responden bisa pilih porsi dengan visual yang jelas

**Acceptance Criteria:**

- Admin dapat buat "As Served Set" untuk makanan
- Admin dapat upload multiple images (1 porsi, 2 porsi, dll)
- Setiap image di-set berat dalam gram
- System generate thumbnail otomatis
- Images tampil di UI responden

---

**Story 7: Set Portion Methods**

- **Sebagai** admin
- **Saya ingin** konfigurasi metode portion per makanan
- **Sehingga** responden bisa pilih cara yang paling sesuai

**Acceptance Criteria:**

- Admin dapat pilih selection type (simple_grid, as_served_quantity, counter, input)
- Admin dapat set allow fractions (ya/tidak)
- Admin dapat set max quantity
- Admin dapat link ke As Served Set

---

## Respondent Stories

### Survey Access

**Story 8: Access Survey via Link**

- **Sebagai** responden
- **Saya ingin** mengakses survey dengan link yang diberikan
- **Sehingga** saya bisa langsung isi survey setelah login

**Acceptance Criteria:**

- Responden login terlebih dahulu
- Responden bisa klik link dan masuk ke survey
- Sistem membuat participant untuk user yang login
- Responden bisa isi alias/nama tampilan (opsional)

---

### Food Selection

**Story 9: Search Food**

- **Sebagai** responden
- **Saya ingin** mencari makanan yang saya makan
- **Sehingga** saya bisa menambahkannya ke recall

**Acceptance Criteria:**

- Search dengan minimal 3 karakter
- Hasil search muncul dalam < 500ms
- Bisa filter berdasarkan kategori
- Tampilkan emoji/icon per kategori

---

**Story 10: Add Food to Meal**

- **Sebagai** responden
- **Saya ingin** menambahkan makanan ke waktu makan tertentu
- **Sehingga** saya bisa membuat daftar makanan per meal

**Acceptance Criteria:**

- Klik makanan dari hasil search untuk add
- Makanan muncul di list "akan di-portion-kan"
- Bisa hapus makanan dari list
- Bisa tambah multiple makanan sebelum ke portion selection

---

### Portion Selection

**Story 11: Select Portion with Images (Simple Grid)**

- **Sebagai** responden
- **Saya ingin** pilih porsi dengan melihat gambar
- **Sehingga** saya bisa estimasi porsi yang saya makan dengan akurat

**Acceptance Criteria:**

- Tampilkan grid gambar porsi (½, 1, 1½, 2, 3 porsi)
- Setiap gambar ada label dan estimasi gram
- Bisa klik gambar untuk select
- Bisa input manual gram sebagai alternatif

---

**Story 12: Select Portion with Quantity (As Served)**

- **Sebagai** responden
- **Saya ingin** pilih porsi dan tentukan jumlahnya
- **Sehingga** saya bisa input "3 setengah pisang"

**Acceptance Criteria:**

- Tampilkan gambar besar porsi yang dipilih
- Ada counter untuk WHOLE (bilangan bulat)
- Ada counter untuk FRACTION (¼, ½, ¾)
- Tampilkan kalkulasi real-time: "3 and ¼ of the largest portion (617.50g)"
- Thumbnail row untuk pilih gambar lain

---

### Review & Submit

**Story 13: Review Meal Data**

- **Sebagai** responden
- **Saya ingin** review semua makanan yang sudah saya input
- **Sehingga** saya bisa pastikan data sudah benar sebelum submit

**Acceptance Criteria:**

- Tampilkan list per waktu makan (Sarapan, Makan Siang, dll)
- Setiap item: nama makanan + porsi gram
- Bisa edit porsi dari halaman review
- Bisa hapus makanan dari halaman review
- Bisa tambah makanan lagi

---

**Story 14: Submit Survey**

- **Sebagai** responden
- **Saya ingin** submit survey setelah selesai mengisi
- **Sehingga** data saya tersimpan dan bisa dianalisis admin

**Acceptance Criteria:**

- Tombol SUBMIT hanya aktif setelah minimal 1 meal diisi
- Setelah submit, tampilkan confirmation message
- Responden tidak bisa edit setelah submit
- Admin langsung bisa lihat submission

---

## System Stories

**Story 15: Calculate Nutrition**

- **Sebagai** system
- **Saya ingin** menghitung nutrisi berdasarkan porsi yang dipilih
- **Sehingga** data gizi akurat per submission

**Acceptance Criteria:**

- Ambil nilai gizi per 100g dari tabel
- Kalkulasi: (nutrient_value / 100) × portionGram
- Bulatkan ke 1 desimal
- Simpan hasil di submission JSON

---

**Story 16: Export Data**

- **Sebagai** system
- **Saya ingin** export submission data ke CSV
- **Sehingga** admin bisa analisis di Excel/SPSS

**Acceptance Criteria:**

- Export semua submissions per survey
- Format: submission_id, respondent_name, meal_name, food_name, portion_gram, nutrients...
- Bisa filter berdasarkan tanggal sebelum export
- Download sebagai file .csv
