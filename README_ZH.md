<div align="center">

# SingUI

**专为 Sing-box 核心打造的下一代高性能 Web 可视化管理面板**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue.js-3.4+-4FC08D?style=flat-square&logo=vue.js&logoColor=white)](https://vuejs.org)
[![Sing-box](https://img.shields.io/badge/Sing--box-1.9+-blue?style=flat-square)](https://github.com/SagerNet/sing-box)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker&logoColor=white)](https://hub.docker.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)](LICENSE)

[**English**](README.md) • [**简体中文**](README_ZH.md)

</div>

---

## 📖 项目简介

**SingUI** 是一款专为 **Sing-box**（v1.9+）深度定制的现代化、高性能 Web 管理面板。系统采用 Go + Vue 3 技术栈构建，将后端守护进程、动态配置引擎与前端静态资源打包为**零外部运行依赖的单一可执行二进制文件**。

SingUI 提供了细粒度的单端口多用户流量计量、X25519 Reality 密钥协商、可视化规则集分流路由及全能自适应订阅转换管线，为 Sing-box 的生产级部署与节点运维提供坚实保障。

---

## 🏗️ 架构与系统能力

```
                  ┌───────────────────────────────┐
                  │       SingUI Web Panel        │
                  │   (Go Backend + Vue 3 SPA)    │
                  └──────────────┬────────────────┘
                                 │
           ┌─────────────────────┼─────────────────────┐
           ▼                     ▼                     ▼
┌────────────────────┐ ┌───────────────────┐ ┌────────────────────┐
│ 入站与多用户管理   │ │  分流路由与 DNS   │ │ 全能订阅分发引擎 │
│ - VLESS Reality    │ │ - Rule-Set 规则集 │ │ - Sing-box JSON    │
│ - Hysteria 2       │ │ - Geosite / GeoIP │ │ - Clash Meta YAML  │
│ - TUIC v5          │ │ - DoH / DoT 分流  │ │ - Base64 节点链    │
│ - Shadowsocks 2022 │ │ - FakeIP 解析引擎 │ │ - Web 用户自服务页 │
└──────────┬─────────┘ └─────────┬─────────┘ └─────────┬──────────┘
           │                     │                     │
           └─────────────────────┼─────────────────────┘
                                 ▼
                     ┌───────────────────────┐
                     │ 原子级配置预检机制    │
                     │  (sing-box check)     │
                     └───────────┬───────────┘
                                 ▼
                     ┌───────────────────────┐
                     │    Sing-box Core      │
                     │ (Supervisor 守护进程) │
                     └───────────────────────┘
```

### 核心特性

- **原生协议栈全覆盖**：
  - **VLESS**：支持 Reality 伪装（自动生成合法 X25519 密钥对与 ShortID）、`xtls-rprx-vision` 流控及 uTLS 客户端指纹模拟。
  - **Hysteria 2**：支持 Salamander 密码混淆、BBR / Brutal 拥塞控制及上下行独立带宽速率限制。
  - **TUIC v5**：基于 QUIC 协议传输，支持 BBR 拥塞控制与 0-RTT 极速握手。
  - **Shadowsocks 2022**：完整支持 Blake3-AEAD 2022 规范（`2022-blake3-aes-128-gcm` 等）与经典加密算法。
  - **Trojan & VMess**：支持 TCP、WebSocket、gRPC、HTTPUpgrade 传输协议与 TLS 加密。
- **单入站多用户 (Multi-Client) 隔离**：
  - 单端口支持承载多个独立客户端凭据。
  - 独立统计已用上行/下行流量、总限额配额、到期时间与同时在线并发 IP 限制。
  - 内置配额熔断调度器，超额或到期毫秒级禁用并热重载核心。
- **可视化规则集分流与 Split-DNS**：
  - 支持按 `geosite` 与 `geoip` 规则集进行动态分流。
  - 精确指定路由动作（`DIRECT` 直连、`BLOCK` 拦截、`DNS-OUT` 路由分流）。
  - 内置独立 Remote DoH 与国内直连 DNS 解析器。
- **配置预检与守护进程**：
  - 应用新配置前自动执行 `sing-box check` 语法校验，杜绝错误配置导致核心中断。
  - 内置 Supervisor 进程管理器，提供 WebSocket 实时终端日志滚屏与崩溃自愈。
  - 集成 Clash API 控制器，提供实时活跃连接分析与源 IP 追踪。

---

## ⚡ 协议与传输支持矩阵

| 协议 | 传输层 (Transport) | 安全与伪装 (Security) | 多用户支持 | 流控与拥塞控制 |
| :--- | :--- | :--- | :---: | :--- |
| **VLESS** | TCP, WS, gRPC, HTTPUpgrade | Reality (X25519), TLS | ✅ | `xtls-rprx-vision` |
| **Hysteria 2** | UDP | Salamander 混淆, TLS | ✅ | BBR / Brutal, 速率限制 |
| **TUIC v5** | UDP (QUIC) | TLS | ✅ | BBR, 0-RTT 握手 |
| **Shadowsocks** | TCP, UDP | Blake3-AEAD (2022) | ✅ | Multiplex |
| **Trojan** | TCP, WS, gRPC | TLS | ✅ | - |
| **VMess** | TCP, WS, gRPC, HTTPUpgrade | TLS, None | ✅ | AlterID 0 |

---

## 🚀 快速上手

### 1. Linux 一键安装脚本

在 Linux 服务器（Debian / Ubuntu / CentOS / Alpine / Arch）上执行：

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Gingtin/SingUI/main/scripts/install.sh)
```

部署完成后访问：
- **面板地址**：`http://<服务器IP>:2096`
- **默认用户名**：`admin`
- **默认密码**：`admin`

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

## 📡 订阅分发规范

通过 `/sub/:token` 统一获取适配的订阅配置：

| 客户端类型 | 查询参数 | 下发格式 |
| :--- | :--- | :--- |
| **Sing-box 客户端** | `?flag=sing-box` | Sing-box 官方客户端标准 JSON 配置 |
| **Clash Meta / Mihomo** | `?flag=clash` | Clash Meta YAML（含 Proxies、自动测速组与规则集） |
| **通用 Base64** | `?flag=base64` | 节点链接列表（`vless://`, `hysteria2://` 等） |
| **用户自服务页面** | `/sub/view/:token` | 网页端信息面板（展示剩余配额、到期日与二维码） |

所有订阅响应均自动附带 `Subscription-Userinfo` 头部，客户端可直观展示流量仪表盘。

---

## 🔧 命令行参数 (CLI)

```bash
singbox-ui -h
  -p string
        面板监听端口 (默认: 2096)
  -d string
        SQLite 数据库路径 (默认: "data/singbox_ui.db")
  -reset-admin
        将管理员账号密码重置为默认值 (admin/admin)
  -v    查看版本信息
```

---

## 💬 问题反馈与建议

如果您在使用中遇到任何 Bug 或有新功能需求，欢迎直接在 [GitHub Issues](https://github.com/Gingtin/SingUI/issues) 中提交议题。

---

## ❤️ 致谢与参考项目 (Acknowledgments)

- **[SagerNet/sing-box](https://github.com/SagerNet/sing-box)**：通用代理核心底座。
- **[XTLS/Xray-core](https://github.com/XTLS/Xray-core)**：Reality 与 XTLS 规范先驱。
- **[MetaCubeX/mihomo](https://github.com/MetaCubeX/mihomo)**：Clash Meta 规则集与生态标准。

---

## 📜 开源许可证

本项目基于 [MIT License](LICENSE) 许可证开源。
