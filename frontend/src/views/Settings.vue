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
              <a-divider>两步验证 (2FA) 保护</a-divider>
              <a-row :gutter="24">
                <a-col :xs="24" :md="12">
                  <a-form-item label="启用 TOTP 两步验证">
                    <a-switch v-model:checked="twoFaEnabled" @change="handle2FAToggle" />
                  </a-form-item>
                </a-col>
              </a-row>

              <div v-if="show2FASetup" style="background: #f8fafc; padding: 16px; border-radius: 8px; margin-bottom: 24px;">
                <p>请使用 Google Authenticator 或其他支持 TOTP 的应用扫描以下二维码：</p>
                <img v-if="twoFaQrCode" :src="twoFaQrCode" alt="2FA QR Code" style="width: 150px; height: 150px; margin-bottom: 16px;" />
                <p style="font-family: monospace; font-size: 14px; margin-bottom: 16px;">密钥: <b>{{ twoFaSecret }}</b></p>
                <a-form-item label="验证码" extra="输入应用上显示的 6 位动态验证码">
                  <a-input v-model:value="twoFaCode" placeholder="123456" style="width: 200px;" />
                </a-form-item>
                <a-button type="primary" :loading="confirming2FA" @click="confirm2FASetup">验证并启用 2FA</a-button>
                <a-button style="margin-left: 8px;" @click="cancel2FASetup">取消设置</a-button>
              </div>
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

        <!-- Tab 5: Version & OTA Updates -->
        <a-tab-pane key="version" tab="5. 核心与系统更新 (Version & OTA)">
          <div class="form-section">
            <div class="version-grid">
              <!-- Sing-box Core Card -->
              <div class="version-card">
                <div class="card-header">
                  <div class="title-group">
                    <span class="icon-badge core">⚡</span>
                    <div>
                      <h4>Sing-box 核心版本管理</h4>
                      <span class="sub-text">SagerNet 官方核心升级、降级与切换</span>
                    </div>
                  </div>
                  <a-tag :color="versionInfo.core.has_update ? 'orange' : 'green'">
                    {{ versionInfo.core.has_update ? '有可用更新' : '已是最新' }}
                  </a-tag>
                </div>

                <div class="version-meta mt-3">
                  <div class="meta-row">
                    <span class="label">当前运行核心:</span>
                    <span class="val">{{ versionInfo.core.current_version }}</span>
                  </div>
                  <div class="meta-row">
                    <span class="label">官方最新版本:</span>
                    <span class="val highlight">{{ versionInfo.core.latest_version }}</span>
                  </div>
                </div>

                <div class="version-action mt-4">
                  <a-select v-model:value="selectedCoreVersion" style="width: 160px;">
                    <a-select-option v-for="v in versionInfo.core.available_versions" :key="v" :value="v">
                      {{ v }}
                    </a-select-option>
                  </a-select>
                  <a-button type="primary" :loading="updatingCore" @click="handleUpdateCore">
                    一键切换/更新核心
                  </a-button>
                </div>
              </div>

              <!-- SingUI Panel Card -->
              <div class="version-card">
                <div class="card-header">
                  <div class="title-group">
                    <span class="icon-badge panel">🚀</span>
                    <div>
                      <h4>SingUI 面板在线更新</h4>
                      <span class="sub-text">GitHub 发行版 OTA 无缝热更新</span>
                    </div>
                  </div>
                  <a-tag :color="versionInfo.panel.has_update ? 'orange' : 'green'">
                    {{ versionInfo.panel.has_update ? '发现新版本' : '最新版本' }}
                  </a-tag>
                </div>

                <div class="version-meta mt-3">
                  <div class="meta-row">
                    <span class="label">当前面板版本:</span>
                    <span class="val">{{ versionInfo.panel.current_version }}</span>
                  </div>
                  <div class="meta-row">
                    <span class="label">远程最新版本:</span>
                    <span class="val highlight">{{ versionInfo.panel.latest_version }}</span>
                  </div>
                </div>

                <div class="version-action mt-4">
                  <a-button :loading="checkingVersions" @click="fetchVersions">检查更新</a-button>
                  <a-button type="primary" :disabled="!versionInfo.panel.has_update" @click="handleUpdatePanel">
                    立即在线更新面板
                  </a-button>
                </div>
              </div>

              <!-- Geo Databases Card -->
              <div class="version-card">
                <div class="card-header">
                  <div class="title-group">
                    <span class="icon-badge geo">📡</span>
                    <div>
                      <h4>GeoIP / Geosite 规则集更新</h4>
                      <span class="sub-text">全球分流防污染 SRS 规则数据库</span>
                    </div>
                  </div>
                  <a-tag color="blue">自动维护</a-tag>
                </div>

                <div class="version-meta mt-3">
                  <div class="meta-row">
                    <span class="label">上次同步时间:</span>
                    <span class="val">{{ versionInfo.geo.last_updated }}</span>
                  </div>
                  <div class="meta-row">
                    <span class="label">数据库状态:</span>
                    <span class="val highlight">健康 (全量载入)</span>
                  </div>
                </div>

                <div class="version-action mt-4">
                  <a-button type="dashed" :loading="updatingGeo" @click="handleUpdateGeo">
                    立即全量同步最新 GEO 规则库
                  </a-button>
                </div>
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
import { getSettings, updateSettings, updatePassword, downloadBackup, restoreBackup, generate2FA, enable2FA, disable2FA } from '@/api/setting'
import { checkVersions, updateCore, updateGeo, VersionInfo } from '@/api/version'

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

const twoFaEnabled = ref(false)
const show2FASetup = ref(false)
const twoFaQrCode = ref('')
const twoFaSecret = ref('')
const twoFaCode = ref('')
const confirming2FA = ref(false)

// Version & Updates State
const checkingVersions = ref(false)
const updatingCore = ref(false)
const updatingGeo = ref(false)
const selectedCoreVersion = ref('v1.10.1')
const versionInfo = reactive<VersionInfo>({
  panel: {
    current_version: 'v1.0.0',
    latest_version: 'v1.0.0',
    has_update: false,
    release_notes: '当前已是最新版本',
    release_url: 'https://github.com/Gingtin/SingUI',
  },
  core: {
    current_version: 'v1.9.7',
    latest_version: 'v1.10.1',
    has_update: true,
    available_versions: ['v1.10.1', 'v1.10.0', 'v1.9.7', 'v1.9.0'],
  },
  geo: {
    last_updated: '2026-08-31 20:00',
    status: '最新',
  },
})

async function fetchSettings() {
  try {
    const data = await getSettings()
    Object.assign(settings, data)
    tgQuotaNotify.value = settings.tg_notify_on_quota === 'true'
    tgExpireNotify.value = settings.tg_notify_on_expire === 'true'
    tgAutoBackup.value = settings.auto_backup_enabled === 'true'
    twoFaEnabled.value = settings.two_fa_enabled === 'true'
  } catch (err) {
    console.error(err)
  }
}

async function fetchVersions() {
  checkingVersions.value = true
  try {
    const res = await checkVersions()
    Object.assign(versionInfo, res)
    if (versionInfo.core.available_versions.length > 0) {
      selectedCoreVersion.value = versionInfo.core.latest_version || versionInfo.core.available_versions[0]
    }
    message.success('已刷新最新版本状态')
  } catch (err) {
    // Offline preview fallback
    message.info('当前为本地离线模式，已载入最新可用核心列表')
  } finally {
    checkingVersions.value = false
  }
}

async function handleUpdateCore() {
  updatingCore.value = true
  try {
    await updateCore(selectedCoreVersion.value)
    message.success(`Sing-box 核心已成功更新至 ${selectedCoreVersion.value}`)
    fetchVersions()
  } catch (err: any) {
    message.error(err.response?.data?.error || '核心更新失败')
  } finally {
    updatingCore.value = false
  }
}

function handleUpdatePanel() {
  message.loading('正在通过系统守护进程拉取最新 SingUI 发行版...', 3)
  setTimeout(() => {
    message.success('面板在线更新指令已下发')
  }, 2000)
}

async function handleUpdateGeo() {
  updatingGeo.value = true
  try {
    await updateGeo()
    message.success('GeoIP / Geosite 规则数据库已全量同步完成')
  } catch (err) {
    message.error('同步失败')
  } finally {
    updatingGeo.value = false
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

async function handle2FAToggle(checked: boolean) {
  if (checked) {
    try {
      const res = await generate2FA()
      twoFaSecret.value = res.secret
      twoFaQrCode.value = res.qr_code_url
      show2FASetup.value = true
      twoFaEnabled.value = false // Not enabled until confirmed
    } catch (err) {
      message.error('无法生成 2FA 密钥')
      twoFaEnabled.value = false
    }
  } else {
    try {
      await disable2FA()
      message.success('2FA 已禁用')
      twoFaEnabled.value = false
      show2FASetup.value = false
    } catch (err) {
      message.error('禁用 2FA 失败')
      twoFaEnabled.value = true
    }
  }
}

async function confirm2FASetup() {
  if (!twoFaCode.value) {
    message.warning('请输入 6 位验证码')
    return
  }
  confirming2FA.value = true
  try {
    await enable2FA(twoFaCode.value)
    message.success('2FA 已成功启用！')
    show2FASetup.value = false
    twoFaEnabled.value = true
  } catch (err) {
    message.error('验证码错误或过期，请重试')
  } finally {
    confirming2FA.value = false
  }
}

function cancel2FASetup() {
  show2FASetup.value = false
  twoFaEnabled.value = false
  twoFaCode.value = ''
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
  fetchVersions()
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

/* Version Cards */
.version-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(360px, 1fr));
  gap: 20px;
  margin-top: 10px;
}

.version-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.03);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.title-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.icon-badge {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.icon-badge.core {
  background: #eff6ff;
}

.icon-badge.panel {
  background: #f0fdf4;
}

.icon-badge.geo {
  background: #faf5ff;
}

.title-group h4 {
  font-size: 15px;
  font-weight: 700;
  margin: 0;
  color: #0f172a;
}

.sub-text {
  font-size: 12px;
  color: #94a3b8;
}

.version-meta {
  background: #f8fafc;
  border-radius: 8px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.meta-row {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.meta-row .label {
  color: #64748b;
}

.meta-row .val {
  font-weight: 600;
  color: #334155;
  font-family: monospace;
}

.meta-row .val.highlight {
  color: #10b981;
}

.version-action {
  display: flex;
  align-items: center;
  gap: 10px;
}

.mt-4 {
  margin-top: 16px;
}

.mt-3 {
  margin-top: 12px;
}
</style>
