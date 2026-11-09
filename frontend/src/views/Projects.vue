<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Calendar from 'primevue/calendar'
import Dropdown from 'primevue/dropdown'

import { GetProjects, CreateProject, DuplicateProject, ArchiveProject, DeleteProject } from '../../wailsjs/go/services/ProjectService'
import { models } from '../../wailsjs/go/models'

const router = useRouter()
const toast = useToast()
const projects = ref<models.Project[]>([])
const loading = ref(true)

const showCreateModal = ref(false)
const createLoading = ref(false)

const presetColors = ['#6366F1', '#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6']

const templateOptions = [
  { label: 'Proyek Kosong (Blank)', value: 'blank' },
  { label: 'General Project', value: 'general-project' },
  { label: 'Software Development', value: 'software-development' },
  { label: 'Marketing Campaign', value: 'marketing-campaign' },
  { label: 'Event Planning', value: 'event-planning' },
  { label: 'Construction / Infrastructure', value: 'construction-infrastructure' }
]

// Form state
const newProject = ref({
  name: '',
  description: '',
  color: '#6366F1', // default indigo
  startDate: null as Date | null,
  endDate: null as Date | null,
  template: 'blank'
})

const fetchProjects = async () => {
  loading.value = true
  try {
    const res = await GetProjects()
    if (res.success) {
      projects.value = res.projects || []
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
    }
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.message || 'Gagal memuat proyek', life: 3000 })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchProjects()
})

const handleCreateProject = async () => {
  if (!newProject.value.name) {
    toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Nama proyek wajib diisi', life: 3000 })
    return
  }

  createLoading.value = true
  try {
    const req = new models.CreateProjectRequest({
      name: newProject.value.name,
      description: newProject.value.description,
      color: newProject.value.color,
      template: newProject.value.template
    })

    if (newProject.value.startDate) {
      req.start_date = newProject.value.startDate.toISOString().split('T')[0]
    }
    if (newProject.value.endDate) {
      req.end_date = newProject.value.endDate.toISOString().split('T')[0]
    }

    const res = await CreateProject(req)
    if (res.success && res.project) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Proyek berhasil dibuat', life: 3000 })
      projects.value.unshift(res.project)
      showCreateModal.value = false
      
      // Reset form
      newProject.value = {
        name: '',
        description: '',
        color: '#6366F1',
        startDate: null,
        endDate: null,
        template: 'blank'
      }
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
    }
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.message || 'Gagal membuat proyek', life: 3000 })
  } finally {
    createLoading.value = false
  }
}

const handleDuplicateProject = async (id: string) => {
  try {
    const res = await DuplicateProject(id)
    if (res.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: res.message, life: 3000 })
      await fetchProjects() // Refresh list
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
    }
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.message || 'Gagal menduplikasi proyek', life: 3000 })
  }
}
</script>

<template>
  <div class="projects-page">
    <div class="page-header">
      <div>
        <h1 class="page-title">Proyek</h1>
        <p class="page-subtitle">Kelola semua proyek dalam workspace Anda.</p>
      </div>
      <Button label="Proyek Baru" icon="pi pi-plus" class="p-button-primary" @click="showCreateModal = true" />
    </div>

    <div v-if="loading" class="projects-grid">
      <div class="skeleton-card" v-for="i in 4" :key="i"></div>
    </div>
    
    <div v-else-if="projects.length === 0" class="empty-state">
      <i class="pi pi-folder-open empty-icon"></i>
      <h3>Belum Ada Proyek</h3>
      <p>Mulai dengan membuat proyek pertama Anda untuk melacak tugas.</p>
      <Button label="Buat Proyek" icon="pi pi-plus" class="p-button-outlined mt-3" @click="showCreateModal = true" />
    </div>

    <div v-else class="projects-grid">
      <div class="project-card" v-for="project in projects" :key="project.id" @click="router.push(`/projects/${project.id}`)">
        <div class="card-header">
          <div class="color-badge" :style="{ backgroundColor: project.color }"></div>
          <div class="status-badge" :class="project.status">{{ project.status }}</div>
        </div>
        <h3 class="project-name">{{ project.name }}</h3>
        <p class="project-desc">{{ project.description || 'Tidak ada deskripsi' }}</p>
        
        <div class="card-footer">
          <div class="stats-mini">
            <span title="Tasks"><i class="pi pi-check-square"></i> 0</span>
            <span title="Members"><i class="pi pi-users"></i> 1</span>
          </div>
          <div style="display: flex;">
            <button class="btn-action" title="Duplikasi" @click.stop="handleDuplicateProject(project.id!)">
              <i class="pi pi-copy"></i>
            </button>
            <button class="btn-action" title="Pengaturan">
              <i class="pi pi-cog"></i>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Project Modal -->
    <Dialog v-model:visible="showCreateModal" modal header="Buat Proyek Baru" :style="{ width: '50vw', minWidth: '400px' }">
      <form @submit.prevent="handleCreateProject" class="project-form">
        <div class="field">
          <label for="name">Nama Proyek *</label>
          <InputText id="name" v-model="newProject.name" placeholder="Contoh: Website Redesign" autofocus />
        </div>
        
        <div class="field">
          <label for="desc">Deskripsi</label>
          <Textarea id="desc" v-model="newProject.description" rows="3" placeholder="Deskripsi singkat tentang proyek ini..." />
        </div>

        <div class="form-row">
          <div class="field flex-1">
            <label for="start">Tanggal Mulai</label>
            <Calendar id="start" v-model="newProject.startDate" dateFormat="yy-mm-dd" showIcon />
          </div>
          <div class="field flex-1">
            <label for="end">Tenggat Waktu</label>
            <Calendar id="end" v-model="newProject.endDate" dateFormat="yy-mm-dd" showIcon />
          </div>
        </div>

        <div class="field">
          <label for="template">Template Proyek</label>
          <Dropdown id="template" v-model="newProject.template" :options="templateOptions" optionLabel="label" optionValue="value" class="w-full" />
        </div>

        <div class="field">
          <label>Warna Proyek</label>
          <div class="color-picker-preset">
            <button 
              type="button" 
              v-for="c in presetColors" 
              :key="c"
              class="color-circle" 
              :class="{ 'selected': newProject.color === c }"
              :style="{ backgroundColor: c }"
              @click="newProject.color = c"
            >
              <i class="pi pi-check" v-if="newProject.color === c"></i>
            </button>
          </div>
        </div>
        
        <div class="form-actions">
          <Button type="button" label="Batal" class="p-button-text" @click="showCreateModal = false" />
          <Button type="submit" label="Buat Proyek" :loading="createLoading" />
        </div>
      </form>
    </Dialog>
  </div>
</template>

<style scoped>
.projects-page {
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

.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--space-6);
}

.project-card {
  background: var(--bg-surface);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  display: flex;
  flex-direction: column;
  box-shadow: var(--shadow-sm);
  transition: transform var(--transition-fast), box-shadow var(--transition-fast);
  cursor: pointer;
}

.project-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--primary-500);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-3);
}

.color-badge {
  width: 32px;
  height: 8px;
  border-radius: var(--radius-full);
}

.status-badge {
  font-size: 0.75rem;
  padding: 2px 8px;
  border-radius: var(--radius-full);
  font-weight: 600;
  text-transform: capitalize;
}

.status-badge.active { background: rgba(52, 211, 153, 0.15); color: #34d399; }
.status-badge.archived { background: rgba(156, 163, 175, 0.15); color: #9ca3af; }

.project-name {
  margin: 0 0 var(--space-2) 0;
  font-size: 1.15rem;
  font-weight: 600;
  color: var(--text-primary);
}

.project-desc {
  margin: 0 0 var(--space-4) 0;
  color: var(--text-secondary);
  font-size: 0.875rem;
  line-height: 1.5;
  flex: 1;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: var(--space-4);
  border-top: 1px solid var(--border-subtle);
}

.stats-mini {
  display: flex;
  gap: var(--space-4);
  color: var(--text-tertiary);
  font-size: 0.875rem;
}

.stats-mini span {
  display: flex;
  align-items: center;
  gap: var(--space-1);
}

.btn-action {
  background: transparent;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  padding: var(--space-1);
  border-radius: var(--radius-sm);
  transition: all var(--transition-fast);
}

.btn-action:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.skeleton-card {
  height: 180px;
  background: var(--bg-surface);
  border-radius: var(--radius-lg);
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: .5; }
}

.empty-state {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: var(--space-12) 0;
  background: var(--bg-surface);
  border: 1px dashed var(--border-default);
  border-radius: var(--radius-lg);
  color: var(--text-secondary);
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: var(--space-4);
  color: var(--border-default);
}

.mt-3 {
  margin-top: var(--space-4);
}

/* Modal Form Styles */
.project-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  margin-top: var(--space-4);
}

.field {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.field label {
  font-weight: 500;
  color: var(--text-secondary);
}

.form-row {
  display: flex;
  gap: var(--space-4);
}

.flex-1 {
  flex: 1;
}

.color-picker-preset {
  display: flex;
  gap: var(--space-3);
  margin-top: var(--space-2);
}

.color-circle {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  transition: transform 0.2s ease;
}

.color-circle.selected {
  transform: scale(1.1);
  border-color: white;
  box-shadow: 0 0 0 2px var(--primary-500);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  margin-top: var(--space-6);
  padding-top: var(--space-4);
  border-top: 1px solid var(--border-subtle);
}

:deep(.p-inputtext), :deep(.p-textarea) {
  width: 100%;
  background: var(--bg-base);
  border: 1px solid var(--border-default);
  color: var(--text-primary);
}

:deep(.p-calendar) {
  width: 100%;
}
</style>
