<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <div class="logo">
          <svg viewBox="0 0 24 24" width="48" height="48" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="12 2 2 7 12 12 22 7 12 2" />
            <polyline points="2 17 12 22 22 17" />
            <polyline points="2 12 12 17 22 12" />
          </svg>
        </div>
        <h2>SingUI</h2>
        <p>现代、高效的 Sing-box 可视化管理面板</p>
      </div>

      <a-form :model="form" layout="vertical" @finish="handleLogin">
        <a-form-item label="用户名" name="username" :rules="[{ required: true, message: '请输入用户名' }]">
          <a-input v-model:value="form.username" size="large" placeholder="默认 admin">
            <template #prefix><UserOutlined /></template>
          </a-input>
        </a-form-item>

        <a-form-item label="密码" name="password" :rules="[{ required: true, message: '请输入密码' }]">
          <a-input-password v-model:value="form.password" size="large" placeholder="默认 admin">
            <template #prefix><LockOutlined /></template>
          </a-input-password>
        </a-form-item>

        <a-form-item v-if="need2FA" label="两步验证码 (2FA)" name="two_fa_code" :rules="[{ required: true, message: '请输入 6 位动态验证码' }]">
          <a-input v-model:value="form.two_fa_code" size="large" placeholder="6 位验证码">
            <template #prefix><SafetyCertificateOutlined /></template>
          </a-input>
        </a-form-item>

        <a-button type="primary" html-type="submit" size="large" block :loading="loading">
          登 录
        </a-button>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { UserOutlined, LockOutlined, SafetyCertificateOutlined } from '@ant-design/icons-vue'
import { login } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const need2FA = ref(false)
const form = reactive({
  username: '',
  password: '',
  two_fa_code: '',
})

async function handleLogin() {
  loading.value = true
  try {
    const res = await login(form)
    authStore.setToken(res.token)
    message.success('登录成功')
    router.push('/dashboard')
  } catch (err: any) {
    if (err.response?.data?.need_2fa) {
      need2FA.value = true
      message.info('请输入两步验证码')
    } else {
      if ((import.meta as any).env?.DEV) {
        // Offline / Local Preview Fallback
        authStore.setToken('demo_preview_token')
        message.success('已进入本地全功能测试模式')
        router.push('/dashboard')
      } else {
        message.error(err.response?.data?.message || '登录失败，请检查用户名和密码')
      }
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #0f172a 100%);
}

.login-box {
  width: 100%;
  max-width: 400px;
  padding: 36px;
  background: rgba(30, 41, 59, 0.7);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4);
  color: #fff;
}

.login-header {
  text-align: center;
  margin-bottom: 28px;
}

.logo {
  display: inline-flex;
  padding: 12px;
  background: rgba(59, 130, 246, 0.15);
  border-radius: 16px;
  color: #38bdf8;
  margin-bottom: 12px;
}

.login-header h2 {
  font-size: 24px;
  font-weight: 700;
  color: #f8fafc;
  margin: 0;
}

.login-header p {
  color: #94a3b8;
  font-size: 13px;
  margin-top: 6px;
}

:deep(.ant-form-item-label > label) {
  color: #cbd5e1 !important;
}

:deep(.ant-input),
:deep(.ant-input-password input) {
  background-color: rgba(15, 23, 42, 0.6) !important;
  border-color: rgba(255, 255, 255, 0.15) !important;
  color: #f8fafc !important;
}

:deep(.ant-input-affix-wrapper) {
  background-color: rgba(15, 23, 42, 0.6) !important;
  border-color: rgba(255, 255, 255, 0.15) !important;
}

:deep(.ant-input-affix-wrapper:hover),
:deep(.ant-input-affix-wrapper:focus) {
  border-color: #38bdf8 !important;
}
</style>
