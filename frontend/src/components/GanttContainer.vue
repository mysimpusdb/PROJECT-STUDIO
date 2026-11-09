<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import Gantt from 'frappe-gantt'
import '../../node_modules/frappe-gantt/dist/frappe-gantt.css'
import { models } from '../../wailsjs/go/models'

const props = defineProps<{
  tasks: models.Task[]
}>()

const emit = defineEmits(['task-click', 'task-date-change', 'task-progress-change', 'task-update', 'task-reorder'])

const ganttTarget = ref<HTMLElement | null>(null)
const taskListRef = ref<HTMLElement | null>(null)
const timelineRef = ref<HTMLElement | null>(null)
let ganttInstance: any = null

let isSyncingLeft = false
let isSyncingRight = false

const handleTaskListScroll = () => {
  if (!taskListRef.value || !timelineRef.value) return
  if (isSyncingLeft) {
    isSyncingLeft = false
    return
  }
  isSyncingRight = true
  timelineRef.value.scrollTop = taskListRef.value.scrollTop
}

const handleTimelineScroll = () => {
  if (!taskListRef.value || !timelineRef.value) return
  if (isSyncingRight) {
    isSyncingRight = false
    return
  }
  isSyncingLeft = true
  taskListRef.value.scrollTop = timelineRef.value.scrollTop
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

const initGantt = () => {
  if (!ganttTarget.value || props.tasks.length === 0) return

  // Format tasks for Frappe Gantt
  const formattedTasks = props.tasks.map(t => ({
    id: t.id,
    name: t.name,
    start: t.start_date || new Date().toISOString().split('T')[0],
    end: t.end_date || new Date(new Date().getTime() + 86400000 * (t.duration || 1)).toISOString().split('T')[0],
    progress: t.progress,
    dependencies: t.dependencies || '',
    custom_class: `priority-${t.priority} ${t.is_milestone ? 'milestone' : ''} ${assigneeClass(t.assignee_id)}`
  }))

  if (ganttInstance) {
    ganttInstance.refresh(formattedTasks)
  } else {
    ganttInstance = new Gantt(ganttTarget.value, formattedTasks, {
      on_click: (task: any) => emit('task-click', task),
      on_date_change: (task: any, start: Date, end: Date) => emit('task-date-change', task, start, end),
      on_progress_change: (task: any, progress: number) => emit('task-progress-change', task, progress),
      on_view_change: (mode: string) => {
        // e.g., Quarter Day, Half Day, Day, Week, Month
      },
      view_mode: 'Day',
      language: 'en',
      custom_popup_html: function(args: any) {
        // frappe-gantt 1.2.2 passes an object { task, chart, ... } or just task depending on version
        const task = args.task ? args.task : args;
        return `
          <div class="custom-gantt-popup">
            <h5>${task.name}</h5>
            <p>Progress: ${task.progress}%</p>
          </div>
        `
      }
    })
  }
  
  // Try to highlight weekends
  setTimeout(() => {
    highlightWeekends()
  }, 100)
}

const highlightWeekends = () => {
  if (!ganttTarget.value) return
  // Depending on Frappe Gantt version, ticks might not be easily parseable for dates.
  // We look for .lower-text which usually holds the day number or date in some modes
  const svg = ganttTarget.value.querySelector('svg')
  if (!svg) return
  // We can't perfectly tint weekends without knowing the exact date of each column.
  // We will leave this stubbed out for now until we use a more robust custom Gantt library
  // or patch frappe-gantt to emit date info on ticks.
}

watch(() => props.tasks, () => {
  initGantt()
}, { deep: true })

import { computed } from 'vue'

// Tree Logic (Expand/Collapse)
const expandedTasks = ref<Set<string>>(new Set())

const toggleExpand = (taskId: string) => {
  const newSet = new Set(expandedTasks.value)
  if (newSet.has(taskId)) {
    newSet.delete(taskId)
  } else {
    newSet.add(taskId)
  }
  expandedTasks.value = newSet
}

// Compute tree properties: depth, hasChildren, isVisible
const visibleTasks = computed(() => {
  const result: any[] = []
  
  // Create a map for quick lookup
  const taskMap = new Map()
  props.tasks.forEach(t => taskMap.set(t.id, { ...t, children: [] }))
  
  // Build tree
  const rootTasks: any[] = []
  props.tasks.forEach(t => {
    if (t.parent_id && taskMap.has(t.parent_id)) {
      taskMap.get(t.parent_id).children.push(taskMap.get(t.id))
    } else {
      rootTasks.push(taskMap.get(t.id))
    }
  })

  // Traverse tree to flatten with depth
  const traverse = (tasks: any[], depth: number, parentExpanded: boolean) => {
    for (const t of tasks) {
      const hasChildren = t.children.length > 0
      const isExpanded = expandedTasks.value.has(t.id)
      
      // If no tasks are explicitly expanded yet, we can default to true by initializing the set differently, 
      // but let's just make it user controlled. Let's auto-expand root nodes by default if they are empty
      
      if (parentExpanded) {
        result.push({
          ...t,
          depth,
          hasChildren,
          isExpanded
        })
        
        if (isExpanded) {
          traverse(t.children, depth + 1, true)
        }
      }
    }
  }

  // To make all root tasks expanded by default on first load:
  // if expandedTasks is empty, we could populate it, but better to just let traverse handle it.
  traverse(rootTasks, 0, true)
  return result
})

// Initialize all root tasks as expanded on first load
watch(() => props.tasks, (newTasks, oldTasks) => {
  if (oldTasks && oldTasks.length === 0 && newTasks.length > 0) {
    const rootIds = newTasks.filter(t => !t.parent_id).map(t => t.id)
    expandedTasks.value = new Set(rootIds as string[])
  }
  initGantt()
}, { deep: true })

// Inline Edit Logic
const editingTaskId = ref<string | null>(null)
const editNameValue = ref('')

const startEditName = (task: any) => {
  editingTaskId.value = task.id
  editNameValue.value = task.name
}

const saveEditName = (task: any) => {
  if (editingTaskId.value === task.id) {
    if (editNameValue.value.trim() !== '' && editNameValue.value !== task.name) {
      const updatedTask = { ...task, name: editNameValue.value.trim() }
      emit('task-update', updatedTask)
    }
    editingTaskId.value = null
  }
}

// Drag and Drop Reorder
const draggedTaskId = ref<string | null>(null)
const dragOverTaskId = ref<string | null>(null)

const onDragStart = (e: DragEvent, task: any) => {
  draggedTaskId.value = task.id
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/plain', task.id)
  }
}

const onDragOver = (e: DragEvent, task: any) => {
  e.preventDefault()
  if (draggedTaskId.value !== task.id) {
    dragOverTaskId.value = task.id
  }
}

const onDrop = (e: DragEvent, targetTask: any) => {
  e.preventDefault()
  dragOverTaskId.value = null
  const sourceId = draggedTaskId.value
  draggedTaskId.value = null

  if (sourceId && sourceId !== targetTask.id) {
    emit('task-reorder', sourceId, targetTask.id)
  }
}

onMounted(() => {
  if (props.tasks.length > 0) {
    const rootIds = props.tasks.filter(t => !t.parent_id).map(t => t.id)
    expandedTasks.value = new Set(rootIds as string[])
  }
  initGantt()
})
</script>

<template>
  <div class="gantt-wrapper">
    <div class="split-pane">
      <!-- Task List Panel -->
      <div class="task-list" ref="taskListRef" @scroll="handleTaskListScroll">
        <div class="task-list-header">
          <div class="col-name">Task Name</div>
          <div class="col-date">Start</div>
          <div class="col-date">End</div>
          <div class="col-progress">Progress</div>
        </div>
        <div class="task-list-body">
          <div v-for="task in visibleTasks" :key="task.id" class="task-row"
               :class="{ 'drag-over': dragOverTaskId === task.id, 'is-dragging': draggedTaskId === task.id }"
               draggable="true"
               @dragstart="onDragStart($event, task)"
               @dragover="onDragOver($event, task)"
               @dragleave="dragOverTaskId = null"
               @drop="onDrop($event, task)">
            
            <div class="drag-handle"><i class="pi pi-bars"></i></div>
            <div class="col-name text-ellipsis" 
                 :title="task.name" 
                 :style="{ paddingLeft: (task.depth * 16) + 'px' }"
                 @dblclick="startEditName(task)">
              
              <span class="expand-icon" :class="{ invisible: !task.hasChildren }" @click.stop="toggleExpand(task.id)">
                <i :class="task.isExpanded ? 'pi pi-chevron-down' : 'pi pi-chevron-right'"></i>
              </span>
              
              <input v-if="editingTaskId === task.id" 
                     v-model="editNameValue" 
                     class="inline-edit-input"
                     @blur="saveEditName(task)"
                     @keyup.enter="saveEditName(task)"
                     @keyup.esc="editingTaskId = null"
                     autofocus
                     type="text" />
              <span v-else>{{ task.name }}</span>
            </div>
            <div class="col-date">{{ task.start_date ? new Date(task.start_date).toLocaleDateString() : '-' }}</div>
            <div class="col-date">{{ task.end_date ? new Date(task.end_date).toLocaleDateString() : '-' }}</div>
            <div class="col-progress">
              <div class="progress-bar-bg">
                <div class="progress-bar-fill" :style="{ width: task.progress + '%' }"></div>
              </div>
            </div>
          </div>
          
          <div v-if="tasks.length === 0" class="empty-state-list">
            Belum ada tugas. Buat tugas baru untuk melihat Gantt.
          </div>
        </div>
      </div>

      <!-- Timeline Panel -->
      <div class="timeline-panel" ref="timelineRef" @scroll="handleTimelineScroll">
        <div ref="ganttTarget"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.gantt-wrapper {
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  overflow: hidden;
  background: var(--bg-surface);
  display: flex;
  flex-direction: column;
  height: 100%;
}

.split-pane {
  display: flex;
  height: 100%;
  overflow: hidden;
}

.task-list {
  width: 400px;
  flex-shrink: 0;
  border-right: 1px solid var(--border-default);
  display: flex;
  flex-direction: column;
  background: var(--bg-base);
  overflow-y: auto;
}

.task-list-header {
  display: flex;
  padding: 12px 16px;
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border-default);
  font-weight: 600;
  font-size: 0.85rem;
  color: var(--text-secondary);
  position: sticky;
  top: 0;
  z-index: 10;
}

.task-list-body {
  flex: 1;
}

.task-row {
  display: flex;
  padding: 10px 16px;
  border-bottom: 1px solid var(--border-subtle);
  font-size: 0.875rem;
  align-items: center;
}

.task-row:hover {
  background: var(--bg-hover);
}

.task-row.is-dragging {
  opacity: 0.5;
}
.task-row.drag-over {
  border-top: 2px solid var(--primary-500);
}

.drag-handle {
  width: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--border-default);
  cursor: grab;
  margin-right: 8px;
}
.task-row:hover .drag-handle {
  color: var(--text-secondary);
}

.col-name {
  flex: 1;
  min-width: 0;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 8px;
}

.expand-icon {
  cursor: pointer;
  width: 16px;
  height: 16px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
}

.expand-icon.invisible {
  visibility: hidden;
}

.inline-edit-input {
  flex: 1;
  min-width: 0;
  background: var(--bg-hover);
  border: 1px solid var(--primary-500);
  color: var(--text-primary);
  border-radius: 4px;
  padding: 2px 4px;
  font-size: 0.875rem;
  outline: none;
}
.col-date {
  width: 80px;
  text-align: center;
  color: var(--text-tertiary);
  font-size: 0.8rem;
}
.col-progress {
  width: 60px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.text-ellipsis {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.progress-bar-bg {
  width: 100%;
  height: 6px;
  background: var(--border-default);
  border-radius: 4px;
  overflow: hidden;
}
.progress-bar-fill {
  height: 100%;
  background: var(--primary-500);
}

.timeline-panel {
  flex: 1;
  overflow-x: auto;
  overflow-y: auto;
  background: white; /* Frappe gantt uses a white theme by default, we can style it later for dark mode */
}

/* Frappe Gantt overrides for dark mode if body is dark */
:global(.gantt .grid-header) {
  fill: var(--bg-surface) !important;
}
:global(.gantt .grid-row) {
  fill: var(--bg-base) !important;
}
:global(.gantt .grid-row:nth-child(even)) {
  fill: var(--bg-surface) !important;
}
:global(.gantt .tick line) {
  stroke: var(--border-default) !important;
}
:global(.gantt .tick text) {
  fill: var(--text-secondary) !important;
}
:global(.gantt .today-highlight) {
  fill: rgba(56, 189, 248, 0.15) !important; /* --info */
}
:global(.gantt .bar-progress) {
  fill: var(--primary-500) !important;
}
:global(.gantt .bar-wrapper:hover .bar-progress) {
  fill: var(--primary-400) !important;
}
:global(.gantt .bar-label) {
  fill: var(--text-primary) !important;
  font-weight: 500;
}

/* Milestone Styling */
:global(.gantt .bar-wrapper.milestone .bar) {
  transform: rotate(45deg);
  transform-origin: center;
  rx: 4;
  ry: 4;
  width: 20px !important;
  height: 20px !important;
  y: 7px !important;
}
:global(.gantt .bar-wrapper.milestone .bar-progress) {
  visibility: hidden;
}

/* Assignee Color Coding */
:global(.gantt .bar-wrapper.assignee-color-0 .bar-progress) { fill: #6366F1 !important; }
:global(.gantt .bar-wrapper.assignee-color-1 .bar-progress) { fill: #3B82F6 !important; }
:global(.gantt .bar-wrapper.assignee-color-2 .bar-progress) { fill: #10B981 !important; }
:global(.gantt .bar-wrapper.assignee-color-3 .bar-progress) { fill: #F59E0B !important; }
:global(.gantt .bar-wrapper.assignee-color-4 .bar-progress) { fill: #EF4444 !important; }
:global(.gantt .bar-wrapper.assignee-color-5 .bar-progress) { fill: #8B5CF6 !important; }

.empty-state-list {
  padding: 32px 16px;
  text-align: center;
  color: var(--text-tertiary);
  font-size: 0.85rem;
}
</style>
