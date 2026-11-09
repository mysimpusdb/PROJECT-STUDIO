<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import InputText from 'primevue/inputtext'
import Avatar from 'primevue/avatar'
import Menu from 'primevue/menu'
import OverlayPanel from 'primevue/overlaypanel'
import Badge from 'primevue/badge'
import Button from 'primevue/button'
import { Logout, CheckAuth } from '../../../wailsjs/go/services/AuthService'
import { GetMyNotifications, MarkAsRead, MarkAllAsRead, ClearAll } from '../../../wailsjs/go/services/NotificationService'
import { EventsOn, EventsOff } from '../../../wailsjs/runtime/runtime'

const router = useRouter()
const toast = useToast()

const searchQuery = ref('')
const menu = ref()
const notifPanel = ref()

const currentUser = ref<any>(null)
const unreadCount = ref(0)
const notifications = ref<any[]>([])

onMounted(async () => {
  try {
    const auth = await CheckAuth()
    if (auth.is_logged_in && auth.user) {
      currentUser.value = auth.user
      
      // Load initial notifications
      await loadNotifications()
      
      // Listen for new notifications
      EventsOn(`rt:notification:new:${auth.user.id}`, (newNotif: any) => {
        notifications.value.unshift(newNotif)
        unreadCount.value++
        
        toast.add({ 
          severity: 'info', 
          summary: newNotif.title, 
          detail: newNotif.message, 
          life: 4000 
        })
      })
      
      EventsOn(`rt:notification:count:${auth.user.id}`, (count: number) => {
        unreadCount.value = count
      })
    }
  } catch (err) {
    console.error("Failed to setup notifications", err)
  }
})

onUnmounted(() => {
  if (currentUser.value) {
    EventsOff(`rt:notification:new:${currentUser.value.id}`)
    EventsOff(`rt:notification:count:${currentUser.value.id}`)
  }
})

const loadNotifications = async () => {
  try {
    const res = await GetMyNotifications(20, 0)
    if (res.success) {
      notifications.value = res.notifications || []
      unreadCount.value = res.unread_count || 0
    }
  } catch (err) {
    console.error(err)
  }
}

const toggleNotifPanel = (event: Event) => {
  notifPanel.value.toggle(event)
}

const onMarkAsRead = async (id: string) => {
  try {
    const res = await MarkAsRead(id)
    if (res.success) {
      const n = notifications.value.find(x => x.id === id)
      if (n && !n.is_read) {
        n.is_read = true
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
    }
  } catch (err) {
    console.error(err)
  }
}

const onMarkAllAsRead = async () => {
  try {
    const res = await MarkAllAsRead()
    if (res.success) {
      notifications.value.forEach(n => n.is_read = true)
      unreadCount.value = 0
    }
  } catch (err) {
    console.error(err)
  }
}

const onClearAll = async () => {
  try {
    const res = await ClearAll()
    if (res.success) {
      notifications.value = []
      unreadCount.value = 0
      notifPanel.value.hide()
    }
  } catch (err) {
    console.error(err)
  }
}

const toggleMenu = (event: Event) => {
  menu.value.toggle(event)
}

const userMenuItems = ref([
  {
    label: 'Profil Saya',
    icon: 'pi pi-user',
    command: () => router.push('/settings/profile')
  },
  {
    label: 'Pengaturan',
    icon: 'pi pi-cog',
    command: () => router.push('/settings')
  },
  {
    separator: true
  },
  {
    label: 'Logout',
    icon: 'pi pi-sign-out',
    command: async () => {
      try {
        await Logout()
        router.push('/login')
      } catch (err) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal logout.', life: 3000 })
      }
    }
  }
])
</script>

<template>
  <header class="topbar">
    <div class="topbar-left">
      <div class="search-container">
        <i class="pi pi-search search-icon"></i>
        <InputText v-model="searchQuery" placeholder="Cari proyek, tugas..." class="search-input" />
      </div>
    </div>
    
    <div class="topbar-right">
      <button class="action-btn" title="Notifikasi" @click="toggleNotifPanel" aria-haspopup="true" aria-controls="notif_panel">
        <i class="pi pi-bell"></i>
        <Badge v-if="unreadCount > 0" :value="unreadCount" severity="danger" class="notif-badge"></Badge>
      </button>

      <OverlayPanel ref="notifPanel" id="notif_panel" style="width: 350px">
        <div class="notif-header flex justify-content-between align-items-center mb-3 pb-2 border-bottom-1 surface-border">
          <span class="font-bold">Notifikasi</span>
          <div class="flex gap-2">
            <Button label="Tandai Dibaca" class="p-button-text p-button-sm p-0 text-xs" @click="onMarkAllAsRead" v-if="unreadCount > 0" />
            <Button label="Hapus Semua" class="p-button-text p-button-danger p-button-sm p-0 text-xs" @click="onClearAll" v-if="notifications.length > 0" />
          </div>
        </div>
        
        <div class="notif-list" style="max-height: 400px; overflow-y: auto;">
          <div v-if="notifications.length === 0" class="text-center p-4 text-color-secondary text-sm">
            Belum ada notifikasi baru.
          </div>
          <div 
            v-else 
            v-for="notif in notifications" 
            :key="notif.id" 
            class="notif-item p-2 mb-2 border-round cursor-pointer hover:surface-hover transition-colors flex align-items-start gap-3"
            :class="{ 'surface-100': !notif.is_read }"
            @click="onMarkAsRead(notif.id)"
          >
            <div class="notif-icon mt-1">
              <i class="pi pi-check-circle text-primary" v-if="notif.type === 'task_assigned'"></i>
              <i class="pi pi-comment text-green-500" v-else-if="notif.type === 'comment'"></i>
              <i class="pi pi-at text-orange-500" v-else-if="notif.type === 'mention'"></i>
              <i class="pi pi-bell text-blue-500" v-else></i>
            </div>
            <div class="notif-content flex-1">
              <div class="notif-title font-semibold text-sm mb-1" :class="{ 'text-color-secondary': notif.is_read }">{{ notif.title }}</div>
              <div class="notif-message text-xs text-color-secondary line-height-2">{{ notif.message }}</div>
              <div class="notif-time text-xs text-color-secondary mt-1">
                {{ new Date(notif.created_at).toLocaleDateString('id-ID', { hour: '2-digit', minute: '2-digit' }) }}
              </div>
            </div>
            <div class="notif-indicator" v-if="!notif.is_read">
              <div class="w-1rem h-1rem border-circle bg-primary-500"></div>
            </div>
          </div>
        </div>
      </OverlayPanel>
      
      <div class="user-profile" @click="toggleMenu" aria-haspopup="true" aria-controls="user_menu">
        <Avatar label="A" class="user-avatar" shape="circle" />
        <span class="user-name">Admin</span>
        <i class="pi pi-angle-down"></i>
      </div>
      
      <Menu ref="menu" id="user_menu" :model="userMenuItems" :popup="true" />
    </div>
  </header>
</template>

<style scoped>
.topbar {
  height: 64px;
  background-color: var(--bg-surface);
  border-bottom: 1px solid var(--border-default);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--space-6);
  z-index: 10;
}

.topbar-left {
  flex: 1;
  max-width: 400px;
}

.search-container {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  color: var(--text-tertiary);
  font-size: 0.9rem;
}

.search-input {
  width: 100%;
  padding-left: 36px;
  background: var(--bg-base);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-full);
  height: 36px;
  font-size: 0.9rem;
  color: var(--text-primary);
  transition: border-color var(--transition-base);
}

.search-input:focus {
  border-color: var(--primary-500);
  outline: none;
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.action-btn {
  background: transparent;
  border: none;
  color: var(--text-secondary);
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  position: relative;
  transition: background var(--transition-fast);
}

.action-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.notif-badge {
  position: absolute;
  top: 0px;
  right: 0px;
  transform: scale(0.75);
}

.notif-item {
  border: 1px solid transparent;
}
.notif-item:hover {
  border-color: var(--border-subtle);
}

.user-profile {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  cursor: pointer;
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-md);
  transition: background var(--transition-fast);
}

.user-profile:hover {
  background: var(--bg-hover);
}

.user-avatar {
  background-color: var(--primary-500);
  color: white;
  width: 32px;
  height: 32px;
  font-size: 0.9rem;
}

.user-name {
  font-weight: 500;
  font-size: 0.9rem;
  color: var(--text-primary);
}

.user-profile i {
  font-size: 0.8rem;
  color: var(--text-tertiary);
}

:deep(.p-menu) {
  background: var(--bg-surface);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-md);
}

:deep(.p-menu .p-menuitem-link) {
  padding: 0.75rem 1rem;
  color: var(--text-secondary);
  transition: background var(--transition-fast);
}

:deep(.p-menu .p-menuitem-link:hover) {
  background: var(--bg-hover);
  color: var(--text-primary);
}

:deep(.p-menu .p-menuitem-icon) {
  color: var(--text-tertiary);
  margin-right: 0.5rem;
}
</style>
