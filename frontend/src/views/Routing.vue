<template>
  <div class="routing-page">
    <div class="page-header">
      <div>
        <h2>分流路由与 DNS 配置</h2>
        <p class="subtitle">基于 Sing-box 1.9+ 架构管理智能分流规则集 (Rule Sets) 与 DNS 策略</p>
      </div>
      <div class="header-actions">
        <a-button v-if="activeTab === 'rules'" type="primary" @click="openCreateRuleModal">
          <template #prefix><PlusOutlined /></template>
          添加路由规则
        </a-button>
        <a-button v-if="activeTab === 'dns'" type="primary" :loading="dnsSaving" @click="handleSaveDNS">
          保存 DNS 配置
        </a-button>
      </div>
    </div>

    <a-card class="mt-4 main-tabs-card" :bordered="false">
      <a-tabs v-model:activeKey="activeTab">
        <!-- Tab 1: Routing Rules -->
        <a-tab-pane key="rules" tab="分流规则 (Routing Rules)">
          <div class="preset-bar mb-3">
            <span class="preset-label">常用规则预设:</span>
            <a-button size="small" @click="addPresetRule('ads')">🚫 广告拦截</a-button>
            <a-button size="small" @click="addPresetRule('cn_domain')">🇨🇳 国内域名直连</a-button>
            <a-button size="small" @click="addPresetRule('cn_ip')">🇨🇳 国内 IP 直连</a-button>
            <a-button size="small" @click="addPresetRule('private_ip')">🏠 私有局域网直连</a-button>
          </div>

          <a-table
            :columns="ruleColumns"
            :data-source="rules"
            :loading="rulesLoading"
            row-key="id"
            :pagination="false"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'order'">
                <a-tag color="geekblue">#{{ record.order }}</a-tag>
              </template>

              <template v-if="column.key === 'tag'">
                <b>{{ record.tag }}</b>
                <div class="text-muted" v-if="record.remark">{{ record.remark }}</div>
              </template>

              <template v-if="column.key === 'match'">
                <div class="match-cell">
                  <span v-if="record.protocol"><a-tag color="purple">协议: {{ record.protocol }}</a-tag></span>
                  <span v-if="record.domain"><a-tag color="cyan">域名: {{ record.domain }}</a-tag></span>
                  <span v-if="record.ip"><a-tag color="blue">IP: {{ record.ip }}</a-tag></span>
                </div>
              </template>

              <template v-if="column.key === 'outbound'">
                <a-tag :color="getOutboundColor(record.outbound)">
                  {{ record.outbound.toUpperCase() }}
                </a-tag>
              </template>

              <template v-if="column.key === 'enable'">
                <a-switch v-model:checked="record.enable" @change="toggleRuleEnable(record)" />
              </template>

              <template v-if="column.key === 'actions'">
                <div class="action-buttons">
                  <a-button size="small" type="link" @click="openEditRuleModal(record)">编辑</a-button>
                  <a-popconfirm title="确定删除该分流规则？" @confirm="handleDeleteRule(record.id)">
                    <a-button size="small" type="link" danger>删除</a-button>
                  </a-popconfirm>
                </div>
              </template>
            </template>
          </a-table>
        </a-tab-pane>

        <!-- Tab 2: DNS Settings -->
        <a-tab-pane key="dns" tab="DNS 解析策略 (DNS Engine)">
          <div class="dns-form-container">
            <a-form layout="vertical" :model="dnsForm">
              <a-row :gutter="24">
                <a-col :xs="24" :md="12">
                  <a-form-item label="国外 / 远程 DNS (Remote DNS)" extra="用于代理外网域名的安全解析 (支持 DoH / DoT)">
                    <a-input v-model:value="dnsForm.remote_dns" placeholder="https://1.1.1.1/dns-query" />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :md="12">
                  <a-form-item label="国内直连 DNS (China DNS)" extra="用于国内域名快速直连解析">
                    <a-input v-model:value="dnsForm.china_dns" placeholder="https://223.5.5.5/dns-query" />
                  </a-form-item>
                </a-col>
              </a-row>

              <a-row :gutter="24">
                <a-col :xs="24" :md="12">
                  <a-form-item label="本地基础 DNS (Local DNS)">
                    <a-input v-model:value="dnsForm.local_dns" placeholder="local 或 223.5.5.5" />
                  </a-form-item>
                </a-col>
                <a-col :xs="24" :md="12">
                  <a-form-item label="DNS 解析策略 (Strategy)">
                    <a-select v-model:value="dnsForm.strategy">
                      <a-select-option value="prefer_ipv4">优先 IPv4 (prefer_ipv4)</a-select-option>
                      <a-select-option value="prefer_ipv6">优先 IPv6 (prefer_ipv6)</a-select-option>
                      <a-select-option value="ipv4_only">仅 IPv4 (ipv4_only)</a-select-option>
                      <a-select-option value="ipv6_only">仅 IPv6 (ipv6_only)</a-select-option>
                    </a-select>
                  </a-form-item>
                </a-col>
              </a-row>
            </a-form>
          </div>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- Rule Create/Edit Modal -->
    <a-modal
      v-model:open="ruleModalVisible"
      :title="isEditRule ? '编辑分流规则' : '添加分流规则'"
      @ok="handleSaveRule"
      :confirmLoading="ruleModalLoading"
      width="600px"
    >
      <a-form :model="ruleForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="14">
            <a-form-item label="规则标签 (Tag)" required>
              <a-input v-model:value="ruleForm.tag" placeholder="如: Direct China" />
            </a-form-item>
          </a-col>
          <a-col :span="10">
            <a-form-item label="执行动作 (Outbound)" required>
              <a-select v-model:value="ruleForm.outbound">
                <a-select-option value="direct">DIRECT (直连)</a-select-option>
                <a-select-option value="block">BLOCK (拦截)</a-select-option>
                <a-select-option value="dns-out">DNS-OUT (DNS路由)</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="匹配域名 / Geosite (多项以逗号分隔)" extra="例如: geosite:cn, geosite:apple, google.com">
          <a-input v-model:value="ruleDomainInput" placeholder="geosite:cn" />
        </a-form-item>

        <a-form-item label="匹配 IP / GeoIP (多项以逗号分隔)" extra="例如: geoip:cn, geoip:private, 192.168.0.0/16">
          <a-input v-model:value="ruleIPInput" placeholder="geoip:cn" />
        </a-form-item>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="匹配协议">
              <a-input v-model:value="ruleForm.protocol" placeholder="如: dns, http, tls" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="规则优先级 (序号越小越优先)">
              <a-input-number v-model:value="ruleForm.order" :min="1" style="width: 100%;" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="备注说明">
          <a-input v-model:value="ruleForm.remark" placeholder="备注用途" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import {
  getRoutingRules,
  createRoutingRule,
  updateRoutingRule,
  deleteRoutingRule,
  getDNSSettings,
  updateDNSSettings,
  RoutingRule,
  DNSSettings,
} from '@/api/routing'

const activeTab = ref('rules')

// Rules
const rules = ref<RoutingRule[]>([])
const rulesLoading = ref(false)
const ruleModalVisible = ref(false)
const isEditRule = ref(false)
const ruleModalLoading = ref(false)

const ruleForm = reactive<Partial<RoutingRule>>({
  tag: '',
  outbound: 'direct',
  protocol: '',
  order: 1,
  remark: '',
  enable: true,
})
const ruleDomainInput = ref('')
const ruleIPInput = ref('')

const ruleColumns = [
  { title: '优先级', key: 'order', width: 90 },
  { title: '规则名称', key: 'tag', width: 200 },
  { title: '匹配条件 (Domain / IP / Protocol)', key: 'match' },
  { title: '路由动作', key: 'outbound', width: 140 },
  { title: '状态', key: 'enable', width: 90 },
  { title: '操作', key: 'actions', width: 160 },
]

// DNS
const dnsForm = reactive<DNSSettings>({
  local_dns: 'local',
  remote_dns: 'https://1.1.1.1/dns-query',
  china_dns: 'https://223.5.5.5/dns-query',
  enable_fakeip: false,
  strategy: 'prefer_ipv4',
})
const dnsSaving = ref(false)

function getOutboundColor(outbound: string) {
  if (outbound === 'direct') return 'green'
  if (outbound === 'block') return 'red'
  if (outbound === 'dns-out') return 'purple'
  return 'blue'
}

async function fetchRules() {
  rulesLoading.value = true
  try {
    rules.value = await getRoutingRules()
  } catch (err) {
    // Fallback Mock Rules for Local Preview
    rules.value = [
      { id: 1, tag: 'DNS Hijack', type: 'field', outbound: 'dns-out', protocol: 'dns', order: 1, enable: true, remark: '劫持 DNS 流量并进入内部分流引擎' },
      { id: 2, tag: 'AdBlock Ads', type: 'field', outbound: 'block', domain: '["geosite:category-ads-all"]', order: 2, enable: true, remark: '🚫 广告与恶意挖矿域名拦截' },
      { id: 3, tag: 'Private Network Direct', type: 'field', outbound: 'direct', ip: '["geoip:private"]', order: 3, enable: true, remark: '🏠 私有局域网 RFC1918 直连' },
      { id: 4, tag: 'China Domain Direct', type: 'field', outbound: 'direct', domain: '["geosite:cn"]', order: 4, enable: true, remark: '🇨🇳 国内主流站点与 CDN 极速直连' },
      { id: 5, tag: 'China IP Direct', type: 'field', outbound: 'direct', ip: '["geoip:cn"]', order: 5, enable: true, remark: '🇨🇳 国内 IP 网段直连' },
      { id: 6, tag: 'Global Proxy Fallback', type: 'field', outbound: 'direct', order: 99, enable: true, remark: '🌐 全球其他流量默认路由' },
    ]
  } finally {
    rulesLoading.value = false
  }
}

async function fetchDNS() {
  try {
    const data = await getDNSSettings()
    Object.assign(dnsForm, data)
  } catch (err) {
    console.error(err)
  }
}

function openCreateRuleModal() {
  isEditRule.value = false
  Object.assign(ruleForm, {
    tag: `rule-${rules.value.length + 1}`,
    outbound: 'direct',
    protocol: '',
    order: rules.value.length + 1,
    remark: '',
    enable: true,
  })
  ruleDomainInput.value = ''
  ruleIPInput.value = ''
  ruleModalVisible.value = true
}

function openEditRuleModal(rule: RoutingRule) {
  isEditRule.value = true
  Object.assign(ruleForm, rule)
  try {
    const d = JSON.parse(rule.domain || '[]')
    ruleDomainInput.value = Array.isArray(d) ? d.join(', ') : ''
  } catch {
    ruleDomainInput.value = ''
  }
  try {
    const ip = JSON.parse(rule.ip || '[]')
    ruleIPInput.value = Array.isArray(ip) ? ip.join(', ') : ''
  } catch {
    ruleIPInput.value = ''
  }
  ruleModalVisible.value = true
}

async function handleSaveRule() {
  ruleModalLoading.value = true
  try {
    const domains = ruleDomainInput.value
      ? ruleDomainInput.value.split(',').map((s) => s.trim()).filter(Boolean)
      : []
    const ips = ruleIPInput.value
      ? ruleIPInput.value.split(',').map((s) => s.trim()).filter(Boolean)
      : []

    const payload = {
      ...ruleForm,
      domain: domains.length > 0 ? JSON.stringify(domains) : '',
      ip: ips.length > 0 ? JSON.stringify(ips) : '',
    }

    if (isEditRule.value && ruleForm.id) {
      await updateRoutingRule(ruleForm.id, payload)
      message.success('分流规则已更新')
    } else {
      await createRoutingRule(payload)
      message.success('分流规则已创建')
    }

    ruleModalVisible.value = false
    fetchRules()
  } finally {
    ruleModalLoading.value = false
  }
}

async function handleDeleteRule(id?: number) {
  if (!id) return
  await deleteRoutingRule(id)
  message.success('规则已删除')
  fetchRules()
}

async function toggleRuleEnable(rule: RoutingRule) {
  if (!rule.id) return
  await updateRoutingRule(rule.id, rule)
  message.success(rule.enable ? '规则已启用' : '规则已禁用')
}

async function addPresetRule(type: string) {
  if (type === 'ads') {
    await createRoutingRule({
      tag: 'Block Ads',
      domain: JSON.stringify(['geosite:category-ads-all']),
      outbound: 'block',
      enable: true,
      order: rules.value.length + 1,
      remark: '拦截广告',
    })
  } else if (type === 'cn_domain') {
    await createRoutingRule({
      tag: 'Direct China Domain',
      domain: JSON.stringify(['geosite:cn']),
      outbound: 'direct',
      enable: true,
      order: rules.value.length + 1,
      remark: '国内主流域名直连',
    })
  } else if (type === 'cn_ip') {
    await createRoutingRule({
      tag: 'Direct China IP',
      ip: JSON.stringify(['geoip:cn']),
      outbound: 'direct',
      enable: true,
      order: rules.value.length + 1,
      remark: '国内 IP 直连',
    })
  } else if (type === 'private_ip') {
    await createRoutingRule({
      tag: 'Direct Private IP',
      ip: JSON.stringify(['geoip:private']),
      outbound: 'direct',
      enable: true,
      order: rules.value.length + 1,
      remark: '私有局域网直连',
    })
  }
  message.success('预设规则已添加')
  fetchRules()
}

async function handleSaveDNS() {
  dnsSaving.value = true
  try {
    await updateDNSSettings(dnsForm)
    message.success('DNS 解析配置已保存并重载')
  } finally {
    dnsSaving.value = false
  }
}

onMounted(() => {
  fetchRules()
  fetchDNS()
})
</script>

<style scoped>
.routing-page {
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

.main-tabs-card {
  border-radius: 12px;
  background: #ffffff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.preset-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  padding: 8px 12px;
  background: #f8fafc;
  border-radius: 8px;
}

.preset-label {
  font-size: 13px;
  font-weight: 600;
  color: #475569;
}

.match-cell {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.action-buttons {
  display: flex;
  gap: 4px;
}

.text-muted {
  font-size: 12px;
  color: #94a3b8;
}

.dns-form-container {
  max-width: 800px;
  padding: 12px 0;
}

.mt-4 {
  margin-top: 16px;
}

.mb-3 {
  margin-bottom: 12px;
}
</style>
