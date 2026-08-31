<template>
  <a-layout class="main-layout">
    <!-- macOS Style Sidebar -->
    <a-layout-sider v-model:collapsed="collapsed" :trigger="null" collapsible class="sider" width="230">
      <!-- macOS Window Traffic Light Controls -->
      <div class="macos-window-header">
        <div class="traffic-lights">
          <span class="light close"></span>
          <span class="light minimize"></span>
          <span class="light maximize"></span>
        </div>
      </div>

      <!-- App Logo -->
      <div class="logo">
        <div class="logo-icon-wrapper">
          <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2.5">
            <polygon points="12 2 2 7 12 12 22 7 12 2" />
            <polyline points="2 17 12 22 22 17" />
            <polyline points="2 12 12 17 22 12" />
          </svg>
        </div>
        <span v-if="!collapsed" class="logo-text">SingUI</span>
        <span v-if="!collapsed" class="logo-badge">Pro</span>
      </div>

      <!-- Menu Items -->
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @select="handleMenuSelect">
        <a-menu-item key="/dashboard">
          <template #icon><DashboardOutlined /></template>
          <span>系统仪表盘</span>
        </a-menu-item>
        <a-menu-item key="/inbounds">
          <template #icon><GlobalOutlined /></template>
          <span>入站节点管理</span>
        </a-menu-item>
        <a-menu-item key="/outbounds">
          <template #icon><SendOutlined /></template>
          <span>出站链路与WARP</span>
        </a-menu-item>
        <a-menu-item key="/routing">
          <template #icon><ForkOutlined /></template>
          <span>规则路由与DNS</span>
        </a-menu-item>
        <a-menu-item key="/subscriptions">
          <template #icon><ShareAltOutlined /></template>
          <span>全能订阅中心</span>
        </a-menu-item>
        <a-menu-item key="/logs">
          <template #icon><CodeOutlined /></template>
          <span>核心配置与日志</span>
        </a-menu-item>
        <a-menu-item key="/settings">
          <template #icon><SettingOutlined /></template>
          <span>面板设置与更新</span>
        </a-menu-item>
        <a-menu-item key="/docs">
          <template #icon><FileTextOutlined /></template>
          <span>REST API 文档</span>
        </a-menu-item>
      </a-menu>

      <!-- Sider Footer Status Pill -->
      <div v-if="!collapsed" class="sider-footer">
        <div class="footer-chip">
          <span class="status-dot"></span>
          <span>{{ coreVersion ? coreVersion.split('\n')[0] : 'Sing-box 未知版本' }}</span>
        </div>
      </div>
    </a-layout-sider>

    <a-layout>
      <!-- macOS Style Frosted Glass Header -->
      <a-layout-header class="header">
        <div class="header-left">
          <div class="collapse-trigger" @click="() => (collapsed = !collapsed)">
            <MenuUnfoldOutlined v-if="collapsed" />
            <MenuFoldOutlined v-else />
          </div>
          <div class="page-breadcrumb">
            <span class="platform-tag">macOS Native Style</span>
          </div>
        </div>

        <div class="header-right">
          <!-- Live RAM Indicator -->
          <div class="system-chip memory">
            <ThunderboltOutlined class="chip-icon" />
            <span>内存常驻: <b>{{ formatBytes(memUsed) }}</b></span>
          </div>

          <!-- Core Status Pill -->
          <div class="system-chip core">
            <span class="pulse-dot" :class="{ stopped: !coreRunning }"></span>
            <span>Sing-box: <b>{{ coreRunning ? '运行中' : '已停止' }}</b></span>
          </div>

          <!-- User Profile Dropdown -->
          <a-dropdown placement="bottomRight">
            <div class="user-pill">
              <div class="user-avatar">A</div>
              <span class="username">{{ authStore.user?.username || 'Administrator' }}</span>
            </div>
            <template #overlay>
              <a-menu class="user-menu">
                <a-menu-item key="settings" @click="router.push('/settings')">
                  <SettingOutlined />
                  <span>面板设置</span>
                </a-menu-item>
                <a-menu-divider />
                <a-menu-item key="logout" @click="handleLogout">
                  <LogoutOutlined />
                  <span>退出登录</span>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>

      <!-- Main Dynamic Content Container -->
      <a-layout-content class="content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  DashboardOutlined,
  GlobalOutlined,
  SendOutlined,
  ForkOutlined,
  ShareAltOutlined,
  CodeOutlined,
  SettingOutlined,
  FileTextOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  ThunderboltOutlined,
} from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { getSystemStatus, getCoreStatus } from '@/api/server'
import { formatBytes } from '@/utils/format'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const collapsed = ref(false)
const selectedKeys = ref<string[]>([route.path])

const memUsed = ref(0)
const coreRunning = ref(false)
const coreVersion = ref('')
let pollTimer: any = null

async function fetchHeaderStatus() {
  try {
    const sys = await getSystemStatus()
    memUsed.value = sys.mem_used || 0
  } catch(e) {}
  try {
    const core = await getCoreStatus()
    coreRunning.value = core.is_running || false
    coreVersion.value = core.version || ''
  } catch(e) {}
}

onMounted(() => {
  fetchHeaderStatus()
  pollTimer = setInterval(fetchHeaderStatus, 5000)
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})

watch(
  () => route.path,
  (path) => {
    selectedKeys.value = [path]
  }
)

function handleMenuSelect({ key }: { key: string }) {
  router.push(key)
}

function handleLogout() {
  authStore.clearToken()
  router.push('/login')
}
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
  background: #f8fafc;
}

.sider {
  background: #0f172a !important;
  box-shadow: 2px 0 12px rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
  z-index: 10;
}

.macos-window-header {
  padding: 14px 16px 6px 16px;
}

.traffic-lights {
  display: flex;
  gap: 7px;
}

.light {
  width: 11px;
  height: 11px;
  border-radius: 50%;
  display: inline-block;
}

.light.close {
  background: #ff5f56;
  box-shadow: 0 0 4px rgba(255, 95, 86, 0.6);
}

.light.minimize {
  background: #ffbd2e;
  box-shadow: 0 0 4px rgba(255, 189, 46, 0.6);
}

.light.maximize {
  background: #27c93f;
  box-shadow: 0 0 4px rgba(39, 201, 63, 0.6);
}

.logo {
  height: 52px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 16px;
  margin-bottom: 8px;
}

.logo-icon-wrapper {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: linear-gradient(135deg, #3b82f6, #6366f1);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  box-shadow: 0 4px 10px rgba(59, 130, 246, 0.4);
}

.logo-text {
  color: #f8fafc;
  font-size: 17px;
  font-weight: 700;
  letter-spacing: 0.5px;
}

.logo-badge {
  font-size: 10px;
  font-weight: 700;
  background: rgba(59, 130, 246, 0.2);
  color: #60a5fa;
  padding: 1px 6px;
  border-radius: 4px;
  border: 1px solid rgba(96, 165, 250, 0.3);
}

.sider-footer {
  position: absolute;
  bottom: 16px;
  left: 16px;
  right: 16px;
}

.footer-chip {
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 20px;
  padding: 6px 12px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #94a3b8;
  font-size: 11px;
}

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #10b981;
  box-shadow: 0 0 6px #10b981;
}

.header {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  padding: 0 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 56px;
  border-bottom: 1px solid rgba(226, 232, 240, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.collapse-trigger {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #64748b;
  transition: all 0.2s;
}

.collapse-trigger:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.platform-tag {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  background: #f1f5f9;
  padding: 3px 8px;
  border-radius: 12px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.system-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 12px;
}

.system-chip.memory {
  background: #f0fdf4;
  color: #15803d;
  border: 1px solid #dcfce7;
}

.chip-icon {
  font-size: 12px;
}

.system-chip.core {
  background: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #dbeafe;
}

.pulse-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #3b82f6;
  box-shadow: 0 0 6px #3b82f6;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(59, 130, 246, 0.7); }
  70% { transform: scale(1); box-shadow: 0 0 0 6px rgba(59, 130, 246, 0); }
  100% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(59, 130, 246, 0); }
}

.user-pill {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 3px 10px 3px 4px;
  border-radius: 20px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  cursor: pointer;
  transition: all 0.2s;
}

.user-pill:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
}

.user-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

.username {
  font-size: 13px;
  font-weight: 600;
  color: #334155;
}

.content {
  margin: 16px 20px;
  padding: 0;
  min-height: 280px;
}
</style>
