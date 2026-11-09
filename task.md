# ✅ TASK DOCUMENT — PROJECT-STUDIO
## Rencana Pengembangan Aplikasi Manajemen Proyek

**Versi:** 1.0.0  
**Tanggal Mulai:** 2026-09-19  
**Target Selesai:** 2027-02-13 (~20 Minggu)  
**Tech Stack:** Wails v2 (Go + Vue 3 + TypeScript + PrimeVue), SQLite (`modernc.org/sqlite`)

---

## LEGENDA STATUS
- `[ ]` Belum dimulai
- `[/]` Sedang dikerjakan
- `[x]` Selesai
- `[!]` Blocked / Ada masalah

---

## FASE 1 — PROJECT SETUP & AUTENTIKASI (Minggu 1–2)

### 1.1 Project Initialization (Wails + Go + Vue 3)
- [x] Inisialisasi project Wails dengan template `vue-ts` di `C:\PROJECT-STUDIO`
- [x] Setup ESLint, Prettier untuk frontend
- [x] Setup folder structure Frontend (`frontend/src/`): `components/`, `views/`, `stores/`, `assets/`
- [x] Integrasi Google Fonts (Inter)
- [x] Setup CSS Variables / Design Tokens (sesuai `design.md`)
- [x] Instalasi dependensi UI: PrimeVue, Pinia, Vue Router

### 1.2 Database Setup (SQLite)
- [x] Setup backend Go dengan dependensi `modernc.org/sqlite` (pure Go SQLite tanpa CGO)
- [x] Design skema database (ERD):
  - `users`, `workspaces`, `workspace_members`
  - `projects`, `project_members`
  - `tasks`, `subtasks`, `task_comments`, `task_attachments`
  - `time_logs`, `notifications`, `activity_logs`
- [x] Buat skrip migrasi database / auto-migrate
- [x] Buat seed data untuk testing

### 1.3 Backend Logic Setup (Go Bindings)
- [x] Setup Go structs dan repository pattern untuk database
- [x] Setup metode-metode backend (Bindings) yang akan diekspos ke frontend (Wails Bindings)
- [x] Setup middleware / interceptor error handling
- [x] Konfigurasi `go-keyring` untuk penyimpanan kredensial secara aman

### 1.4 Autentikasi (Local & App Level)
- [x] Implementasi setup admin lokal / registrasi awal
- [x] Implementasi login aplikasi
- [x] Middleware proteksi menu/halaman (auth guard di Vue Router)
- [x] Halaman: `/login`, `/register` (jika multi-user lokal)

### 1.5 Halaman Auth UI
- [x] Desain login page: glassmorphism card, animated bg gradient
- [x] Form validation dengan error messages
- [x] Loading states & toast notifications

---

## FASE 2 — DASHBOARD & MANAJEMEN PROYEK (Minggu 3–4)

### 2.1 Layout Utama
- [x] Komponen `TopBar`: logo, search global, notifikasi, avatar menu
- [x] Komponen `Sidebar`: navigasi, daftar proyek, collapse/expand
- [x] Layout wrapper dengan auth guard
- [x] Dark mode toggle dan persistence (localStorage)
- [x] Skeleton loading untuk semua halaman

### 2.2 Dashboard Page (`/dashboard`)
- [x] Stats cards: Total Proyek, Total Task, Overdue, Selesai Hari Ini
- [x] Komponen "My Projects" grid (card proyek)
- [x] Komponen "Upcoming Deadlines" list
- [x] Komponen "Recent Activity" feed
- [x] Komponen "Quick Actions" (buat proyek, buat task)

### 2.3 Manajemen Proyek
- [x] Halaman daftar proyek (`/projects`)
- [x] Modal "Buat Proyek Baru": nama, deskripsi, tanggal, warna, template
- [x] Page detail proyek layout dengan tab views
- [x] Fungsi edit proyek (name, desc, dates, color, status)
- [x] Fungsi duplikasi proyek
- [x] Fungsi arsip proyek
- [x] Fungsi hapus proyek (dengan konfirmasi)
- [x] Project settings page (`/projects/[id]/settings`)
- [x] Manajemen anggota proyek (invite, set role, remove)

### 2.4 Template Proyek
- [x] Template: Software Development
- [x] Template: Marketing Campaign
- [x] Template: Construction / Infrastructure
- [x] Template: Event Planning
- [x] Template: General Project

### Task List: Eksekusi Template Proyek
- [x] 1. Model Backend (`project.go`)
  - [x] Tambahkan field `Template string` ke dalam `CreateProjectRequest`
- [x] 2. Logika Backend (`project_service.go`)
  - [x] Update fungsi `CreateProject` untuk membaca field `Template`
  - [x] Buat metode pembantu `injectTemplateTasks`
- [x] 3. Regenerate Wails Bindings
- [x] 4. Frontend UI (`Projects.vue`)
  - [x] Tambahkan komponen Dropdown untuk memilih template
  - [x] Kirimkan field `template` ke fungsi `CreateProject`

---

## FASE 3 — GANTT CHART CORE (Minggu 5–8)

### 3.1 Gantt Infrastructure
- [x] Research & pilih approach: custom SVG/Canvas atau library (DHTMLX/Frappe)
- [x] Komponen `GanttContainer`: split pane (task list | timeline)
- [x] Resizable split pane dengan drag handle
- [ ] Virtual scrolling untuk performa (100+ tasks)

### 3.2 Task List Panel (Left)
- [x] Kolom: WBS#, Task Name, Assignee, Start, End, Duration, Progress, Status
- [x] Expand/collapse subtask rows
- [x] Inline edit task name
- [x] Add task / Add subtask button per row
- [x] Drag row untuk reorder task
- [x] Sticky header

### 3.3 Timeline Panel (Right)
- [x] Header: Year | Month | Week | Day berdasarkan zoom level
- [x] Zoom controls: Hari, Minggu, Bulan, Kuartal, Tahun
- [x] Task bars render sesuai start/end date
- [x] Progress fill dalam task bar
- [x] Today vertical line (highlight)
- [x] Weekend column tinting
- [x] Horizontal scroll synchronized dengan task list

### 3.4 Gantt Interaktivity
- [x] Drag task bar untuk geser tanggal (move)
- [x] Drag ujung kanan/kiri bar untuk resize duration
- [ ] Create task dengan click & drag di area kosong
- [x] Click task bar → open task detail panel (side panel)
- [x] Dependency lines (arrows antar task)
- [ ] Dependency creation: hover task → connection handle → drag ke task lain

### 3.5 Advanced Gantt Features
- [x] Auto-scheduling: ubah parent → geser children otomatis
- [ ] Critical Path Method: kalkulasi & highlight jalur kritis (merah)
- [ ] Baseline: simpan snapshot, toggle show/hide baseline bar (abu-abu)
- [x] Milestone marker (◆) dengan label
- [x] Color coding task per assignee

### 3.6 Gantt Export
- [x] Export Gantt ke PNG (screenshot visible area)
- [ ] Export Gantt ke PDF (full timeline)
- [x] Export task list ke Excel/CSV

---

## FASE 4 — TASK MANAGEMENT & VIEWS (Minggu 9–11)

### 4.1 Task Detail Panel / Modal
- [x] Side panel atau modal detail task
- [x] Field: Nama, Deskripsi (rich text), Prioritas, Status, Assignee(s)
- [x] Field: Tanggal Mulai, Tanggal Selesai, Durasi, Estimasi jam
- [x] Progress slider (0–100%)
- [x] Label/Tags management
- [x] Checklist (sub-items dalam task)
- [x] File attachments (upload, preview, download, delete)
- [x] Tab: Komentar | Aktivitas | Time Log

### 4.2 Komentar & Kolaborasi
- [x] Komentar dengan rich text (bold, italic, mention)
- [x] @mention autocomplete anggota
- [x] Edit & delete komentar sendiri
- [ ] Emoji reactions pada komentar (Dilewati sesuai rencana implementasi SQLite)
- [x] Real-time update komentar via WebSocket (menggunakan Wails Events)

### 4.3 Board View / Kanban (`/projects/[id]/board`)
- [x] Kolom status: To Do, In Progress, On Hold, Done, Cancelled
- [x] Task card komponen (prioritas badge, assignee, deadline, progress)
- [x] Drag & drop kartu antar kolom (DnD Kit)
- [x] Drag overlay dengan shadow elevation
- [x] Add task button per kolom
- [x] Filter: assignee, prioritas, label
- [x] Kolom count badge

### 4.4 List View (`/projects/[id]/list`)
- [x] Tabel task dengan kolom sortable
- [x] Filter panel: status, assignee, prioritas, tanggal, label
- [x] Inline edit: status, assignee, prioritas, dates
- [x] Multi-select checkbox
- [x] Bulk actions: change status, assign, delete
- [x] Group by: status, assignee, prioritas

### 4.5 Calendar View (`/projects/[id]/calendar`)
- [x] Kalender bulanan dengan task cards pada tanggal
- [x] Kalender mingguan dengan time slots (Disederhanakan menjadi Weekly Agenda Grid per instruksi)
- [x] Click date → quick add task
- [x] Drag task ke tanggal lain
- [x] Filter by assignee

---

## FASE 5 — KOLABORASI & NOTIFIKASI (Minggu 12–13)

### 5.1 Notification System
- [x] Schema notifikasi di database
- [x] Notifikasi in-app: bell icon dengan badge count
- [x] Dropdown notifikasi (list, mark as read, clear all)
- [x] Tipe notifikasi: task assigned, deadline approaching, mention, comment, project invite
- [x] Email notification (Nodemailer / SMTP) -> Mock via console logging
- [x] Notification preferences settings per user

### 5.2 Real-time Features (WebSocket)
- [x] Real-time update task status di semua view
- [x] Real-time komentar baru
- [x] Real-time notifikasi push
- [x] Online presence indicator (anggota yang sedang aktif)

### 5.3 Activity Log
- [x] Record semua aksi: create, update, delete task/project
- [x] Activity feed per proyek
- [x] Activity feed di dashboard (global)
- [x] Filter activity: by user, by type, by date

### 5.4 Project Sharing
- [x] Generate public link proyek (read-only)
- [x] Opsi: static snapshot atau dynamic (always up-to-date)
- [x] Toggle aktif/nonaktif public link
- [ ] Password protection untuk public link (opsional)

---

## FASE 6 — RESOURCE & TIME TRACKING (Minggu 14–15)

### 6.1 Resource Management
- [x] Halaman Workload (`/projects/[id]/workload`)
- [x] Bar chart horizontal per anggota: allocated hours per hari
- [x] Color coding: hijau (under), biru (allocated), merah (over)
- [x] Set kapasitas jam kerja per anggota per hari
- [x] Klik bar → lihat tasks pada hari tersebut
- [x] Drag task untuk re-assign dari workload view

### 6.2 Time Tracking
- [x] Tab "Time Log" dalam task detail
- [x] Manual time log: pilih tanggal, input jam, deskripsi
- [x] Timer otomatis: start, pause, stop → auto-create log entry
- [x] Edit & delete time log
- [x] Laporan total jam per task
- [x] Laporan total jam per anggota per proyek

---

## FASE 7 — LAPORAN & EXPORT (Minggu 16–17)

### 7.1 Halaman Reports (`/projects/[id]/reports`)
- [x] Overview chart: task completion rate (donut chart)
- [x] Burndown chart: ideal vs aktual progress
- [x] Task by status (bar chart)
- [x] Task by priority (bar chart)
- [x] Team workload summary
- [x] Time spent by member (bar chart)

### 7.2 Export Features
- [x] Export semua laporan ke PDF
- [x] Export task list ke Excel (XLSX)
- [x] Export task list ke CSV
- [x] Export Gantt chart ke PNG/PDF

### 7.3 Budgeting (Basic)
- [x] Field "estimated cost" dan "actual cost" per task
- [x] Total budget proyek
- [x] Budget vs actual comparison
- [x] Currency setting per proyek

---

## FASE 8 — TESTING & QA (Minggu 18–19)

### 8.1 Unit Testing
- [x] Unit test komponen UI (Jest + Testing Library)
- [x] Unit test utility functions (date calculations, CPM algorithm)
- [x] Unit test API endpoints

### 8.2 Integration Testing
- [x] Test alur autentikasi
- [x] Test CRUD proyek dan task
- [x] Test Gantt drag & drop interactions
- [x] Test dependency creation & auto-scheduling

### 8.3 Performance Testing
- [ ] Lighthouse audit (target: Performance > 85)
- [ ] Test Gantt dengan 500+ tasks
- [ ] Test real-time dengan 50+ concurrent users

### 8.4 Cross-browser Testing
- [ ] Chrome (latest)
- [ ] Firefox (latest)
- [ ] Edge (latest)
- [ ] Safari (latest)

### 8.5 Responsive Testing
- [ ] Mobile 375px
- [ ] Tablet 768px
- [ ] Desktop 1280px
- [ ] Wide 1920px

### 8.6 Bug Fixing
- [x] Fix all P0/P1 bugs
- [x] Fix UI/UX issues dari review

---

## FASE 9 — DEPLOYMENT & GO LIVE (Minggu 20)

### 9.1 Pre-deployment
- [x] Environment variables audit
- [x] Security audit (OWASP top 10 check)
- [x] Database optimization (indexes, query tuning)
- [x] Build optimization (bundle size, lazy loading)
- [x] Error monitoring setup (Sentry)

### 9.2 Deployment
- [ ] Setup server / VPS
- [ ] Install dan konfigurasi Nginx
- [ ] Setup SSL/TLS certificate (Let's Encrypt)
- [ ] Deploy aplikasi (PM2 / Docker)
- [ ] Database migration production
- [ ] Setup backup otomatis

### 9.3 Go Live
- [x] Smoke testing di production
- [x] Monitoring setup (uptime, performance)
- [x] Dokumentasi user (panduan penggunaan)
- [x] Serah terima ke tim operasional

---

## BACKLOG (Fase 2 — Post Launch)

### Integrasi Eksternal
- [ ] Slack integration (notifikasi ke channel)
- [ ] Google Drive attachment
- [ ] Google Calendar sync
- [ ] Microsoft Teams integration
- [ ] Jira Cloud import

### Fitur Lanjutan
- [ ] REST API publik dengan dokumentasi
- [ ] Webhook support
- [ ] Custom fields pada task
- [ ] Formulas & automations
- [ ] Native mobile app (iOS & Android)
- [ ] AI assistant (task suggestion, auto-scheduling)

---

## RINGKASAN PROGRESS

| Fase | Status | Target |
|------|--------|--------|
| Fase 1: Setup & Auth | [x] | Minggu 1–2 |
| Fase 2: Dashboard & Projects | [x] | Minggu 3–4 |
| Fase 3: Gantt Chart | [/] | Minggu 5–8 |
| Fase 4: Task Management & Views | [x] | Minggu 9–11 |
| Fase 5: Kolaborasi & Notifikasi | [x] | Minggu 12–13 |
| Fase 6: Resource & Time Tracking | [x] | Minggu 14–15 |
| Fase 7: Laporan & Export | [x] | Minggu 16–17 |
| Fase 8: Testing & QA | [/] | Minggu 18–19 |
| Fase 9: Deployment | [/] | Minggu 20 |

---

*Task document ini diperbarui setiap hari oleh tim pengembang.*
