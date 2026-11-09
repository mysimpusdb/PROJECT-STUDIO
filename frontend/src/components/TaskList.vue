<script setup lang="ts">
import { ref, computed } from 'vue'
import { models } from '../../wailsjs/go/models'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import ProgressBar from 'primevue/progressbar'
import Tag from 'primevue/tag'
import Toolbar from 'primevue/toolbar'
import Dropdown from 'primevue/dropdown'
import Button from 'primevue/button'

const props = defineProps<{
  tasks: models.Task[]
  members: models.ProjectMember[]
}>()

const emit = defineEmits(['task-click', 'task-update', 'bulk-delete', 'bulk-update'])

const selectedTasks = ref<models.Task[]>([])

const getAssigneeName = (assigneeId: string | undefined | null) => {
  if (!assigneeId) return '-'
  const member = props.members.find(m => m.user_id === assigneeId)
  return member ? member.user_full_name : 'Tidak diketahui'
}

const getStatusSeverity = (status: string) => {
  switch(status) {
    case 'todo': return 'secondary'
    case 'in-progress': return 'info'
    case 'done': return 'success'
    default: return 'secondary'
  }
}

const getStatusLabel = (status: string) => {
  switch(status) {
    case 'todo': return 'To Do'
    case 'in-progress': return 'In Progress'
    case 'done': return 'Done'
    default: return status
  }
}

const getPrioritySeverity = (priority: string) => {
  switch(priority) {
    case 'urgent': return 'danger'
    case 'high': return 'warning'
    case 'medium': return 'info'
    case 'low': return 'success'
    default: return 'info'
  }
}

const formatDate = (dateStr: any) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

const filters = ref({
  global: { value: null, matchMode: 'contains' },
})

const handleRowClick = (event: any) => {
  emit('task-click', event.data)
}

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

const onCellEditComplete = (event: any) => {
  let { data, newValue, field } = event;
  
  if (data[field] !== newValue) {
    const updatedTask = { ...data, [field]: newValue }
    emit('task-update', updatedTask)
  }
}

const emitBulkDelete = () => {
  emit('bulk-delete', selectedTasks.value)
  selectedTasks.value = []
}

const emitBulkStatusUpdate = (newStatus: string) => {
  if (newStatus) {
    emit('bulk-update', { tasks: selectedTasks.value, status: newStatus })
  }
}

const bulkStatus = ref(null)
const handleBulkStatusChange = (e: any) => {
  emitBulkStatusUpdate(e.value)
  bulkStatus.value = null // reset dropdown
}

const groupByField = ref<string | null>(null)
const groupByOptions = [
  { label: 'None', value: null },
  { label: 'Status', value: 'status' },
  { label: 'Prioritas', value: 'priority' },
  { label: 'Assignee', value: 'assignee_id' }
]

const getGroupHeaderName = (field: string, data: any) => {
  const val = data[field]
  if (!val) return 'Unassigned / Kosong'
  if (field === 'status') return getStatusLabel(val)
  if (field === 'priority') return val.toUpperCase()
  if (field === 'assignee_id') return getAssigneeName(val)
  return val
}
</script>

<template>
  <div class="task-list-container">
    <div class="flex justify-content-between align-items-center p-3 border-bottom-1 surface-border">
      <div class="flex align-items-center gap-2">
        <span class="font-semibold text-color-secondary">Group By:</span>
        <Dropdown 
          v-model="groupByField" 
          :options="groupByOptions" 
          optionLabel="label" 
          optionValue="value" 
          class="p-dropdown-sm w-12rem"
        />
      </div>
    </div>

    <Toolbar v-if="selectedTasks.length > 0" class="bulk-action-toolbar p-2 border-none">
      <template #start>
        <span class="mr-3 ml-2 font-semibold text-sm">{{ selectedTasks.length }} tugas dipilih</span>
        <Dropdown 
          v-model="bulkStatus"
          placeholder="Ubah Status" 
          :options="statusOptions" 
          optionLabel="label" 
          optionValue="value" 
          class="p-dropdown-sm mr-2"
          @change="handleBulkStatusChange" 
        />
      </template>
      <template #end>
        <Button icon="pi pi-trash" class="p-button-sm p-button-danger p-button-text" label="Hapus" @click="emitBulkDelete" />
      </template>
    </Toolbar>

    <DataTable 
      :value="tasks" 
      responsiveLayout="scroll" 
      class="p-datatable-sm" 
      :paginator="true" 
      :rows="15"
      v-model:filters="filters"
      v-model:selection="selectedTasks"
      filterDisplay="menu"
      :globalFilterFields="['name', 'description', 'status', 'priority']"
      @row-click="handleRowClick"
      editMode="cell"
      @cell-edit-complete="onCellEditComplete"
      rowHover
      stripedRows
      :rowGroupMode="groupByField ? 'subheader' : undefined"
      :groupRowsBy="groupByField || undefined"
      sortMode="single"
      :sortField="groupByField || undefined"
      :sortOrder="1"
    >
      <template #groupheader="slotProps" v-if="groupByField">
        <div class="flex align-items-center gap-2 font-bold py-2">
          <span>{{ getGroupHeaderName(groupByField, slotProps.data) }}</span>
        </div>
      </template>

      <Column selectionMode="multiple" headerStyle="width: 3rem"></Column>

      <Column field="name" header="Nama Tugas" sortable style="min-width: 250px">
        <template #body="slotProps">
          <div class="font-semibold">{{ slotProps.data.name }}</div>
          <div class="text-xs text-gray-500 truncate max-w-xs" v-if="slotProps.data.description">
            {{ slotProps.data.description.substring(0, 50) }}{{ slotProps.data.description.length > 50 ? '...' : '' }}
          </div>
        </template>
      </Column>
      
      <Column field="status" header="Status" sortable style="width: 150px">
        <template #body="slotProps">
          <Tag :severity="getStatusSeverity(slotProps.data.status)" :value="getStatusLabel(slotProps.data.status)"></Tag>
        </template>
        <template #editor="{ data, field }">
          <Dropdown v-model="data[field]" :options="statusOptions" optionLabel="label" optionValue="value" class="w-full" autofocus />
        </template>
      </Column>

      <Column field="priority" header="Prioritas" sortable style="width: 150px">
        <template #body="slotProps">
          <Tag :severity="getPrioritySeverity(slotProps.data.priority)" :value="slotProps.data.priority.toUpperCase()"></Tag>
        </template>
        <template #editor="{ data, field }">
          <Dropdown v-model="data[field]" :options="priorityOptions" optionLabel="label" optionValue="value" class="w-full" autofocus />
        </template>
      </Column>

      <Column field="assignee_id" header="Assignee" sortable style="width: 150px">
        <template #body="slotProps">
          <div class="flex align-items-center gap-2">
            <span v-if="slotProps.data.assignee_id" class="member-avatar-small">
              {{ getAssigneeName(slotProps.data.assignee_id)?.charAt(0).toUpperCase() }}
            </span>
            <span class="text-sm">{{ getAssigneeName(slotProps.data.assignee_id) }}</span>
          </div>
        </template>
      </Column>

      <Column field="start_date" header="Mulai" sortable style="width: 120px">
        <template #body="slotProps">
          <span class="text-sm">{{ formatDate(slotProps.data.start_date) }}</span>
        </template>
      </Column>

      <Column field="end_date" header="Tenggat" sortable style="width: 120px">
        <template #body="slotProps">
          <span class="text-sm">{{ formatDate(slotProps.data.end_date) }}</span>
        </template>
      </Column>

      <Column field="progress" header="Progress" sortable style="width: 150px">
        <template #body="slotProps">
          <div class="flex align-items-center gap-2" style="display:flex; align-items:center; gap:0.5rem">
            <ProgressBar :value="slotProps.data.progress" :showValue="false" style="height: 6px; flex: 1" />
            <span class="text-xs">{{ slotProps.data.progress }}%</span>
          </div>
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<style scoped>
.task-list-container {
  background: var(--bg-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-subtle);
  overflow: hidden;
}

.bulk-action-toolbar {
  background: rgba(99, 102, 241, 0.1);
  border-bottom: 1px solid var(--border-subtle) !important;
}

:deep(.p-datatable-tbody > tr) {
  cursor: pointer;
}

.member-avatar-small {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--primary-500);
  color: white;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.75rem;
}

.text-sm {
  font-size: 0.875rem;
}

.text-xs {
  font-size: 0.75rem;
}

.text-gray-500 {
  color: var(--text-secondary);
}

.truncate {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
