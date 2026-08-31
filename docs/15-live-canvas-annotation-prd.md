# 🎨 Product Requirement Document (PRD) & Technical Specification
## Fitur: Live Canvas Annotation (Coret-Coret Layar Real-Time)

> **Project:** Atlas Food — Food Recall Survey Platform  
> **Module:** Real-Time Collaboration (`internal/domain/collab`)  
> **Status:** Draft / Technical Specification  
> **Author:** Antigravity AI Engineering Team  
> **Branch:** `feature/live-canvas-annotation`  

---

## 📋 Table of Contents
1. [Ringkasan Eksekutif & Tujuan](#1-ringkasan-eksekutif--tujuan)
2. [Skenario Pengguna & Workflow UX/UI](#2-skenario-pengguna--workflow-uxui)
3. [Arsitektur & Diagram Alur Data](#3-arsitektur--diagram-alur-data)
4. [Normalisasi Koordinat (Cross-Device Scaling)](#4-normalisasi-koordinat-cross-device-scaling)
5. [Spesifikasi Protokol WebSocket](#5-spesifikasi-protokol-websocket)
6. [Manajemen State & Siklus Hidup Canvas](#6-manajemen-state--siklus-hidup-canvas)
7. [Teknik Optimasi Kinerja & Throughput](#7-teknik-optimasi-kinerja--throughput)
8. [Keamanan & Kontrol Akses Berperan (RBAC)](#8-keamanan--kontrol-akses-berperan-rbac)
9. [Handling Edge Cases & Graceful Degradation](#9-handling-edge-cases--graceful-degradation)
10. [Rencana Implementasi Backend (Go)](#10-rencana-implementasi-backend-go)
11. [Rencana Implementasi Frontend (Next.js / React)](#11-rencana-implementasi-frontend-nextjs--react)
12. [Kriteria Penerimaan & Plan Pengujian](#12-kriteria-penerimaan--plan-pengujian)

---

## 1. Ringkasan Eksekutif & Tujuan

### 1.1 Masalah
Saat sesi recall berpendamping jarak jauh, enumerator/ahli gizi sering kali perlu menunjukkan area porsi spesifik pada foto porsi *As Served* (misalnya: *"Porsi nasi yang dimakan seberapa piring?", "Telurnya bagian putihnya saja atau utuh?"*). Tanpa alat bantu visual langsung, komunikasi verbal saja sering kali membingungkan responden dan memicu bias estimasi porsi.

### 1.2 Solusi
Fitur **Live Canvas Annotation** memungkinkan Enumerator (atau pengguna berizin `editor`/`owner`) menggambar, melingkari, menunjuk dengan *laser pointer*, atau memberi catatan transparan secara *real-time* di atas canvas foto porsi yang sedang dibuka di layar responden.

### 1.3 Tujuan Utama
* **Akurasi Estimasi Porsi:** Menurunkan kesalahan komunikasi visual antara enumerator dan responden.
* **Low Latency:** Latensi coretan $< 100\text{ ms}$ antar perangkat.
* **Responsive Scaling:** Coretan yang dibuat di monitor Desktop ($1920\times 1080$) harus tampil di posisi yang persis sama pada smartphone Responden ($390\times 844$).
* **Ringan & Non-blocking:** Tidak mengganggu kinerja memori peramban atau jalur pengisian survey.

---

## 2. Skenario Pengguna & Workflow UX/UI

### 2.1 Peran Pengguna
* **Enumerator / Pendamping (`editor` / `owner`):** Memiliki *Toolbar Canvas* (Pensil, Laser Pointer, Lingkaran/Kotak, Penghapus, Warna, Clear All).
* **Responden (`owner` / `editor` / `viewer`):** Melihat coretan secara *real-time* di atas gambar foto porsi yang sedang aktif. Jika bertindak sebagai `viewer`, kontrol menggambar dinonaktifkan (*fail-closed*).

### 2.2 UI Toolbar Component (`CanvasToolbar.tsx`)
Toolbar melayang (*floating widget*) yang muncul hanya saat ruang kolaborasi aktif dan mode canvas dinyalakan:
* 🖌️ **Pencil Tool:** Goresan bebas (*freehand drawing*).
* 🔴 **Laser Pointer:** Titik sorot sementara yang memudar otomatis dalam 1.5 detik (*vanishing trail*).
* ⭕ **Shape Tool:** Lingkaran/Oval & Persegi Panjang untuk menunjuk area makanan.
* 🎨 **Color Picker:** Pilihan warna kontras tinggi (Merah Kritis, Kuning Sorot, Hijau Petunjuk, Biru).
* 🗑️ **Clear Canvas:** Menghapus seluruh coretan di ruang aktif.
* 👁️ **Toggle Canvas Visibility:** Menyembunyikan/menampilkan layer coretan tanpa menghapus state.

---

## 3. Arsitektur & Diagram Alur Data

```
┌─────────────────────────┐               ┌─────────────────────────┐
│ Klien Enumerator (Goes) │               │ Klien Responden (Views) │
│  - Mouse / Touch Move   │               │  - Receives Normalized  │
│  - Local Render         │               │    Coordinates          │
│  - Throttled Sync       │               │  - Denormalize & Draw   │
└────────────┬────────────┘               └────────────▲────────────┘
             │                                         │
             │ WS Payload (Normalized x,y: 0.0-1.0)    │ WS Broadcast
             ▼                                         │
┌──────────────────────────────────────────────────────┴────────────┐
│                    Go WebSocket Hub Server                        │
│  1. Check Client Role (canEdit() == true?)                       │
│  2. Room Broadcast Manager (Coalescing & Rate Limit)             │
│  3. Canvas Ring Buffer (Simpan 50 stroke terakhir untuk latejoin)│
└───────────────────────────────────────────────────────────────────┘
```

---

## 4. Normalisasi Koordinat (Cross-Device Scaling)

### 4.1 Tantangan
Perangkat Enumerator dan Responden memiliki ukuran layar, rasio aspek, dan resolusi pixel (DPR - *Device Pixel Ratio*) yang berbeda.

### 4.2 Solusi: Koordinat Relatif $0.0 \rightarrow 1.0$
Seluruh koordinat $(X, Y)$ titik coretan dipetakan ke rentang persentase relatif terhadap elemen gambar target (`ImageContainer`):

$$x_{norm} = \frac{x_{pixel}}{Width_{container}}$$

$$y_{norm} = \frac{y_{pixel}}{Height_{container}}$$

Saat klien penerima menerima koordinat $(x_{norm}, y_{norm})$, klien melakukan denormalisasi sesuai ukuran kontainer gambar lokalnya:

$$x_{local} = x_{norm} \times Width_{local\_container}$$

$$y_{local} = y_{norm} \times Height_{local\_container}$$

---

## 5. Spesifikasi Protokol WebSocket

### 5.1 Pesan dari Klien ke Server (K $\rightarrow$ S)

#### A. Inisiasi Stroke / Coretan (`canvas_draw_start`)
```json
{
  "type": "canvas_draw_start",
  "payload": {
    "stroke_id": "str_9a8b7c6d",
    "tool": "pencil",
    "color": "#EF4444",
    "width": 3,
    "x": 0.4521,
    "y": 0.3120,
    "target_image_id": "img_as_served_102"
  }
}
```

#### B. Pergerakan / Penambahan Titik (`canvas_draw_move`)
```json
{
  "type": "canvas_draw_move",
  "payload": {
    "stroke_id": "str_9a8b7c6d",
    "points": [
      [0.4530, 0.3135],
      [0.4545, 0.3150]
    ]
  }
}
```

#### C. Penutupan Stroke (`canvas_draw_end`)
```json
{
  "type": "canvas_draw_end",
  "payload": {
    "stroke_id": "str_9a8b7c6d"
  }
}
```

#### D. Laser Pointer Update (`canvas_laser_move`)
```json
{
  "type": "canvas_laser_move",
  "payload": {
    "x": 0.6120,
    "y": 0.4510,
    "color": "#FF0055"
  }
}
```

#### E. Hapus Seluruh Coretan (`canvas_clear`)
```json
{
  "type": "canvas_clear",
  "payload": {
    "target_image_id": "img_as_served_102"
  }
}
```

---

### 5.2 Pesan Siaran dari Server ke Klien (S $\rightarrow$ K)
Server memvalidasi role pengirim (`canEdit() == true`), lalu menyiarkan pesan ke seluruh peserta di ruang yang sama:

* `canvas_stroke_started`
* `canvas_stroke_updated`
* `canvas_stroke_ended`
* `canvas_laser_updated`
* `canvas_cleared`
* `canvas_state_sync` *(dikirim ke klien baru yang bergabung untuk merender coretan yang sudah ada)*

---

## 6. Manajemen State & Siklus Hidup Canvas

### 6.1 State Model (`canvasStore.ts`)
```typescript
export interface Point {
  x: number; // 0.0 - 1.0
  y: number; // 0.0 - 1.0
}

export interface CanvasStroke {
  id: string;
  userId: string;
  tool: 'pencil' | 'circle' | 'rectangle' | 'laser';
  color: string;
  width: number;
  points: Point[];
  targetImageId?: string;
  createdAt: number;
}

export interface CanvasState {
  activeStrokes: Map<string, CanvasStroke>;
  laserPoint: (Point & { userId: string; color: string; timestamp: number }) | null;
  isVisible: boolean;
  activeTool: 'pencil' | 'circle' | 'rectangle' | 'laser' | 'eraser';
  activeColor: string;
  activeWidth: number;
}
```

---

## 7. Teknik Optimasi Kinerja & Throughput

Untuk mencegah lonjakan pesan WebSocket dan lag UI (*jank*):

1. **Point Throttling & Batching (Client Side):**
   * Peristiwa `mousemove` / `touchmove` dipicu hingga 120 FPS.
   * Klien mengumpulkan titik setiap **16 ms (60 FPS)** dan mengirim batch `points` dalam 1 payload `canvas_draw_move`.
2. **Bezier Curve / Smooth Interpolation (Render Side):**
   * Menggunakan interpolasi *Quadratic Curve* (`quadraticCurveTo`) saat merender garis agar garis tidak patah-patah meskipun titik yang dikirim berjarak.
3. **Offscreen Canvas Rendering:**
   * Coretan stasioner yang sudah selesai (`canvas_draw_end`) dirender ke *offscreen buffer layer* agar canvas utama tidak perlu melakukan *redraw* seluruh stroke dari awal pada tiap frame.
4. **Ring Buffer History di Server (Late-Joiner Sync):**
   * Server menyimpan **50 stroke terakhir** per `room` di memori (`in-memory ring buffer`). Klien yang baru masuk/rekonek langsung menerima `canvas_state_sync` tanpa perlu request ulang.

---

## 8. Keamanan & Kontrol Akses Berperan (RBAC)

Penegakan kontrol keamanan dilakukan dengan prinsip **Fail-Closed 3 Lapis**:

1. **Lapis 1 (UI Level):** Toolbar canvas disembunyikan total jika `role === 'viewer'`.
2. **Lapis 2 (Client Level):** Function `sendCanvasMessage()` memblokir event menggambar jika status role di Zustand belum terverifikasi sebagai `editor`/`owner`.
3. **Lapis 3 (Server Level Go Handler):**
   ```go
   func (c *Client) handleCanvasDraw(msg Message) {
       if !c.canEdit() {
           c.sendError("UNAUTHORIZED", "Viewer role cannot draw on canvas")
           return
       }
       // Broadcast message
   }
   ```

---

## 9. Handling Edge Cases & Graceful Degradation

| Edge Case | Penanganan / Recovery |
| --- | --- |
| **Koneksi Terputus di Tengah Menggambar** | `canvas_draw_end` sintetis otomatis dipicu oleh klien untuk menutup stroke aktif agar garis tidak tersambung liar. |
| **Responden Mengganti Foto Porsi** | Event ganti foto memicu `canvas_clear` lokal dan server untuk menyesuaikan canvas dengan foto baru. |
| **Layar Diputar (Orientation Change)** | Canvas di-resize secara otomatis; fungsi denormalisasi menyesuaikan skala koordinat sesuai rasio kontainer baru. |
| **High Latency / Slow Network** | Render lokal instan di sisi pengirim (*Optimistic Local Rendering*). Pengikut merender pesan saat tiba tanpa memblokir thread. |

---

## 10. Rencana Implementasi Backend (Go)

### File yang Ditambahkan / Dimodifikasi:
1. `internal/domain/collab/dto.go`
   * Menambahkan struct pesan: `CanvasDrawPayload`, `CanvasStrokeMessage`, `CanvasSyncPayload`.
2. `internal/domain/collab/room.go`
   * Menambahkan `canvasHistory` ring buffer di struct `Room`.
3. `internal/domain/collab/client.go`
   * Menambahkan handler WebSocket event `canvas_*` dengan validasi `c.canEdit()`.

---

## 11. Rencana Implementasi Frontend (Next.js / React)

### File yang Ditambahkan / Dimodifikasi:
1. `internal/domain/collab/store/canvasStore.ts` *(Baru)*
   * Zustand store untuk mengelola state alat canvas, koordinat stroke, laser pointer, dan sejarah menggambar.
2. `internal/domain/collab/hooks/useLiveCanvas.ts` *(Baru)*
   * Custom hook penanganan mouse/touch event, normalisasi koordinat, throttling 16ms, dan integrasi WebSocket.
3. `internal/domain/collab/components/LiveCanvasOverlay.tsx` *(Baru)*
   * Komponen Canvas transparan (`<canvas>`) di atas foto porsi makanan dengan auto-resize observer.
4. `internal/domain/collab/components/CanvasToolbar.tsx` *(Baru)*
   * Toolbar floating alat gambar (Pencil, Laser, Color, Clear) khusus untuk `editor`/`owner`.

---

## 12. Kriteria Penerimaan & Plan Pengujian

### 12.1 Acceptance Criteria
* [ ] Enumerator dapat mengaktifkan mode gambar dan melingkari area makanan pada foto porsi.
* [ ] Coretan muncul di layar responden dalam waktu $< 100\text{ ms}$.
* [ ] Coretan yang dibuat di Desktop tampil secara presisi di smartphone tanpa tergeser.
* [ ] Pengguna berizin `viewer` **dilarang keras** menggambar (tombol disembunyikan & pesan diblokir server).
* [ ] Tombol *Clear All* menghapus coretan di semua peramban peserta secara serentak.
* [ ] Laser pointer memudar otomatis dalam 1.5 detik.

### 12.2 Test Plan
1. **Unit Test (FE):** Testing utility fungsi normalisasi `normalizeCoordinates()` dan `denormalizeCoordinates()`.
2. **Integration Test (BE):** Verification websocket handler Go menolak payload `canvas_draw_start` dari client dengan role `viewer`.
3. **End-to-End Test (Multi-browser):** Uji coba 2 peramban bersisian (Desktop vs Mobile Responsive Viewport).
