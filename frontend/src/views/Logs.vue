<template>
  <div class="logs-page">
    <div class="page-header">
      <div>
        <h2>系统日志与连接诊断</h2>
        <p class="subtitle">实时捕获 Sing-box 内核运行日志与 Clash API 活动连接</p>
      </div>
      <div class="header-actions">
        <a-button @click="clearLogs">清屏</a-button>
        <a-button type="primary" @click="fetchConnections">刷新活动连接</a-button>
      </div>
    </div>

    <!-- Terminal Log View -->
    <a-card class="log-card mt-4" :bordered="false">
      <div class="log-toolbar">
        <div class="log-status">
          <span class="status-dot" :class="{ connected: wsConnected }"></span>
          {{ wsConnected ? 'WebSocket 日志流已连接' : 'WebSocket 断开重连中...' }}
        </div>
        <div class="log-controls">
          <a-switch v-model:checked="autoScroll" checked-children="自动滚屏" un-checked-children="停止滚屏" />
        </div>
      </div>
      <div ref="terminalRef" class="terminal-body">
        <div v-for="(line, idx) in logs" :key="idx" class="log-line">
          {{ line }}
        </div>
      </div>
    </a-card>

    <!-- Active Connections Table -->
    <a-card class="connections-card mt-4" title="活动连接列表 (Clash API)" :bordered="false">
      <a-table
        :columns="connColumns"
        :data-source="connections"
        :loading="connLoading"
        row-key="id"
        size="small"
        :pagination="{ pageSize: 10 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'source'">
            <code>{{ record.metadata?.sourceIP }}:{{ record.metadata?.sourcePort }}</code>
          </template>
          <template v-if="column.key === 'destination'">
            <b>{{ record.metadata?.host || record.metadata?.destinationIP }}</b>:{{ record.metadata?.destinationPort }}
          </template>
          <template v-if="column.key === 'network'">
            <a-tag>{{ record.metadata?.network?.toUpperCase() }}</a-tag>
          </template>
          <template v-if="column.key === 'traffic'">
            <span>⬇️ {{ formatBytes(record.download) }} | ⬆️ {{ formatBytes(record.upload) }}</span>
          </template>
          <template v-if="column.key === 'rule'">
            <a-tag color="blue">{{ record.rule }}</a-tag>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { getActiveConnections } from '@/api/server'
import { formatBytes } from '@/utils/format'

const logs = ref<string[]>([])
const autoScroll = ref(true)
const wsConnected = ref(false)
const terminalRef = ref<HTMLDivElement>()
let ws: WebSocket | null = null

// Connections
const connections = ref<any[]>([])
const connLoading = ref(false)

const connColumns = [
  { title: '源地址', key: 'source' },
  { title: '目标域名/IP', key: 'destination' },
  { title: '网络', key: 'network', width: 90 },
  { title: '流量', key: 'traffic' },
  { title: '分流规则', key: 'rule' },
]

function connectWS() {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const url = `${protocol}//${window.location.host}/api/server/logs/ws`
  ws = new WebSocket(url)

  ws.onopen = () => {
    wsConnected.value = true
  }

  ws.onmessage = (event) => {
    logs.value.push(event.data)
    if (logs.value.length > 500) {
      logs.value.shift()
    }
    if (autoScroll.value) {
      nextTick(() => {
        if (terminalRef.value) {
          terminalRef.value.scrollTop = terminalRef.value.scrollHeight
        }
      })
    }
  }

  ws.onclose = () => {
    wsConnected.value = false
    setTimeout(connectWS, 3000)
  }
}

function clearLogs() {
  logs.value = []
}

async function fetchConnections() {
  connLoading.value = true
  try {
    const res = await getActiveConnections()
    connections.value = res.connections || []
  } finally {
    connLoading.value = false
  }
}

onMounted(() => {
  connectWS()
  fetchConnections()
})

onUnmounted(() => {
  if (ws) ws.close()
})
</script>

<style scoped>
.logs-page {
  padding: 4px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.page-header h2 {
  font-size: 20px;
  font-weight: 700;
  margin: 0;
  color: #0f172a;
}

.subtitle {
  color: #64748b;
  font-size: 13px;
  margin: 4px 0 0 0;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.log-card {
  border-radius: 12px;
  background: #0f172a;
  color: #f8fafc;
  overflow: hidden;
}

.log-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  margin-bottom: 12px;
}

.log-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #94a3b8;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #ef4444;
}

.status-dot.connected {
  background: #10b981;
}

.terminal-body {
  height: 380px;
  overflow-y: auto;
  font-family: 'Cascadia Code', 'Fira Code', Consolas, Monaco, monospace;
  font-size: 12px;
  line-height: 1.6;
  background: #020617;
  padding: 12px;
  border-radius: 8px;
  color: #38bdf8;
}

.log-line {
  word-break: break-all;
  white-space: pre-wrap;
}

.connections-card {
  border-radius: 12px;
  background: #ffffff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.mt-4 {
  margin-top: 16px;
}
</style>
