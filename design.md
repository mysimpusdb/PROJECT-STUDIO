# 🎨 DESIGN DOCUMENT — PROJECT-STUDIO
## Aplikasi Manajemen Proyek (GanttPRO-Inspired)

**Versi:** 1.0.0  
**Tanggal:** 2026-09-19  
**Project:** PROJECT-STUDIO

---

## 1. DESIGN PHILOSOPHY

PROJECT-STUDIO menganut prinsip desain modern dengan pendekatan:
- **Clarity First**: Informasi proyek yang kompleks disajikan secara jelas dan terstruktur
- **Dark Mode Premium**: Tema gelap sebagai default untuk mengurangi eye strain saat kerja lama
- **Glassmorphism + Depth**: Layering visual dengan blur dan transparansi untuk kesan premium
- **Micro-interactions**: Animasi halus pada setiap interaksi untuk feedback responsif
- **Data Density**: Menampilkan banyak informasi tanpa terasa padat

---

## 2. COLOR PALETTE

### Primary Colors
| Token | Hex | Penggunaan |
|-------|-----|------------|
| --primary-500 | #6366F1 | Indigo - Accent utama, CTA button |
| --primary-400 | #818CF8 | Hover states, link |
| --primary-300 | #A5B4FC | Secondary accent |
| --primary-600 | #4F46E5 | Active states, pressed |

### Neutral (Dark Mode Base)
| Token | Hex | Penggunaan |
|-------|-----|------------|
| --bg-base | #0F0F14 | Background halaman utama |
| --bg-surface | #16161E | Card, panel, sidebar |
| --bg-elevated | #1E1E2A | Modal, dropdown, tooltip |
| --bg-hover | #252535 | Hover state pada item |
| --border-subtle | #2A2A3E | Border tipis |
| --border-default | #383850 | Border normal |

### Text Colors
| Token | Hex | Penggunaan |
|-------|-----|------------|
| --text-primary | #F0F0FF | Heading, konten utama |
| --text-secondary | #9494B8 | Label, deskripsi sekunder |
| --text-muted | #6B6B8A | Placeholder, disabled |
| --text-inverse | #0F0F14 | Teks di atas background terang |

### Semantic Colors
| Token | Hex | Penggunaan |
|-------|-----|------------|
| --success | #22C55E | Task selesai, status sukses |
| --warning | #F59E0B | Overdue approaching, peringatan |
| --danger | #EF4444 | Critical path, deadline terlewat |
| --info | #38BDF8 | Informasi, milestone |

### Priority Colors (Task)
| Priority | Color |
|----------|-------|
| Critical | #EF4444 (red) |
| High | #F97316 (orange) |
| Normal | #6366F1 (indigo) |
| Low | #6B7280 (gray) |

### Project Colors (untuk diferensiasi proyek di Gantt)
Palet 12 warna: Indigo, Violet, Sky, Emerald, Amber, Rose, Teal, Fuchsia, Lime, Cyan, Orange, Pink

---

## 3. TYPOGRAPHY

### Font Family
```css
--font-sans: 'Inter', 'Segoe UI', system-ui, sans-serif;
--font-mono: 'JetBrains Mono', 'Fira Code', monospace;
```
*(Import dari Google Fonts: Inter 300, 400, 500, 600, 700)*

### Type Scale
| Token | Size | Weight | Line Height | Penggunaan |
|-------|------|--------|-------------|------------|
| --text-xs | 11px | 400 | 16px | Badge, label kecil |
| --text-sm | 13px | 400 | 20px | Meta info, caption |
| --text-base | 14px | 400 | 22px | Body text default |
| --text-md | 15px | 500 | 24px | Task nama |
| --text-lg | 17px | 600 | 26px | Section title |
| --text-xl | 20px | 600 | 28px | Panel heading |
| --text-2xl | 24px | 700 | 32px | Page title |
| --text-3xl | 30px | 700 | 40px | Hero heading |

---

## 4. SPACING & LAYOUT

### Spacing Scale
```css
--space-1: 4px   --space-2: 8px   --space-3: 12px
--space-4: 16px  --space-5: 20px  --space-6: 24px
--space-8: 32px  --space-10: 40px --space-12: 48px
--space-16: 64px
```

### Border Radius
```css
--radius-sm: 4px    /* Badge, small chips */
--radius-md: 8px    /* Card, input field */
--radius-lg: 12px   /* Modal, panel */
--radius-xl: 16px   /* Large card */
--radius-full: 9999px /* Pill, avatar */
```

### Shadow
```css
--shadow-sm: 0 1px 3px rgba(0,0,0,0.4);
--shadow-md: 0 4px 16px rgba(0,0,0,0.5);
--shadow-lg: 0 8px 32px rgba(0,0,0,0.6);
--shadow-glow: 0 0 20px rgba(99,102,241,0.3); /* Primary glow */
```

---

## 5. LAYOUT ARCHITECTURE

### 5.1 Struktur Halaman Utama
```
┌─────────────────────────────────────────────────────┐
│  TOP BAR (56px)  - Logo | Search | Notif | Avatar   │
├──────────┬──────────────────────────────────────────┤
│          │  CONTENT AREA                             │
│ SIDEBAR  │  (dinamis berdasarkan active view)        │
│ (240px)  │                                           │
│          │                                           │
│          │                                           │
└──────────┴──────────────────────────────────────────┘
```

### 5.2 Sidebar Navigation
- **Fixed width**: 240px (collapsible ke 60px icon-only)
- **Sections**:
  - My Space (Dashboard, My Tasks, Inbox)
  - Projects (list proyek dengan color dot)
  - Workspace (Members, Settings)
- Transition: slide + fade 200ms ease

### 5.3 Gantt Chart Layout
```
┌─────────────────────────────────────────────────────┐
│ TOOLBAR: Views | Filter | Group | Zoom | Export     │
├────────────────────┬────────────────────────────────┤
│  TASK LIST (40%)   │  TIMELINE (60%)                │
│  - WBS #           │  Month/Week headers             │
│  - Task Name       │  ████░░░░░░░ Task bars          │
│  - Assignee        │  →→→→ Dependency arrows         │
│  - Start/End       │  ◆ Milestones                  │
│  - Progress        │  🔴 Critical path               │
│  - Status          │                                 │
└────────────────────┴────────────────────────────────┘
```

---

## 6. KOMPONEN UI

### 6.1 Button System
```
Variants: primary | secondary | ghost | danger | success
Sizes: xs | sm | md | lg
States: default | hover | active | disabled | loading
```

**Primary Button:**
- Background: `--primary-500`
- Border-radius: `--radius-md`
- Hover: `--primary-400` + `--shadow-glow`
- Transition: background 150ms, shadow 200ms

**Ghost Button:**
- Background: transparent
- Border: 1px solid `--border-default`
- Hover: `--bg-elevated`

### 6.2 Task Card (Kanban Board)
```
┌────────────────────────────────────┐
│ [■ High]          [Assignee Avatar] │
│ Task Name                           │
│ Deskripsi singkat task...           │
│ ─────────────────────────────────  │
│ 📎 3  💬 5  ⏱ 2h   [████░░] 70%   │
│ Deadline: 25 Sep 2026               │
└────────────────────────────────────┘
```
- Card background: `--bg-elevated` dengan glassmorphism (backdrop-blur: 8px)
- Border: 1px solid `--border-subtle`
- Hover: translate-Y(-2px) + `--shadow-md`
- Left border accent: warna prioritas

### 6.3 Gantt Bar
- Bar height: 28px, border-radius: 4px
- Progress fill: warna lebih gelap dari bar utama
- Critical path bar: gradient merah `#EF4444 → #DC2626`
- Milestone: diamond shape ◆ 14px
- Hover tooltip: nama task, tanggal, progress, assignee

### 6.4 Input Fields
- Background: `rgba(255,255,255,0.04)`
- Border: 1px solid `--border-default`
- Focus: border `--primary-500` + outer glow
- Placeholder: `--text-muted`
- Height: 36px (sm), 40px (md), 44px (lg)

### 6.5 Badges / Pills
```
Status badges:
• To Do      → --bg-elevated, text-secondary
• In Progress → blue bg (rgba 20%), blue text
• On Hold    → amber bg, amber text
• Done       → green bg, green text
• Cancelled  → red bg (muted), red text
```

### 6.6 Modal / Dialog
- Backdrop: `rgba(0,0,0,0.7)` dengan backdrop-blur: 4px
- Container: `--bg-elevated`, border-radius: `--radius-xl`
- Max-width: 560px (standard), 800px (wide)
- Animation: fade-in + scale(0.96→1.0) 200ms

### 6.7 Avatar & User Chips
- Size: 24px (xs), 32px (sm), 40px (md), 48px (lg)
- Fallback: initial huruf pada background gradient warna unik per user
- Border: 2px solid `--bg-surface`
- Group avatar: overlap -8px

### 6.8 Progress Bar
- Height: 6px (thin), 10px (normal)
- Background: `rgba(255,255,255,0.1)`
- Fill: gradient `--primary-400 → --primary-600`
- Animation: width transition 400ms ease-out

---

## 7. PAGE DESIGNS

### 7.1 Landing/Login Page
- Full-screen dark background dengan animated gradient mesh
- Glassmorphism login card di tengah
- Logo PROJECT-STUDIO dengan glow effect
- Google OAuth button + Email/Password form

### 7.2 Dashboard Page
```
┌──────────────────────────────────────────────────┐
│ Good Morning, [Name] 👋     Today: 19 Sep 2026   │
├──────┬──────┬──────┬──────────────────────────── │
│ 12   │ 48   │ 7    │ 3                            │
│Proj  │Tasks │Overdue│ Done Today                  │
├──────┴──────┴──────┴──────────────────────────── │
│ MY PROJECTS                    [+ New Project]    │
│ [Card] [Card] [Card] [Card]                       │
├───────────────────────────────────────────────── │
│ UPCOMING DEADLINES          RECENT ACTIVITY       │
│ ████ Task 1 - 2 days        [Feed items]          │
│ ████ Task 2 - 5 days                              │
└──────────────────────────────────────────────────┘
```

### 7.3 Gantt Chart Page
- Split pane: Task list (resizable) | Timeline (scrollable horizontal)
- Sticky header row dengan zoom controls
- Zebra striping pada rows (alternating opacity)
- Weekend columns: darker background tint
- Today line: vertical bright blue `--info` dengan subtle glow
- Sticky first column (task name)

### 7.4 Board (Kanban) Page
- Horizontal scrollable columns
- Column header: count badge, add task button
- Card drag dengan drag overlay shadow
- Empty state per kolom: illustration + add button

### 7.5 My Tasks / List View
- Sticky table header
- Sortable columns (click header)
- Row hover highlight
- Inline status dropdown
- Multi-select dengan checkbox

---

## 8. ANIMASI & MICRO-INTERACTIONS

### Transition Defaults
```css
--transition-fast: 100ms ease;
--transition-base: 200ms ease;
--transition-slow: 350ms ease;
--transition-spring: 400ms cubic-bezier(0.34, 1.56, 0.64, 1);
```

### Key Animations
| Element | Animation |
|---------|-----------|
| Page route change | Fade + slide-up 200ms |
| Modal open | Scale 0.96→1 + fade 200ms |
| Card hover | TranslateY -2px + shadow 150ms |
| Button click | Scale 0.97 press effect 100ms |
| Task drag | Scale 1.05 + shadow elevation |
| Sidebar toggle | Width slide 250ms ease |
| Progress update | Bar width animate 400ms |
| Notification | Slide-in from right 300ms |
| Skeleton loading | Shimmer gradient animation |

---

## 9. RESPONSIVE DESIGN

| Breakpoint | Width | Layout |
|-----------|-------|--------|
| Mobile | < 640px | Sidebar hidden, hamburger menu; views simplified |
| Tablet | 640–1024px | Sidebar collapsible; Gantt: list-only mode |
| Desktop | 1024–1440px | Full layout |
| Wide | > 1440px | Max content width 1800px |

---

## 10. ICONOGRAPHY

- **Icon Library**: Lucide Icons (stroke-based, consistent weight)
- **Size standard**: 14px (xs), 16px (sm), 20px (md), 24px (lg)
- **Stroke width**: 1.5px (default), 2px (emphasis)
- Konsisten menggunakan outline style, bukan filled

---

## 11. ACCESSIBILITY

- Color contrast minimum 4.5:1 (WCAG AA)
- Focus ring visible: 2px outline `--primary-400`, offset 2px
- Keyboard navigation support (Tab, Arrow keys, Enter, Escape)
- ARIA labels pada semua interactive elements
- Screen reader-friendly semantic HTML
- Reduced motion support: `@media (prefers-reduced-motion: reduce)`

---

## 12. DESIGN TOKENS SUMMARY (CSS Variables)

```css
:root {
  /* Colors - Dark Mode (default) */
  --bg-base: #0F0F14;
  --bg-surface: #16161E;
  --bg-elevated: #1E1E2A;
  --bg-hover: #252535;
  --border-subtle: #2A2A3E;
  --border-default: #383850;

  --text-primary: #F0F0FF;
  --text-secondary: #9494B8;
  --text-muted: #6B6B8A;

  --primary-300: #A5B4FC;
  --primary-400: #818CF8;
  --primary-500: #6366F1;
  --primary-600: #4F46E5;

  --success: #22C55E;
  --warning: #F59E0B;
  --danger: #EF4444;
  --info: #38BDF8;

  /* Spacing */
  --space-1: 4px; --space-2: 8px; --space-3: 12px;
  --space-4: 16px; --space-6: 24px; --space-8: 32px;

  /* Border Radius */
  --radius-sm: 4px; --radius-md: 8px;
  --radius-lg: 12px; --radius-xl: 16px; --radius-full: 9999px;

  /* Shadows */
  --shadow-sm: 0 1px 3px rgba(0,0,0,0.4);
  --shadow-md: 0 4px 16px rgba(0,0,0,0.5);
  --shadow-lg: 0 8px 32px rgba(0,0,0,0.6);
  --shadow-glow: 0 0 20px rgba(99,102,241,0.3);

  /* Typography */
  --font-sans: 'Inter', system-ui, sans-serif;
  --font-mono: 'JetBrains Mono', monospace;

  /* Transitions */
  --transition-fast: 100ms ease;
  --transition-base: 200ms ease;
  --transition-slow: 350ms ease;
}
```

---

*Design document ini menjadi acuan utama dalam implementasi UI PROJECT-STUDIO.*
