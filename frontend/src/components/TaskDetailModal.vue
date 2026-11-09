<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import Calendar from 'primevue/calendar'
import Chips from 'primevue/chips'
import Slider from 'primevue/slider'
import Checkbox from 'primevue/checkbox'
import Avatar from 'primevue/avatar'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { marked } from 'marked'
import { models } from '../../wailsjs/go/models'
import { GetTaskComments, AddComment, UpdateComment, DeleteComment, GetTaskAttachments, SelectAndUploadAttachment, DeleteAttachment, OpenAttachment } from '../../wailsjs/go/services/TaskService'
import { CheckAuth } from '../../wailsjs/go/services/AuthService'

const props = defineProps<{
  visible: boolean
  task: models.Task | null
  members?: models.ProjectMember[]
}>()

const emit = defineEmits(['update:visible', 'save', 'close'])

const toast = useToast()
const confirm = useConfirm()
const currentUserId = ref<string | null>(null)

const activeTab = ref(0)
const editForm = ref({
  description: '',
  progress: 0,
  status: 'todo',
  priority: 'medium',
  assignee_id: null as string | null,
  start_date: null as Date | null,
  end_date: null as Date | null,
  labels: [] as string[],
  checklists: [] as { text: string; done: boolean }[],
  estimated_hours: 0,
  estimated_cost: 0,
  actual_cost: 0
})
const newChecklistItem = ref('')

const statusOptions = [
  { label: 'To Do', value: 'todo' },
  { label: 'In Progress', value: 'in-progress' },
  { label: 'On Hold', value: 'on-hold' },
  { label: 'Done', value: 'done' },
  { label: 'Cancelled', value: 'cancelled' }
]

const priorityOptions = [
  { label: 'Urgent', value: 'urgent' },
  { label: 'High', value: 'high' },
  { label: 'Medium', value: 'medium' },
  { label: 'Low', value: 'low' }
]

const comments = ref<models.TaskComment[]>([])
const newComment = ref('')
const loadingComments = ref(false)
const submittingComment = ref(false)

const attachments = ref<models.TaskAttachment[]>([])
const loadingAttachments = ref(false)
const uploadingAttachment = ref(false)

watch(() => props.visible, (newVal) => {
  if (newVal && props.task) {
    activeTab.value = 0
    editForm.value.description = props.task.description || ''
    editForm.value.progress = props.task.progress || 0
    editForm.value.status = props.task.status || 'todo'
    editForm.value.priority = props.task.priority || 'medium'
    editForm.value.assignee_id = props.task.assignee_id || null
    editForm.value.start_date = props.task.start_date ? new Date(props.task.start_date) : null
    editForm.value.end_date = props.task.end_date ? new Date(props.task.end_date) : null
    editForm.value.estimated_hours = props.task.estimated_hours || (props.task.duration ? props.task.duration * 8 : 0)
    editForm.value.estimated_cost = props.task.estimated_cost || 0
    editForm.value.actual_cost = props.task.actual_cost || 0
    
    try {
      editForm.value.labels = props.task.labels ? JSON.parse(props.task.labels) : []
    } catch { editForm.value.labels = [] }
    
    try {
      editForm.value.checklists = props.task.checklists ? JSON.parse(props.task.checklists) : []
    } catch { editForm.value.checklists = [] }

    CheckAuth().then((res: any) => {
      if (res && res.user) {
        currentUserId.value = res.user.id
      }
    })
    
    loadComments()
    loadAttachments()
  }
})

const addChecklist = () => {
  if (newChecklistItem.value.trim()) {
    editForm.value.checklists.push({ text: newChecklistItem.value.trim(), done: false })
    newChecklistItem.value = ''
  }
}

const removeChecklist = (index: number) => {
  editForm.value.checklists.splice(index, 1)
}

const toDateStr = (d: Date | null) => {
  if (!d) return ''
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

const handleSave = () => {
  if (!props.task) return
  
  const updatedTask = {
    ...props.task,
    description: editForm.value.description,
    progress: editForm.value.progress,
    status: editForm.value.status,
    priority: editForm.value.priority,
    assignee_id: editForm.value.assignee_id,
    start_date: editForm.value.start_date ? toDateStr(editForm.value.start_date) : undefined,
    end_date: editForm.value.end_date ? toDateStr(editForm.value.end_date) : undefined,
    estimated_hours: Number(editForm.value.estimated_hours || 0),
    estimated_cost: Number(editForm.value.estimated_cost || 0),
    actual_cost: Number(editForm.value.actual_cost || 0),
    labels: JSON.stringify(editForm.value.labels),
    checklists: JSON.stringify(editForm.value.checklists)
  }
  
  emit('save', updatedTask)
}

const handleClose = () => {
  emit('update:visible', false)
  emit('close')
}

const loadComments = async () => {
  if (!props.task) return
  loadingComments.value = true
  try {
    const res = await GetTaskComments(props.task.id)
    if (res.success) {
      comments.value = res.comments || []
    }
  } catch (err) {
    console.error("Gagal memuat komentar", err)
  } finally {
    loadingComments.value = false
  }
}

const submitComment = async () => {
  if (!newComment.value.trim() || !props.task) return
  
  submittingComment.value = true
  try {
    const req = new models.CreateCommentRequest({
      task_id: props.task.id,
      content: newComment.value.trim()
    })
    
    const res = await AddComment(req)
    if (res.success && res.comment) {
      comments.value.push(res.comment)
      newComment.value = ''
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Komentar ditambahkan', life: 2000 })
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mengirim komentar', life: 3000 })
  } finally {
    submittingComment.value = false
  }
}

const editingCommentId = ref<string | null>(null)
const editingCommentContent = ref('')

const sanitizeHtml = (raw: string) => {
  if (!raw || typeof DOMParser === 'undefined') return raw
  const doc = new DOMParser().parseFromString(raw, 'text/html')
  doc.querySelectorAll('script, iframe, object, embed, link, meta, style, form').forEach(el => el.remove())
  doc.querySelectorAll('*').forEach(el => {
    for (const attr of Array.from(el.attributes)) {
      const name = attr.name.toLowerCase()
      const val = attr.value.trim().toLowerCase()
      if (name.startsWith('on') || ((name === 'href' || name === 'src' || name === 'xlink:href') && (val.startsWith('javascript:') || val.startsWith('data:text/html')))) {
        el.removeAttribute(attr.name)
      }
    }
  })
  return doc.body.innerHTML
}

const renderMarkdown = (text: string) => {
  if (!text) return ''
  return sanitizeHtml(String(marked.parse(text)))
}

const startEditComment = (c: any) => {
  editingCommentId.value = c.id
  editingCommentContent.value = c.content
}

const cancelEditComment = () => {
  editingCommentId.value = null
  editingCommentContent.value = ''
}

const saveEditComment = async () => {
  if (!editingCommentContent.value.trim() || !editingCommentId.value) return
  try {
    const req = new models.UpdateCommentRequest({
      id: editingCommentId.value,
      content: editingCommentContent.value.trim()
    })
    const res = await UpdateComment(req)
    if (res.success) {
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Komentar diperbarui', life: 2000 })
      cancelEditComment()
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal memperbarui komentar', life: 3000 })
  }
}

const deleteComment = (id: string) => {
  confirm.require({
    message: 'Apakah Anda yakin ingin menghapus komentar ini?',
    header: 'Konfirmasi Hapus',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        const res = await DeleteComment(id)
        if (res.success) {
          toast.add({ severity: 'success', summary: 'Sukses', detail: 'Komentar dihapus', life: 2000 })
        } else {
          toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
        }
      } catch (err) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal menghapus komentar', life: 3000 })
      }
    }
  })
}

// Mention autocomplete logic
const showMentionOverlay = ref(false)
const mentionFilter = ref('')
const mentionOverlayStyles = ref({ top: '0px', left: '0px' })

const handleCommentInput = (e: any) => {
  const val = newComment.value
  const cursorPos = e.target.selectionStart
  const textBeforeCursor = val.substring(0, cursorPos)
  const match = textBeforeCursor.match(/@(\w*)$/)
  
  if (match) {
    showMentionOverlay.value = true
    mentionFilter.value = match[1].toLowerCase()
    
    // Position overlay roughly based on cursor position (simplified)
    mentionOverlayStyles.value = {
      top: '-150px',
      left: '10px'
    }
  } else {
    showMentionOverlay.value = false
  }
}

const filteredMembers = computed(() => {
  if (!props.members) return []
  return props.members.filter(m => 
    (m.user_full_name || '').toLowerCase().includes(mentionFilter.value) ||
    (m.user_full_name || '').replace(/\s/g, '').toLowerCase().includes(mentionFilter.value)
  )
})

const selectMention = (member: any) => {
  const val = newComment.value
  // We need to replace the @xxx part before the cursor with @MemberName
  const match = val.match(/@(\w*)$/)
  if (match) {
    const start = val.lastIndexOf('@' + match[1])
    const newText = val.substring(0, start) + '@' + (member.user_full_name || '').replace(/\s/g, '') + ' ' + val.substring(start + match[0].length)
    newComment.value = newText
  }
  showMentionOverlay.value = false
}

// Wails Event Listeners
let cancelEventsOn: (() => void) | null = null

const subscribeCommentEvents = (taskId?: string) => {
  if (cancelEventsOn && typeof cancelEventsOn === 'function') {
    cancelEventsOn()
    cancelEventsOn = null
  }
  if (!taskId) return
  if ((window as any).runtime && (window as any).runtime.EventsOn) {
    cancelEventsOn = (window as any).runtime.EventsOn('rt:comment:' + taskId, (action: string, data: any) => {
      if (action === 'add') {
        if (!comments.value.find(c => c.id === data.id)) {
          comments.value.push(data)
        }
      } else if (action === 'update') {
        const idx = comments.value.findIndex(c => c.id === data.id)
        if (idx !== -1) {
          comments.value[idx].content = data.content
        }
      } else if (action === 'delete') {
        comments.value = comments.value.filter(c => c.id !== data)
      }
    })
  }
}

watch(() => props.task?.id, (id) => {
  subscribeCommentEvents(id)
}, { immediate: true })

onUnmounted(() => {
  if (cancelEventsOn && typeof cancelEventsOn === 'function') {
    cancelEventsOn()
  }
})

const formatDate = (dateString: any) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('id-ID', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
}

const formatBytes = (bytes: number, decimals = 2) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
}

const loadAttachments = async () => {
  if (!props.task) return
  loadingAttachments.value = true
  try {
    const res = await GetTaskAttachments(props.task.id)
    if (res.success) {
      attachments.value = res.attachments || []
    }
  } catch (err) {
    console.error("Gagal memuat lampiran", err)
  } finally {
    loadingAttachments.value = false
  }
}

const handleUploadAttachment = async () => {
  if (!props.task) return
  uploadingAttachment.value = true
  try {
    const res = await SelectAndUploadAttachment(props.task.id)
    if (res.success && res.attachment) {
      attachments.value.unshift(res.attachment)
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Lampiran berhasil diunggah', life: 2000 })
    } else if (res.message !== "Tidak ada file yang dipilih.") {
      toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal mengunggah lampiran', life: 3000 })
  } finally {
    uploadingAttachment.value = false
  }
}

const handleDeleteAttachment = async (id: string) => {
  try {
    const res = await DeleteAttachment(id)
    if (res.success) {
      attachments.value = attachments.value.filter(a => a.id !== id)
      toast.add({ severity: 'success', summary: 'Sukses', detail: 'Lampiran dihapus', life: 2000 })
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal menghapus lampiran', life: 3000 })
  }
}

const handleOpenAttachment = async (id: string) => {
  try {
    const res = await OpenAttachment(id)
    if (!res.success) {
      toast.add({ severity: 'error', summary: 'Error', detail: res.message, life: 3000 })
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Gagal membuka lampiran', life: 3000 })
  }
}
</script>

<template>
  <Dialog 
    :visible="visible" 
    @update:visible="$emit('update:visible', $event)"
    modal 
    :header="task ? task.name : 'Detail Tugas'" 
    :style="{ width: '60vw', minWidth: '500px' }"
    @hide="handleClose"
  >
    <div v-if="task" class="task-modal-container">

      <TabView v-model:activeIndex="activeTab" class="mt-2">
        <!-- Tab 1: Detail & Metadata -->
        <TabPanel value="0" header="Detail">
          <div class="detail-form">
            <div class="form-grid">
              <div class="field">
                <label>Status</label>
                <Dropdown v-model="editForm.status" :options="statusOptions" optionLabel="label" optionValue="value" class="w-full" />
              </div>
              <div class="field">
                <label>Prioritas</label>
                <Dropdown v-model="editForm.priority" :options="priorityOptions" optionLabel="label" optionValue="value" class="w-full" />
              </div>
              <div class="field">
                <label>Assignee</label>
                <Dropdown v-model="editForm.assignee_id" :options="members || []" optionLabel="user_full_name" optionValue="user_id" placeholder="Pilih Assignee" class="w-full" showClear />
              </div>
              <div class="field">
                <label>Progress ({{ editForm.progress }}%)</label>
                <Slider v-model="editForm.progress" :step="5" class="mt-2" />
              </div>
            </div>

            <div class="form-grid mt-4">
              <div class="field">
                <label>Tanggal Mulai</label>
                <Calendar v-model="editForm.start_date" dateFormat="yy-mm-dd" showIcon class="w-full" />
              </div>
              <div class="field">
                <label>Tenggat Waktu</label>
                <Calendar v-model="editForm.end_date" dateFormat="yy-mm-dd" showIcon class="w-full" />
              </div>
            </div>

            <div class="form-grid mt-4">
              <div class="field">
                <label>Estimasi Jam Kerja</label>
                <InputNumber v-model="editForm.estimated_hours" suffix=" jam" :min="0" class="w-full" />
              </div>
              <div class="field">
                <label>Estimated Cost</label>
                <InputNumber v-model="editForm.estimated_cost" class="w-full" />
              </div>
              <div class="field">
                <label>Actual Cost</label>
                <InputNumber v-model="editForm.actual_cost" class="w-full" />
              </div>
            </div>

            <div class="field mt-4">
              <label>Deskripsi</label>
              <Textarea v-model="editForm.description" rows="4" placeholder="Tambahkan deskripsi detail..." autoResize class="w-full" />
            </div>
            


            <div class="field mt-4">
              <label>Label / Tags</label>
              <Chips v-model="editForm.labels" placeholder="Ketik label lalu tekan Enter..." class="w-full" />
            </div>

            <div class="field mt-4">
              <label>Checklist</label>
              <div class="checklist-container">
                <div v-for="(item, idx) in editForm.checklists" :key="idx" class="checklist-item">
                  <Checkbox v-model="item.done" :binary="true" />
                  <InputText v-model="item.text" class="flex-1 ml-2 p-inputtext-sm" />
                  <Button icon="pi pi-times" class="p-button-rounded p-button-text p-button-danger ml-2" @click="removeChecklist(idx)" />
                </div>
                
                <div class="checklist-add mt-2 flex">
                  <InputText v-model="newChecklistItem" placeholder="Tambah item checklist..." class="flex-1 p-inputtext-sm" @keyup.enter="addChecklist" />
                  <Button icon="pi pi-plus" label="Tambah" class="ml-2 p-button-sm" @click="addChecklist" />
                </div>
              </div>
            </div>
          </div>
        </TabPanel>

        <!-- Tab 2: Comments -->
        <TabPanel value="1" header="Komentar">
          <div class="comments-section">
            <div v-if="loadingComments" class="loading-state">
              <i class="pi pi-spin pi-spinner text-2xl"></i>
              <p>Memuat komentar...</p>
            </div>
            
            <div v-else class="comments-list">
              <div v-if="comments.length === 0" class="empty-comments">
                Belum ada komentar.
              </div>
              <div v-for="c in comments" :key="c.id" class="comment-item">
                <Avatar :label="c.user_name ? c.user_name.charAt(0).toUpperCase() : 'U'" shape="circle" class="comment-avatar" />
                <div class="comment-body">
                  <div class="comment-header">
                    <span class="comment-author">{{ c.user_name }}</span>
                    <div class="comment-meta flex align-items-center">
                      <span class="comment-time">{{ formatDate(c.created_at) }}</span>
                      <div class="comment-actions ml-2" v-if="currentUserId === c.user_id">
                        <Button icon="pi pi-pencil" class="p-button-text p-button-sm p-button-secondary py-0 px-2" @click="startEditComment(c)" />
                        <Button icon="pi pi-trash" class="p-button-text p-button-sm p-button-danger py-0 px-2" @click="deleteComment(c.id)" />
                      </div>
                    </div>
                  </div>
                  
                  <div v-if="editingCommentId === c.id" class="comment-edit">
                    <Textarea v-model="editingCommentContent" rows="2" autoResize class="w-full mb-2" />
                    <div class="flex justify-content-end gap-2">
                      <Button label="Batal" class="p-button-text p-button-sm" @click="cancelEditComment" />
                      <Button label="Simpan" class="p-button-sm" @click="saveEditComment" />
                    </div>
                  </div>
                  
                  <div v-else class="comment-content markdown-body" v-html="renderMarkdown(c.content)"></div>
                </div>
              </div>
            </div>

            <div class="comment-input mt-4" style="position: relative;">
              <div v-if="showMentionOverlay && filteredMembers.length > 0" class="mention-overlay shadow-2 border-round surface-overlay absolute z-5" :style="mentionOverlayStyles">
                <ul class="list-none p-0 m-0">
                  <li v-for="m in filteredMembers" :key="m.user_id" class="p-2 cursor-pointer hover:surface-hover flex align-items-center gap-2" @click="selectMention(m)">
                    <Avatar :label="(m.user_full_name || 'U').charAt(0).toUpperCase()" shape="circle" size="normal" />
                    <span>{{ m.user_full_name || 'User' }}</span>
                  </li>
                </ul>
              </div>
              
              <Textarea 
                v-model="newComment" 
                rows="2" 
                placeholder="Tulis komentar (gunakan Markdown, ketik @ untuk mention)..." 
                autoResize 
                class="w-full" 
                @input="handleCommentInput"
              />
              <div class="flex justify-content-end mt-2">
                <Button label="Kirim Komentar" icon="pi pi-send" @click="submitComment" :loading="submittingComment" :disabled="!newComment.trim()" />
              </div>
            </div>
          </div>
        </TabPanel>
        <!-- Tab 3: Lampiran -->
        <TabPanel value="2" header="Lampiran">
          <div class="attachments-section">
            <div class="flex justify-content-between align-items-center mb-3">
              <h3 class="m-0 text-lg">Daftar Lampiran</h3>
              <Button label="Unggah File" icon="pi pi-upload" class="p-button-sm" @click="handleUploadAttachment" :loading="uploadingAttachment" />
            </div>

            <div v-if="loadingAttachments" class="loading-state">
              <i class="pi pi-spin pi-spinner text-2xl"></i>
              <p>Memuat lampiran...</p>
            </div>
            
            <div v-else class="attachments-list">
              <div v-if="attachments.length === 0" class="empty-state text-center py-4 text-color-secondary italic">
                Belum ada lampiran.
              </div>
              
              <div v-for="att in attachments" :key="att.id" class="attachment-card flex align-items-center p-3 mb-2 border-round surface-hover border-1 border-300">
                <i class="pi pi-file text-3xl text-primary mr-3"></i>
                <div class="flex-1">
                  <div class="font-semibold text-primary cursor-pointer hover:underline" @click="handleOpenAttachment(att.id)">
                    {{ att.file_name }}
                  </div>
                  <div class="text-sm text-color-secondary mt-1">
                    {{ formatBytes(att.file_size) }} • Diunggah oleh {{ att.user_name }} pada {{ formatDate(att.uploaded_at) }}
                  </div>
                </div>
                <div class="actions ml-3 flex gap-2">
                  <Button icon="pi pi-external-link" class="p-button-rounded p-button-text p-button-sm" title="Buka File" @click="handleOpenAttachment(att.id)" />
                  <Button icon="pi pi-trash" class="p-button-rounded p-button-text p-button-danger p-button-sm" title="Hapus File" @click="handleDeleteAttachment(att.id)" />
                </div>
              </div>
            </div>
          </div>
        </TabPanel>
      </TabView>

      <div class="form-actions mt-4">
        <Button label="Batal" class="p-button-text" @click="handleClose" />
        <Button label="Simpan Perubahan" icon="pi pi-save" @click="handleSave" />
      </div>
    </div>
  </Dialog>
</template>

<style scoped>
.mention-overlay {
  background: var(--surface-overlay, #1e1e1e);
  border: 1px solid var(--surface-border, #333);
  max-height: 200px;
  overflow-y: auto;
  width: 250px;
  bottom: 100%;
  margin-bottom: 8px;
}

.comment-actions {
  display: flex;
  opacity: 0;
  transition: opacity 0.2s;
}

.comment-item:hover .comment-actions {
  opacity: 1;
}

:deep(.markdown-body p) { margin: 0 0 0.5rem 0; }
:deep(.markdown-body p:last-child) { margin-bottom: 0; }
:deep(.markdown-body strong) { font-weight: 600; color: var(--text-color); }
:deep(.markdown-body em) { font-style: italic; }

.task-detail-modal :deep(.p-dialog-content) { gap: var(--space-4);
  background: var(--bg-surface-hover);
  padding: var(--space-3);
  border-radius: var(--radius-md);
  margin-bottom: var(--space-4);
}
.task-meta {
  display: flex;
  gap: var(--space-4);
  background: var(--bg-surface-hover);
  padding: var(--space-3);
  border-radius: var(--radius-md);
  margin-bottom: var(--space-4);
}
.meta-item {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  color: var(--text-secondary);
  font-size: 0.9rem;
  font-weight: 500;
}
.detail-form {
  display: flex;
  flex-direction: column;
}
.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--space-4);
}
@media (max-width: 600px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
.field label {
  display: block;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: var(--space-2);
}
.checklist-item {
  display: flex;
  align-items: center;
  margin-bottom: var(--space-2);
}
.empty-comments {
  text-align: center;
  color: var(--text-tertiary);
  padding: var(--space-4) 0;
  font-style: italic;
}
.comments-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  max-height: 400px;
  overflow-y: auto;
  padding-right: var(--space-2);
}
.comment-item {
  display: flex;
  gap: var(--space-3);
}
.comment-avatar {
  background-color: var(--primary-100);
  color: var(--primary-700);
}
.comment-body {
  flex: 1;
  background: var(--bg-surface-hover);
  padding: var(--space-3);
  border-radius: var(--radius-md);
}
.comment-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: var(--space-2);
}
.comment-author {
  font-weight: 600;
  color: var(--text-primary);
}
.comment-time {
  font-size: 0.8rem;
  color: var(--text-tertiary);
}
.comment-content {
  color: var(--text-secondary);
  line-height: 1.5;
  white-space: pre-wrap;
}
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  border-top: 1px solid var(--border-subtle);
  padding-top: var(--space-4);
}
.flex { display: flex; }
.flex-1 { flex: 1; }
.ml-2 { margin-left: 0.5rem; }
.mt-2 { margin-top: 0.5rem; }
.mt-4 { margin-top: 1rem; }
.w-full { width: 100%; }
.justify-content-end { justify-content: flex-end; }
:deep(.p-calendar) { width: 100%; }
</style>
