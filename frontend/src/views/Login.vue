<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import { Login, Register } from '../../wailsjs/go/services/AuthService'
import { models } from '../../wailsjs/go/models'

const router = useRouter()
const toast = useToast()

const isRegister = ref(false)
const fullName = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)

const handleSubmit = async () => {
  if (!email.value || !password.value || (isRegister.value && !fullName.value)) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Semua kolom wajib diisi.', life: 3000 })
    return
  }

  loading.value = true
  try {
    const res = isRegister.value
      ? await Register(new models.RegisterRequest({
          full_name: fullName.value,
          email: email.value,
          password: password.value
        }))
      : await Login(new models.LoginRequest({
          email: email.value,
          password: password.value
        }))
    
    if (res.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: res.message, life: 2000 })
      setTimeout(() => {
        router.push('/')
      }, 400)
    } else {
      toast.add({ severity: 'error', summary: 'Gagal', detail: res.message, life: 3000 })
    }
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.message || 'Terjadi kesalahan sistem', life: 3000 })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-mesh"></div>
    <div class="login-card">
      <div class="login-header">
        <h1>PROJECT STUDIO</h1>
        <p>{{ isRegister ? 'Buat Akun Baru di Workspace' : 'Login ke Workspace Anda' }}</p>
      </div>

      <form @submit.prevent="handleSubmit" class="login-form">
        <div v-if="isRegister" class="field">
          <label for="fullName">Nama Lengkap</label>
          <InputText id="fullName" v-model="fullName" type="text" placeholder="Nama lengkap" class="w-full" />
        </div>

        <div class="field">
          <label for="email">Email</label>
          <InputText id="email" v-model="email" type="email" placeholder="admin@dbstudio.local" class="w-full" />
        </div>
        
        <div class="field">
          <label for="password">Kata Sandi</label>
          <InputText id="password" v-model="password" type="password" placeholder="••••••••" class="w-full" />
        </div>

        <Button type="submit" :label="isRegister ? 'Daftar' : 'Masuk'" :loading="loading" class="w-full btn-primary" />

        <div class="mode-switch">
          <span>{{ isRegister ? 'Sudah punya akun?' : 'Belum punya akun?' }}</span>
          <a href="#" @click.prevent="isRegister = !isRegister">
            {{ isRegister ? 'Masuk di sini' : 'Daftar anggota baru' }}
          </a>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--bg-base);
  position: relative;
  overflow: hidden;
}

.login-mesh {
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle at 50% 50%, rgba(99, 102, 241, 0.15) 0%, transparent 40%),
              radial-gradient(circle at 80% 20%, rgba(56, 189, 248, 0.1) 0%, transparent 30%);
  z-index: 0;
  animation: float 20s infinite alternate linear;
}

@keyframes float {
  0% { transform: rotate(0deg) translate(0, 0); }
  100% { transform: rotate(10deg) translate(-2%, 2%); }
}

.login-card {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 420px;
  background: rgba(30, 30, 42, 0.7);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-xl);
  padding: var(--space-8);
  box-shadow: var(--shadow-lg);
}

.login-header {
  text-align: center;
  margin-bottom: var(--space-8);
}

.login-header h1 {
  margin: 0;
  font-size: var(--text-2xl);
  font-weight: 700;
  background: linear-gradient(to right, #F0F0FF, #818CF8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  text-shadow: var(--shadow-glow);
}

.login-header p {
  margin: var(--space-2) 0 0 0;
  color: var(--text-secondary);
  font-size: var(--text-sm);
}

.login-form .field {
  margin-bottom: var(--space-4);
  display: flex;
  flex-direction: column;
}

.login-form label {
  margin-bottom: var(--space-2);
  color: var(--text-primary);
  font-size: var(--text-sm);
  font-weight: 500;
}

.w-full {
  width: 100%;
}

:deep(.p-inputtext) {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-default);
  color: var(--text-primary);
  padding: 10px 14px;
  border-radius: var(--radius-md);
  font-family: var(--font-sans);
  transition: var(--transition-base);
}

:deep(.p-inputtext:focus) {
  border-color: var(--primary-500);
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.2);
  outline: none;
}

:deep(.btn-primary) {
  background: var(--primary-500);
  border: none;
  color: #fff;
  padding: 12px;
  border-radius: var(--radius-md);
  font-weight: 600;
  font-family: var(--font-sans);
  cursor: pointer;
  transition: var(--transition-fast);
  margin-top: var(--space-4);
  display: flex;
  justify-content: center;
}

:deep(.btn-primary:hover) {
  background: var(--primary-400);
  box-shadow: var(--shadow-glow);
}

:deep(.btn-primary:active) {
  transform: scale(0.98);
}

.mode-switch {
  margin-top: var(--space-4);
  text-align: center;
  font-size: var(--text-sm);
  color: var(--text-secondary);
}

.mode-switch a {
  margin-left: 6px;
  color: var(--primary-400);
  text-decoration: none;
  font-weight: 600;
}

.mode-switch a:hover {
  text-decoration: underline;
}
</style>
