<template>
  <div class="subs-page">
    <div class="page-header">
      <div>
        <h2>订阅中心</h2>
        <p class="subtitle">管理并分发各客户端订阅链接 (Sing-box / Clash Meta / Base64 / 自服务页)</p>
      </div>
    </div>

    <!-- Client Subscriptions Table -->
    <a-table
      :columns="columns"
      :data-source="clientSubs"
      :loading="loading"
      row-key="id"
      class="subs-table mt-4"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'email'">
          <b>{{ record.email }}</b>
          <div class="node-tag">节点: {{ record.inbound_tag }}</div>
        </template>

        <template v-if="column.key === 'traffic'">
          <span>{{ formatBytes(record.up + record.down) }} / {{ record.total > 0 ? formatBytes(record.total) : '无限制' }}</span>
        </template>

        <template v-if="column.key === 'expiry'">
          <span>{{ formatTime(record.expiry_time) }}</span>
        </template>

        <template v-if="column.key === 'sub_token'">
          <code class="token-badge">{{ record.sub_token }}</code>
        </template>

        <template v-if="column.key === 'actions'">
          <div class="action-links">
            <a-dropdown>
              <template #overlay>
                <a-menu @click="({ key }) => handleCopySub(record.sub_token, key as string)">
                  <a-menu-item key="sing-box">Sing-box JSON 订阅</a-menu-item>
                  <a-menu-item key="clash">Clash Meta (Mihomo) 订阅</a-menu-item>
                  <a-menu-item key="base64">Base64 通用订阅</a-menu-item>
                </a-menu>
              </template>
              <a-button size="small" type="primary">
                复制订阅链接 <DownOutlined />
              </a-button>
            </a-dropdown>

            <a-button size="small" @click="showSubQR(record)">
              二维码
            </a-button>
          </div>
        </template>
      </template>
    </a-table>

    <!-- QR Modal -->
    <a-modal
      v-model:open="qrVisible"
      :title="`订阅二维码 - [${currentRecord?.email}]`"
      :footer="null"
      width="400px"
    >
      <div class="qr-box">
        <qrcode-vue :value="qrSubUrl" :size="240" level="H" />
        <p class="qr-tip mt-3">支持 Sing-box、Clash Verge、Shadowrocket 扫码导入</p>
        <a-input v-model:value="qrSubUrl" readonly class="mt-2" />
        <a-button type="primary" block class="mt-2" @click="copyDirectUrl">复制链接</a-button>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { DownOutlined } from '@ant-design/icons-vue'
import QrcodeVue from 'qrcode.vue'
import { getInbounds, Inbound, Client } from '@/api/inbound'
import { formatBytes, formatTime } from '@/utils/format'

interface ClientSubItem extends Client {
  inbound_tag: string
}

const inbounds = ref<Inbound[]>([])
const loading = ref(false)
const qrVisible = ref(false)
const currentRecord = ref<ClientSubItem | null>(null)
const qrSubUrl = ref('')

const columns = [
  { title: '用户/节点', key: 'email' },
  { title: '流量使用', key: 'traffic' },
  { title: '到期时间', key: 'expiry' },
  { title: '安全 Token', key: 'sub_token' },
  { title: '操作', key: 'actions', width: 260 },
]

const clientSubs = computed(() => {
  const list: ClientSubItem[] = []
  for (const inb of inbounds.value) {
    if (inb.clients) {
      for (const c of inb.clients) {
        list.push({
          ...c,
          inbound_tag: inb.tag,
        })
      }
    }
  }
  return list
})

async function fetchSubs() {
  loading.value = true
  try {
    const data = await getInbounds()
    inbounds.value = data
  } catch (err) {
    inbounds.value = [
      {
        id: 1,
        tag: 'HK-VLESS-Reality-01',
        protocol: 'vless',
        port: 443,
        listen: '0.0.0.0',
        network: 'tcp',
        security: 'reality',
        settings: '{}',
        stream_settings: '{}',
        sniffing: '{}',
        enable: true,
        remark: '🇭🇰 香港 BGP 专线',
        clients: [
          { id: 101, email: 'alice', uuid: '6ba7b810-9dad-11d1-80b4-00c04fd430c8', up: 1024 * 1024 * 1024 * 5, down: 1024 * 1024 * 1024 * 28, total: 1024 * 1024 * 1024 * 100, expiry_time: Date.now() + 86400000 * 30, enable: true, sub_token: 'sub-alice-01' },
          { id: 102, email: 'bob', uuid: '7ca8b820-9dad-11d1-80b4-00c04fd430c9', up: 1024 * 1024 * 1024 * 1, down: 1024 * 1024 * 1024 * 8, total: 1024 * 1024 * 1024 * 50, expiry_time: Date.now() + 86400000 * 15, enable: true, sub_token: 'sub-bob-02' },
        ],
      },
      {
        id: 2,
        tag: 'US-Hysteria-2-Fast',
        protocol: 'hysteria2',
        port: 8443,
        listen: '0.0.0.0',
        network: 'udp',
        security: 'none',
        settings: '{}',
        stream_settings: '{}',
        sniffing: '{}',
        enable: true,
        remark: '🇺🇸 美国 洛杉矶 极速 UDP',
        clients: [
          { id: 103, email: 'carol', uuid: '8da9b830-9dad-11d1-80b4-00c04fd430ca', up: 1024 * 1024 * 1024 * 12, down: 1024 * 1024 * 1024 * 65, total: 0, expiry_time: 0, enable: true, sub_token: 'sub-carol-03' },
        ],
      },
    ]
  } finally {
    loading.value = false
  }
}

function getSubBaseUrl(token: string, flag = '') {
  const origin = window.location.origin
  let url = `${origin}/sub/${token}`
  if (flag) {
    url += `?flag=${flag}`
  }
  return url
}

function handleCopySub(token: string, type: string) {
  const url = getSubBaseUrl(token, type)
  navigator.clipboard.writeText(url)
  message.success(`${type} 订阅链接已复制`)
}

function showSubQR(record: ClientSubItem) {
  currentRecord.value = record
  qrSubUrl.value = getSubBaseUrl(record.sub_token || '', 'sing-box')
  qrVisible.value = true
}

function copyDirectUrl() {
  navigator.clipboard.writeText(qrSubUrl.value)
  message.success('订阅链接已复制')
}

onMounted(() => {
  fetchSubs()
})
</script>

<style scoped>
.subs-page {
  padding: 4px;
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

.subs-table {
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.node-tag {
  font-size: 12px;
  color: #94a3b8;
}

.token-badge {
  font-family: monospace;
  background: #f1f5f9;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
}

.action-links {
  display: flex;
  gap: 8px;
}

.qr-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 0;
}

.qr-tip {
  color: #64748b;
  font-size: 12px;
}

.mt-4 {
  margin-top: 16px;
}

.mt-3 {
  margin-top: 12px;
}

.mt-2 {
  margin-top: 8px;
}
</style>
