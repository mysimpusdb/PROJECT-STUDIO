<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'

import Button from 'primevue/button'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Calendar from 'primevue/calendar'
import ConfirmDialog from 'primevue/confirmdialog'
import InputNumber from 'primevue/inputnumber'
import { marked } from 'marked'
import GanttContainer from '../components/GanttContainer.vue'
import KanbanBoard from '../components/KanbanBoard.vue'
import TaskList from '../components/TaskList.vue'
import TaskDetailModal from '../components/TaskDetailModal.vue'
import CalendarBoard from '../components/CalendarBoard.vue'
import ProjectReport from '../components/ProjectReport.vue'
import html2canvas from 'html2canvas'
import { jsPDF } from 'jspdf'

import Dropdown from 'primevue/dropdown'
import Dialog from 'primevue/dialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Checkbox from 'primevue/checkbox'
import { GetProjectByID, UpdateProject, ArchiveProject, DeleteProject, GetProjectMembers, AddProjectMember, RemoveProjectMember, UpdateProjectMemberRole, GetWorkspaceUsers } from '../../wailsjs/go/services/ProjectService'
import { GetTasksByProject, CreateTask, UpdateTask, DeleteTask, ExportProjectTasksToExcel } from '../../wailsjs/go/services/TaskService'
import { models } from '../../wailsjs/go/models'

const route = useRoute()
const router = useRouter()
const toast = useToast()
const confirm = useConfirm()

const projectId = route.params.id as string
const project = ref<models.Project | null>(null)
const tasks = ref<models.Task[]>([])
const members = ref<models.ProjectMember[]>([])
const workspaceUsers = ref<models.User[]>([])

const loading = ref(true)
const saveLoading = ref(false)

const editForm = ref({
  name: '',
  description: '',
  color: '',
  startDate: null as Date | null,
  endDate: null as Date | null,
  totalBudget: 0,
  currency: 'IDR'
})
const presetColors = ['#6366F1', '#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6']

const showCreateTaskModal = ref(false)
const createTaskLoading = ref(false)
const newTask = ref({
  name: '',
  description: '',
  startDate: null as Date | null,
  duration: 1
})

const showTaskDetailModal = ref(false)
const selectedTask = ref<models.Task | null>(null)
const isEditingTask = ref(false)
const taskEditForm = ref({
  dependencies: '',
  is_milestone: false,
  description: ''
})

const showInviteModal = ref(false)
const inviteLoading = ref(false)
const inviteForm = ref({
  userId: '',
  role: 'member'
})
const roleOptions = [
  { label: 'Manajer (Manager)', value: 'manager' },
  { label: 'Anggota (Member)', value: 'member' }
]

const handleTaskClick = (task: any) => {
  const localTask = tasks.value.find(t => t.id === task.id)
  if (localTask) {
    selectedTask.value = localTask
    showTaskDetailModal.value = true
  }
}



const loadProject = async () => {
  try {
    const res = await GetProjectByID(projectId)
    if (res.success && res.project) {
      project.value = res.project
      
      editForm.value.name = res.project.name
      editForm.value.description = res.project.description || ''
      editForm.value.color = res.project.color
      editForm.value.totalBudget = res.project.total_budget || 0
      editForm.value.currency = res.project.currency || 'IDR'
      
      if (res.project.start_date) {
        editForm.value.startDate = new Date(res.project.start_date)
      }
      if (res.project.end_date) {
        editForm.value.endDate = new Date(res.project.end_date)
      }
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
      router.push('/projects')
      return
    }

    // Load tasks
    const taskRes = await GetTasksByProject(projectId)
    if (taskRes.success) {
      tasks.value = taskRes.tasks || []
    }
    
    // Load members
    const memberRes = await GetProjectMembers(projectId)
    if (memberRes.success) {
      members.value = memberRes.members || []
    }

  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memuat detail proyek', life: 3000 })
    router.push('/projects')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadProject()
})

const handleInvite = async () => {
  if (!inviteForm.value.userId) return
  inviteLoading.value = true
  try {
    const req = new models.ProjectMemberRequest({
      project_id: projectId,
      user_id: inviteForm.value.userId,
      role: inviteForm.value.role
    })
    const res = await AddProjectMember(req)
    if (res.success && res.member) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: res.message, life: 3000 })
      members.value.push(res.member)
      showInviteModal.value = false
      inviteForm.value.userId = ''
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mengundang anggota', life: 3000 })
  } finally {
    inviteLoading.value = false
  }
}

const openInviteModal = async () => {
  showInviteModal.value = true
  try {
    const res = await GetWorkspaceUsers()
    if (res.success) {
      // Filter out users who are already in the project
      const existingIds = new Set(members.value.map(m => m.user_id))
      workspaceUsers.value = (res.users || []).filter(u => !existingIds.has(u.id))
    }
  } catch (err) {
    console.error(err)
  }
}

const handleRoleChange = async (member: models.ProjectMember, newRole: string) => {
  try {
    const req = new models.ProjectMemberRequest({
      project_id: projectId,
      user_id: member.user_id,
      role: newRole
    })
    const res = await UpdateProjectMemberRole(req)
    if (res.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: res.message, life: 2000 })
      member.role = newRole
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mengubah peran', life: 3000 })
  }
}

const handleRemoveMember = (member: models.ProjectMember) => {
  confirm.require({
    message: `Hapus ${member.user_full_name} dari proyek ini?`,
    header: 'Konfirmasi',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        const res = await RemoveProjectMember(projectId, member.user_id)
        if (res.success) {
          toast.add({ severity: 'success', summary: 'Sukses', detail: res.message, life: 3000 })
          members.value = members.value.filter(m => m.user_id !== member.user_id)
        } else {
          toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
        }
      } catch (err) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal menghapus anggota', life: 3000 })
      }
    }
  })
}

const handleCreateTask = async () => {
  if (!newTask.value.name) {
    toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Nama tugas wajib diisi', life: 3000 })
    return
  }
  
  createTaskLoading.value = true
  try {
    const req = new models.CreateTaskRequest({
      project_id: projectId,
      name: newTask.value.name,
      description: newTask.value.description,
      duration: newTask.value.duration,
      status: 'todo',
      priority: 'medium'
    })

    if (newTask.value.startDate) {
      req.start_date = newTask.value.startDate.toISOString().split('T')[0]
      // End date based on duration
      const endDate = new Date(newTask.value.startDate)
      endDate.setDate(endDate.getDate() + newTask.value.duration)
      req.end_date = endDate.toISOString().split('T')[0]
    }

    const res = await CreateTask(req)
    if (res.success && res.task) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Tugas berhasil dibuat', life: 3000 })
      tasks.value.push(res.task)
      showCreateTaskModal.value = false
      newTask.value = { name: '', description: '', startDate: null, duration: 1 }
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
    }
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal membuat tugas', life: 3000 })
  } finally {
    createTaskLoading.value = false
  }
}

const handleTaskUpdate = async (updatedTask: any) => {
  try {
    const req = new models.UpdateTaskRequest({
      id: updatedTask.id,
      name: updatedTask.name,
      description: updatedTask.description,
      status: updatedTask.status,
      priority: updatedTask.priority,
      progress: updatedTask.progress,
      duration: updatedTask.duration,
      estimated_hours: updatedTask.estimated_hours,
      // '' = unassign; omitting the field keeps the current assignee
      assignee_id: updatedTask.assignee_id ?? '',
      sort_order: updatedTask.sort_order,
      dependencies: updatedTask.dependencies,
      is_milestone: updatedTask.is_milestone,
      labels: updatedTask.labels,
      checklists: updatedTask.checklists,
      estimated_cost: updatedTask.estimated_cost,
      actual_cost: updatedTask.actual_cost
    })
    
    if (updatedTask.start_date) {
      req.start_date = updatedTask.start_date.split('T')[0]
    }
    if (updatedTask.end_date) {
      req.end_date = updatedTask.end_date.split('T')[0]
    }

    const res = await UpdateTask(req)
    if (res.success && res.task) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Tugas diperbarui', life: 2000 })
      // Update local task list
      const idx = tasks.value.findIndex(t => t.id === updatedTask.id)
      if (idx !== -1) {
        tasks.value[idx] = res.task
      }
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
      // Revert if error? For now just fetch tasks again or leave as is.
      await loadProject() 
    }
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memperbarui tugas', life: 3000 })
  }
}

const formatLocalDate = (date: Date) => {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

const handleTaskDateChange = async (task: any, start: Date, end: Date) => {
  const localTask = tasks.value.find(t => t.id === task.id)
  if (localTask) {
    // Calculate duration in days (ensure at least 1)
    const durationMs = end.getTime() - start.getTime()
    const durationDays = Math.max(1, Math.round(durationMs / (1000 * 60 * 60 * 24)))
    
    // Auto-scheduling calculation: delta days
    let oldStart = localTask.start_date ? new Date(localTask.start_date) : null
    let deltaDays = 0
    if (oldStart) {
      deltaDays = Math.round((start.getTime() - oldStart.getTime()) / (1000 * 60 * 60 * 24))
    }

    const updatedTask = {
      ...localTask,
      start_date: formatLocalDate(start),
      end_date: formatLocalDate(end),
      duration: durationDays
    }
    
    await handleTaskUpdate(updatedTask)

    // Recursive function to shift all successors
    const shiftSuccessors = async (taskId: string, deltaDays: number) => {
      const successors = tasks.value.filter(t => {
        if (!t.dependencies) return false
        const deps = t.dependencies.split(',').map(s => s.trim())
        return deps.includes(taskId)
      })

      for (const succ of successors) {
        if (succ.start_date && succ.end_date) {
          const succStart = new Date(succ.start_date)
          const succEnd = new Date(succ.end_date)
          succStart.setDate(succStart.getDate() + deltaDays)
          succEnd.setDate(succEnd.getDate() + deltaDays)
          
          const updatedSucc = {
            ...succ,
            start_date: formatLocalDate(succStart),
            end_date: formatLocalDate(succEnd)
          }
          await handleTaskUpdate(updatedSucc)
          // Recursively shift successors of this successor
          await shiftSuccessors(succ.id, deltaDays)
        }
      }
    }

    if (deltaDays !== 0) {
      await shiftSuccessors(localTask.id, deltaDays)
    }
  }
}

const handleTaskProgressChange = async (task: any, progress: number) => {
  const localTask = tasks.value.find(t => t.id === task.id)
  if (localTask) {
    const updatedTask = {
      ...localTask,
      progress: progress
    }
    await handleTaskUpdate(updatedTask)
  }
}

const handleTaskReorder = async (sourceId: string, targetId: string) => {
  const sourceIndex = tasks.value.findIndex(t => t.id === sourceId)
  const targetIndex = tasks.value.findIndex(t => t.id === targetId)
  if (sourceIndex === -1 || targetIndex === -1) return

  const [movedTask] = tasks.value.splice(sourceIndex, 1)
  tasks.value.splice(targetIndex, 0, movedTask)

  // Update sort_order locally
  const updates = []
  for (let i = 0; i < tasks.value.length; i++) {
    if (tasks.value[i].sort_order !== i) {
      tasks.value[i].sort_order = i
      updates.push(tasks.value[i])
    }
  }

  // Update to backend in background so UI doesn't freeze
  Promise.all(updates.map(t => {
    const req = new models.UpdateTaskRequest({
      id: t.id,
      name: t.name,
      description: t.description,
      status: t.status,
      priority: t.priority,
      progress: t.progress,
      duration: t.duration,
      assignee_id: t.assignee_id ?? '',
      sort_order: t.sort_order
    })
    if (t.start_date) req.start_date = t.start_date.split('T')[0]
    if (t.end_date) req.end_date = t.end_date.split('T')[0]
    return UpdateTask(req)
  })).then(() => {
    toast.add({ severity: 'success', summary: 'Sukses', detail: 'Urutan tugas diperbarui', life: 2000 })
  }).catch(() => {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memperbarui urutan', life: 3000 })
    loadProject() // Reload on error
  })
}




const handleBulkDelete = async (selectedTasks: any[]) => {
  confirm.require({
    message: `Apakah Anda yakin ingin menghapus ${selectedTasks.length} tugas yang dipilih?`,
    header: 'Konfirmasi Hapus',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      let successCount = 0
      for (const task of selectedTasks) {
        try {
          const res = await DeleteTask(task.id)
          if (res.success) successCount++
        } catch (err) {
          console.error(err)
        }
      }
      toast.add({ severity: 'success', summary: 'Sukses', detail: `${successCount} tugas berhasil dihapus`, life: 3000 })
      await loadProject() // Reload list
    }
  })
}

const handleBulkUpdate = async (data: { tasks: any[], status: string }) => {
  let successCount = 0
  for (const task of data.tasks) {
    try {
      const updatedTask = { ...task, status: data.status }
      const req = new models.UpdateTaskRequest({
        id: updatedTask.id,
        name: updatedTask.name,
        description: updatedTask.description,
        status: updatedTask.status,
        priority: updatedTask.priority,
        progress: updatedTask.progress,
        duration: updatedTask.duration,
        assignee_id: updatedTask.assignee_id ?? '',
        sort_order: updatedTask.sort_order,
        dependencies: updatedTask.dependencies,
        is_milestone: updatedTask.is_milestone
      })
      if (updatedTask.start_date) req.start_date = updatedTask.start_date.split('T')[0]
      if (updatedTask.end_date) req.end_date = updatedTask.end_date.split('T')[0]

      const res = await UpdateTask(req)
      if (res.success) successCount++
    } catch (err) {
      console.error(err)
    }
  }
  toast.add({ severity: 'success', summary: 'Sukses', detail: `${successCount} tugas berhasil diperbarui`, life: 3000 })
  await loadProject()
}

const saveTaskDetail = async (updatedTask: any) => {
  if (updatedTask) {
    await handleTaskUpdate(updatedTask)
    showTaskDetailModal.value = false
  }
}


const handleSaveSettings = async () => {
  if (!editForm.value.name) return
  
  saveLoading.value = true
  try {
    const req = new models.UpdateProjectRequest({
      id: projectId,
      name: editForm.value.name,
      description: editForm.value.description,
      color: editForm.value.color,
      status: project.value?.status || 'active',
      total_budget: editForm.value.totalBudget,
      currency: editForm.value.currency
    })

    if (editForm.value.startDate) {
      req.start_date = editForm.value.startDate.toISOString().split('T')[0]
    }
    if (editForm.value.endDate) {
      req.end_date = editForm.value.endDate.toISOString().split('T')[0]
    }

    const res = await UpdateProject(req)
    if (res.success && res.project) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Pengaturan proyek disimpan', life: 3000 })
      project.value = res.project
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal menyimpan proyek', life: 3000 })
  } finally {
    saveLoading.value = false
  }
}

const handleArchive = async () => {
  try {
    const res = await ArchiveProject(projectId)
    if (res.success) {
      toast.add({ severity: 'info', summary: 'Sukses', detail: 'Proyek diarsipkan', life: 3000 })
      if (project.value) project.value.status = 'archived'
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mengarsipkan proyek', life: 3000 })
  }
}

const ganttWrapperRef = ref<HTMLElement | null>(null)

const exportToCSV = () => {
  if (tasks.value.length === 0) return
  const headers = ['ID', 'Name', 'Status', 'Priority', 'Start Date', 'End Date', 'Progress', 'Assignee', 'Dependencies', 'Milestone']
  const rows = tasks.value.map(t => [
    t.id, 
    `"${t.name.replace(/"/g, '""')}"`, 
    t.status, 
    t.priority, 
    t.start_date ? t.start_date.split('T')[0] : '', 
    t.end_date ? t.end_date.split('T')[0] : '', 
    t.progress, 
    t.assignee_id || '', 
    `"${t.dependencies || ''}"`,
    t.is_milestone ? 'Yes' : 'No'
  ])
  
  const csvContent = "data:text/csv;charset=utf-8," + [headers.join(','), ...rows.map(r => r.join(','))].join('\n')
  const encodedUri = encodeURI(csvContent)
  const link = document.createElement('a')
  link.setAttribute('href', encodedUri)
  link.setAttribute('download', `${project.value?.name}_tasks.csv`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

const exportToPNG = async () => {
  // Try to find the gantt element
  const el = document.querySelector('.timeline-panel') as HTMLElement
  if (!el) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Elemen Gantt tidak ditemukan', life: 3000 })
    return
  }
  try {
    const canvas = await html2canvas(el, {
      backgroundColor: '#1E293B', // dark mode bg or fallback
      scale: 2 // better resolution
    })
    const imgData = canvas.toDataURL('image/png')
    const link = document.createElement('a')
    link.download = `${project.value?.name}_gantt.png`
    link.href = imgData
    link.click()
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mengekspor PNG', life: 3000 })
  }
}

const exportTasksToExcel = async () => {
  try {
    const res = await ExportProjectTasksToExcel(projectId)
    if (res.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: res.message, life: 3000 })
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mengekspor tugas ke Excel', life: 3000 })
  }
}

const exportToPDF = async () => {
  const el = document.querySelector('.report-container') as HTMLElement
  if (!el) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Tidak ada laporan untuk diekspor', life: 3000 })
    return
  }
  
  toast.add({ severity: 'info', summary: 'Info', detail: 'Mengekspor laporan ke PDF...', life: 2000 })
  
  try {
    const canvas = await html2canvas(el, {
      scale: 2,
      useCORS: true,
      backgroundColor: '#ffffff'
    })
    
    const imgData = canvas.toDataURL('image/png')
    const pdf = new jsPDF('p', 'mm', 'a4')
    const pdfWidth = pdf.internal.pageSize.getWidth()
    const pdfHeight = (canvas.height * pdfWidth) / canvas.width
    
    pdf.text(`Laporan Proyek: ${project.value?.name}`, 10, 10)
    pdf.addImage(imgData, 'PNG', 10, 20, pdfWidth - 20, pdfHeight - 20)
    pdf.save(`${project.value?.name}_report.pdf`)
    
    toast.add({ severity: 'success', summary: 'Sukses', detail: 'Laporan berhasil diekspor ke PDF', life: 3000 })
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mengekspor PDF', life: 3000 })
  }
}

const confirmDelete = () => {
  confirm.require({
    message: 'Apakah Anda yakin ingin menghapus proyek ini secara permanen? Semua tugas di dalamnya akan ikut terhapus.',
    header: 'Konfirmasi Hapus Proyek',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    acceptLabel: 'Hapus Permanen',
    rejectLabel: 'Batal',
    accept: async () => {
      try {
        const res = await DeleteProject(projectId)
        if (res.success) {
          toast.add({ severity: 'success', summary: 'Sukses', detail: 'Proyek dihapus', life: 3000 })
          router.push('/projects')
        } else {
          toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
        }
      } catch (err) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal menghapus proyek', life: 3000 })
      }
    }
  })
}
</script>

<template>
  <div class="project-detail" v-if="project">
    <div class="header-section">
      <div class="header-left">
        <Button icon="pi pi-arrow-left" class="p-button-text p-button-rounded mr-2" @click="router.push('/projects')" />
        <div class="color-indicator" :style="{ backgroundColor: project.color }"></div>
        <div>
          <h1 class="project-title">{{ project.name }}</h1>
          <div class="status-badge" :class="project.status">{{ project.status }}</div>
        </div>
      </div>
      <div class="header-actions">
        <Button label="Export Excel" icon="pi pi-file-excel" class="p-button-outlined p-button-sm mr-2" @click="exportTasksToExcel" />
        <Button label="Export PDF" icon="pi pi-file-pdf" class="p-button-outlined p-button-sm" @click="exportToPDF" />
      </div>
    </div>

    <TabView class="project-tabs">
      <TabPanel value="0" header="Kanban Board">
        <div class="tab-content" style="background:transparent; border:none; padding-left:0; padding-right:0;">
          <KanbanBoard 
            :tasks="tasks" 
            @task-update="handleTaskUpdate" 
            @task-click="handleTaskClick" 
          />
        </div>
      </TabPanel>
      
      <TabPanel value="1" header="List">
        <div class="tab-content" style="padding: 1rem 0;">
          <div class="flex justify-between items-center mb-4 px-4" style="display:flex; justify-content:space-between; align-items:center; margin-bottom: 1rem; padding: 0 1rem;">
            <h3 style="margin:0; font-size:1.25rem;">Daftar Tugas</h3>
            <Button label="Tambah Tugas" icon="pi pi-plus" class="p-button-sm" @click="showCreateTaskModal = true" />
          </div>
          <TaskList 
            :tasks="tasks" 
            :members="members" 
            @task-click="handleTaskClick" 
            @task-update="handleTaskUpdate"
            @bulk-delete="handleBulkDelete"
            @bulk-update="handleBulkUpdate"
          />
        </div>
      </TabPanel>
      <TabPanel value="2" header="Kalender">
        <div class="tab-content" style="padding: 1rem; height: calc(100vh - 180px);">
          <CalendarBoard 
            :tasks="tasks"
            :members="members"
            @task-click="handleTaskClick"
            @date-click="(d) => { newTask.startDate = d; showCreateTaskModal = true }"
            @task-date-change="handleTaskDateChange"
          />
        </div>
      </TabPanel>


      <TabPanel value="3" header="Gantt Chart">
        <div class="tab-content no-padding">
          <div class="gantt-controls" style="justify-content: space-between;">
            <div class="left-controls">
              <Button label="Export CSV" icon="pi pi-file-excel" class="p-button-sm p-button-outlined p-button-secondary mr-2" @click="exportToCSV" />
              <Button label="Export PNG" icon="pi pi-image" class="p-button-sm p-button-outlined p-button-secondary" @click="exportToPNG" />
            </div>
            <Button label="Tugas Baru" icon="pi pi-plus" class="p-button-sm" @click="showCreateTaskModal = true" />
          </div>
          <GanttContainer 
            :tasks="tasks" 
            @task-update="handleTaskUpdate" 
            @task-date-change="handleTaskDateChange"
            @task-progress-change="handleTaskProgressChange"
            @task-click="handleTaskClick"
            @task-reorder="handleTaskReorder"
          />
        </div>
      </TabPanel>
      
      <TabPanel value="4" header="Anggota">
        <div class="tab-content">
          <div class="flex justify-between items-center mb-4" style="display:flex; justify-content:space-between; align-items:center; margin-bottom: 1rem;">
            <h3 style="margin:0; font-size:1.25rem;">Anggota Proyek</h3>
            <Button label="Undang Anggota" icon="pi pi-user-plus" class="p-button-sm" @click="openInviteModal" />
          </div>
          
          <DataTable :value="members" responsiveLayout="scroll" class="p-datatable-sm" :emptyMessage="'Belum ada anggota yang ditambahkan.'">
            <Column field="user_full_name" header="Nama">
              <template #body="slotProps">
                <div style="display:flex; align-items:center; gap: 0.5rem;">
                  <div class="member-avatar">
                    {{ slotProps.data.user_full_name.charAt(0).toUpperCase() }}
                  </div>
                  <div>
                    <div style="font-weight:600">{{ slotProps.data.user_full_name }}</div>
                    <div style="font-size:0.875rem; color:var(--text-secondary)">{{ slotProps.data.user_email }}</div>
                  </div>
                </div>
              </template>
            </Column>
            <Column field="role" header="Peran">
              <template #body="slotProps">
                <Dropdown 
                  :modelValue="slotProps.data.role" 
                  :options="roleOptions" 
                  optionLabel="label" 
                  optionValue="value" 
                  @change="(e) => handleRoleChange(slotProps.data, e.value)"
                  class="p-dropdown-sm w-full"
                />
              </template>
            </Column>
            <Column header="Aksi" style="width: 100px; text-align:right;">
              <template #body="slotProps">
                <Button icon="pi pi-trash" class="p-button-rounded p-button-danger p-button-text" @click="handleRemoveMember(slotProps.data)" />
              </template>
            </Column>
          </DataTable>
        </div>
      </TabPanel>

      <TabPanel value="5" header="Laporan">
        <div class="tab-content">
          <div class="flex justify-between items-center mb-4" style="display:flex; justify-content:space-between; align-items:center; margin-bottom: 1rem;">
            <h3 style="margin:0; font-size:1.25rem;">Laporan & Analitik</h3>
          </div>
          <ProjectReport :project="project" :tasks="tasks" />
        </div>
      </TabPanel>
      
      <TabPanel value="6" header="Pengaturan">
        <div class="tab-content settings-content">
          <div class="settings-grid">
            <!-- Form Pengaturan Umum -->
            <div class="settings-card">
              <h3>Pengaturan Umum</h3>
              <form @submit.prevent="handleSaveSettings" class="settings-form">
                <div class="field">
                  <label>Nama Proyek</label>
                  <InputText v-model="editForm.name" />
                </div>
                <div class="field">
                  <label>Deskripsi</label>
                  <Textarea v-model="editForm.description" rows="3" />
                </div>
                
                <div class="form-row">
                  <div class="field flex-1">
                    <label>Tanggal Mulai</label>
                    <Calendar v-model="editForm.startDate" dateFormat="yy-mm-dd" showIcon />
                  </div>
                  <div class="field flex-1">
                    <label>Tenggat Waktu</label>
                    <Calendar v-model="editForm.endDate" dateFormat="yy-mm-dd" showIcon />
                  </div>
                </div>

                <div class="form-row mt-4">
                  <div class="field flex-1">
                    <label>Total Budget</label>
                    <InputNumber v-model="editForm.totalBudget" class="w-full" />
                  </div>
                  <div class="field flex-1">
                    <label>Mata Uang</label>
                    <InputText type="text" v-model="editForm.currency" placeholder="IDR / USD" />
                  </div>
                </div>

                <div class="field">
                  <label>Warna Proyek</label>
                  <div class="color-picker-preset">
                    <button 
                      type="button" 
                      v-for="c in presetColors" 
                      :key="c"
                      class="color-circle" 
                      :class="{ 'selected': editForm.color === c }"
                      :style="{ backgroundColor: c }"
                      @click="editForm.color = c"
                    >
                      <i class="pi pi-check" v-if="editForm.color === c"></i>
                    </button>
                  </div>
                </div>

                <Button type="submit" label="Simpan Perubahan" :loading="saveLoading" />
              </form>
            </div>

            <!-- Zona Bahaya (Danger Zone) -->
            <div class="settings-card danger-zone">
              <h3 class="danger-title">Zona Bahaya</h3>
              <p class="danger-desc">Tindakan di bawah ini bersifat destruktif dan akan mempengaruhi proyek.</p>
              
              <div class="danger-action">
                <div>
                  <strong>Arsipkan Proyek</strong>
                  <p>Sembunyikan proyek dari daftar aktif. Anda masih bisa melihatnya di daftar arsip.</p>
                </div>
                <Button label="Arsipkan" class="p-button-outlined p-button-warning" 
                  :disabled="project.status === 'archived'" @click="handleArchive" />
              </div>

              <div class="danger-action border-top">
                <div>
                  <strong>Hapus Proyek</strong>
                  <p>Hapus proyek beserta semua tugas di dalamnya secara permanen.</p>
                </div>
                <Button label="Hapus Permanen" class="p-button-danger" @click="confirmDelete" />
              </div>
            </div>
          </div>
        </div>
      </TabPanel>
    </TabView>

    <!-- Create Task Modal -->
    <Dialog v-model:visible="showCreateTaskModal" modal header="Buat Tugas Baru" :style="{ width: '400px' }">
      <form @submit.prevent="handleCreateTask" class="project-form">
        <div class="field">
          <label>Nama Tugas *</label>
          <InputText v-model="newTask.name" autofocus />
        </div>
        <div class="field">
          <label>Tanggal Mulai</label>
          <Calendar v-model="newTask.startDate" dateFormat="yy-mm-dd" showIcon />
        </div>
        <div class="field">
          <label>Durasi (Hari)</label>
          <input type="number" v-model.number="newTask.duration" min="1" class="p-inputtext p-component w-full" />
        </div>
        <div class="form-actions">
          <Button type="button" label="Batal" class="p-button-text" @click="showCreateTaskModal = false" />
          <Button type="submit" label="Simpan" :loading="createTaskLoading" />
        </div>
      </form>
    </Dialog>

    <!-- Task Detail Modal -->
    <TaskDetailModal 
      v-model:visible="showTaskDetailModal" 
      :task="selectedTask" 
      :members="members"
      @save="saveTaskDetail" 
    />

    <!-- Invite Member Modal -->
    <Dialog v-model:visible="showInviteModal" modal header="Undang Anggota" :style="{ width: '400px' }">
      <form @submit.prevent="handleInvite" class="project-form">
        <div class="field">
          <label>Pilih Pengguna</label>
          <Dropdown v-model="inviteForm.userId" :options="workspaceUsers" optionLabel="full_name" optionValue="id" placeholder="Pilih anggota workspace" class="w-full" filter />
        </div>
        <div class="field">
          <label>Peran</label>
          <Dropdown v-model="inviteForm.role" :options="roleOptions" optionLabel="label" optionValue="value" class="w-full" />
        </div>
        <div class="form-actions">
          <Button type="button" label="Batal" class="p-button-text" @click="showInviteModal = false" />
          <Button type="submit" label="Undang" :loading="inviteLoading" :disabled="!inviteForm.userId" />
        </div>
      </form>
    </Dialog>

    <ConfirmDialog></ConfirmDialog>
  </div>
  <div v-else-if="loading" class="loading-state">
    <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
    <p>Memuat proyek...</p>
  </div>
</template>

<style scoped>
.project-detail {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
  height: 100%;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.mr-2 { margin-right: 0.5rem; }

.color-indicator {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
}

.project-title {
  margin: 0 0 4px 0;
  font-size: 1.75rem;
  color: var(--text-primary);
}

.status-badge {
  display: inline-block;
  font-size: 0.75rem;
  padding: 2px 10px;
  border-radius: var(--radius-full);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.status-badge.active { background: rgba(52, 211, 153, 0.15); color: #34d399; }
.status-badge.archived { background: rgba(156, 163, 175, 0.15); color: #9ca3af; }

.project-tabs {
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* Override PrimeVue tabview base to fill height */
:deep(.p-tabview-panels) {
  flex: 1;
  background: transparent;
  padding: var(--space-6) 0;
  display: flex;
  flex-direction: column;
}
:deep(.p-tabview-panel) {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.tab-content {
  flex: 1;
  background: var(--bg-surface);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  padding: var(--space-6);
  min-height: 400px;
  display: flex;
  flex-direction: column;
}

.tab-content.no-padding {
  padding: 0;
  overflow: hidden;
}

.gantt-controls {
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--border-default);
  background: var(--bg-base);
  display: flex;
  justify-content: flex-end;
}

.empty-content {
  align-items: center;
  justify-content: center;
  text-align: center;
  color: var(--text-secondary);
}
.empty-icon {
  font-size: 3rem;
  color: var(--border-hover);
  margin-bottom: var(--space-4);
}

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
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  margin-top: var(--space-6);
  padding-top: var(--space-4);
  border-top: 1px solid var(--border-subtle);
}

/* Settings Layout */
.settings-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: var(--space-8);
}
@media (max-width: 900px) {
  .settings-grid { grid-template-columns: 1fr; }
}

.settings-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}
.settings-card h3 {
  margin: 0;
  font-size: 1.25rem;
  color: var(--text-primary);
  border-bottom: 1px solid var(--border-subtle);
  padding-bottom: var(--space-2);
}

.settings-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  max-width: 600px;
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
.flex-1 { flex: 1; }

.color-picker-preset {
  display: flex;
  gap: var(--space-3);
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
}
.color-circle.selected {
  transform: scale(1.1);
  border-color: white;
  box-shadow: 0 0 0 2px var(--primary-500);
}

/* Danger Zone */
.danger-zone {
  border: 1px solid rgba(239, 68, 68, 0.3);
  background: rgba(239, 68, 68, 0.05);
  padding: var(--space-4);
  border-radius: var(--radius-lg);
}
.danger-title {
  color: #ef4444 !important;
  border-bottom-color: rgba(239, 68, 68, 0.2) !important;
}
.danger-desc {
  color: var(--text-secondary);
  font-size: 0.9rem;
}
.danger-action {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: var(--space-3);
}
.danger-action.border-top {
  border-top: 1px solid rgba(239, 68, 68, 0.2);
  margin-top: var(--space-2);
}
.danger-action p {
  margin: 4px 0 0 0;
  font-size: 0.85rem;
  color: var(--text-tertiary);
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--text-secondary);
  gap: var(--space-4);
}

:deep(.p-inputtext), :deep(.p-textarea) { width: 100%; }
:deep(.p-calendar) {
  width: 100%;
}

.member-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--primary-500);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.875rem;
}

.task-detail-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  padding-top: var(--space-2);
}
.detail-row {
  display: flex;
  justify-content: space-between;
  border-bottom: 1px solid var(--border-subtle);
  padding-bottom: var(--space-2);
}
.detail-row.description-row {
  flex-direction: column;
  border-bottom: none;
  gap: var(--space-2);
}
.detail-label {
  color: var(--text-secondary);
  font-weight: 500;
}
.description-text {
  margin: 0;
  color: var(--text-primary);
  line-height: 1.5;
  white-space: pre-wrap;
}
.mt-3 { margin-top: 1rem; }
.ml-2 { margin-left: 0.5rem; }
.flex { display: flex; }
.align-items-center { align-items: center; }
</style>
