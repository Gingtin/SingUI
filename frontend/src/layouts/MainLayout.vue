<template>
  <a-layout class="main-layout">
    <a-layout-sider v-model:collapsed="collapsed" :trigger="null" collapsible class="sider">
      <div class="logo">
        <svg viewBox="0 0 24 24" width="28" height="28" fill="none" stroke="currentColor" stroke-width="2">
          <polygon points="12 2 2 7 12 12 22 7 12 2" />
          <polyline points="2 17 12 22 22 17" />
          <polyline points="2 12 12 17 22 12" />
        </svg>
        <span v-if="!collapsed" class="logo-text">SingUI</span>
      </div>
      <a-menu v-model:selectedKeys="selectedKeys" theme="dark" mode="inline" @select="handleMenuSelect">
        <a-menu-item key="/dashboard">
          <template #icon><DashboardOutlined /></template>
          <span>系统仪表盘</span>
        </a-menu-item>
        <a-menu-item key="/inbounds">
          <template #icon><GlobalOutlined /></template>
          <span>入站节点</span>
        </a-menu-item>
        <a-menu-item key="/subscriptions">
          <template #icon><ShareAltOutlined /></template>
          <span>订阅中心</span>
        </a-menu-item>
        <a-menu-item key="/logs">
          <template #icon><CodeOutlined /></template>
          <span>日志诊断</span>
        </a-menu-item>
        <a-menu-item key="/settings">
          <template #icon><SettingOutlined /></template>
          <span>面板设置</span>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <a-layout-header class="header">
        <div class="header-left">
          <MenuUnfoldOutlined v-if="collapsed" class="trigger" @click="() => (collapsed = !collapsed)" />
          <MenuFoldOutlined v-else class="trigger" @click="() => (collapsed = !collapsed)" />
        </div>
        <div class="header-right">
          <a-dropdown>
            <span class="user-dropdown">
              <UserOutlined class="user-icon" />
              <span class="username">{{ authStore.user?.username || 'Admin' }}</span>
            </span>
            <template #overlay>
              <a-menu>
                <a-menu-item key="logout" @click="handleLogout">
                  <LogoutOutlined />
                  <span>退出登录</span>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>

      <a-layout-content class="content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  DashboardOutlined,
  GlobalOutlined,
  ShareAltOutlined,
  CodeOutlined,
  SettingOutlined,
  MenuUnfoldOutlined,
  MenuFoldOutlined,
  UserOutlined,
  LogoutOutlined,
} from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const collapsed = ref(false)
const selectedKeys = ref<string[]>([route.path])

watch(
  () => route.path,
  (val) => {
    selectedKeys.value = [val]
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
}

.sider {
  background: #0f172a;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 20px;
  color: #38bdf8;
  font-weight: 700;
  font-size: 18px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.logo-text {
  color: #f8fafc;
  letter-spacing: 0.5px;
}

.header {
  background: #ffffff;
  padding: 0 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  z-index: 10;
}

.trigger {
  font-size: 18px;
  cursor: pointer;
  transition: color 0.3s;
}

.trigger:hover {
  color: #3b82f6;
}

.user-dropdown {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 8px;
  border-radius: 6px;
  transition: background 0.2s;
}

.user-dropdown:hover {
  background: #f1f5f9;
}

.content {
  margin: 20px 24px;
  min-height: 280px;
}
</style>
