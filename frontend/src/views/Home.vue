<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { CheckAuth } from '../../wailsjs/go/services/AuthService'
import { models } from '../../wailsjs/go/models'

const user = ref<models.User | null>(null)

// Mock data for Dashboard stats (to be replaced with real backend data later)
const stats = ref([
  { title: 'Total Proyek', value: '0', icon: 'pi pi-folder', color: 'blue' },
  { title: 'Total Tugas', value: '0', icon: 'pi pi-check-square', color: 'green' },
  { title: 'Overdue', value: '0', icon: 'pi pi-exclamation-circle', color: 'red' },
  { title: 'Selesai Hari Ini', value: '0', icon: 'pi pi-star', color: 'yellow' }
])

onMounted(async () => {
  try {
    const authState = await CheckAuth()
    if (authState.is_logged_in && authState.user) {
      user.value = authState.user
    }
  } catch (err) {
    console.error(err)
  }
})
</script>

<template>
  <div class="dashboard-page">
    <div class="page-header">
      <div>
        <h1 class="page-title">Dashboard</h1>
        <p class="page-subtitle" v-if="user">Selamat datang kembali, {{ user.full_name }}!</p>
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="stats-grid">
      <div class="stat-card" v-for="stat in stats" :key="stat.title">
        <div class="stat-icon" :class="'icon-' + stat.color">
          <i :class="stat.icon"></i>
        </div>
        <div class="stat-content">
          <div class="stat-title">{{ stat.title }}</div>
          <div class="stat-value">{{ stat.value }}</div>
        </div>
      </div>
    </div>

    <!-- Layout Utama Bawah -->
    <div class="dashboard-content">
      <div class="main-column">
        <!-- My Projects Section -->
        <div class="section-container">
          <div class="section-header">
            <h2 class="section-title">Proyek Saya</h2>
            <a href="#" class="view-all" @click.prevent="$router.push('/projects')">Lihat Semua</a>
          </div>
          
          <div class="projects-grid empty-state">
            <div class="empty-message">
              <i class="pi pi-folder-open empty-icon"></i>
              <p>Belum ada proyek.</p>
              <button class="btn-primary-outline" @click="$router.push('/projects')">Buat Proyek Baru</button>
            </div>
          </div>
        </div>

        <!-- Recent Activity Section -->
        <div class="section-container">
          <div class="section-header">
            <h2 class="section-title">Aktivitas Terbaru</h2>
          </div>
          <div class="activity-feed empty-state">
            <p>Belum ada aktivitas.</p>
          </div>
        </div>
      </div>

      <div class="side-column">
        <!-- Quick Actions -->
        <div class="section-container">
          <div class="section-header">
            <h2 class="section-title">Tindakan Cepat</h2>
          </div>
          <div class="quick-actions">
            <button class="quick-btn" @click="$router.push('/projects')">
              <i class="pi pi-plus"></i> Proyek Baru
            </button>
            <button class="quick-btn">
              <i class="pi pi-check-square"></i> Tugas Baru
            </button>
            <button class="quick-btn">
              <i class="pi pi-users"></i> Undang Anggota
            </button>
          </div>
        </div>

        <!-- Upcoming Deadlines -->
        <div class="section-container">
          <div class="section-header">
            <h2 class="section-title">Tenggat Mendekat</h2>
          </div>
          <div class="deadlines-list empty-state">
            <p>Tidak ada tenggat waktu dalam waktu dekat.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
  gap: var(--space-8);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.page-title {
  margin: 0;
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text-primary);
}

.page-subtitle {
  margin: var(--space-1) 0 0 0;
  color: var(--text-secondary);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: var(--space-6);
}

.stat-card {
  background: var(--bg-surface);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  display: flex;
  align-items: center;
  gap: var(--space-4);
  box-shadow: var(--shadow-sm);
  transition: transform var(--transition-fast), box-shadow var(--transition-fast);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.25rem;
}

.icon-blue { background: rgba(56, 189, 248, 0.15); color: #38bdf8; }
.icon-green { background: rgba(52, 211, 153, 0.15); color: #34d399; }
.icon-red { background: rgba(248, 113, 113, 0.15); color: #f87171; }
.icon-yellow { background: rgba(251, 191, 36, 0.15); color: #fbbf24; }

.stat-content {
  flex: 1;
}

.stat-title {
  color: var(--text-secondary);
  font-size: 0.875rem;
  font-weight: 500;
  margin-bottom: var(--space-1);
}

.stat-value {
  color: var(--text-primary);
  font-size: 1.5rem;
  font-weight: 700;
}

.dashboard-content {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: var(--space-8);
}

@media (max-width: 1024px) {
  .dashboard-content {
    grid-template-columns: 1fr;
  }
}

.main-column, .side-column {
  display: flex;
  flex-direction: column;
  gap: var(--space-8);
}

.section-container {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  margin: 0;
  font-size: 1.15rem;
  font-weight: 600;
  color: var(--text-primary);
}

.view-all {
  color: var(--primary-400);
  text-decoration: none;
  font-size: 0.875rem;
  font-weight: 500;
}

.view-all:hover {
  text-decoration: underline;
}

.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--space-4);
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: var(--space-10) 0;
  background: var(--bg-surface);
  border: 1px dashed var(--border-default);
  border-radius: var(--radius-lg);
  grid-column: 1 / -1;
  color: var(--text-tertiary);
  font-size: 0.9rem;
}

.empty-message {
  text-align: center;
}

.empty-icon {
  font-size: 2.5rem;
  margin-bottom: var(--space-4);
  color: var(--border-default);
}

.btn-primary-outline {
  margin-top: var(--space-4);
  background: transparent;
  border: 1px solid var(--primary-500);
  color: var(--primary-400);
  padding: 8px 16px;
  border-radius: var(--radius-md);
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.btn-primary-outline:hover {
  background: rgba(99, 102, 241, 0.1);
}

.quick-actions {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--space-3);
}

.quick-btn {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  background: var(--bg-surface);
  border: 1px solid var(--border-default);
  color: var(--text-primary);
  padding: var(--space-3) var(--space-4);
  border-radius: var(--radius-md);
  cursor: pointer;
  font-size: 0.95rem;
  font-weight: 500;
  transition: all var(--transition-fast);
}

.quick-btn i {
  color: var(--primary-400);
}

.quick-btn:hover {
  background: var(--bg-hover);
  border-color: var(--border-hover);
  transform: translateX(4px);
}
</style>
