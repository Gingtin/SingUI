<template>
  <div class="api-docs-page">
    <div class="page-header">
      <div>
        <h2>REST API 开发者接口文档</h2>
        <p class="subtitle">SingUI 标准 RESTful API 接口规范，便于第三方自动化、机器人集成与集群节点调用</p>
      </div>
    </div>

    <a-card class="mt-4 main-card" :bordered="false">
      <div class="auth-box mb-4">
        <h4>🔑 鉴权说明 (Authentication)</h4>
        <p>所有受保护的 API 均需要在 HTTP Header 中携带 JWT Token：</p>
        <code>Authorization: Bearer &lt;YOUR_JWT_TOKEN&gt;</code>
      </div>

      <a-collapse v-model:activeKey="activeApiKeys">
        <!-- 1. Auth API -->
        <a-collapse-panel key="auth" header="1. 身份认证 (Auth API)">
          <div class="endpoint-item">
            <span class="method post">POST</span>
            <span class="path">/api/auth/login</span>
            <span class="desc">用户登录并获取 JWT Token</span>
          </div>
          <div class="endpoint-item">
            <span class="method get">GET</span>
            <span class="path">/api/auth/info</span>
            <span class="desc">获取当前登录管理员身份信息</span>
          </div>
        </a-collapse-panel>

        <!-- 2. Inbound API -->
        <a-collapse-panel key="inbounds" header="2. 入站节点与用户 (Inbounds & Clients API)">
          <div class="endpoint-item">
            <span class="method get">GET</span>
            <span class="path">/api/inbounds</span>
            <span class="desc">获取全部入站节点列表与挂载的客户端用户</span>
          </div>
          <div class="endpoint-item">
            <span class="method post">POST</span>
            <span class="path">/api/inbounds</span>
            <span class="desc">创建新的入站节点（支持 VLESS Reality, AnyTLS, Hysteria 2 等）</span>
          </div>
          <div class="endpoint-item">
            <span class="method put">PUT</span>
            <span class="path">/api/inbounds/:id</span>
            <span class="desc">修改指定入站节点配置</span>
          </div>
          <div class="endpoint-item">
            <span class="method delete">DELETE</span>
            <span class="path">/api/inbounds/:id</span>
            <span class="desc">删除指定入站节点</span>
          </div>
          <div class="endpoint-item">
            <span class="method post">POST</span>
            <span class="path">/api/inbounds/:id/clients</span>
            <span class="desc">在指定入站节点下添加客户端（UUID/配额/到期时间）</span>
          </div>
          <div class="endpoint-item">
            <span class="method post">POST</span>
            <span class="path">/api/inbounds/:id/clients/:clientId/reset</span>
            <span class="desc">重置指定用户的已用流量</span>
          </div>
        </a-collapse-panel>

        <!-- 3. Outbounds API -->
        <a-collapse-panel key="outbounds" header="3. 出站链路 (Outbounds API)">
          <div class="endpoint-item">
            <span class="method get">GET</span>
            <span class="path">/api/outbounds</span>
            <span class="desc">获取全部出站链路列表（Direct, Block, DNS, WARP）</span>
          </div>
          <div class="endpoint-item">
            <span class="method post">POST</span>
            <span class="path">/api/outbounds</span>
            <span class="desc">添加新的出站链路或 WARP 落地</span>
          </div>
        </a-collapse-panel>

        <!-- 4. Subscriptions API -->
        <a-collapse-panel key="subscriptions" header="4. 通用订阅分发 (Subscriptions API - Public)">
          <div class="endpoint-item">
            <span class="method get">GET</span>
            <span class="path">/sub/:token?flag=sing-box</span>
            <span class="desc">获取 Sing-box 官方客户端格式配置 JSON</span>
          </div>
          <div class="endpoint-item">
            <span class="method get">GET</span>
            <span class="path">/sub/:token?flag=clash</span>
            <span class="desc">获取 Clash Meta / Mihomo 格式 YAML 订阅</span>
          </div>
          <div class="endpoint-item">
            <span class="method get">GET</span>
            <span class="path">/sub/:token?flag=base64</span>
            <span class="desc">获取通用 Base64 节点 URI 链接列表</span>
          </div>
          <div class="endpoint-item">
            <span class="method get">GET</span>
            <span class="path">/sub/view/:token</span>
            <span class="desc">用户免登录自服务信息门户（流量进度条与二维码）</span>
          </div>
        </a-collapse-panel>
      </a-collapse>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const activeApiKeys = ref(['auth', 'inbounds', 'outbounds', 'subscriptions'])
</script>

<style scoped>
.api-docs-page {
  padding: 4px;
}

.page-header {
  margin-bottom: 16px;
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

.auth-box {
  padding: 12px 16px;
  background: #f8fafc;
  border-left: 4px solid #3b82f6;
  border-radius: 6px;
}

.auth-box h4 {
  margin: 0 0 4px 0;
  font-size: 14px;
  font-weight: 600;
}

.auth-box p {
  margin: 0 0 6px 0;
  font-size: 13px;
  color: #64748b;
}

.auth-box code {
  background: #0f172a;
  color: #38bdf8;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.endpoint-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid #f1f5f9;
}

.endpoint-item:last-child {
  border-bottom: none;
}

.method {
  font-family: monospace;
  font-weight: 700;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
  min-width: 60px;
  text-align: center;
}

.method.get {
  background: #dcfce7;
  color: #15803d;
}

.method.post {
  background: #dbeafe;
  color: #1d4ed8;
}

.method.put {
  background: #fef3c7;
  color: #b45309;
}

.method.delete {
  background: #fee2e2;
  color: #b91c1c;
}

.path {
  font-family: monospace;
  font-weight: 600;
  color: #0f172a;
  font-size: 13px;
  min-width: 280px;
}

.desc {
  font-size: 13px;
  color: #64748b;
}

.mt-4 {
  margin-top: 16px;
}

.mb-4 {
  margin-bottom: 16px;
}
</style>
