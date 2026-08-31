<template>
  <div class="settings-page">
    <div class="page-header">
      <div>
        <h2>面板全局设置 (Panel Settings)</h2>
        <p class="subtitle">深度管理 Web 服务、订阅分发策略、Telegram 自动化告警与 SQLite 备份</p>
      </div>
      <div class="header-actions">
        <a-button type="primary" :loading="saving" @click="handleSaveSettings">
          保存全部配置
        </a-button>
      </div>
    </div>

    <a-card class="mt-4 main-card" :bordered="false">
      <a-tabs v-model:activeKey="activeTab">
        <!-- Tab 1: Web Server & Security -->
        <a-tab-pane key="web" tab="1. 面板与安全 (Web & Auth)">
          <div class="form-section">
            <a-form layout="vertical">
              <a-row :gutter="24">
                <a-col :xs="24" :md="12">
                  <a-form-item label="面板监听端口 (Web Port)" extra="修改后需重启面板生效">
                    <a-input v-model:value="settings.web_port" placeholder="2096" />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :md="12">
                  <a-form-item label="安全访问根路径 (Web Base Path)" extra="如: /my_secret_panel/">
                    <a-input v-model:value="settings.web_base_path" placeholder="/" />
                  </a-form-item>
                </a-col>
              </a-row>

              <a-divider>管理员密码修改</a-divider>
              <a-row :gutter="24">
                <a-col :xs="24" :md="12">
                  <a-form-item label="当前原密码">
                    <a-input-password v-model:value="pwdForm.old_password" placeholder="原密码" />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :md="12">
                  <a-form-item label="设置新密码">
                    <a-input-password v-model:value="pwdForm.new_password" placeholder="新密码" />
                  </a-form-item>
                </a-col>
              </a-row>
              <a-button type="dashed" @click="handleChangePassword">确认修改密码</a-button>
            </a-form>
          </div>
        </a-tab-pane>

        <!-- Tab 2: Subscription Hub -->
        <a-tab-pane key="sub" tab="2. 订阅分发配置 (Subscriptions)">
          <div class="form-section">
            <a-form layout="vertical">
              <a-row :gutter="24">
                <a-col :xs="24" :md="12">
                  <a-form-item label="订阅下发域名 / IP (Sub Domain)" extra="生成订阅链接和节点链接时使用的公网连接地址（留空则默认使用访问面板的域名/IP）">
                    <a-input v-model:value="settings.sub_domain" placeholder="node.yourdomain.com 或 1.2.3.4" />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :md="12">
                  <a-form-item label="订阅路径前缀 (Sub Path)">
                    <a-input v-model:value="settings.sub_path" placeholder="/sub" />
                  </a-form-item>
                </a-col>
              </a-row>

              <a-row :gutter="24">
                <a-col :xs="24" :md="12">
                  <a-form-item label="订阅节点名称前缀 (Sub Title)">
                    <a-input v-model:value="settings.sub_title" placeholder="SingUI Nodes" />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :md="12">
                  <a-form-item label="流量统计保留天数 (Traffic Log Days)">
                    <a-input v-model:value="settings.traffic_log_days" placeholder="30" />
                  </a-form-item>
                </a-col>
              </a-row>
            </a-form>
          </div>
        </a-tab-pane>

        <!-- Tab 3: Telegram Bot & Alerts -->
        <a-tab-pane key="telegram" tab="3. Telegram 机器人通知 (Telegram Bot)">
          <div class="form-section">
            <a-form layout="vertical">
              <a-row :gutter="24">
                <a-col :xs="24" :md="12">
                  <a-form-item label="Telegram Bot Token" extra="从 @BotFather 获取的机器令牌">
                    <a-input v-model:value="settings.tg_bot_token" placeholder="123456789:ABCdefGhIJKlmNoPQRstuVWXyz" />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :md="12">
                  <a-form-item label="管理员 Chat ID" extra="从 @userinfobot 获取的接收通知 ID">
                    <a-input v-model:value="settings.tg_chat_id" placeholder="123456789" />
                  </a-form-item>
                </a-col>
              </a-row>

              <a-row :gutter="24">
                <a-col :xs="24" :md="8">
                  <a-form-item label="流量耗尽自动熔断告警">
                    <a-switch v-model:checked="tgQuotaNotify" />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :md="8">
                  <a-form-item label="用户到期提前提醒">
                    <a-switch v-model:checked="tgExpireNotify" />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :md="8">
                  <a-form-item label="每周定时备份推送">
                    <a-switch v-model:checked="tgAutoBackup" />
                  </a-form-item>
                </a-col>
              </a-row>
            </a-form>
          </div>
        </a-tab-pane>

        <!-- Tab 4: Backup & Restore -->
        <a-tab-pane key="backup" tab="4. 数据库备份与灾备 (Backup & Restore)">
          <div class="form-section">
            <div class="backup-grid">
              <div class="backup-card">
                <h4>💾 导出 SQLite 数据库热备份</h4>
                <p>一键下载当前面板所有入站节点、客户端用户、路由规则与系统设置的完整快照。</p>
                <a-button type="primary" @click="handleDownloadBackup">立即下载备份文件 (.db)</a-button>
              </div>

              <div class="backup-card">
                <h4>🔄 上传并恢复数据库</h4>
                <p>选择之前导出的 <code>.db</code> 备份文件上传，系统将自动验证并热恢复数据。</p>
                <a-upload :before-upload="handleUploadBackup" :show-upload-list="false">
                  <a-button>选择 .db 备份文件上传恢复</a-button>
                </a-upload>
              </div>
            </div>
          </div>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { getSettings, updateSettings, updatePassword, downloadBackup, restoreBackup } from '@/api/setting'

const activeTab = ref('web')
const saving = ref(false)

const settings = reactive<Record<string, string>>({
  web_port: '2096',
  web_base_path: '/',
  sub_domain: '',
  sub_path: '/sub',
  sub_title: 'SingUI Nodes',
  traffic_log_days: '30',
  tg_bot_token: '',
  tg_chat_id: '',
  tg_notify_on_quota: 'true',
  tg_notify_on_expire: 'true',
  auto_backup_enabled: 'false',
})

const tgQuotaNotify = ref(true)
const tgExpireNotify = ref(true)
const tgAutoBackup = ref(false)

const pwdForm = reactive({
  old_password: '',
  new_password: '',
})

async function fetchSettings() {
  try {
    const data = await getSettings()
    Object.assign(settings, data)
    tgQuotaNotify.value = settings.tg_notify_on_quota === 'true'
    tgExpireNotify.value = settings.tg_notify_on_expire === 'true'
    tgAutoBackup.value = settings.auto_backup_enabled === 'true'
  } catch (err) {
    console.error(err)
  }
}

async function handleSaveSettings() {
  saving.value = true
  try {
    settings.tg_notify_on_quota = tgQuotaNotify.value ? 'true' : 'false'
    settings.tg_notify_on_expire = tgExpireNotify.value ? 'true' : 'false'
    settings.auto_backup_enabled = tgAutoBackup.value ? 'true' : 'false'

    await updateSettings(settings)
    message.success('面板配置已保存')
  } finally {
    saving.value = false
  }
}

async function handleChangePassword() {
  if (!pwdForm.old_password || !pwdForm.new_password) {
    message.warning('请输入原密码和新密码')
    return
  }
  try {
    await updatePassword(pwdForm.old_password, pwdForm.new_password)
    message.success('密码修改成功，请牢记新密码')
    pwdForm.old_password = ''
    pwdForm.new_password = ''
  } catch (err: any) {
    message.error(err.response?.data?.error || '密码修改失败')
  }
}

async function handleDownloadBackup() {
  try {
    const res = await downloadBackup()
    const blob = new Blob([res.data])
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `singbox_ui_backup_${new Date().toISOString().slice(0, 10)}.db`
    link.click()
    window.URL.revokeObjectURL(url)
    message.success('备份文件下载完成')
  } catch (err) {
    message.error('下载备份失败')
  }
}

async function handleUploadBackup(file: File) {
  const formData = new FormData()
  formData.append('backup', file)
  try {
    await restoreBackup(formData)
    message.success('数据库已恢复成功，请刷新页面')
    fetchSettings()
  } catch (err) {
    message.error('数据库恢复失败')
  }
  return false
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

.main-card {
  border-radius: 14px;
  background: #ffffff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.form-section {
  max-width: 860px;
  padding: 12px 0;
}

.backup-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 20px;
  margin-top: 10px;
}

.backup-card {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.backup-card h4 {
  font-size: 15px;
  font-weight: 700;
  margin-bottom: 8px;
  color: #0f172a;
}

.backup-card p {
  color: #64748b;
  font-size: 13px;
  margin-bottom: 16px;
  line-height: 1.5;
}

.mt-4 {
  margin-top: 16px;
}
</style>
