<template>
  <div class="inbounds-page">
    <!-- Header -->
    <div class="page-header">
      <div>
        <h2>入站节点管理</h2>
        <p class="subtitle">高品质 Sing-box 1.9+ 原生多协议节点与多用户配额调度</p>
      </div>
      <div class="header-actions">
        <a-radio-group v-model:value="viewMode" button-style="solid">
          <a-radio-button value="cards"><AppstoreOutlined /> 卡片视图</a-radio-button>
          <a-radio-button value="table"><BarsOutlined /> 列表视图</a-radio-button>
        </a-radio-group>
        <a-button @click="handleResetAllTraffic">重置全部流量</a-button>
        <a-button type="primary" @click="openCreateInboundModal">
          <template #prefix><PlusOutlined /></template>
          添加入站节点
        </a-button>
      </div>
    </div>

    <!-- Filter & Batch Toolbar -->
    <a-card class="toolbar-card mt-3" :bordered="false">
      <div class="toolbar-wrapper">
        <div class="search-filters">
          <a-input-search
            v-model:value="searchKeyword"
            placeholder="搜索节点名称、端口、备注..."
            style="width: 260px;"
            allow-clear
          />
          <a-select v-model:value="protocolFilter" style="width: 140px;" placeholder="全部协议">
            <a-select-option value="all">全部协议</a-select-option>
            <a-select-option value="vless">VLESS Reality</a-select-option>
            <a-select-option value="hysteria2">Hysteria 2</a-select-option>
            <a-select-option value="tuic">TUIC v5</a-select-option>
            <a-select-option value="shadowsocks">Shadowsocks</a-select-option>
            <a-select-option value="trojan">Trojan</a-select-option>
            <a-select-option value="vmess">VMess</a-select-option>
          </a-select>
        </div>

        <div v-if="selectedRowKeys.length > 0" class="batch-actions">
          <span class="selected-text">已选中 {{ selectedRowKeys.length }} 项</span>
          <a-button size="small" @click="handleBatchToggle(true)">批量启用</a-button>
          <a-button size="small" @click="handleBatchToggle(false)">批量禁用</a-button>
          <a-popconfirm title="确定批量删除选中的节点吗？" @confirm="handleBatchDelete">
            <a-button size="small" danger>批量删除</a-button>
          </a-popconfirm>
        </div>
      </div>
    </a-card>

    <!-- 1. Card View -->
    <div v-if="viewMode === 'cards'" class="cards-container mt-4">
      <a-row :gutter="[16, 16]">
        <a-col v-for="inbound in filteredInbounds" :key="inbound.id" :xs="24" :sm="12" :lg="8" :xl="6">
          <div class="node-card" :class="{ disabled: !inbound.enable }">
            <div class="node-card-header">
              <div class="node-title">
                <span class="protocol-icon" :style="{ backgroundColor: getProtocolBg(inbound.protocol) }">
                  {{ inbound.protocol.substring(0, 2).toUpperCase() }}
                </span>
                <div class="node-info">
                  <span class="node-tag">{{ inbound.tag }}</span>
                  <span class="node-remark">{{ inbound.remark || '无备注' }}</span>
                </div>
              </div>
              <a-switch :checked="inbound.enable" @change="() => toggleInboundEnable(inbound)" size="small" />
            </div>

            <div class="node-card-body">
              <div class="node-meta">
                <div class="meta-item">
                  <span class="label">端口</span>
                  <span class="val port">{{ inbound.port }}</span>
                </div>
                <div class="meta-item">
                  <span class="label">传输/安全</span>
                  <span class="val">{{ inbound.network }} / {{ inbound.security }}</span>
                </div>
                <div class="meta-item">
                  <span class="label">用户数</span>
                  <span class="val">{{ inbound.clients?.length || 0 }}</span>
                </div>
              </div>

              <div class="traffic-section">
                <div class="traffic-labels">
                  <span>总流量</span>
                  <span>{{ formatBytes(calcInboundTraffic(inbound).up + calcInboundTraffic(inbound).down) }}</span>
                </div>
                <div class="traffic-detail">
                  <span>⬆️ {{ formatBytes(calcInboundTraffic(inbound).up) }}</span>
                  <span>⬇️ {{ formatBytes(calcInboundTraffic(inbound).down) }}</span>
                </div>
              </div>
            </div>

            <div class="node-card-footer">
              <a-button type="link" size="small" @click="openClientDrawer(inbound)">
                <UserOutlined /> 用户 ({{ inbound.clients?.length || 0 }})
              </a-button>
              <a-button type="link" size="small" @click="openEditInboundModal(inbound)">
                <EditOutlined /> 编辑
              </a-button>
              <a-popconfirm title="确定删除该节点？" @confirm="() => handleDeleteInbound(inbound.id)">
                <a-button type="link" size="small" danger>
                  <DeleteOutlined />
                </a-button>
              </a-popconfirm>
            </div>
          </div>
        </a-col>
      </a-row>
    </div>

    <!-- 2. Table View -->
    <a-table
      v-else
      :columns="columns"
      :data-source="filteredInbounds"
      :loading="loading"
      row-key="id"
      :row-selection="{ selectedRowKeys: selectedRowKeys, onChange: onSelectChange }"
      :pagination="false"
      class="inbound-table mt-4"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'tag'">
          <div class="tag-cell">
            <span class="tag-title">{{ record.tag }}</span>
            <span class="tag-remark" v-if="record.remark">{{ record.remark }}</span>
          </div>
        </template>

        <template v-if="column.key === 'protocol'">
          <a-tag :color="getProtocolColor(record.protocol)" class="protocol-tag">
            {{ record.protocol.toUpperCase() }}
          </a-tag>
        </template>

        <template v-if="column.key === 'port'">
          <span class="port-badge">{{ record.port }}</span>
        </template>

        <template v-if="column.key === 'transport'">
          <span class="transport-badge">{{ record.network }} / {{ record.security }}</span>
        </template>

        <template v-if="column.key === 'traffic'">
          <div class="traffic-cell">
            <span>⬆️ {{ formatBytes(calcInboundTraffic(record).up) }}</span>
            <span>⬇️ {{ formatBytes(calcInboundTraffic(record).down) }}</span>
          </div>
        </template>

        <template v-if="column.key === 'clients'">
          <a-badge :count="record.clients?.length || 0" :number-style="{ backgroundColor: '#3b82f6' }" />
        </template>

        <template v-if="column.key === 'enable'">
          <a-switch v-model:checked="record.enable" @change="toggleInboundEnable(record)" />
        </template>

        <template v-if="column.key === 'actions'">
          <div class="action-buttons">
            <a-button size="small" type="link" @click="openClientDrawer(record)">
              用户 ({{ record.clients?.length || 0 }})
            </a-button>
            <a-button size="small" type="link" @click="openEditInboundModal(record)">
              编辑
            </a-button>
            <a-popconfirm title="确定删除该入站节点吗？" @confirm="handleDeleteInbound(record.id)">
              <a-button size="small" type="link" danger>删除</a-button>
            </a-popconfirm>
          </div>
        </template>
      </template>
    </a-table>

    <!-- Inbound Create/Edit Wizard Modal -->
    <a-modal
      v-model:open="inboundModalVisible"
      :title="isEditInbound ? '编辑入站节点' : '添加入站节点'"
      width="720px"
      @ok="handleSaveInbound"
      :confirmLoading="modalLoading"
    >
      <a-form :model="inboundForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="协议类型" required>
              <a-select v-model:value="inboundForm.protocol" @change="onProtocolChange" :disabled="isEditInbound">
                <a-select-option value="vless">VLESS (支持 Reality / Vision)</a-select-option>
                <a-select-option value="hysteria2">Hysteria 2 (极速 UDP / Salamander)</a-select-option>
                <a-select-option value="tuic">TUIC v5 (QUIC / BBR)</a-select-option>
                <a-select-option value="shadowsocks">Shadowsocks 2022</a-select-option>
                <a-select-option value="trojan">Trojan</a-select-option>
                <a-select-option value="vmess">VMess</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="监听端口" required>
              <a-input-number v-model:value="inboundForm.port" :min="1" :max="65535" style="width: 100%;" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="节点标签 (Tag)">
              <a-input v-model:value="inboundForm.tag" placeholder="如: vless-reality-01" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="备注说明">
              <a-input v-model:value="inboundForm.remark" placeholder="如: 香港优化节点" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="传输协议 (Network)">
              <a-select v-model:value="inboundForm.network">
                <a-select-option value="tcp">TCP</a-select-option>
                <a-select-option value="udp">UDP</a-select-option>
                <a-select-option value="ws">WebSocket</a-select-option>
                <a-select-option value="grpc">gRPC</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="安全协议 (Security)">
              <a-select v-model:value="inboundForm.security">
                <a-select-option value="none">None</a-select-option>
                <a-select-option value="tls">TLS</a-select-option>
                <a-select-option value="reality">Reality (仅 VLESS)</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <!-- Reality Settings Section -->
        <div v-if="inboundForm.security === 'reality'" class="config-section">
          <div class="section-title">
            <span>Reality 伪装配置</span>
            <a-button size="small" type="link" @click="generateRealityKeys">一键生成密钥对</a-button>
          </div>
          <a-form-item label="目标伪装域名 (SNI)">
            <a-input v-model:value="realityForm.server_name" placeholder="www.apple.com / addons.mozilla.org / www.cloudflare.com" />
          </a-form-item>
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="私钥 (Private Key)">
                <a-input v-model:value="realityForm.private_key" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="公钥 (Public Key)">
                <a-input v-model:value="realityForm.public_key" />
              </a-form-item>
            </a-col>
          </a-row>
          <a-form-item label="Short ID (8位 Hex)">
            <a-input v-model:value="realityForm.short_id" />
          </a-form-item>
        </div>

        <!-- Hysteria 2 Settings -->
        <div v-if="inboundForm.protocol === 'hysteria2'" class="config-section">
          <div class="section-title">Hysteria 2 速率与混淆</div>
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="上行速率上限 (Mbps)">
                <a-input-number v-model:value="hyForm.up_mbps" :min="0" placeholder="0 为不限制" style="width: 100%;" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="下行速率上限 (Mbps)">
                <a-input-number v-model:value="hyForm.down_mbps" :min="0" placeholder="0 为不限制" style="width: 100%;" />
              </a-form-item>
            </a-col>
          </a-row>
          <a-form-item label="Salamander 混淆密码 (可选)">
            <a-input v-model:value="hyForm.obfs_password" placeholder="留空表示不开启混淆" />
          </a-form-item>
        </div>

        <!-- Shadowsocks 2022 Settings -->
        <div v-if="inboundForm.protocol === 'shadowsocks'" class="config-section">
          <div class="section-title">Shadowsocks 加密配置</div>
          <a-form-item label="加密方法 (Cipher)">
            <a-select v-model:value="ssForm.method">
              <a-select-option value="2022-blake3-aes-128-gcm">2022-blake3-aes-128-gcm (推荐)</a-select-option>
              <a-select-option value="2022-blake3-aes-256-gcm">2022-blake3-aes-256-gcm</a-select-option>
              <a-select-option value="aes-128-gcm">aes-128-gcm</a-select-option>
              <a-select-option value="chacha20-ietf-poly1305">chacha20-ietf-poly1305</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="服务器主密码">
            <a-input v-model:value="ssForm.password" placeholder="主密码" />
          </a-form-item>
        </div>
      </a-form>
    </a-modal>

    <!-- Client Drawer -->
    <a-drawer
      v-model:open="clientDrawerVisible"
      :title="`用户管理 - [${currentInbound?.tag}]`"
      width="820px"
    >
      <div class="drawer-header-actions mb-4">
        <a-button type="primary" size="small" @click="openAddClientModal">
          <template #prefix><PlusOutlined /></template>
          添加用户
        </a-button>
        <a-button v-if="selectedClientKeys.length > 0" danger size="small" @click="handleBatchDeleteClients">
          批量删除 ({{ selectedClientKeys.length }})
        </a-button>
      </div>

      <a-table
        :columns="clientColumns"
        :data-source="currentInbound?.clients || []"
        row-key="id"
        :row-selection="{ selectedRowKeys: selectedClientKeys, onChange: onClientSelectChange }"
        :pagination="false"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'email'">
            <b>{{ record.email }}</b>
          </template>

          <template v-if="column.key === 'traffic'">
            <span>{{ formatBytes(record.up + record.down) }} / {{ record.total > 0 ? formatBytes(record.total) : '无限制' }}</span>
          </template>

          <template v-if="column.key === 'expiry'">
            <span>{{ formatTime(record.expiry_time) }}</span>
          </template>

          <template v-if="column.key === 'enable'">
            <a-switch v-model:checked="record.enable" @change="toggleClientEnable(record)" />
          </template>

          <template v-if="column.key === 'actions'">
            <div class="action-buttons">
              <a-button size="small" type="link" @click="showClientQR(record)">二维码</a-button>
              <a-button size="small" type="link" @click="openEditClientModal(record)">编辑</a-button>
              <a-button size="small" type="link" @click="handleResetClientTraffic(record.id)">清零</a-button>
              <a-popconfirm title="确定删除该用户？" @confirm="handleDeleteClient(record.id)">
                <a-button size="small" type="link" danger>删除</a-button>
              </a-popconfirm>
            </div>
          </template>
        </template>
      </a-table>
    </a-drawer>

    <!-- Client Modal (Add/Edit) -->
    <a-modal
      v-model:open="clientModalVisible"
      :title="isEditClient ? '编辑用户' : '添加用户'"
      @ok="handleSaveClient"
      :confirmLoading="modalLoading"
    >
      <a-form :model="clientForm" layout="vertical">
        <a-form-item label="用户备注 / Email" required>
          <a-input v-model:value="clientForm.email" placeholder="如: user1" />
        </a-form-item>

        <a-form-item label="UUID / 认证密钥">
          <a-input v-model:value="clientForm.uuid">
            <template #addonAfter>
              <span style="cursor: pointer;" @click="generateClientUUID">自动生成</span>
            </template>
          </a-input>
        </a-form-item>

        <a-form-item v-if="currentInbound?.protocol === 'vless' && currentInbound?.security === 'reality'" label="流控 (Flow)">
          <a-select v-model:value="clientForm.flow">
            <a-select-option value="xtls-rprx-vision">xtls-rprx-vision (推荐)</a-select-option>
            <a-select-option value="">无</a-select-option>
          </a-select>
        </a-form-item>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="流量限额 (GB)">
              <a-input-number v-model:value="clientQuotaGB" :min="0" placeholder="0 为不限制" style="width: 100%;" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="到期时间">
              <a-date-picker v-model:value="clientExpiryDayjs" show-time placeholder="留空为永久有效" style="width: 100%;" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <!-- QR Code / Link Modal -->
    <a-modal
      v-model:open="qrModalVisible"
      title="节点分享与二维码"
      :footer="null"
      width="400px"
    >
      <div class="qr-container">
        <qrcode-vue :value="qrValue" :size="240" level="H" />
        <div class="link-box mt-4">
          <a-textarea v-model:value="qrValue" :rows="3" readonly />
          <a-button type="primary" block class="mt-2" @click="copyLink">复制节点链接</a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import {
  PlusOutlined,
  AppstoreOutlined,
  BarsOutlined,
  UserOutlined,
  EditOutlined,
  DeleteOutlined,
} from '@ant-design/icons-vue'
import QrcodeVue from 'qrcode.vue'
import dayjs, { Dayjs } from 'dayjs'
import {
  getInbounds,
  createInbound,
  updateInbound,
  deleteInbound,
  batchDeleteInbounds,
  batchToggleInbounds,
  addClient,
  updateClient,
  deleteClient,
  batchDeleteClients,
  resetClientTraffic,
  resetAllTraffic,
  getRealityKeypair,
  getRandomUUID,
  Inbound,
  Client,
} from '@/api/inbound'
import { formatBytes, formatTime } from '@/utils/format'

const viewMode = ref<'cards' | 'table'>('cards')
const inbounds = ref<Inbound[]>([])
const loading = ref(false)
const modalLoading = ref(false)

// Filters & Selection
const searchKeyword = ref('')
const protocolFilter = ref('all')
const selectedRowKeys = ref<number[]>([])
const selectedClientKeys = ref<number[]>([])

const filteredInbounds = computed(() => {
  return inbounds.value.filter((inb) => {
    const matchProto = protocolFilter.value === 'all' || inb.protocol === protocolFilter.value
    const kw = searchKeyword.value.toLowerCase().trim()
    const matchKeyword =
      !kw ||
      inb.tag.toLowerCase().includes(kw) ||
      (inb.remark && inb.remark.toLowerCase().includes(kw)) ||
      String(inb.port).includes(kw)
    return matchProto && matchKeyword
  })
})

const columns = [
  { title: '节点名称/备注', key: 'tag' },
  { title: '协议', key: 'protocol', width: 120 },
  { title: '端口', key: 'port', width: 90 },
  { title: '传输 / 安全', key: 'transport', width: 160 },
  { title: '已用流量', key: 'traffic', width: 180 },
  { title: '用户数', key: 'clients', width: 90 },
  { title: '状态', key: 'enable', width: 80 },
  { title: '操作', key: 'actions', width: 220 },
]

const clientColumns = [
  { title: '用户', key: 'email' },
  { title: '已用 / 总限额', key: 'traffic' },
  { title: '到期时间', key: 'expiry' },
  { title: '启用', key: 'enable', width: 80 },
  { title: '操作', key: 'actions', width: 220 },
]

// Inbound Modals State
const inboundModalVisible = ref(false)
const isEditInbound = ref(false)
const inboundForm = reactive<Partial<Inbound>>({
  tag: '',
  protocol: 'vless',
  port: 443,
  listen: '0.0.0.0',
  network: 'tcp',
  security: 'reality',
  settings: '{}',
  stream_settings: '{}',
  enable: true,
  remark: '',
})

const realityForm = reactive({
  server_name: 'www.apple.com',
  private_key: '',
  public_key: '',
  short_id: '',
})

const hyForm = reactive({
  up_mbps: 0,
  down_mbps: 0,
  obfs_password: '',
})

const ssForm = reactive({
  method: '2022-blake3-aes-128-gcm',
  password: '',
})

// Client Drawer & Modal State
const clientDrawerVisible = ref(false)
const currentInbound = ref<Inbound | null>(null)
const clientModalVisible = ref(false)
const isEditClient = ref(false)
const clientForm = reactive<Partial<Client>>({
  email: '',
  uuid: '',
  flow: 'xtls-rprx-vision',
  total: 0,
  expiry_time: 0,
  enable: true,
})
const clientQuotaGB = ref<number>(0)
const clientExpiryDayjs = ref<Dayjs | null>(null)

// QR Modal
const qrModalVisible = ref(false)
const qrValue = ref('')

async function fetchInbounds() {
  loading.value = true
  try {
    const data = await getInbounds()
    inbounds.value = data
    if (currentInbound.value) {
      currentInbound.value = data.find((i) => i.id === currentInbound.value?.id) || null
    }
  } finally {
    loading.value = false
  }
}

function getProtocolBg(protocol: string) {
  const map: Record<string, string> = {
    vless: '#3b82f6',
    hysteria2: '#8b5cf6',
    tuic: '#06b6d4',
    shadowsocks: '#f59e0b',
    trojan: '#10b981',
    vmess: '#ef4444',
  }
  return map[protocol] || '#64748b'
}

function getProtocolColor(protocol: string) {
  const map: Record<string, string> = {
    vless: 'blue',
    hysteria2: 'purple',
    tuic: 'cyan',
    shadowsocks: 'orange',
    trojan: 'green',
    vmess: 'volcano',
  }
  return map[protocol] || 'default'
}

function calcInboundTraffic(inbound: Inbound) {
  let up = 0
  let down = 0
  if (inbound.clients) {
    for (const c of inbound.clients) {
      up += c.up || 0
      down += c.down || 0
    }
  }
  return { up, down }
}

function onSelectChange(keys: any[]) {
  selectedRowKeys.value = keys
}

function onClientSelectChange(keys: any[]) {
  selectedClientKeys.value = keys
}

async function handleBatchToggle(enable: boolean) {
  if (selectedRowKeys.value.length === 0) return
  await batchToggleInbounds(selectedRowKeys.value, enable)
  message.success(enable ? '已批量启用节点' : '已批量禁用节点')
  selectedRowKeys.value = []
  fetchInbounds()
}

async function handleBatchDelete() {
  if (selectedRowKeys.value.length === 0) return
  await batchDeleteInbounds(selectedRowKeys.value)
  message.success('已批量删除节点')
  selectedRowKeys.value = []
  fetchInbounds()
}

async function handleBatchDeleteClients() {
  if (selectedClientKeys.value.length === 0) return
  await batchDeleteClients(selectedClientKeys.value)
  message.success('已批量删除用户')
  selectedClientKeys.value = []
  fetchInbounds()
}

async function onProtocolChange(val: string) {
  if (val === 'vless') {
    inboundForm.security = 'reality'
    inboundForm.network = 'tcp'
    await generateRealityKeys()
  } else if (val === 'hysteria2') {
    inboundForm.security = 'none'
    inboundForm.network = 'udp'
  } else if (val === 'tuic') {
    inboundForm.security = 'none'
    inboundForm.network = 'udp'
  } else if (val === 'shadowsocks') {
    inboundForm.security = 'none'
    inboundForm.network = 'tcp'
    ssForm.password = Math.random().toString(36).substring(2, 18)
  }
}

async function generateRealityKeys() {
  try {
    const res = await getRealityKeypair()
    realityForm.private_key = res.private_key
    realityForm.public_key = res.public_key
    realityForm.short_id = res.short_id
  } catch (err) {
    console.error(err)
  }
}

function openCreateInboundModal() {
  isEditInbound.value = false
  Object.assign(inboundForm, {
    tag: `inbound-${Math.floor(1000 + Math.random() * 9000)}`,
    protocol: 'vless',
    port: 443,
    listen: '0.0.0.0',
    network: 'tcp',
    security: 'reality',
    enable: true,
    remark: '',
  })
  generateRealityKeys()
  inboundModalVisible.value = true
}

function openEditInboundModal(record: Inbound) {
  isEditInbound.value = true
  Object.assign(inboundForm, record)

  try {
    const stream = JSON.parse(record.stream_settings || '{}')
    if (stream.reality) {
      realityForm.server_name = stream.reality.server_names?.[0] || 'www.apple.com'
      realityForm.private_key = stream.reality.private_key || ''
      realityForm.public_key = stream.reality.public_key || ''
      realityForm.short_id = stream.reality.short_ids?.[0] || ''
    }
  } catch {}

  try {
    const settings = JSON.parse(record.settings || '{}')
    if (record.protocol === 'hysteria2') {
      hyForm.up_mbps = settings.up_mbps || 0
      hyForm.down_mbps = settings.down_mbps || 0
      hyForm.obfs_password = settings.obfs_password || ''
    } else if (record.protocol === 'shadowsocks') {
      ssForm.method = settings.method || '2022-blake3-aes-128-gcm'
      ssForm.password = settings.password || ''
    }
  } catch {}

  inboundModalVisible.value = true
}

async function handleSaveInbound() {
  modalLoading.value = true
  try {
    const streamSettingsObj: any = {
      network: inboundForm.network,
      security: inboundForm.security,
    }
    if (inboundForm.security === 'reality') {
      streamSettingsObj.reality = {
        server_names: [realityForm.server_name],
        private_key: realityForm.private_key,
        public_key: realityForm.public_key,
        short_ids: [realityForm.short_id],
      }
    }

    let settingsObj: any = {}
    if (inboundForm.protocol === 'hysteria2') {
      settingsObj = {
        up_mbps: hyForm.up_mbps,
        down_mbps: hyForm.down_mbps,
        obfs_type: hyForm.obfs_password ? 'salamander' : '',
        obfs_password: hyForm.obfs_password,
      }
    } else if (inboundForm.protocol === 'shadowsocks') {
      settingsObj = {
        method: ssForm.method,
        password: ssForm.password,
      }
    }

    const payload = {
      ...inboundForm,
      stream_settings: JSON.stringify(streamSettingsObj),
      settings: JSON.stringify(settingsObj),
    }

    if (isEditInbound.value && inboundForm.id) {
      await updateInbound(inboundForm.id, payload)
      message.success('节点修改成功')
    } else {
      await createInbound(payload)
      message.success('节点创建成功')
    }

    inboundModalVisible.value = false
    fetchInbounds()
  } finally {
    modalLoading.value = false
  }
}

async function handleDeleteInbound(id?: number) {
  if (!id) return
  await deleteInbound(id)
  message.success('节点已删除')
  fetchInbounds()
}

async function toggleInboundEnable(record: Inbound) {
  if (!record.id) return
  await updateInbound(record.id, record)
  message.success(record.enable ? '节点已启用' : '节点已禁用')
}

// Client operations
function openClientDrawer(record: Inbound) {
  currentInbound.value = record
  selectedClientKeys.value = []
  clientDrawerVisible.value = true
}

async function generateClientUUID() {
  const res = await getRandomUUID()
  clientForm.uuid = res.uuid
}

async function openAddClientModal() {
  isEditClient.value = false
  const res = await getRandomUUID()
  Object.assign(clientForm, {
    email: `user-${Math.floor(100 + Math.random() * 900)}`,
    uuid: res.uuid,
    flow: currentInbound.value?.security === 'reality' ? 'xtls-rprx-vision' : '',
    enable: true,
  })
  clientQuotaGB.value = 0
  clientExpiryDayjs.value = null
  clientModalVisible.value = true
}

function openEditClientModal(client: Client) {
  isEditClient.value = true
  Object.assign(clientForm, client)
  clientQuotaGB.value = client.total ? client.total / (1024 * 1024 * 1024) : 0
  clientExpiryDayjs.value = client.expiry_time ? dayjs(client.expiry_time) : null
  clientModalVisible.value = true
}

async function handleSaveClient() {
  if (!currentInbound.value?.id) return
  modalLoading.value = true
  try {
    const payload = {
      ...clientForm,
      total: clientQuotaGB.value ? clientQuotaGB.value * 1024 * 1024 * 1024 : 0,
      expiry_time: clientExpiryDayjs.value ? clientExpiryDayjs.value.valueOf() : 0,
    }

    if (isEditClient.value && clientForm.id) {
      await updateClient(currentInbound.value.id, clientForm.id, payload)
      message.success('用户更新成功')
    } else {
      await addClient(currentInbound.value.id, payload)
      message.success('用户添加成功')
    }

    clientModalVisible.value = false
    fetchInbounds()
  } finally {
    modalLoading.value = false
  }
}

async function handleDeleteClient(clientId?: number) {
  if (!currentInbound.value?.id || !clientId) return
  await deleteClient(currentInbound.value.id, clientId)
  message.success('用户已删除')
  fetchInbounds()
}

async function handleResetClientTraffic(clientId?: number) {
  if (!currentInbound.value?.id || !clientId) return
  await resetClientTraffic(currentInbound.value.id, clientId)
  message.success('用户流量已清零')
  fetchInbounds()
}

async function toggleClientEnable(client: Client) {
  if (!currentInbound.value?.id || !client.id) return
  await updateClient(currentInbound.value.id, client.id, client)
  message.success(client.enable ? '用户已启用' : '用户已禁用')
}

async function handleResetAllTraffic() {
  await resetAllTraffic()
  message.success('全部流量已重置')
  fetchInbounds()
}

function showClientQR(client: Client) {
  const host = window.location.hostname
  const port = currentInbound.value?.port || 443
  const tag = currentInbound.value?.tag || ''

  if (currentInbound.value?.protocol === 'vless') {
    let realityOpts = ''
    try {
      const stream = JSON.parse(currentInbound.value.stream_settings || '{}')
      if (stream.reality) {
        const sni = stream.reality.server_names?.[0] || 'www.apple.com'
        const pbk = stream.reality.public_key || ''
        const sid = stream.reality.short_ids?.[0] || ''
        realityOpts = `&security=reality&sni=${sni}&pbk=${pbk}&sid=${sid}&flow=${client.flow || 'xtls-rprx-vision'}`
      }
    } catch {}
    qrValue.value = `vless://${client.uuid}@${host}:${port}?type=tcp${realityOpts}#${encodeURIComponent(tag + '-' + client.email)}`
  } else if (currentInbound.value?.protocol === 'hysteria2') {
    qrValue.value = `hysteria2://${client.password || client.uuid}@${host}:${port}/?sni=${host}#${encodeURIComponent(tag + '-' + client.email)}`
  } else {
    qrValue.value = `Sub Token: ${client.sub_token}`
  }

  qrModalVisible.value = true
}

function copyLink() {
  navigator.clipboard.writeText(qrValue.value)
  message.success('链接已复制到剪贴板')
}

onMounted(() => {
  fetchInbounds()
})
</script>

<style scoped>
.inbounds-page {
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
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.toolbar-card {
  border-radius: 12px;
  background: #ffffff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.toolbar-wrapper {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.search-filters {
  display: flex;
  gap: 12px;
}

.batch-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.selected-text {
  font-size: 13px;
  font-weight: 600;
  color: #3b82f6;
  margin-right: 4px;
}

/* Card View Styling */
.node-card {
  background: #ffffff;
  border-radius: 14px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid #f1f5f9;
  overflow: hidden;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
}

.node-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.08);
}

.node-card.disabled {
  opacity: 0.6;
}

.node-card-header {
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  border-bottom: 1px solid #f8fafc;
}

.node-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.protocol-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #ffffff;
  font-weight: 700;
  font-size: 13px;
}

.node-info {
  display: flex;
  flex-direction: column;
}

.node-tag {
  font-weight: 700;
  font-size: 14px;
  color: #0f172a;
}

.node-remark {
  font-size: 12px;
  color: #94a3b8;
}

.node-card-body {
  padding: 16px;
  flex: 1;
}

.node-meta {
  display: flex;
  justify-content: space-between;
  background: #f8fafc;
  padding: 10px 12px;
  border-radius: 8px;
  margin-bottom: 12px;
}

.meta-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.meta-item .label {
  font-size: 11px;
  color: #94a3b8;
}

.meta-item .val {
  font-size: 13px;
  font-weight: 600;
  color: #334155;
}

.meta-item .val.port {
  font-family: monospace;
  color: #3b82f6;
}

.traffic-section {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.traffic-labels {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  font-weight: 600;
  color: #475569;
}

.traffic-detail {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: #94a3b8;
}

.node-card-footer {
  padding: 8px 12px;
  background: #fcfcfd;
  border-top: 1px solid #f1f5f9;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* Table View */
.inbound-table {
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.tag-cell {
  display: flex;
  flex-direction: column;
}

.tag-title {
  font-weight: 600;
  color: #0f172a;
}

.tag-remark {
  font-size: 12px;
  color: #94a3b8;
}

.port-badge {
  font-family: monospace;
  font-weight: 600;
  color: #3b82f6;
  background: #eff6ff;
  padding: 2px 8px;
  border-radius: 4px;
}

.transport-badge {
  font-size: 12px;
  color: #475569;
  background: #f8fafc;
  padding: 2px 6px;
  border-radius: 4px;
}

.traffic-cell {
  display: flex;
  flex-direction: column;
  font-size: 12px;
  color: #64748b;
}

.config-section {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 16px;
}

.section-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  color: #334155;
  margin-bottom: 12px;
}

.qr-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 0;
}

.link-box {
  width: 100%;
}

.action-buttons {
  display: flex;
  gap: 4px;
}

.drawer-header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.mt-4 {
  margin-top: 16px;
}

.mt-3 {
  margin-top: 12px;
}

.mb-4 {
  margin-bottom: 16px;
}

.mt-2 {
  margin-top: 8px;
}
</style>
