<script setup lang="ts">
import { ref, computed } from 'vue'
import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'
import SelectButton from 'primevue/selectbutton'

const props = defineProps<{
  tasks: any[]
  members?: any[]
}>()

const emit = defineEmits<{
  (e: 'task-click', task: any): void
  (e: 'date-click', date: Date): void
  (e: 'task-date-change', task: any, start: Date, end: Date): void
}>()

const currentDate = ref(new Date())
const selectedAssignee = ref<string | null>(null)

const currentMonth = computed(() => currentDate.value.getMonth())
const currentYear = computed(() => currentDate.value.getFullYear())

const viewMode = ref('month')
const viewOptions = [
  { label: 'Bulanan', value: 'month' },
  { label: 'Mingguan', value: 'week' }
]

const monthNames = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni', 
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
]

const dayNames = ['Min', 'Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab']

const assigneeOptions = computed(() => {
  const opts = [{ label: 'Semua Assignee', value: null }]
  if (props.members) {
    props.members.forEach(m => {
      opts.push({ label: m.user_full_name, value: m.user_id })
    })
  }
  return opts
})

const prevPeriod = () => {
  if (viewMode.value === 'month') {
    currentDate.value = new Date(currentYear.value, currentMonth.value - 1, 1)
  } else {
    const d = new Date(currentDate.value)
    d.setDate(d.getDate() - 7)
    currentDate.value = d
  }
}

const nextPeriod = () => {
  if (viewMode.value === 'month') {
    currentDate.value = new Date(currentYear.value, currentMonth.value + 1, 1)
  } else {
    const d = new Date(currentDate.value)
    d.setDate(d.getDate() + 7)
    currentDate.value = d
  }
}

const goToToday = () => {
  currentDate.value = new Date()
}

// Generate calendar grid
const monthlyDays = computed(() => {
  const days = []
  const firstDay = new Date(currentYear.value, currentMonth.value, 1)
  const lastDay = new Date(currentYear.value, currentMonth.value + 1, 0)
  
  const startingDayOfWeek = firstDay.getDay()
  const daysInMonth = lastDay.getDate()
  
  // Previous month padding
  const prevMonthLastDay = new Date(currentYear.value, currentMonth.value, 0).getDate()
  for (let i = 0; i < startingDayOfWeek; i++) {
    const day = prevMonthLastDay - startingDayOfWeek + i + 1
    const date = new Date(currentYear.value, currentMonth.value - 1, day)
    days.push({
      date,
      dayNumber: day,
      isCurrentMonth: false,
      isToday: isSameDate(date, new Date()),
      tasks: getTasksForDate(date)
    })
  }
  
  // Current month days
  for (let i = 1; i <= daysInMonth; i++) {
    const date = new Date(currentYear.value, currentMonth.value, i)
    days.push({
      date,
      dayNumber: i,
      isCurrentMonth: true,
      isToday: isSameDate(date, new Date()),
      tasks: getTasksForDate(date)
    })
  }
  
  // Next month padding to fill exactly 6 rows (42 cells)
  const remainingCells = 42 - days.length
  for (let i = 1; i <= remainingCells; i++) {
    const date = new Date(currentYear.value, currentMonth.value + 1, i)
    days.push({
      date,
      dayNumber: i,
      isCurrentMonth: false,
      isToday: isSameDate(date, new Date()),
      tasks: getTasksForDate(date)
    })
  }
  
  return days
})

const weeklyDays = computed(() => {
  const days = []
  const current = new Date(currentDate.value)
  const first = current.getDate() - current.getDay() // Sunday as first day
  
  for (let i = 0; i < 7; i++) {
    const d = new Date(current.getFullYear(), current.getMonth(), first + i)
    days.push({
      date: d,
      dayNumber: d.getDate(),
      dayName: dayNames[d.getDay()],
      isToday: isSameDate(d, new Date()),
      tasks: getTasksForDate(d)
    })
  }
  return days
})

const currentDisplayPeriodText = computed(() => {
  if (viewMode.value === 'month') {
    return `${monthNames[currentMonth.value]} ${currentYear.value}`
  } else {
    if (weeklyDays.value.length > 0) {
      const first = weeklyDays.value[0].date
      const last = weeklyDays.value[6].date
      if (first.getMonth() === last.getMonth()) {
        return `${first.getDate()} - ${last.getDate()} ${monthNames[first.getMonth()]} ${first.getFullYear()}`
      } else {
        return `${first.getDate()} ${monthNames[first.getMonth()].substr(0,3)} - ${last.getDate()} ${monthNames[last.getMonth()].substr(0,3)} ${last.getFullYear()}`
      }
    }
    return ''
  }
})

function isSameDate(d1: Date, d2: Date) {
  return d1.getFullYear() === d2.getFullYear() &&
         d1.getMonth() === d2.getMonth() &&
         d1.getDate() === d2.getDate()
}

function getTasksForDate(date: Date) {
  return props.tasks.filter(t => {
    if (!t.start_date) return false
    
    // Filter by assignee if set
    if (selectedAssignee.value && t.assignee_id !== selectedAssignee.value) return false
    
    const taskStart = new Date(t.start_date)
    
    // Simplest logic: show task on its start date
    return isSameDate(taskStart, date)
  })
}

const getPriorityColor = (priority: string) => {
  switch (priority?.toLowerCase()) {
    case 'urgent':
    case 'critical': return 'var(--danger)'
    case 'high': return 'var(--warning)'
    case 'medium':
    case 'normal': return 'var(--primary-500)'
    case 'low': return 'var(--text-muted)'
    default: return 'var(--primary-500)'
  }
}

const onDateClick = (date: Date) => {
  emit('date-click', date)
}

const onTaskClick = (e: Event, task: any) => {
  e.stopPropagation()
  emit('task-click', task)
}

const onDragStart = (e: DragEvent, task: any) => {
  if (e.dataTransfer) {
    e.dataTransfer.setData('taskId', task.id)
    e.dataTransfer.effectAllowed = 'move'
  }
}

const onDrop = (e: DragEvent, date: Date) => {
  const taskId = e.dataTransfer?.getData('taskId')
  if (taskId) {
    const task = props.tasks.find(t => t.id === taskId)
    if (task && task.start_date) {
      const oldStart = new Date(task.start_date)
      const oldEnd = task.end_date ? new Date(task.end_date) : new Date(task.start_date)
      
      const durationMs = oldEnd.getTime() - oldStart.getTime()
      
      const newStart = new Date(date)
      const newEnd = new Date(newStart.getTime() + durationMs)
      
      if (!isSameDate(oldStart, newStart)) {
        emit('task-date-change', task, newStart, newEnd)
      }
    }
  }
}
</script>

<template>
  <div class="calendar-board">
    <div class="calendar-header">
      <div class="header-left" style="display: flex; align-items: center; gap: 1rem;">
        <h2 class="month-title">{{ currentDisplayPeriodText }}</h2>
        <Dropdown 
          v-model="selectedAssignee" 
          :options="assigneeOptions" 
          optionLabel="label" 
          optionValue="value" 
          placeholder="Filter Assignee" 
          class="p-dropdown-sm ml-4"
          style="width: 200px"
        />
      </div>
      <div class="header-right flex align-items-center gap-3">
        <SelectButton v-model="viewMode" :options="viewOptions" optionLabel="label" optionValue="value" aria-labelledby="basic" class="p-button-sm mr-2" />
        <Button label="Hari Ini" class="p-button-outlined p-button-sm mr-2" @click="goToToday" />
        <Button icon="pi pi-chevron-left" class="p-button-rounded p-button-text p-button-sm mr-1" @click="prevPeriod" />
        <Button icon="pi pi-chevron-right" class="p-button-rounded p-button-text p-button-sm" @click="nextPeriod" />
      </div>
    </div>
    
    <div v-if="viewMode === 'month'" class="calendar-grid">
      <!-- Weekday Headers -->
      <div class="weekday-header" v-for="day in dayNames" :key="day">
        {{ day }}
      </div>
      
      <!-- Calendar Days -->
      <div 
        v-for="(day, index) in monthlyDays" 
        :key="index"
        class="calendar-day"
        :class="{
          'not-current-month': !day.isCurrentMonth,
          'today': day.isToday
        }"
        @click="onDateClick(day.date)"
        @dragover.prevent
        @drop="(e) => onDrop(e, day.date)"
      >
        <div class="day-number">{{ day.dayNumber }}</div>
        
        <div class="day-tasks">
          <div 
            v-for="task in day.tasks" 
            :key="task.id"
            class="task-item"
            :style="{ borderLeftColor: getPriorityColor(task.priority) }"
            draggable="true"
            @dragstart="(e) => onDragStart(e, task)"
            @click="(e) => onTaskClick(e, task)"
          >
            <span class="task-name">{{ task.name }}</span>
            <span class="task-status" :class="task.status"></span>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="weekly-grid">
      <div 
        v-for="(day, index) in weeklyDays" 
        :key="index" 
        class="weekly-column"
        :class="{ 'today-column': day.isToday }"
        @click="onDateClick(day.date)"
        @dragover.prevent
        @drop="(e) => onDrop(e, day.date)"
      >
        <div class="weekly-header">
          <div class="day-name">{{ day.dayName }}</div>
          <div class="day-number">{{ day.dayNumber }}</div>
        </div>
        <div class="weekly-tasks">
          <div 
            v-for="task in day.tasks" 
            :key="task.id"
            class="task-item-weekly"
            :style="{ borderLeftColor: getPriorityColor(task.priority) }"
            draggable="true"
            @dragstart="(e) => onDragStart(e, task)"
            @click="(e) => onTaskClick(e, task)"
          >
            <span class="task-name">{{ task.name }}</span>
            <span class="task-status" :class="task.status"></span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.calendar-board {
  display: flex;
  flex-direction: column;
  background-color: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  overflow: hidden;
  height: 100%;
}

.calendar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-4) var(--space-6);
  border-bottom: 1px solid var(--border-subtle);
  background-color: var(--bg-elevated);
}

.month-title {
  margin: 0;
  font-size: var(--text-2xl);
  font-weight: 700;
  color: var(--text-primary);
}

.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  grid-auto-rows: minmax(120px, 1fr);
  flex: 1;
  background-color: var(--border-subtle);
  gap: 1px; /* Creates borders between cells */
}

.weekday-header {
  background-color: var(--bg-elevated);
  text-align: center;
  padding: var(--space-2);
  font-weight: 600;
  color: var(--text-secondary);
  font-size: var(--text-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  height: 40px;
}

.calendar-day {
  background-color: var(--bg-surface);
  padding: var(--space-2);
  display: flex;
  flex-direction: column;
  transition: background-color var(--transition-fast);
  cursor: pointer;
  overflow: hidden;
}

.calendar-day:hover {
  background-color: var(--bg-hover);
}

.calendar-day.not-current-month {
  background-color: var(--bg-base);
  opacity: 0.7;
}

.calendar-day.not-current-month:hover {
  background-color: var(--bg-surface);
}

.day-number {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: var(--space-2);
  align-self: flex-end;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-full);
}

.calendar-day.today .day-number {
  background-color: var(--primary-500);
  color: #fff;
  font-weight: 700;
}

.day-tasks {
  display: flex;
  flex-direction: column;
  gap: 4px;
  overflow-y: auto;
  flex: 1;
}

/* Hide scrollbar for clean look */
.day-tasks::-webkit-scrollbar {
  width: 4px;
}
.day-tasks::-webkit-scrollbar-thumb {
  background-color: var(--border-default);
  border-radius: 4px;
}

.task-item {
  background-color: var(--bg-elevated);
  border: 1px solid var(--border-subtle);
  border-left-width: 3px;
  border-left-style: solid;
  border-radius: var(--radius-sm);
  padding: 4px 8px;
  font-size: 11px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: all var(--transition-fast);
  white-space: nowrap;
}

.task-item:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
  background-color: var(--bg-hover);
}

.task-name {
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--text-primary);
  flex: 1;
}

/* Status Indicators */
.task-status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.task-status.todo { background-color: var(--text-muted); }
.task-status.in-progress { background-color: var(--primary-500); }
.task-status.done { background-color: var(--success); }
.task-status.on-hold { background-color: var(--warning); }
.task-status.cancelled { background-color: var(--danger); }

/* Weekly View Styles */
.weekly-grid {
  display: flex;
  flex: 1;
  overflow-x: auto;
  border-top: 1px solid var(--border-subtle);
}

.weekly-column {
  flex: 1;
  min-width: 150px;
  border-right: 1px solid var(--border-subtle);
  display: flex;
  flex-direction: column;
  background-color: var(--bg-surface);
}

.weekly-column:last-child {
  border-right: none;
}

.weekly-column.today-column {
  background-color: rgba(99, 102, 241, 0.05); /* Highlight today's column slightly */
}

.weekly-header {
  padding: var(--space-3);
  text-align: center;
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.weekly-header .day-name {
  font-size: 0.8rem;
  color: var(--text-tertiary);
  text-transform: uppercase;
  font-weight: 600;
  margin-bottom: 4px;
}

.weekly-header .day-number {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-primary);
}

.today-column .day-number {
  color: var(--primary-500);
}

.weekly-tasks {
  flex: 1;
  padding: var(--space-2);
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.task-item-weekly {
  background-color: var(--bg-elevated);
  border: 1px solid var(--border-subtle);
  border-left: 3px solid var(--primary-500);
  border-radius: var(--radius-sm);
  padding: var(--space-2) var(--space-3);
  font-size: 0.8rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: var(--shadow-sm);
}

.task-item-weekly:hover {
  border-color: var(--primary-400);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.mr-1 { margin-right: 0.25rem; }
.mr-2 { margin-right: 0.5rem; }
</style>
