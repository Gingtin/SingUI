<div align="center">

# SingUI

**专为 Sing-box 打造的下一代高性能可视化管理面板**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue.js-3.4+-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white)](https://vuejs.org)
[![Sing-box Core](https://img.shields.io/badge/Sing--box-1.9+-blue?style=for-the-badge)](https://github.com/SagerNet/sing-box)
[![Docker Support](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://hub.docker.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)

[**English**](README.md) • [**简体中文**](README_ZH.md)

</div>

---

## 📖 项目简介

**SingUI** 是一款专为 **Sing-box** 通用代理平台量身定制的现代化、高性能 Web 管理面板。系统采用 Go + Vue 3 技术栈开发，具备零依赖单文件交付、原生多用户流量精确计量、一键生成 Reality 密钥对、全协议可视化配置以及自适应多格式订阅分发能力。

无论是个人多节点管理还是团队与社区配额分发，SingUI 都能提供极致的性能、灵活性与安全保障。

---

## ✨ 核心特性

- ⚡ **深度集成 Sing-box 特色协议栈**：
  - **VLESS**：原生 Reality 伪装（一键生成合法 X25519 密钥对与 ShortID）、`xtls-rprx-vision` 流控，支持 TCP、WebSocket、gRPC、HTTPUpgrade 传输。
  - **Hysteria 2**：Salamander 混淆密码配置、BBR / Brutal 拥塞控制、独立上/下行速率限速。
  - **TUIC v5**：原生 QUIC 传输、BBR 拥塞控制与安全 Token 鉴权。
  - **Shadowsocks 2022**：完整支持 2022 规范（`2022-blake3-aes-128-gcm` 等）与经典加密算法。
  - **Trojan & VMess**：支持 TLS、WebSocket、gRPC 传输协议。
- 👥 **单入站多用户 (Multi-Client) 架构**：
  - 单个监听端口支持挂载多个独立 Client。
  - 各用户拥有独立已用上行、下行流量统计、总配额限额与到期时间。
  - 内置流量超额与到期自动熔断调度器，毫秒级触发 Sing-box 核心热重载。
  - 支持限制用户同时在线并发 IP 数。
- 📡 **全能自适应多格式订阅引擎**：
  - **Sing-box 官方客户端 JSON 配置** (`flag=sing-box`)
  - **Clash Meta / Mihomo YAML 配置** (`flag=clash` / `flag=meta`)，自动生成代理组与规则集
  - **通用 Base64 / URI 链接列表** (`vless://`, `hysteria2://`, `tuic://`, `ss://`, `trojan://`)
  - **用户自服务门户** (`/sub/view/:token`)：提供剩余流量、到期天数与二维码一键导入。
- 📊 **系统级监控与实时诊断**：
  - CPU、内存、磁盘实时占用监控指标。
  - ECharts 动态折线图实时呈现网络上行/下行速率。
  - Sing-box 核心进程 Supervisor（自动拉起、崩溃自愈、WebSocket 实时终端日志）。
  - Clash API 实时活跃连接列表与源 IP 追踪。
- 🤖 **运维与自动化**：
  - Telegram 机器人告警（流量耗尽、账号到期通知）与 SQLite 数据库定期自动备份推送。
  - SQLite 数据库一键下载备份与上传恢复。
  - ACME SSL 证书管理与面板安全路径自定义。
- 📦 **零外部依赖单二进制分发**：
  - Vue 3 前端静态产物在编译时直接内嵌至 Go 可执行文件（`embed.FS`），生产服务器无需安装 Node.js 或任何运行环境。

---

## ⚡ 协议支持与功能矩阵

| 协议 | 传输层 (Transport) | 安全与伪装 (Security) | 多用户支持 | 流控与拥塞控制 |
| :--- | :--- | :--- | :---: | :--- |
| **VLESS** | TCP, WS, gRPC, HTTPUpgrade | Reality, TLS, None | ✅ | `xtls-rprx-vision` |
| **Hysteria 2** | UDP | Salamander 混淆, TLS | ✅ | Brutal / BBR, 带宽速率限制 |
| **TUIC v5** | UDP (QUIC) | TLS | ✅ | BBR 拥塞控制 |
| **Shadowsocks** | TCP, UDP | None (2022 Blake3 / AEAD) | ✅ | Multiplex |
| **Trojan** | TCP, WS, gRPC | TLS | ✅ | - |
| **VMess** | TCP, WS, gRPC, HTTPUpgrade | TLS, None | ✅ | AlterID 0 |

---

## 🚀 快速开始

### 方式一：Linux 一键安装脚本（推荐）

在你的 Linux VPS 服务器（Debian / Ubuntu / CentOS / Alpine / Arch）上执行：

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Gingtin/SingUI/main/scripts/install.sh)
```

安装完成后在浏览器中打开：
- **访问地址**：`http://<你的服务器IP>:2096`
- **默认用户名**：`admin`
- **默认密码**：`admin`

> ⚠️ **安全提示**：初次登录后，请第一时间进入 **面板设置** -> **安全设置** 中修改默认管理员密码！

---

### 方式二：Docker Compose 容器部署

```bash
git clone https://github.com/Gingtin/SingUI.git
cd SingUI
docker compose up -d
```

---

### 方式三：从源码编译

#### 环境要求
- **Go**：1.22 或更高版本
- **Node.js**：18+ 与 npm / pnpm

```bash
# 1. 编译前端静态资源
cd frontend
npm install
npm run build

# 2. 编译 Go 后端可执行文件
cd ../backend
go build -ldflags="-s -w" -o ../singbox-ui ./cmd/server

# 3. 启动面板
../singbox-ui -p 2096 -d data/singbox_ui.db
```

---

## 🖥️ 界面功能概览

- **系统仪表盘**：实时展示 CPU、内存、硬盘占用及网络上下行吞吐动态折线图。
- **入站节点管理**：可视化节点配置向导，支持一键生成 Reality X25519 密钥对，多用户独立分配配额。
- **订阅中心**：一键导出 Sing-box JSON、Clash Meta YAML、Base64 订阅及扫码二维码。
- **运行日志与连接**：基于 WebSocket 的实时终端日志滚屏与 Clash API 活跃连接监控。
- **系统设置**：支持自定义面板端口、密码修改、Telegram 机器人预警与 SQLite 备份恢复。

---

## 📡 订阅接口规范

SingUI 会根据请求参数或客户端 User-Agent 自动分发适配的配置格式：

| 订阅格式 | URL 路径规范 | 适配客户端 |
| :--- | :--- | :--- |
| **Sing-box JSON** | `/sub/:token?flag=sing-box` | Sing-box 官方客户端, Box4, SFA |
| **Clash Meta YAML** | `/sub/:token?flag=clash` | Clash Verge Rev, Mihomo Party, Stash, Flclash |
| **Base64 节点链** | `/sub/:token?flag=base64` | Shadowrocket, v2rayN, Quantumult X, Loon |
| **用户自服务页面** | `/sub/view/:token` | 网页端自服务页（展示剩余流量、到期时间与二维码） |

所有订阅响应均自动携带标准的 `Subscription-Userinfo` 头部，客户端可直接显示流量使用仪表盘。

---

## 🔧 命令行参数 (CLI)

```bash
singbox-ui -h
  -p string
        面板监听端口 (默认从数据库读取或 2096)
  -d string
        SQLite 数据库路径 (默认: "data/singbox_ui.db")
  -reset-admin
        将管理员密码重置为默认值 'admin'
  -v    打印版本信息
```

---

## 🔒 生产环境安全建议

1. **防火墙保护**：妥善保护面板端口（默认 2096），建议通过 Nginx / Caddy 反向代理配置 HTTPS 访问。
2. **Reality SNI 选择**：配置 VLESS Reality 时，推荐使用大型正规网站的 SNI 伪装域名（如 `www.apple.com`、`addons.mozilla.org`、`www.cloudflare.com`）。
3. **定期数据备份**：定期在设置中下载 SQLite 备份文件，或启用 Telegram Bot 自动推送每周备份。

---

## 🤝 参与贡献

欢迎社区开发者参与贡献！请查阅 [CONTRIBUTING.md](CONTRIBUTING.md) 了解 Issue 提交规范与 Pull Request 流程。

---

## 📜 开源许可证

本项目基于 [MIT License](LICENSE) 许可证开源。

---

## ❤️ 致谢与参考项目 (Acknowledgments)

我们向以下优秀的开源项目与社区先驱致以崇高的敬意，它们为 **SingUI** 的设计与实现提供了重要灵感与基础支撑：

- **[Sing-box](https://github.com/SagerNet/sing-box)**：驱动 SingUI 核心的高性能下一代代理核心平台。
- **[3x-ui](https://github.com/MHSanaei/3x-ui)**：在多用户节点交互设计与面板体验方面的经典开创。
- **[Marzban](https://github.com/Gozargah/Marzban)**：在多协议统一订阅与多用户配额管控上的优秀设计理念。
- **[s-ui](https://github.com/alireza0/s-ui)**：在 Sing-box 可视化界面早期的探索与实践。
- **[Clash Meta / Mihomo](https://github.com/MetaCubeX/mihomo)**：成熟规范的规则集与代理提供者生态标准。
