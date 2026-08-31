<template>
  <div class="logs-page">
    <div class="page-header">
      <div>
        <h2>日志与系统诊断 (Diagnostics & Config)</h2>
        <p class="subtitle">实时 WebSocket 日志滚屏、Clash API 活动连接与 Sing-box 原生配置查看</p>
      </div>
      <div class="header-actions">
        <a-button v-if="activeTab === 'logs'" @click="clearLogs">清空屏幕</a-button>
        <a-button v-if="activeTab === 'connections'" @click="fetchConnections">刷新连接</a-button>
        <a-button v-if="activeTab === 'config'" @click="fetchRawConfig">刷新配置</a-button>
      </div>
    </div>

    <a-card class="mt-4 main-card" :bordered="false">
      <a-tabs v-model:activeKey="activeTab">
        <!-- Tab 1: Terminal Logs -->
        <a-tab-pane key="logs" tab="实时终端日志 (Live Logs)">
          <div class="terminal-toolbar">
            <div class="terminal-status">
              <span class="status-dot" :class="{ connected: isConnected }"></span>
              <span>{{ isConnected ? 'WebSocket 已连接' : '正在重新连接...' }}</span>
            </div>
            <a-checkbox v-model:checked="autoScroll">自动滚屏</a-checkbox>
          </div>
          <div ref="terminalRef" class="terminal-window">
            <div v-for="(log, idx) in logs" :key="idx" class="log-line">
              <span class="log-index">#{{ idx + 1 }}</span>
              <span class="log-text" :class="getLogLevelClass(log)">{{ log }}</span>
            </div>
            <div v-if="logs.length === 0" class="log-empty">暂无运行日志...</div>
          </div>
        </a-tab-pane>

        <!-- Tab 2: Active Connections -->
        <a-tab-pane key="connections" tab="活跃网络连接 (Active Connections)">
          <div class="connections-summary mb-3">
            <span>当前活跃连接: <b>{{ connections.length }}</b> 个</span>
            <span>总上传: <b>{{ formatBytes(uploadTotal) }}</b></span>
            <span>总下载: <b>{{ formatBytes(downloadTotal) }}</b></span>
          </div>

          <a-table
            :columns="connColumns"
            :data-source="connections"
            :loading="connLoading"
            row-key="id"
            :pagination="{ pageSize: 15 }"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'host'">
                <span>{{ record.metadata?.host || record.metadata?.destinationIP }}:{{ record.metadata?.destinationPort }}</span>
              </template>
              <template v-if="column.key === 'sourceIP'">
                <span>{{ record.metadata?.sourceIP }}:{{ record.metadata?.sourcePort }}</span>
              </template>
              <template v-if="column.key === 'traffic'">
                <span>⬆️ {{ formatBytes(record.upload) }} | ⬇️ {{ formatBytes(record.download) }}</span>
              </template>
              <template v-if="column.key === 'chains'">
                <a-tag color="blue">{{ record.chains?.[0] || 'DIRECT' }}</a-tag>
              </template>
            </template>
          </a-table>
        </a-tab-pane>

        <!-- Tab 3: Raw config.json Viewer -->
        <a-tab-pane key="config" tab="Sing-box 核心配置 (config.json)">
          <div class="config-viewer-wrapper">
            <div class="config-toolbar mb-2">
              <a-button size="small" type="primary" @click="copyConfig">复制完整配置</a-button>
            </div>
            <pre class="config-code"><code>{{ rawConfig }}</code></pre>
          </div>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { message } from 'ant-design-vue'
import { getActiveConnections } from '@/api/server'
import { getRawConfig } from '@/api/outbound'
import { formatBytes } from '@/utils/format'

const activeTab = ref('logs')

// Logs
const logs = ref<string[]>([])
const isConnected = ref(false)
const autoScroll = ref(true)
const terminalRef = ref<HTMLDivElement>()
let ws: WebSocket | null = null

// Connections
const connections = ref<any[]>([])
const uploadTotal = ref(0)
const downloadTotal = ref(0)
const connLoading = ref(false)

const connColumns = [
  { title: '目标地址 (Host / IP)', key: 'host' },
  { title: '源地址 (Source IP)', key: 'sourceIP', width: 180 },
  { title: '传输量 (Up / Down)', key: 'traffic', width: 220 },
  { title: '出站链路 (Detour)', key: 'chains', width: 140 },
]

// Raw Config
const rawConfig = ref('')

function connectWS() {
  const token = localStorage.getItem('token') || ''
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = `${protocol}//${window.location.host}/api/server/logs/ws?token=${token}`

  ws = new WebSocket(wsUrl)

  ws.onopen = () => {
    isConnected.value = true
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
    isConnected.value = false
    setTimeout(connectWS, 3000)
  }

  ws.onerror = () => {
    isConnected.value = false
  }
}

function clearLogs() {
  logs.value = []
}

function getLogLevelClass(line: string) {
  if (line.includes('ERROR') || line.includes('FATAL') || line.includes('panic')) return 'text-error'
  if (line.includes('WARN')) return 'text-warn'
  if (line.includes('INFO')) return 'text-info'
  return ''
}

async function fetchConnections() {
  connLoading.value = true
  try {
    const res = await getActiveConnections()
    connections.value = res.connections || []
    uploadTotal.value = res.uploadTotal || 0
    downloadTotal.value = res.downloadTotal || 0
  } finally {
    connLoading.value = false
  }
}

async function fetchRawConfig() {
  try {
    const res = await getRawConfig()
    rawConfig.value = res.config
  } catch (err) {
    console.error(err)
  }
}

function copyConfig() {
  navigator.clipboard.writeText(rawConfig.value)
  message.success('配置已复制到剪贴板')
}

onMounted(() => {
  connectWS()
  fetchConnections()
  fetchRawConfig()
})

onUnmounted(() => {
  if (ws) {
    ws.close()
  }
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

.main-card {
  border-radius: 14px;
  background: #ffffff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.terminal-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.terminal-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #64748b;
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

.terminal-window {
  background: #0f172a;
  border-radius: 10px;
  padding: 16px;
  height: 520px;
  overflow-y: auto;
  font-family: 'JetBrains Mono', 'Fira Code', Menlo, monospace;
  font-size: 13px;
  line-height: 1.6;
  color: #e2e8f0;
}

.log-line {
  display: flex;
  gap: 12px;
  word-break: break-all;
}

.log-index {
  color: #475569;
  user-select: none;
  min-width: 40px;
}

.text-info {
  color: #38bdf8;
}

.text-warn {
  color: #fbbf24;
}

.text-error {
  color: #f87171;
  font-weight: bold;
}

.log-empty {
  color: #475569;
  text-align: center;
  padding: 40px 0;
}

.connections-summary {
  display: flex;
  gap: 24px;
  font-size: 13px;
  color: #475569;
  padding: 10px 14px;
  background: #f8fafc;
  border-radius: 8px;
}

.config-code {
  background: #0f172a;
  color: #38bdf8;
  padding: 16px;
  border-radius: 10px;
  max-height: 520px;
  overflow-y: auto;
  font-family: monospace;
  font-size: 13px;
}

.mt-4 {
  margin-top: 16px;
}

.mb-3 {
  margin-bottom: 12px;
}

.mb-2 {
  margin-bottom: 8px;
}
</style>
