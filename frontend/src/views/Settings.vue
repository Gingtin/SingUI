<template>
  <div class="settings-page">
    <div class="page-header">
      <div>
        <h2>面板与系统设置</h2>
        <p class="subtitle">配置 Web 端口、Sing-box 核心参数、Telegram 机器人告警与数据库备份</p>
      </div>
      <a-button type="primary" :loading="saving" @click="handleSaveSettings">保存全局设置</a-button>
    </div>

    <a-row :gutter="[16, 16]" class="mt-4">
      <!-- Web & Sub Settings -->
      <a-col :xs="24" :md="12">
        <a-card title="基础与订阅配置" :bordered="false" class="settings-card">
          <a-form layout="vertical">
            <a-form-item label="面板监听端口">
              <a-input v-model:value="settings.web_port" placeholder="2096" />
            </a-form-item>
            <a-form-item label="面板基础路径 (Base Path)">
              <a-input v-model:value="settings.web_base_path" placeholder="/" />
            </a-form-item>
            <a-form-item label="订阅外部域名 / IP">
              <a-input v-model:value="settings.sub_domain" placeholder="留空则自动使用当前访问的主机名" />
            </a-form-item>
            <a-form-item label="订阅标题">
              <a-input v-model:value="settings.sub_title" placeholder="SingBox UI Nodes" />
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>

      <!-- Sing-box Core Settings -->
      <a-col :xs="24" :md="12">
        <a-card title="Sing-box 核心与 Clash API" :bordered="false" class="settings-card">
          <a-form layout="vertical">
            <a-form-item label="Sing-box 可执行文件路径">
              <a-input v-model:value="settings.singbox_bin_path" placeholder="sing-box 或 /usr/local/bin/sing-box" />
            </a-form-item>
            <a-form-item label="配置文件生成路径">
              <a-input v-model:value="settings.singbox_config_path" placeholder="config/singbox_config.json" />
            </a-form-item>
            <a-form-item label="Clash API 端口">
              <a-input v-model:value="settings.clash_api_port" placeholder="9090" />
            </a-form-item>
            <a-form-item label="Clash API Secret 密钥">
              <a-input v-model:value="settings.clash_api_secret" placeholder="API 访问密钥" />
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>

      <!-- Telegram Bot Settings -->
      <a-col :xs="24" :md="12">
        <a-card title="Telegram 监控与告警机器人" :bordered="false" class="settings-card">
          <a-form layout="vertical">
            <a-form-item label="Bot Token">
              <a-input v-model:value="settings.tg_bot_token" placeholder="从 @BotFather 获取" />
            </a-form-item>
            <a-form-item label="Chat ID">
              <a-input v-model:value="settings.tg_chat_id" placeholder="管理员或群组 Chat ID" />
            </a-form-item>
            <a-form-item label="到期及超额提醒">
              <a-switch v-model:checked="tgNotifyOnExpire" />
              <span class="ml-2 text-muted">当用户流量超额或到期时推送通知</span>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>

      <!-- Security & Backup Settings -->
      <a-col :xs="24" :md="12">
        <a-card title="安全设置与数据备份" :bordered="false" class="settings-card">
          <div class="section-subtitle">修改管理员密码</div>
          <a-form layout="vertical" :model="pwdForm" @finish="handleChangePassword">
            <a-form-item label="当前密码" name="old_password" required>
              <a-input-password v-model:value="pwdForm.old_password" />
            </a-form-item>
            <a-form-item label="新密码" name="new_password" required>
              <a-input-password v-model:value="pwdForm.new_password" />
            </a-form-item>
            <a-button type="primary" html-type="submit" :loading="pwdLoading">更新密码</a-button>
          </a-form>

          <a-divider />

          <div class="section-subtitle">数据库备份与恢复</div>
          <div class="backup-actions mt-3">
            <a-button @click="downloadBackupFile">
              <template #prefix><DownloadOutlined /></template>
              下载 SQLite 备份
            </a-button>
            <a-upload
              name="backup_file"
              action="/api/settings/restore"
              :headers="uploadHeaders"
              :show-upload-list="false"
              @change="handleRestoreUpload"
            >
              <a-button>
                <template #prefix><UploadOutlined /></template>
                上传并恢复备份
              </a-button>
            </a-upload>
          </div>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { DownloadOutlined, UploadOutlined } from '@ant-design/icons-vue'
import { getSettings, updateSettings, updatePassword } from '@/api/setting'

const settings = reactive<Record<string, string>>({})
const saving = ref(false)
const pwdLoading = ref(false)

const tgNotifyOnExpire = ref(true)

const pwdForm = reactive({
  old_password: '',
  new_password: '',
})

const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${localStorage.getItem('token')}`,
}))

async function fetchSettings() {
  try {
    const data = await getSettings()
    Object.assign(settings, data)
    tgNotifyOnExpire.value = settings.tg_notify_on_expire === 'true'
  } catch (err) {
    console.error(err)
  }
}

async function handleSaveSettings() {
  saving.value = true
  try {
    settings.tg_notify_on_expire = tgNotifyOnExpire.value ? 'true' : 'false'
    await updateSettings(settings)
    message.success('设置保存成功')
  } finally {
    saving.value = false
  }
}

async function handleChangePassword() {
  pwdLoading.value = true
  try {
    await updatePassword(pwdForm)
    message.success('管理员密码修改成功')
    pwdForm.old_password = ''
    pwdForm.new_password = ''
  } finally {
    pwdLoading.value = false
  }
}

function downloadBackupFile() {
  const token = localStorage.getItem('token')
  window.open(`/api/settings/backup?token=${token}`, '_blank')
}

function handleRestoreUpload(info: any) {
  if (info.file.status === 'done') {
    message.success('备份恢复成功，请重启面板')
  } else if (info.file.status === 'error') {
    message.error('备份恢复失败')
  }
}

onMounted(() => {
  fetchSettings()
})
</script>

<style scoped>
.settings-page {
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

.settings-card {
  border-radius: 12px;
  background: #ffffff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  height: 100%;
}

.section-subtitle {
  font-weight: 600;
  color: #334155;
  margin-bottom: 12px;
}

.backup-actions {
  display: flex;
  gap: 12px;
}

.ml-2 {
  margin-left: 8px;
}

.text-muted {
  font-size: 12px;
  color: #94a3b8;
}

.mt-4 {
  margin-top: 16px;
}

.mt-3 {
  margin-top: 12px;
}
</style>
