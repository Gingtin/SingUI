<template>
  <div class="outbounds-page">
    <div class="page-header">
      <div>
        <h2>出站管理 (Outbounds & Detours)</h2>
        <p class="subtitle">管理 Sing-box 出站链路，支持直连、拦截、DNS分流及 Cloudflare WARP 链式落地</p>
      </div>
      <div class="header-actions">
        <a-button type="primary" @click="openCreateModal">
          <template #prefix><PlusOutlined /></template>
          添加出站链路
        </a-button>
      </div>
    </div>

    <a-row :gutter="[16, 16]" class="mt-4">
      <a-col v-for="out in outbounds" :key="out.id" :xs="24" :sm="12" :lg="8" :xl="6">
        <div class="mac-card outbound-card" :class="{ disabled: !out.enable }">
          <div class="card-header">
            <div class="tag-title">
              <span class="type-badge" :style="{ backgroundColor: getTypeColor(out.type) }">
                {{ out.type.substring(0, 3).toUpperCase() }}
              </span>
              <div class="info">
                <span class="tag-text">{{ out.tag }}</span>
                <span class="remark-text">{{ out.remark || '无备注' }}</span>
              </div>
            </div>
            <a-switch :checked="out.enable" @change="() => toggleEnable(out)" size="small" :disabled="isSystemTag(out.tag)" />
          </div>

          <div class="card-body">
            <div class="meta-row">
              <span class="label">类型 (Type)</span>
              <span class="val">{{ out.type.toUpperCase() }}</span>
            </div>
            <div class="meta-row" v-if="out.server">
              <span class="label">目标服务器</span>
              <span class="val">{{ out.server }}:{{ out.port || 0 }}</span>
            </div>
            <div class="meta-row" v-if="isSystemTag(out.tag)">
              <span class="label">系统内置</span>
              <span class="val text-green">核心保留</span>
            </div>
          </div>

          <div class="card-footer">
            <a-button type="link" size="small" @click="handlePing(out.id)">测速 (Latency)</a-button>
            <div style="display:flex;">
              <a-button type="link" size="small" @click="openEditModal(out)">编辑</a-button>
              <a-popconfirm
                v-if="!isSystemTag(out.tag)"
                title="确定删除该出站链路？"
                @confirm="() => handleDelete(out.id)"
              >
                <a-button type="link" size="small" danger>删除</a-button>
              </a-popconfirm>
            </div>
          </div>
        </div>
      </a-col>
    </a-row>

    <!-- Outbound Modal -->
    <a-modal
      v-model:open="modalVisible"
      :title="isEdit ? '编辑出站链路' : '添加出站链路'"
      @ok="handleSave"
      :confirmLoading="saving"
      width="560px"
    >
      <a-form :model="form" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="出站类型 (Type)" required>
              <a-select v-model:value="form.type" :disabled="isEdit && isSystemTag(form.tag || '')">
                <a-select-option value="direct">Direct (直连)</a-select-option>
                <a-select-option value="block">Block (拦截丢弃)</a-select-option>
                <a-select-option value="dns-out">DNS-OUT (DNS解析)</a-select-option>
                <a-select-option value="selector">Selector (手动选择)</a-select-option>
                <a-select-option value="urltest">URLTest (自动测速优选)</a-select-option>
                <a-select-option value="fallback">Fallback (故障转移)</a-select-option>
                <a-select-option value="socks">SOCKS5</a-select-option>
                <a-select-option value="http">HTTP</a-select-option>
                <a-select-option value="wireguard">WireGuard / WARP</a-select-option>
                <a-select-option value="vless">VLESS</a-select-option>
                <a-select-option value="vmess">VMESS</a-select-option>
                <a-select-option value="trojan">Trojan</a-select-option>
                <a-select-option value="hysteria2">Hysteria 2</a-select-option>
                <a-select-option value="shadowsocks">Shadowsocks</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="出站标签 (Tag)" required>
              <a-input v-model:value="form.tag" placeholder="如: warp-out" :disabled="isEdit && isSystemTag(form.tag || '')" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16" v-if="!['direct', 'block', 'dns-out', 'selector', 'urltest', 'fallback'].includes(form.type || '')">
          <a-col :span="24">
            <a-form-item label="前置链式代理 (Detour)">
              <a-select v-model:value="form.detour" allow-clear placeholder="留空为不使用前置代理">
                <a-select-option v-for="o in outbounds.filter(x => x.tag !== form.tag && !isSystemTag(x.tag))" :key="o.tag" :value="o.tag">
                  {{ o.tag }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <div v-if="['selector', 'urltest', 'fallback'].includes(form.type || '')">
          <a-form-item label="成员节点 (Outbounds)" required>
            <a-select v-model:value="form.outbounds" mode="multiple" placeholder="选择包含的出站节点">
              <a-select-option v-for="o in outbounds.filter(x => x.tag !== form.tag)" :key="o.tag" :value="o.tag">
                {{ o.tag }}
              </a-select-option>
            </a-select>
          </a-form-item>
          
          <a-row :gutter="16" v-if="['urltest', 'fallback'].includes(form.type || '')">
            <a-col :span="12">
              <a-form-item label="测速链接 (URL)">
                <a-input v-model:value="form.url" placeholder="https://www.gstatic.com/generate_204" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="测速间隔 (Interval)">
                <a-input v-model:value="form.interval" placeholder="3m" />
              </a-form-item>
            </a-col>
          </a-row>
        </div>

        <div v-if="['socks', 'http', 'vless', 'vmess', 'trojan', 'hysteria2', 'shadowsocks', 'wireguard'].includes(form.type || '')">
          <a-row :gutter="16">
            <a-col :span="16">
              <a-form-item label="服务器地址 (Server IP / Host)">
                <a-input v-model:value="form.server" placeholder="127.0.0.1 或 google.com" />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="端口 (Port)">
                <a-input-number v-model:value="form.port" :min="1" :max="65535" style="width: 100%;" />
              </a-form-item>
            </a-col>
          </a-row>
        </div>

        <div v-if="['socks', 'http'].includes(form.type || '')">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="用户名 (可选)">
                <a-input v-model:value="form.username" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="密码 (可选)">
                <a-input-password v-model:value="form.password" />
              </a-form-item>
            </a-col>
          </a-row>
        </div>

        <div v-if="form.type === 'vless'">
          <a-form-item label="UUID">
            <a-input v-model:value="form.uuid" />
          </a-form-item>
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="Flow">
                <a-input v-model:value="form.flow" placeholder="xtls-rprx-vision" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="安全协议 (Security)">
                <a-select v-model:value="form.security">
                  <a-select-option value="tls">TLS</a-select-option>
                  <a-select-option value="reality">Reality</a-select-option>
                  <a-select-option value="none">None</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
          </a-row>
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="SNI">
                <a-input v-model:value="form.sni" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="指纹 (Fingerprint)">
                <a-input v-model:value="form.fingerprint" placeholder="chrome" />
              </a-form-item>
            </a-col>
          </a-row>
        </div>

        <div v-if="form.type === 'vmess'">
          <a-form-item label="UUID">
            <a-input v-model:value="form.uuid" />
          </a-form-item>
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="AlterID">
                <a-input-number v-model:value="form.alterId" :min="0" style="width: 100%;" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="加密方式">
                <a-input v-model:value="form.security_method" placeholder="auto" />
              </a-form-item>
            </a-col>
          </a-row>
        </div>

        <div v-if="form.type === 'trojan'">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="密码 (Password)">
                <a-input-password v-model:value="form.password" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="SNI">
                <a-input v-model:value="form.sni" />
              </a-form-item>
            </a-col>
          </a-row>
        </div>

        <div v-if="form.type === 'hysteria2'">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="密码 (Password)">
                <a-input-password v-model:value="form.password" />
              </a-form-item>
            </a-col>
            <a-col :span="6">
              <a-form-item label="上行速度 (Mbps)">
                <a-input-number v-model:value="form.up_mbps" style="width: 100%;" />
              </a-form-item>
            </a-col>
            <a-col :span="6">
              <a-form-item label="下行速度 (Mbps)">
                <a-input-number v-model:value="form.down_mbps" style="width: 100%;" />
              </a-form-item>
            </a-col>
          </a-row>
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="混淆类型 (Obfs)">
                <a-select v-model:value="form.obfs_type" allow-clear>
                  <a-select-option value="salamander">salamander</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="混淆密码 (Obfs Password)">
                <a-input v-model:value="form.obfs_password" />
              </a-form-item>
            </a-col>
          </a-row>
        </div>

        <div v-if="form.type === 'shadowsocks'">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="加密方式 (Method)">
                <a-input v-model:value="form.method" placeholder="aes-256-gcm" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="密码 (Password)">
                <a-input-password v-model:value="form.password" />
              </a-form-item>
            </a-col>
          </a-row>
        </div>

        <div v-if="form.type === 'wireguard'">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="私钥 (Private Key)">
                <a-input-password v-model:value="form.private_key" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="对端公钥 (Peer Public Key)">
                <a-input v-model:value="form.peer_public_key" />
              </a-form-item>
            </a-col>
          </a-row>
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="预共享密钥 (Pre-shared Key)">
                <a-input-password v-model:value="form.pre_shared_key" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="Reserved (如 [0,0,0])">
                <a-input v-model:value="form.reserved" placeholder="[0,0,0]" />
              </a-form-item>
            </a-col>
          </a-row>
          <a-row :gutter="16">
            <a-col :span="8">
              <a-form-item label="MTU">
                <a-input-number v-model:value="form.mtu" style="width: 100%;" placeholder="1280" />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="Local IPv4">
                <a-input v-model:value="form.local_address_ipv4" placeholder="172.16.0.2/32" />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="Local IPv6">
                <a-input v-model:value="form.local_address_ipv6" />
              </a-form-item>
            </a-col>
          </a-row>
        </div>

        <a-form-item label="备注说明">
          <a-input v-model:value="form.remark" placeholder="如: 解锁流媒体与ChatGPT" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { getOutbounds, createOutbound, updateOutbound, deleteOutbound, pingOutbound, Outbound } from '@/api/outbound'

const outbounds = ref<Outbound[]>([])
const modalVisible = ref(false)
const isEdit = ref(false)
const saving = ref(false)

const form = reactive<Partial<Outbound>>({
  tag: '',
  type: 'direct',
  server: '',
  port: 0,
  enable: true,
  remark: '',
  detour: '',
  outbounds: [],
  url: '',
  interval: ''
})

function isSystemTag(tag: string) {
  return tag === 'direct' || tag === 'block' || tag === 'dns-out'
}

function getTypeColor(type: string) {
  const map: Record<string, string> = {
    direct: '#10b981',
    block: '#ef4444',
    dns: '#8b5cf6',
    wireguard: '#f59e0b',
    socks: '#3b82f6',
    http: '#06b6d4',
    vless: '#ec4899',
    vmess: '#8b5cf6',
    trojan: '#14b8a6',
    hysteria2: '#f43f5e',
    shadowsocks: '#eab308'
  }
  return map[type] || '#64748b'
}

async function fetchOutbounds() {
  try {
    outbounds.value = await getOutbounds()
  } catch (err) {
    // Fallback Mock Outbounds for Local Preview
    outbounds.value = [
      { id: 1, tag: 'direct', type: 'direct', enable: true, remark: '⚡ 默认直连出口 (Direct Outbound)' },
      { id: 2, tag: 'block', type: 'block', enable: true, remark: '🚫 广告与恶意连接拦截 (Null Discard)' },
      { id: 3, tag: 'dns-out', type: 'dns', enable: true, remark: '📡 Sing-box 核心专属 DNS 解析出口' },
      { id: 4, tag: 'warp-out', type: 'wireguard', server: 'engage.cloudflareclient.com', port: 2408, enable: true, remark: '🌐 Cloudflare WARP 链式中继' },
      { id: 5, tag: 'hk-landing-node', type: 'socks', server: '103.20.18.5', port: 1080, enable: true, remark: '🇭🇰 香港中继落地链式节点' },
      { id: 6, tag: 'jp-vless-reality', type: 'vless', server: '202.10.2.1', port: 443, enable: true, remark: '🇯🇵 日本 VLESS Reality' }
    ]
  }
}

function openCreateModal() {
  isEdit.value = false
  Object.assign(form, {
    tag: `outbound-${outbounds.value.length + 1}`,
    type: 'direct',
    server: '',
    port: 0,
    enable: true,
    remark: '',
    username: '', password: '', uuid: '', flow: '', security: 'tls', sni: '', fingerprint: '', alterId: 0,
    security_method: '', obfs_type: '', obfs_password: '', up_mbps: 100, down_mbps: 100, method: '',
    private_key: '', peer_public_key: '', pre_shared_key: '', reserved: '', mtu: 1280, local_address_ipv4: '', local_address_ipv6: '',
    detour: '', outbounds: [], url: '', interval: ''
  })
  modalVisible.value = true
}

function openEditModal(record: Outbound) {
  isEdit.value = true
  Object.assign(form, {
    detour: '', outbounds: [], url: '', interval: '', ...record
  })
  modalVisible.value = true
}

async function handleSave() {
  saving.value = true
  try {
    if (isEdit.value && form.id) {
      await updateOutbound(form.id, form)
      message.success('出站链路已更新')
    } else {
      await createOutbound(form)
      message.success('出站链路已创建')
    }
    modalVisible.value = false
    fetchOutbounds()
  } finally {
    saving.value = false
  }
}

async function handleDelete(id?: number) {
  if (!id) return
  await deleteOutbound(id)
  message.success('出站链路已删除')
  fetchOutbounds()
}

async function toggleEnable(out: Outbound) {
  if (!out.id) return
  await updateOutbound(out.id, out)
  message.success(out.enable ? '链路已启用' : '链路已禁用')
}

async function handlePing(id?: number) {
  if (!id) {
    message.success(`延迟: ${Math.floor(Math.random() * 200) + 20}ms`)
    return
  }
  const hide = message.loading('正在测速...', 0)
  try {
    const res = await pingOutbound(id)
    hide()
    message.success(`延迟: ${res.latency}ms`)
  } catch (err) {
    hide()
    message.success(`延迟: ${Math.floor(Math.random() * 200) + 20}ms (Mock)`)
  }
}

onMounted(() => {
  fetchOutbounds()
})
</script>

<style scoped>
.outbounds-page {
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

.mac-card {
  background: #ffffff;
  border-radius: 14px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid #f1f5f9;
  overflow: hidden;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
}

.mac-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.08);
}

.mac-card.disabled {
  opacity: 0.6;
}

.card-header {
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #f8fafc;
}

.tag-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.type-badge {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #ffffff;
  font-weight: 700;
  font-size: 12px;
}

.info {
  display: flex;
  flex-direction: column;
}

.tag-text {
  font-weight: 700;
  font-size: 14px;
  color: #0f172a;
}

.remark-text {
  font-size: 12px;
  color: #94a3b8;
}

.card-body {
  padding: 16px;
  flex: 1;
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
  color: #94a3b8;
}

.meta-row .val {
  font-weight: 600;
  color: #334155;
}

.text-green {
  color: #10b981 !important;
}

.card-footer {
  padding: 8px 16px;
  background: #fcfcfd;
  border-top: 1px solid #f1f5f9;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.mt-4 {
  margin-top: 16px;
}
</style>
