<script setup lang="ts">
import { computed } from 'vue'
import { models } from '../../wailsjs/go/models'
import { Pie, Bar } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale, BarElement } from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale, BarElement)

const props = defineProps<{
  project: models.Project
  tasks: models.Task[]
}>()

const totalTasks = computed(() => props.tasks.length)
const completedTasks = computed(() => props.tasks.filter(t => t.status === 'done').length)
const inProgressTasks = computed(() => props.tasks.filter(t => t.status === 'in-progress').length)
const todoTasks = computed(() => props.tasks.filter(t => t.status === 'todo').length)

const progressPercent = computed(() => {
  if (totalTasks.value === 0) return 0
  return Math.round((completedTasks.value / totalTasks.value) * 100)
})

const totalEstimatedCost = computed(() => {
  return props.tasks.reduce((sum, t) => sum + (t.estimated_cost || 0), 0)
})

const totalActualCost = computed(() => {
  return props.tasks.reduce((sum, t) => sum + (t.actual_cost || 0), 0)
})

const formatCurrency = (val: number) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: props.project.currency || 'IDR' }).format(val)
}

const statusChartData = computed(() => ({
  labels: ['To Do', 'In Progress', 'Done'],
  datasets: [
    {
      backgroundColor: ['#94a3b8', '#3b82f6', '#10b981'],
      data: [todoTasks.value, inProgressTasks.value, completedTasks.value]
    }
  ]
}))

const costChartData = computed(() => ({
  labels: ['Cost Comparison'],
  datasets: [
    {
      label: 'Budget',
      backgroundColor: '#f59e0b',
      data: [props.project.total_budget || 0]
    },
    {
      label: 'Estimated Cost',
      backgroundColor: '#6366f1',
      data: [totalEstimatedCost.value]
    },
    {
      label: 'Actual Cost',
      backgroundColor: '#ef4444',
      data: [totalActualCost.value]
    }
  ]
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false
}
</script>

<template>
  <div class="report-container">
    <div class="summary-cards">
      <div class="card">
        <h4>Progress Proyek</h4>
        <div class="value">{{ progressPercent }}%</div>
        <div class="subtitle">{{ completedTasks }} dari {{ totalTasks }} tugas selesai</div>
      </div>
      <div class="card">
        <h4>Total Budget</h4>
        <div class="value">{{ formatCurrency(props.project.total_budget || 0) }}</div>
      </div>
      <div class="card">
        <h4>Estimated Cost</h4>
        <div class="value">{{ formatCurrency(totalEstimatedCost) }}</div>
      </div>
      <div class="card" :class="{ 'over-budget': totalActualCost > (props.project.total_budget || 0) }">
        <h4>Actual Cost</h4>
        <div class="value">{{ formatCurrency(totalActualCost) }}</div>
      </div>
    </div>

    <div class="charts-grid mt-6">
      <div class="chart-box">
        <h3>Status Tugas</h3>
        <div class="chart-wrapper">
          <Pie :data="statusChartData" :options="chartOptions" />
        </div>
      </div>
      <div class="chart-box">
        <h3>Anggaran vs Pengeluaran</h3>
        <div class="chart-wrapper">
          <Bar :data="costChartData" :options="chartOptions" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.report-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}
.summary-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
}
.card {
  background: var(--bg-base);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.card h4 {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.9rem;
  font-weight: 500;
}
.card .value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-primary);
}
.card .subtitle {
  font-size: 0.8rem;
  color: var(--text-tertiary);
}
.card.over-budget .value {
  color: #ef4444;
}
.charts-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}
.chart-box {
  background: var(--bg-base);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  padding: 1.5rem;
}
.chart-box h3 {
  margin: 0 0 1rem 0;
  font-size: 1.1rem;
}
.chart-wrapper {
  height: 300px;
  position: relative;
}
</style>
