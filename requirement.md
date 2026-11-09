# 📋 REQUIREMENT DOCUMENT — PROJECT-STUDIO
## Aplikasi Manajemen Proyek (Terinspirasi dari GanttPRO)

**Versi:** 1.0.0  
**Tanggal:** 2026-09-19  
**Project:** PROJECT-STUDIO  
**Referensi:** [app.ganttpro.com](https://app.ganttpro.com)

---

## 1. 📌 RINGKASAN PROYEK

PROJECT-STUDIO adalah aplikasi manajemen proyek berbasis web yang memungkinkan tim untuk merencanakan, menjadwalkan, dan melacak proyek secara visual menggunakan Gantt Chart interaktif. Aplikasi ini dirancang untuk meningkatkan produktivitas tim, transparansi progres, dan efisiensi alokasi sumber daya.

---

## 2. 🎯 TUJUAN BISNIS

| No | Tujuan |
|----|--------|
| 1 | Memungkinkan manajer proyek membuat dan mengelola rencana proyek secara visual |
| 2 | Memberikan visibilitas real-time terhadap progres proyek kepada seluruh anggota tim |
| 3 | Mengoptimalkan alokasi sumber daya dan menghindari beban kerja berlebih |
| 4 | Mendukung kolaborasi tim secara efektif dengan fitur komunikasi terintegrasi |
| 5 | Menghasilkan laporan dan analitik proyek yang akurat |

---

## 3. 👥 PENGGUNA (USER ROLES)

### 3.1 Owner / Admin
- Membuat dan mengelola workspace
- Mengelola semua proyek dan anggota
- Mengakses semua fitur termasuk billing dan pengaturan
- Menetapkan peran pengguna

### 3.2 Project Manager
- Membuat, mengedit, dan menghapus proyek
- Mengelola tugas dan anggota dalam proyek
- Melihat laporan dan analitik
- Mengatur baseline dan jalur kritis

### 3.3 Team Member
- Melihat dan mengupdate status tugas yang ditugaskan
- Menambahkan komentar dan lampiran
- Mencatat waktu pengerjaan (time log)
- Melihat proyek yang bersangkutan

### 3.4 Stakeholder / Viewer
- Hanya melihat proyek melalui public link
- Tidak dapat mengedit apapun

---

## 4. 🔧 KEBUTUHAN FUNGSIONAL

### 4.1 Modul Autentikasi & Manajemen Akun
- [ ] Registrasi dengan email & password
- [ ] Login dengan email, Google OAuth
- [ ] Lupa password & reset via email
- [ ] Manajemen profil pengguna (nama, foto, zona waktu)
- [ ] Multi-workspace support

### 4.2 Modul Dashboard
- [ ] Dashboard ringkasan semua proyek aktif
- [ ] Widget statistik: proyek aktif, tugas tertunda, overdue, selesai
- [ ] Aktivitas terbaru (activity feed)
- [ ] Shortcut ke proyek favorit
- [ ] Notifikasi dan pengingat terpusat

### 4.3 Modul Proyek
- [ ] Buat proyek baru (dari blank / dari template)
- [ ] Pengaturan proyek: nama, deskripsi, tanggal mulai, tanggal selesai, warna, status
- [ ] Duplikasi proyek
- [ ] Arsip proyek
- [ ] Hapus proyek (dengan konfirmasi)
- [ ] Template proyek siap pakai (Software Dev, Marketing, Construction, dll.)

### 4.4 Modul Gantt Chart (Fitur Utama)
- [ ] Tampilan Gantt Chart interaktif dengan timeline horizontal
- [ ] Pembuatan task/subtask dengan drag-and-drop di timeline
- [ ] Resize duration task dengan drag handle
- [ ] Dependency antar task: FS, SS, FF, SF
- [ ] Auto-scheduling: perubahan satu task otomatis menggeser yang dependen
- [ ] Critical Path Method (CPM): highlight jalur kritis
- [ ] Baseline: simpan snapshot rencana awal & bandingkan dengan aktual
- [ ] Milestone marker di timeline
- [ ] Zoom level: hari, minggu, bulan, kuartal, tahun
- [ ] Progress bar per task (0-100%)
- [ ] Color coding per task/assignee
- [ ] Export Gantt ke PNG, PDF, Excel

### 4.5 Modul Tugas (Task Management)
- [ ] Buat, edit, hapus task dan subtask
- [ ] Atribut task: nama, deskripsi, tanggal mulai, selesai, durasi, prioritas, status, assignee, progress
- [ ] Prioritas: Critical, High, Normal, Low
- [ ] Status: To Do, In Progress, On Hold, Done, Cancelled
- [ ] Lampiran file (gambar, dokumen, PDF)
- [ ] Komentar & diskusi per task
- [ ] Label/tag task
- [ ] Checklist dalam task
- [ ] Recurring task (pengulangan)
- [ ] Notifikasi deadline approaching

### 4.6 Modul Board View (Kanban)
- [ ] Tampilan Kanban board per proyek
- [ ] Kolom berdasarkan status task
- [ ] Drag-and-drop antar kolom
- [ ] Filter dan sort task
- [ ] Quick edit task dari card

### 4.7 Modul List View
- [ ] Tampilan tabel/list semua task
- [ ] Filter multi-kolom (status, assignee, prioritas, tanggal)
- [ ] Sort per kolom
- [ ] Inline edit
- [ ] Bulk action (assign, ubah status, hapus)

### 4.8 Modul Calendar View
- [ ] Tampilan kalender bulanan / mingguan
- [ ] Task ditampilkan pada tanggal deadline/start
- [ ] Quick-add task dari kalender
- [ ] Export ke Google Calendar / iCal

### 4.9 Modul Resource Management
- [ ] Tampilan beban kerja per anggota (workload view)
- [ ] Grafik alokasi waktu: under/over-allocated
- [ ] Re-assign task dari workload view
- [ ] Kapasitas anggota per hari (jam kerja)
- [ ] Laporan utilisasi sumber daya

### 4.10 Modul Time Tracking
- [ ] Log waktu manual per task
- [ ] Timer otomatis (start/stop)
- [ ] Riwayat time log
- [ ] Laporan total jam per task/proyek/anggota
- [ ] Export laporan waktu ke CSV/Excel

### 4.11 Modul Budgeting
- [ ] Estimasi biaya per task
- [ ] Biaya aktual vs estimasi
- [ ] Total anggaran proyek
- [ ] Laporan pengeluaran proyek
- [ ] Currency support (IDR, USD, EUR)

### 4.12 Modul Laporan & Analitik
- [ ] Laporan progres proyek (grafik burn-down, velocity)
- [ ] Laporan per anggota tim
- [ ] Laporan penggunaan waktu
- [ ] Laporan overdue tasks
- [ ] Export laporan ke PDF, Excel

### 4.13 Modul Kolaborasi & Notifikasi
- [ ] Komentar real-time di task
- [ ] Mention anggota (@username)
- [ ] Notifikasi in-app dan email
- [ ] Notification settings (customizable)
- [ ] Riwayat aktivitas proyek (audit log)
- [ ] Share proyek via link publik (read-only)

### 4.14 Modul Integrasi
- [ ] Export ke Excel/CSV
- [ ] Export Gantt ke PDF/PNG
- [ ] Import dari MS Project (.xml)
- [ ] Import dari Excel/CSV
- [ ] Webhook support (fase 2)
- [ ] REST API untuk integrasi pihak ketiga (fase 2)

### 4.15 Modul Pengaturan
- [ ] Pengaturan workspace (nama, logo, timezone)
- [ ] Manajemen anggota & undangan
- [ ] Pengaturan notifikasi
- [ ] Pengaturan tema (light/dark mode)
- [ ] Billing & subscription management (fase 2)

---

## 5. KEBUTUHAN NON-FUNGSIONAL

| Kategori | Requirement |
|----------|-------------|
| Performa | Load < 3 detik, Gantt render < 2 detik (500 task) |
| Keamanan | HTTPS, bcrypt, JWT, RBAC, XSS/SQLi protection |
| Skalabilitas | 10.000 pengguna aktif, 10.000 task/proyek |
| Ketersediaan | Uptime 99.5%, backup harian |
| Kompatibilitas | Chrome, Firefox, Edge, Safari; Desktop + Mobile |

---

## 6. TECH STACK & SOFTWARE

### Framework Utama
- **Wails (v2.13.0):** Digunakan sebagai framework utama untuk membangun aplikasi desktop cross-platform (menggabungkan backend Go dan frontend web).

### Backend (Go)
- **Bahasa Pemrograman:** Go (versi 1.25.0)
- **Database Driver:**
  - `modernc.org/sqlite`: Driver untuk SQLite murni menggunakan Go (tanpa CGO).
- **Library Pendukung Utama:**
  - `xuri/excelize`: Digunakan untuk membaca/menulis file Excel (mungkin untuk fitur import/export data).
  - `zalando/go-keyring`: Digunakan untuk menyimpan kredensial atau password secara aman di OS keychain/credential manager.
  - `google/uuid`: Pembuat Unique ID.

### Frontend (Web Technologies)
- **Core Framework & Build Tool:**
  - **Vue.js 3 (v3.5.39):** Framework JavaScript utama.
  - **Vite (v8.1.1):** Build tool dan development server yang sangat cepat.
  - **TypeScript (~v6.0.2):** Digunakan untuk pengetikan statis (static typing) pada JavaScript.
  - **Pinia (v4.0.2):** State management standar untuk Vue 3.
- **UI & Komponen Visual:**
  - **PrimeVue (v4.2.5) & PrimeIcons:** Library komponen UI utama yang digunakan untuk mendesain antarmuka aplikasi.
- **Library Fitur Khusus:**
  - **Monaco Editor (v0.56.0):** Digunakan sebagai editor kode di dalam aplikasi (seperti yang digunakan di VS Code, biasanya untuk editor query SQL).
  - **Vue Flow (v1.48.2):** Digunakan untuk membuat antarmuka berbasis node/diagram interaktif (mungkin untuk visualisasi relasi antar tabel/ERD).
  - **Chart.js (v4.5.1):** Digunakan untuk membuat visualisasi data/grafik.
  - **Leaflet (v1.9.4):** Digunakan untuk menampilkan peta interaktif (berkaitan dengan data spasial).
  - **Marked (v18.0.9):** Digunakan untuk merender teks Markdown menjadi HTML (misalnya untuk dokumentasi atau catatan).

> **Ringkasan:** Secara garis besar, project ini adalah aplikasi Database Management/Studio modern berbentuk Desktop App yang memiliki fitur komprehensif seperti eksekusi query dengan Monaco Editor, visualisasi data (Chart.js & Leaflet), visualisasi relasi (Vue Flow), hingga fitur export/import (Excelize) yang mendukung berbagai jenis database (MySQL, PostgreSQL, SQLite).

---

## 7. PRIORITAS FITUR (MoSCoW)

### Must Have (Fase 1)
- Autentikasi & manajemen pengguna
- Dashboard proyek
- Gantt Chart interaktif (core)
- Task management CRUD
- Board view (Kanban)
- Notifikasi dasar
- Export PDF/Excel

### Should Have (Fase 1-2)
- Resource management
- Time tracking
- List & Calendar view
- Laporan & analitik dasar

### Could Have (Fase 2)
- Budgeting
- Import MS Project
- Integrasi eksternal (Slack, Google Drive)
- REST API publik

### Won't Have (Out of Scope)
- Native mobile app
- AI/ML features
- Advanced billing system

---

## 8. ESTIMASI TIMELINE

| Fase | Deskripsi | Durasi |
|------|-----------|--------|
| 1 | Setup & Autentikasi | 2 minggu |
| 2 | Dashboard & Manajemen Proyek | 2 minggu |
| 3 | Gantt Chart Core | 4 minggu |
| 4 | Task Management & Views (Board, List, Calendar) | 3 minggu |
| 5 | Kolaborasi & Notifikasi | 2 minggu |
| 6 | Resource & Time Tracking | 2 minggu |
| 7 | Laporan & Export | 2 minggu |
| 8 | Testing & QA | 2 minggu |
| 9 | Deployment & Go Live | 1 minggu |
| **Total** | | **~20 minggu** |

---

*Dokumen ini adalah living document dan akan diperbarui seiring perkembangan proyek.*
