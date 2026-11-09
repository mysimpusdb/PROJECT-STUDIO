<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import Sidebar from './Sidebar.vue'
import TopBar from './TopBar.vue'
import { Heartbeat } from '../../../wailsjs/go/services/PresenceService'

let heartbeatTimer: ReturnType<typeof setInterval> | null = null

const sendHeartbeat = () => {
  Heartbeat().catch(() => {})
}

onMounted(() => {
  sendHeartbeat()
  heartbeatTimer = setInterval(sendHeartbeat, 15000)
})

onUnmounted(() => {
  if (heartbeatTimer) {
    clearInterval(heartbeatTimer)
    heartbeatTimer = null
  }
})
</script>

<template>
  <div class="app-layout">
    <Sidebar />
    <div class="main-container">
      <TopBar />
      <main class="content-area">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  background-color: var(--bg-base);
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0; /* Important for preventing flex children from overflowing */
}

.content-area {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: var(--space-6);
  background-color: var(--bg-base);
}

/* Route transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
