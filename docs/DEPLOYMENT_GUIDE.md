# 🚀 Panduan Deployment Atlas Food (Vercel + Railway)

> **Stack Deployment:**  
> - **Frontend (Next.js):** Vercel  
> - **Backend (Go API & WebSocket):** Railway  
> - **Database (MySQL 8):** Railway MySQL Plugin / PlanetScale  

---

## 📑 Ringkasan Arsitektur Production

```
                  ┌──────────────────────────────┐
                  │    Responden / Admin (Browser)│
                  └──────────────┬───────────────┘
                                 │
                   HTTPS / WSS   │
            ┌────────────────────┴────────────────────┐
            │                                         │
            ▼                                         ▼
┌──────────────────────┐                  ┌──────────────────────┐
│   Vercel (Frontend)  │                  │  Railway (Backend)   │
│ Next.js App Router   │                  │ Go API + WebSocket   │
│ https://atlas-food.vercel.app           │ https://backend.up.railway.app
└──────────────────────┘                  └──────────┬───────────┘
                                                     │
                                                     ▼
                                          ┌──────────────────────┐
                                          │ Railway MySQL DB 8.0 │
                                          └──────────────────────┘
```

---

## 1. 🗄️ Langkah 1: Setup Database MySQL di Railway

1. Buka [Railway.app](https://railway.app) dan login dengan GitHub.
2. Buat **New Project** $\rightarrow$ Pilih **Provision MySQL**.
3. Setelah MySQL dibuat, klik pada service MySQL $\rightarrow$ Buka tab **Variables**.
4. Catat kredensial berikut untuk variabel backend:
   * `MYSQLHOST` (atau `MYSQL_HOST`)
   * `MYSQLPORT` (atau `MYSQL_PORT`)
   * `MYSQLUSER`
   * `MYSQLPASSWORD`
   * `MYSQLDATABASE`

---

## 2. ⚙️ Langkah 2: Deploy Backend Go di Railway

1. Di Project Railway yang sama, klik **+ New** $\rightarrow$ **GitHub Repo** $\rightarrow$ Pilih repository `atlas_food_backend`.
2. Railway akan mendeteksi `Dockerfile` yang telah kita buat secara otomatis.
3. Buka tab **Variables** di Service Backend Railway, lalu tambahkan Environment Variables:

```env
# Database (Hubungkan ke MySQL Railway)
DB_HOST=${{MySQL.MYSQLHOST}}
DB_PORT=${{MySQL.MYSQLPORT}}
DB_USER=${{MySQL.MYSQLUSER}}
DB_PASSWORD=${{MySQL.MYSQLPASSWORD}}
DB_NAME=${{MySQL.MYSQLDATABASE}}

# JWT Secret
JWT_SECRET=rahasia-super-aman-jwt-minimal-32-karakter-acak
JWT_EXPIRATION=24h
REFRESH_TOKEN_EXPIRATION=168h

# Server
SERVER_MODE=release

# Frontend Domain (URL Vercel kamu)
FRONTEND_URL=https://atlas-food.vercel.app

# Groq AI
GROQ_API_KEY=gsk_your_groq_api_key_here
GROQ_MODEL=llama-3.3-70b-versatile
```

4. Buka tab **Settings** $\rightarrow$ Di bagian **Networking**, klik **Generate Domain** (contoh URL: `https://atlas-food-backend-production.up.railway.app`).
5. Backend akan otomatis mengompilasi Go, mengeksekusi Auto Migration 17 tabel DB, dan seeding data awal makanan Indonesia!

---

## 3. 🌐 Langkah 3: Deploy Frontend Next.js di Vercel

1. Buka [Vercel.com](https://vercel.com) dan login dengan akun GitHub kamu.
2. Klik **Add New...** $\rightarrow$ **Project**.
3. Import repository `atlas_food_frontend`.
4. Di bagian **Environment Variables**, tambahkan:

```env
# API Backend di Railway (NEXT_PUBLIC_API_URL & NEXT_PUBLIC_API_BASE_URL didukung otomatis)
NEXT_PUBLIC_API_URL=https://atlas-food-backend-production.up.railway.app/api/v1

# WebSocket Backend di Railway (wajib gunakan wss:// untuk HTTPS)
NEXT_PUBLIC_WS_BASE_URL=wss://atlas-food-backend-production.up.railway.app/api/v1
```

5. Klik **Deploy**! Vercel akan mem-build aplikasi Next.js dalam 1-2 menit.

---

## 4. ✅ Verifikasi & Uji Coba Production

* **Health Check API:**  
  Buka `https://atlas-food-backend-production.up.railway.app/health`  
  *Harus mengembalikan: `{"service":"atlas-food-api","status":"ok"}`*

* **Login Admin:**  
  * Email: `admin@mail.com`  
  * Password: `Password123!`  

* **Fitur Real-Time Kolaborasi & WSS:**  
  Fitur WebSocket secara otomatis menggunakan `wss://` aman tanpa masalah Mixed Content SSL!
