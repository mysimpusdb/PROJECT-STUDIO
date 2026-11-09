<script setup lang="ts">
import { ref, computed } from 'vue'
import { models } from '../../wailsjs/go/models'

const props = defineProps<{
  tasks: models.Task[]
}>()

const emit = defineEmits(['task-update', 'task-click'])

const columns = [
  { id: 'todo', title: 'To Do', icon: 'pi pi-list' },
  { id: 'in-progress', title: 'In Progress', icon: 'pi pi-spin pi-spinner' },
  { id: 'done', title: 'Done', icon: 'pi pi-check-circle' }
]

// HTML5 Drag and Drop state
const draggedTaskId = ref<string | null>(null)
const dragOverColId = ref<string | null>(null)

const onDragStart = (e: DragEvent, task: models.Task) => {
  draggedTaskId.value = task.id
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/plain', task.id)
  }
}

const onDragOver = (e: DragEvent, colId: string) => {
  e.preventDefault() // Necessary to allow dropping
  dragOverColId.value = colId
}

const onDragLeave = () => {
  dragOverColId.value = null
}

const onDrop = (e: DragEvent, colId: string) => {
  e.preventDefault()
  dragOverColId.value = null
  
  if (draggedTaskId.value) {
    const task = props.tasks.find(t => t.id === draggedTaskId.value)
    if (task && task.status !== colId) {
      const updatedTask = { ...task, status: colId }
      // If moving to done, auto set progress to 100
      if (colId === 'done') updatedTask.progress = 100
      else if (colId === 'todo') updatedTask.progress = 0
      
      emit('task-update', updatedTask)
    }
  }
  draggedTaskId.value = null
}

const getTasksByStatus = (status: string) => {
  return props.tasks.filter(t => t.status === status).sort((a, b) => a.sort_order - b.sort_order)
}

const assigneeClass = (id: string | null | undefined) => {
  if (!id) return ''
  let hash = 0
  for (let i = 0; i < id.length; i++) {
    hash = id.charCodeAt(i) + ((hash << 5) - hash)
  }
  const index = Math.abs(hash) % 6
  return `assignee-color-${index}`
}
</script>

<template>
  <div class="kanban-board">
    <div 
      v-for="col in columns" 
      :key="col.id" 
      class="kanban-col"
      :class="{ 'drag-over': dragOverColId === col.id }"
      @dragover="onDragOver($event, col.id)"
      @dragleave="onDragLeave"
      @drop="onDrop($event, col.id)"
    >
      <div class="col-header">
        <div class="col-title">
          <i :class="col.icon"></i>
          <span>{{ col.title }}</span>
          <span class="count-badge">{{ getTasksByStatus(col.id).length }}</span>
        </div>
      </div>
      
      <div class="col-body">
        <div 
          v-for="task in getTasksByStatus(col.id)" 
          :key="task.id"
          class="task-card"
          :class="{ 'is-dragging': draggedTaskId === task.id }"
          draggable="true"
          @dragstart="onDragStart($event, task)"
          @dragend="draggedTaskId = null"
          @click="emit('task-click', task)"
        >
          <div class="card-header">
            <span class="priority-badge" :class="task.priority">{{ task.priority }}</span>
            <span class="milestone-icon" v-if="task.is_milestone" title="Milestone">◆</span>
          </div>
          
          <h4 class="task-name">{{ task.name }}</h4>
          
          <div class="card-meta">
            <div class="date-info" v-if="task.end_date">
              <i class="pi pi-calendar"></i>
              <span>{{ new Date(task.end_date).toLocaleDateString() }}</span>
            </div>
            
            <div class="assignee-avatar" v-if="task.assignee_id" :class="assigneeClass(task.assignee_id)" :title="task.assignee_id">
              {{ task.assignee_id.substring(0, 2).toUpperCase() }}
            </div>
          </div>
          
          <div class="progress-container">
            <div class="progress-bar-bg">
              <div class="progress-bar-fill" :style="{ width: task.progress + '%' }"></div>
            </div>
            <span class="progress-text">{{ task.progress }}%</span>
          </div>
        </div>
        
        <div v-if="getTasksByStatus(col.id).length === 0" class="empty-col">
          Tarik tugas ke sini
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.kanban-board {
  display: flex;
  gap: var(--space-6);
  height: 100%;
  overflow-x: auto;
  padding-bottom: var(--space-4);
}

.kanban-col {
  flex: 0 0 320px;
  background: var(--bg-surface);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  display: flex;
  flex-direction: column;
  transition: all 0.2s ease;
}

.kanban-col.drag-over {
  background: var(--bg-hover);
  border-color: var(--primary-400);
}

.col-header {
  padding: var(--space-4);
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.col-title {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-weight: 600;
  color: var(--text-primary);
}

.count-badge {
  background: var(--bg-hover);
  color: var(--text-secondary);
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 0.75rem;
  margin-left: 4px;
}

.col-body {
  padding: var(--space-4);
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

/* Task Card */
.task-card {
  background: var(--bg-base);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  padding: var(--space-4);
  cursor: grab;
  box-shadow: var(--shadow-sm);
  transition: transform 0.1s ease, box-shadow 0.2s ease;
}

.task-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--border-hover);
}

.task-card:active {
  cursor: grabbing;
}

.task-card.is-dragging {
  opacity: 0.4;
  transform: scale(0.95);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-2);
}

.priority-badge {
  font-size: 0.65rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 600;
}
.priority-badge.urgent { background: rgba(239, 68, 68, 0.15); color: #ef4444; }
.priority-badge.high { background: rgba(245, 158, 11, 0.15); color: #f59e0b; }
.priority-badge.medium { background: rgba(59, 130, 246, 0.15); color: #3b82f6; }
.priority-badge.low { background: rgba(16, 185, 129, 0.15); color: #10b981; }

.milestone-icon {
  color: var(--primary-500);
  font-size: 0.8rem;
}

.task-name {
  margin: 0 0 var(--space-3) 0;
  font-size: 0.95rem;
  color: var(--text-primary);
  line-height: 1.4;
}

.card-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-3);
}

.date-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.75rem;
  color: var(--text-tertiary);
}

.assignee-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.6rem;
  font-weight: bold;
  color: white;
}

/* Colors */
.assignee-color-0 { background-color: #6366F1; }
.assignee-color-1 { background-color: #3B82F6; }
.assignee-color-2 { background-color: #10B981; }
.assignee-color-3 { background-color: #F59E0B; }
.assignee-color-4 { background-color: #EF4444; }
.assignee-color-5 { background-color: #8B5CF6; }

.progress-container {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.progress-bar-bg {
  flex: 1;
  height: 4px;
  background: var(--border-default);
  border-radius: 2px;
  overflow: hidden;
}

.progress-bar-fill {
  height: 100%;
  background: var(--primary-500);
  border-radius: 2px;
}

.progress-text {
  font-size: 0.7rem;
  color: var(--text-secondary);
  font-weight: 500;
  width: 28px;
  text-align: right;
}

.empty-col {
  border: 2px dashed var(--border-default);
  border-radius: var(--radius-md);
  padding: var(--space-6);
  text-align: center;
  color: var(--text-tertiary);
  font-size: 0.85rem;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100px;
}
</style>
