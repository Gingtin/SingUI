<div align="center">

# SingUI

**专为 Sing-box 核心打造的下一代高性能 Web 可视化管理平台**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue.js-3.4+-4FC08D?style=flat-square&logo=vue.js&logoColor=white)](https://vuejs.org)
[![Sing-box](https://img.shields.io/badge/Sing--box-1.9+-blue?style=flat-square)](https://github.com/SagerNet/sing-box)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker&logoColor=white)](https://hub.docker.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)](LICENSE)

[**English**](README.md) • [**简体中文**](README_ZH.md)

</div>

---

## 📖 项目简介

**SingUI** 是一款专为 **Sing-box** 通用代理平台深度定制的现代化、高性能 Web 管理系统。系统采用 Go + Vue 3 技术栈构建，将后端守护进程、动态配置引擎、多格式订阅转换器与 macOS 极简美学仪表盘打包为**零外部运行依赖的单一可执行文件**。

SingUI 拥有极致的内存控制能力（**VPS 常驻内存 < 25MB**），提供单端口多用户流量精确计量、一键 X25519 Reality / AnyTLS 密钥协商、Cloudflare WARP 链式出站解锁、多国规则集分流与 Linux 原生 BBR / UDP 缓冲区调优。

---

## 🏗️ 架构与系统能力

```
                  ┌───────────────────────────────────────────────┐
                  │            SingUI Web 面板 (Vue 3)            │
                  │              (macOS 极简美学设计)             │
                  └───────────────────────┬───────────────────────┘
                                          │
       ┌──────────────┬───────────────────┼───────────────────┬──────────────┐
       ▼              ▼                   ▼                   ▼              ▼
┌─────────────┐ ┌─────────────┐   ┌───────────────┐   ┌───────────────┐ ┌────────────┐
│ 入站节点管理│ │ 出站链路管理│   │ 规则路由与 DNS│   │ 全能订阅中心  │ │ 诊断与配置 │
│ - Reality   │ │ - Direct    │   │ - 中/伊/俄 SRS│   │ - Sing-box    │ │ - 实时 WS  │
│ - AnyTLS    │ │ - Block     │   │ - Split-DNS   │   │ - Clash Meta  │ │ - Clash API│
│ - Hy2/TUIC  │ │ - WARP / WG │   │ - 防 DPI 分片 │   │ - Base64 链接 │ │ - 原生配置 │
│ - SS2022    │ │ - 链式代理  │   │ - FakeIP 引擎 │   │ - 用户自服务页│ │ - 备份恢复 │
└──────┬──────┘ └──────┬──────┘   └───────┬───────┘   └───────┬───────┘ └─────┬──────┘
       │               │                  │                   │               │
       └───────────────┴──────────────────┼───────────────────┴───────────────┘
                                          ▼
                              ┌───────────────────────┐
                              │ 原子级配置语法预检    │
                              │   (sing-box check)    │
                              └───────────┬───────────┘
                                          ▼
                              ┌───────────────────────┐
                              │  Sing-box Core 1.9+   │
                              │ (Supervisor 守护进程) │
                              └───────────────────────┘
```

---

## ✨ 核心特性矩阵

### 1. 原生协议栈全覆盖
- **VLESS Reality & Vision**：免证书伪装部署，自动生成 X25519 密钥对与 ShortID，支持 `xtls-rprx-vision` 流控及 uTLS 客户端指纹模拟（`chrome`, `firefox`, `safari`, `ios`）。
- **AnyTLS (官方原生集成)**：通过数据包动态填充（Padding）与会话多路复用，彻底解决 TLS-in-TLS 嵌套特征阻断。
- **Hysteria 2**：支持 Salamander 密码混淆、BBR / Brutal 拥塞控制与上下行独立限速。
- **TUIC v5**：基于 QUIC 协议传输，支持 BBR 拥塞控制与 0-RTT 极速握手。
- **Shadowsocks 2022**：完整支持 Blake3-AEAD 2022 规范（`2022-blake3-aes-128-gcm` 等）。
- **Trojan & VMess**：支持 TLS、WebSocket、gRPC、HTTPUpgrade 传输协议。

### 2. 单入站多用户 (Multi-Client) 架构
- 单个监听端口支持挂载多个独立客户端。
- 独立统计已用上行、下行流量、总限额配额、到期时间与并发 IP 限制。
- 超额或到期自动熔断并毫秒级触发核心热重载。

### 3. 出站链路与 Cloudflare WARP 链式落地
- 可视化管理 `direct`、`block`、`dns-out` 及自定义代理中继。
- **Cloudflare WARP / WireGuard 链式出站**：一键让出口流量通过 WARP 转发，完美解锁 ChatGPT、Claude、Netflix、Disney+。

### 4. 多国家分流规则集与 Split-DNS
- 完整支持 Sing-box 1.9+ 原生二进制 `.srs` 规则集。
- **一键国家优化预设**：
  - 🇨🇳 **中国优化预设**：`geosite:cn` 与 `geoip:cn` 直连 + 广告拦截 + AliDNS / DNSPod 直连解析 + Remote DoH。
  - 🇮🇷 **伊朗优化预设**：`geosite:ir` 与 `geoip:ir` 直连 + Shecan 本地反审查 DNS + TLS 分片穿透 DPI。
  - 🇷🇺 **俄罗斯优化预设**：`geosite:ru` 与 `geoip:ru` 直连 + Yandex / Quad9 DNS。
  - 🌐 **全球通用预设**：私有局域网直连 + 全局代理。

### 5. 全能自适应多格式订阅管线
- **Sing-box 官方客户端 JSON** (`flag=sing-box`)
- **Clash Meta / Mihomo YAML** (`flag=clash`)，自动生成代理组与规则集
- **通用 Base64 / URI 链接** (`vless://`, `hysteria2://`, `tuic://`, `ss://`, `trojan://`)
- **用户自服务门户** (`/sub/view/:token`)：提供剩余配额进度环、到期天数与客户端一键导入。

### 6. 高可靠性与极简内存常驻
- **SQLite WAL 并发引擎**：开启 `PRAGMA journal_mode=WAL`，彻底杜绝数据库锁死。
- **配置语法沙盒预检**：应用配置前自动调用 `sing-box check`，保障核心零中断。
- **内存极低开销**：VPS 常驻内存小于 **25MB**。

---

## ⚡ 协议与传输支持矩阵

| 协议 | 传输层 (Transport) | 安全与伪装 (Security) | 多用户支持 | 流控与拥塞控制 |
| :--- | :--- | :--- | :---: | :--- |
| **VLESS** | TCP, WS, gRPC, HTTPUpgrade | Reality (X25519), TLS | ✅ | `xtls-rprx-vision` |
| **AnyTLS** | TCP | Standard TLS, Padding | ✅ | 会话多路复用 (Mux) |
| **Hysteria 2** | UDP | Salamander 混淆, TLS | ✅ | BBR / Brutal, 速率限制 |
| **TUIC v5** | UDP (QUIC) | TLS | ✅ | BBR, 0-RTT 握手 |
| **Shadowsocks** | TCP, UDP | Blake3-AEAD (2022) | ✅ | Multiplex |
| **Trojan** | TCP, WS, gRPC | TLS | ✅ | - |
| **VMess** | TCP, WS, gRPC, HTTPUpgrade | TLS, None | ✅ | AlterID 0 |

---

## 🚀 快速上手

### 1. Linux 一键安装脚本（推荐）

在 Linux 服务器（Debian / Ubuntu / CentOS / Alpine / Arch）上执行：

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Gingtin/SingUI/main/scripts/install.sh)
```

脚本会自动启用 **Linux 原生 BBR** 与 **UDP 内核缓冲区调优**，释放 Hysteria 2 / TUIC 的极致吞吐。

- **面板地址**：`http://<服务器IP>:2096`
- **默认用户名**：`admin`
- **默认密码**：`admin`
- **管理快捷指令**：`sing-ui {start|stop|restart|status|reset-admin}`

---

### 2. Docker 部署

```bash
git clone https://github.com/Gingtin/SingUI.git
cd SingUI
docker compose up -d
```

---

### 3. 从源码编译

#### 前置要求
- Go 1.22+
- Node.js 20+

```bash
# 1. 编译前端资源
cd frontend
npm install
npm run build

# 2. 编译 Go 单二进制可执行文件
cd ../backend
go build -ldflags="-s -w" -o ../singbox-ui ./cmd/server

# 3. 运行面板
../singbox-ui -p 2096 -d data/singbox_ui.db
```

---

## 🔧 命令行参数 (CLI)

```bash
singbox-ui -h
  -p string
        面板监听端口 (默认: 2096)
  -d string
        SQLite 数据库路径 (默认: "data/singbox_ui.db")
  -reset-admin
        重置管理员账号密码为默认值 (admin/admin)
  -v    查看版本信息
```

---

## 💬 问题反馈与建议

如果您在使用中遇到任何 Bug 或有新功能需求，欢迎在 [GitHub Issues](https://github.com/Gingtin/SingUI/issues) 中提交议题。

---

## ❤️ 致谢与参考项目 (Acknowledgments)

- **[SagerNet/sing-box](https://github.com/SagerNet/sing-box)**：驱动 SingUI 核心的高性能下一代代理核心平台。
- **[XTLS/Xray-core](https://github.com/XTLS/Xray-core)**：Reality 与 XTLS 规范先驱。
- **[MetaCubeX/mihomo](https://github.com/MetaCubeX/mihomo)**：Clash Meta 规则集与生态标准。

---

## 📜 开源许可证

本项目基于 [MIT License](LICENSE) 许可证开源。
