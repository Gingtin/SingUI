# SingUI

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Vue 3](https://img.shields.io/badge/Vue-3.4+-4FC08D?style=flat&logo=vue.js)](https://vuejs.org)
[![Sing-box](https://img.shields.io/badge/Sing--box-1.9+-blue?style=flat)](https://github.com/SagerNet/sing-box)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

**SingUI** 是一个结合 **3x-ui** 经典交互与 **Sing-box** 现代协议栈的高性能代理面板。单二进制交付、零外部依赖、原生支持多用户流量配额管控、全协议节点生成与一键多格式订阅转换。

---

## ✨ 核心特性

- 🚀 **全面支持 Sing-box 协议栈**：
  - **VLESS**：Reality 伪装（一键生成 X25519 密钥对与 ShortID）、XTLS-Vision、TCP、WebSocket、gRPC、HTTPUpgrade。
  - **Hysteria 2**：Salamander 混淆、BBR/Brutal 拥塞控制、上下行速率限速。
  - **TUIC v5**：原生 QUIC / BBR 拥塞控制。
  - **Shadowsocks 2022**：2022 规范（`2022-blake3-aes-128-gcm` 等）与经典加密。
  - **Trojan / VMess**：TLS、WebSocket、gRPC 传输。
- 👥 **单入站多用户模型 (Multi-Client)**：
  - 每个入站支持多 Client 独立统计已用上行/下行流量。
  - 独立到期时间、总流量限额、超额/到期自动熔断与配置热重载。
  - 单用户并发在线 IP 限制。
- 📡 **全客户端全能订阅引擎**：
  - **Sing-box 官方客户端格式** (`json`)
  - **Clash Meta / Mihomo 格式** (`yaml`，自动包含 Proxies、Proxy Groups、规则分流)
  - **通用 Base64 / URI 链接列表** (`vless://`, `hysteria2://`, `tuic://`, `ss://`, `trojan://`)
  - **用户自服务门户**：提供剩余流量、到期天数、扫码一键导入。
- 📊 **系统监控与诊断**：
  - CPU、内存、磁盘实时使用率。
  - ECharts 实时网络上/下行速率动态折线图。
  - Sing-box 核心进程 Supervisor（自动拉起、崩溃自愈、WebSocket 实时终端日志）。
  - Clash API 实时活动连接列表与源 IP 追踪。
- 🤖 **Telegram 机器人与运维辅助**：
  - 用户流量耗尽 / 到期自动告警。
  - 每日/每周 SQLite 数据库自动备份推送到 Telegram。
  - `/status`, `/traffic`, `/restart` 远程交互指令。
- 📦 **轻量化单文件分发**：
  - 基于 Go `embed.FS` 将 Vue 3 前端完整打包至单个二进制文件中，开箱即用。

---

## 🚀 快速开始

### 方式一：Linux 一键安装脚本（推荐）

在你的 VPS 服务器（Debian / Ubuntu / CentOS / Alpine / Arch）上执行：

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Gingtin/SingUI/main/scripts/install.sh)
```

安装完成后访问：`http://<你的服务器IP>:2096`
- **默认用户名**：`admin`
- **默认密码**：`admin`

### 方式二：Docker Compose 部署

```bash
git clone https://github.com/Gingtin/SingUI.git
cd SingUI
docker compose up -d
```

### 方式三：源码编译

#### 前置要求
- Go 1.22+
- Node.js 18+ & npm / pnpm

```bash
# 1. 构建前端资源
cd frontend
npm install
npm run build

# 2. 编译 Go 单二进制文件
cd ../backend
go build -ldflags="-s -w" -o ../bin/singbox-ui ./cmd/server

# 3. 运行面板
../bin/singbox-ui -p 2096 -d data/singbox_ui.db
```

---

## 📖 使用指南

### 1. 添加入站节点 (以 VLESS Reality 为例)
1. 进入 **入站节点** 页面，点击 **添加入站节点**。
2. 协议选择 `VLESS`，安全协议选择 `Reality`。
3. 点击 **一键生成密钥对**，系统将自动生成 X25519 `私钥`、`公钥` 与 `ShortID`。
4. 伪装域名 (SNI) 推荐填写：`www.apple.com` 或 `addons.mozilla.org`。
5. 监听端口填写 `443`（或其他端口），点击保存。
6. 面板将自动更新并热重载 Sing-box 核心。

### 2. 获取客户端订阅
1. 进入 **订阅中心** 页面。
2. 找到对应用户的操作列，点击 **复制订阅链接**：
   - 选择 **Sing-box JSON 订阅**：直接导入 Sing-box 客户端。
   - 选择 **Clash Meta (Mihomo) 订阅**：直接导入 Clash Verge / Mihomo Party / Stash。
   - 选择 **Base64 通用订阅**：适用于 Shadowrocket、v2rayN 等通用客户端。
3. 也可点击 **二维码** 使用手机客户端直接扫码导入。

---

## 🔒 生产环境安全建议

1. 登录面板后，请第一时间在 **面板设置** -> **安全设置** 中修改默认管理员密码。
2. 建议通过 Nginx 反向代理配置 HTTPS 证书，或修改默认 `2096` 端口。
3. 可在 **系统设置** 中配置 **Telegram 机器人**，当用户流量超限或节点异常时第一时间收到推送。

---

## 📄 开源许可证

本项目基于 [MIT License](LICENSE) 开源。
