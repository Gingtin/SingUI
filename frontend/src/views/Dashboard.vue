<template>
  <div class="dashboard-page">
    <!-- Top Core Status Card -->
    <a-card class="status-card" :bordered="false">
      <div class="status-wrapper">
        <div class="status-info">
          <div class="status-badge" :class="{ active: coreStatus.is_running }">
            <span class="dot"></span>
            {{ coreStatus.is_running ? 'Sing-box 运行中' : 'Sing-box 已停止' }}
          </div>
          <div class="status-details">
            <span v-if="coreStatus.pid"><b>PID:</b> {{ coreStatus.pid }}</span>
            <span><b>运行时间:</b> {{ formatDuration(sysStatus.uptime) }}</span>
            <span><b>系统平台:</b> {{ sysStatus.platform || sysStatus.os }}</span>
            <span v-if="coreStatus.version"><b>版本:</b> {{ coreStatus.version.split('\n')[0] }}</span>
          </div>
        </div>
        <div class="status-actions">
          <a-button v-if="!coreStatus.is_running" type="primary" :loading="coreLoading" @click="handleStartCore">
            启动内核
          </a-button>
          <a-button v-else danger :loading="coreLoading" @click="handleStopCore">
            停止内核
          </a-button>
          <a-button :loading="coreLoading" @click="handleRestartCore">
            重启内核
          </a-button>
        </div>
      </div>
    </a-card>

    <!-- Metrics Row -->
    <a-row :gutter="[16, 16]" class="mt-4">
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="metric-card" :bordered="false">
          <div class="metric-header">
            <span class="title">CPU 使用率</span>
            <span class="value">{{ sysStatus.cpu_percent?.toFixed(1) || 0 }}%</span>
          </div>
          <a-progress :percent="Number(sysStatus.cpu_percent?.toFixed(1)) || 0" :stroke-color="getProgressColor(sysStatus.cpu_percent)" :show-info="false" />
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="metric-card" :bordered="false">
          <div class="metric-header">
            <span class="title">内存占用</span>
            <span class="value">{{ sysStatus.mem_percent?.toFixed(1) || 0 }}%</span>
          </div>
          <a-progress :percent="Number(sysStatus.mem_percent?.toFixed(1)) || 0" :stroke-color="getProgressColor(sysStatus.mem_percent)" :show-info="false" />
          <div class="metric-sub">{{ formatBytes(sysStatus.mem_used) }} / {{ formatBytes(sysStatus.mem_total) }}</div>
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="metric-card" :bordered="false">
          <div class="metric-header">
            <span class="title">硬盘空间</span>
            <span class="value">{{ sysStatus.disk_percent?.toFixed(1) || 0 }}%</span>
          </div>
          <a-progress :percent="Number(sysStatus.disk_percent?.toFixed(1)) || 0" :stroke-color="getProgressColor(sysStatus.disk_percent)" :show-info="false" />
          <div class="metric-sub">{{ formatBytes(sysStatus.disk_used) }} / {{ formatBytes(sysStatus.disk_total) }}</div>
        </a-card>
      </a-col>

      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="metric-card" :bordered="false">
          <div class="metric-header">
            <span class="title">实时网速</span>
            <span class="speed-indicator">
              <span>⬇️ {{ formatBytes(sysStatus.net_download_rate) }}/s</span>
              <span>⬆️ {{ formatBytes(sysStatus.net_upload_rate) }}/s</span>
            </span>
          </div>
          <div class="metric-sub mt-2">总流量: ⬇️ {{ formatBytes(sysStatus.net_total_recv) }} | ⬆️ {{ formatBytes(sysStatus.net_total_sent) }}</div>
        </a-card>
      </a-col>
    </a-row>

    <!-- Charts Row -->
    <a-row :gutter="[16, 16]" class="mt-4">
      <!-- Speed Chart (16 cols) -->
      <a-col :xs="24" :lg="16">
        <a-card class="chart-card" title="实时网络速率动态监控 (KB/s)" :bordered="false">
          <div ref="chartRef" style="height: 320px; width: 100%;"></div>
        </a-card>
      </a-col>

      <!-- Protocol Distribution (8 cols) -->
      <a-col :xs="24" :lg="8">
        <a-card class="chart-card" title="入站协议分布" :bordered="false">
          <div ref="pieChartRef" style="height: 320px; width: 100%;"></div>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { message } from 'ant-design-vue'
import * as echarts from 'echarts'
import { getSystemStatus, getCoreStatus, startCore, stopCore, restartCore, SystemStatus, CoreStatus } from '@/api/server'
import { getInbounds } from '@/api/inbound'
import { formatBytes, formatDuration } from '@/utils/format'

const sysStatus = ref<Partial<SystemStatus>>({})
const coreStatus = ref<Partial<CoreStatus>>({})
const coreLoading = ref(false)

const chartRef = ref<HTMLDivElement>()
const pieChartRef = ref<HTMLDivElement>()
let chartInstance: echarts.ECharts | null = null
let pieInstance: echarts.ECharts | null = null

const downloadHistory: number[] = []
const uploadHistory: number[] = []
const timeLabels: string[] = []
const maxPoints = 20

let pollTimer: any = null

function getProgressColor(percent: number | undefined) {
  if (!percent) return '#10b981'
  if (percent > 85) return '#ef4444'
  if (percent > 65) return '#f59e0b'
  return '#3b82f6'
}

async function fetchStatus() {
  try {
    const [sys, core] = await Promise.all([getSystemStatus(), getCoreStatus()])
    sysStatus.value = sys
    coreStatus.value = core

    // Update Chart
    const nowStr = new Date().toLocaleTimeString()
    timeLabels.push(nowStr)
    downloadHistory.push(sys.net_download_rate / 1024) // KB/s
    uploadHistory.push(sys.net_upload_rate / 1024)     // KB/s

    if (timeLabels.length > maxPoints) {
      timeLabels.shift()
      downloadHistory.shift()
      uploadHistory.shift()
    }

    if (chartInstance) {
      chartInstance.setOption({
        xAxis: { data: timeLabels },
        series: [
          { name: '下载速率 (KB/s)', data: downloadHistory },
          { name: '上传速率 (KB/s)', data: uploadHistory },
        ],
      })
    }
  } catch (err) {
    console.error('Failed to fetch status', err)
  }
}

async function fetchProtocolDistribution() {
  try {
    const inbounds = await getInbounds()
    const counts: Record<string, number> = {}
    for (const inb of inbounds) {
      counts[inb.protocol] = (counts[inb.protocol] || 0) + 1
    }

    const pieData = Object.keys(counts).map((k) => ({
      name: k.toUpperCase(),
      value: counts[k],
    }))

    if (pieData.length === 0) {
      pieData.push({ name: '暂无节点', value: 1 })
    }

    if (pieInstance) {
      pieInstance.setOption({
        series: [{ data: pieData }],
      })
    }
  } catch (err) {
    console.error(err)
  }
}

async function handleStartCore() {
  coreLoading.value = true
  try {
    await startCore()
    message.success('Sing-box 内核已启动')
    fetchStatus()
  } finally {
    coreLoading.value = false
  }
}

async function handleStopCore() {
  coreLoading.value = true
  try {
    await stopCore()
    message.success('Sing-box 内核已停止')
    fetchStatus()
  } finally {
    coreLoading.value = false
  }
}

async function handleRestartCore() {
  coreLoading.value = true
  try {
    await restartCore()
    message.success('Sing-box 内核已重启')
    fetchStatus()
  } finally {
    coreLoading.value = false
  }
}

function initCharts() {
  if (chartRef.value) {
    chartInstance = echarts.init(chartRef.value)
    chartInstance.setOption({
      tooltip: { trigger: 'axis' },
      legend: { data: ['下载速率 (KB/s)', '上传速率 (KB/s)'], textStyle: { color: '#64748b' } },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: timeLabels,
        axisLine: { lineStyle: { color: '#cbd5e1' } },
      },
      yAxis: {
        type: 'value',
        axisLine: { lineStyle: { color: '#cbd5e1' } },
        splitLine: { lineStyle: { color: '#f1f5f9' } },
      },
      series: [
        {
          name: '下载速率 (KB/s)',
          type: 'line',
          smooth: true,
          showSymbol: false,
          areaStyle: { opacity: 0.15, color: '#10b981' },
          lineStyle: { color: '#10b981', width: 2 },
          itemStyle: { color: '#10b981' },
          data: downloadHistory,
        },
        {
          name: '上传速率 (KB/s)',
          type: 'line',
          smooth: true,
          showSymbol: false,
          areaStyle: { opacity: 0.15, color: '#3b82f6' },
          lineStyle: { color: '#3b82f6', width: 2 },
          itemStyle: { color: '#3b82f6' },
          data: uploadHistory,
        },
      ],
    })
  }

  if (pieChartRef.value) {
    pieInstance = echarts.init(pieChartRef.value)
    pieInstance.setOption({
      tooltip: { trigger: 'item' },
      legend: { bottom: '5%', left: 'center' },
      series: [
        {
          name: '节点协议',
          type: 'pie',
          radius: ['45%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: { borderRadius: 8, borderColor: '#fff', borderWidth: 2 },
          label: { show: false, position: 'center' },
          emphasis: {
            label: { show: true, fontSize: 16, fontWeight: 'bold' },
          },
          data: [],
        },
      ],
    })
  }
}

onMounted(() => {
  initCharts()
  fetchStatus()
  fetchProtocolDistribution()
  pollTimer = setInterval(fetchStatus, 2000)
  window.addEventListener('resize', () => {
    chartInstance?.resize()
    pieInstance?.resize()
  })
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
  if (chartInstance) chartInstance.dispose()
  if (pieInstance) pieInstance.dispose()
})
</script>

<style scoped>
.dashboard-page {
  padding: 4px;
}

.status-card {
  border-radius: 12px;
  background: #ffffff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.status-wrapper {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 14px;
  border-radius: 20px;
  background: #fee2e2;
  color: #ef4444;
  font-weight: 600;
  font-size: 14px;
}

.status-badge.active {
  background: #dcfce7;
  color: #10b981;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
}

.status-details {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  color: #64748b;
  font-size: 13px;
  margin-top: 8px;
}

.status-actions {
  display: flex;
  gap: 10px;
}

.metric-card {
  border-radius: 12px;
  background: #ffffff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.metric-header .title {
  color: #64748b;
  font-size: 13px;
}

.metric-header .value {
  font-size: 18px;
  font-weight: 700;
  color: #0f172a;
}

.speed-indicator {
  display: flex;
  gap: 12px;
  font-size: 13px;
  font-weight: 600;
  color: #0f172a;
}

.metric-sub {
  font-size: 12px;
  color: #94a3b8;
  margin-top: 6px;
}

.chart-card {
  border-radius: 12px;
  background: #ffffff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.mt-4 {
  margin-top: 16px;
}
</style>
