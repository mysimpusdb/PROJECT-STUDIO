<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const isCollapsed = ref(false)

const menuItems = [
  { label: 'Dashboard', icon: 'pi pi-home', route: '/' },
  { label: 'Proyek', icon: 'pi pi-folder', route: '/projects' },
  { label: 'Tugas', icon: 'pi pi-check-square', route: '/tasks' },
  { label: 'Gantt Chart', icon: 'pi pi-chart-bar', route: '/gantt' },
  { label: 'Pengaturan', icon: 'pi pi-cog', route: '/settings' }
]

const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
}

const navigate = (path: string) => {
  router.push(path)
}
</script>

<template>
  <aside :class="['sidebar', { 'collapsed': isCollapsed }]">
    <div class="sidebar-header">
      <div class="logo-container" v-if="!isCollapsed">
        <span class="logo-icon pi pi-database"></span>
        <span class="logo-text">DB-STUDIO</span>
      </div>
      <div class="logo-container" v-else>
        <span class="logo-icon pi pi-database"></span>
      </div>
      <button class="toggle-btn" @click="toggleSidebar">
        <i :class="isCollapsed ? 'pi pi-angle-right' : 'pi pi-angle-left'"></i>
      </button>
    </div>

    <div class="sidebar-content">
      <nav class="nav-menu">
        <div class="nav-section">
          <div class="nav-section-title" v-if="!isCollapsed">MENU UTAMA</div>
          <ul>
            <li v-for="item in menuItems" :key="item.route">
              <a 
                href="#" 
                :class="['nav-item', { 'active': route.path === item.route }]"
                @click.prevent="navigate(item.route)"
                :title="isCollapsed ? item.label : ''"
              >
                <i :class="item.icon"></i>
                <span class="nav-label" v-if="!isCollapsed">{{ item.label }}</span>
              </a>
            </li>
          </ul>
        </div>
      </nav>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  width: 260px;
  background-color: var(--bg-surface);
  border-right: 1px solid var(--border-default);
  display: flex;
  flex-direction: column;
  transition: width var(--transition-base);
  height: 100%;
}

.sidebar.collapsed {
  width: 72px;
}

.sidebar-header {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--space-4);
  border-bottom: 1px solid var(--border-default);
}

.logo-container {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  color: var(--primary-400);
}

.logo-icon {
  font-size: 1.25rem;
}

.logo-text {
  font-weight: 700;
  font-size: 1.1rem;
  letter-spacing: 0.5px;
}

.toggle-btn {
  background: transparent;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  padding: var(--space-1);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
}

.toggle-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.sidebar-content {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-4) 0;
}

.nav-section-title {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-tertiary);
  margin: 0 var(--space-4) var(--space-2);
  letter-spacing: 1px;
}

.nav-menu ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: var(--space-3) var(--space-4);
  color: var(--text-secondary);
  text-decoration: none;
  transition: background var(--transition-fast), color var(--transition-fast);
  gap: var(--space-3);
  position: relative;
}

.nav-item i {
  font-size: 1.1rem;
}

.nav-item:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.nav-item.active {
  background: rgba(99, 102, 241, 0.1);
  color: var(--primary-400);
  font-weight: 500;
}

.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background-color: var(--primary-500);
  border-radius: 0 4px 4px 0;
}

.sidebar.collapsed .nav-item {
  justify-content: center;
  padding: var(--space-3) 0;
}

.sidebar.collapsed .nav-item i {
  margin: 0;
}
</style>
