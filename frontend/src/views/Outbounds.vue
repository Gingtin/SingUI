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
                <a-select-option value="dns">DNS (专属解析)</a-select-option>
                <a-select-option value="socks">SOCKS5 链式中继</a-select-option>
                <a-select-option value="http">HTTP 链式中继</a-select-option>
                <a-select-option value="wireguard">WireGuard / WARP</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="出站标签 (Tag)" required>
              <a-input v-model:value="form.tag" placeholder="如: warp-out" :disabled="isEdit && isSystemTag(form.tag || '')" />
            </a-form-item>
          </a-col>
        </a-row>

        <div v-if="form.type === 'socks' || form.type === 'http'">
          <a-row :gutter="16">
            <a-col :span="16">
              <a-form-item label="服务器地址 (Server IP / Host)">
                <a-input v-model:value="form.server" placeholder="127.0.0.1" />
              </a-form-item>
            </a-col>
            <a-col :span="8">
              <a-form-item label="端口 (Port)">
                <a-input-number v-model:value="form.port" :min="1" :max="65535" style="width: 100%;" />
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
import { getOutbounds, createOutbound, updateOutbound, deleteOutbound, Outbound } from '@/api/outbound'

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
  }
  return map[type] || '#64748b'
}

async function fetchOutbounds() {
  try {
    outbounds.value = await getOutbounds()
  } catch (err) {
    console.error(err)
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
  })
  modalVisible.value = true
}

function openEditModal(record: Outbound) {
  isEdit.value = true
  Object.assign(form, record)
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
