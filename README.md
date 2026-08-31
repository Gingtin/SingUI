<div align="center">

# SingUI

**Next-Generation, High-Performance Web Management Panel for Sing-box**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue.js-3.4+-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white)](https://vuejs.org)
[![Sing-box Core](https://img.shields.io/badge/Sing--box-1.9+-blue?style=for-the-badge)](https://github.com/SagerNet/sing-box)
[![Docker Support](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://hub.docker.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)

[**English**](README.md) • [**简体中文**](README_ZH.md)

</div>

---

## 📖 Overview

**SingUI** is a modern, feature-rich web management panel designed specifically for the **Sing-box** universal proxy platform. Engineered with Go and Vue 3, SingUI delivers a single-binary deployment with zero external dependencies, robust multi-client traffic accounting, one-click Reality keypair generation, comprehensive protocol configuration, and an adaptive universal subscription engine.

Whether you are managing personal nodes or provisioning access for teams and communities, SingUI provides the performance, agility, and security you need.

---

## ✨ Key Features

- ⚡ **Native Sing-box Protocol Suite**:
  - **VLESS**: Reality camouflage (one-click X25519 keypair & ShortID generation), `xtls-rprx-vision` flow, TCP, WebSocket, gRPC, HTTPUpgrade.
  - **Hysteria 2**: Salamander obfuscation password, BBR / Brutal congestion control, and independent upstream/downstream bandwidth limits.
  - **TUIC v5**: Native QUIC transport, BBR congestion control, and secure token authentication.
  - **Shadowsocks 2022**: Full support for 2022 specifications (`2022-blake3-aes-128-gcm`, etc.) and classic ciphers.
  - **Trojan & VMess**: Support for TLS, WebSocket, and gRPC transports.
- 👥 **Multi-Client per Inbound Architecture**:
  - Single listening port supporting multiple independent clients.
  - Granular traffic accounting (Upload, Download, Total Quota) and expiration timers per user.
  - Automatic quota exhaustion and expiration enforcement with instant core hot-reloading.
  - Simultaneous IP concurrency limits.
- 📡 **Universal Multi-Format Subscription Engine**:
  - **Sing-box Official Client JSON** (`flag=sing-box`)
  - **Clash Meta / Mihomo YAML** (`flag=clash` / `flag=meta`) with auto-generated Proxy Groups and routing rule-sets.
  - **Standard Base64 / URI Links** (`vless://`, `hysteria2://`, `tuic://`, `ss://`, `trojan://`)
  - **Client User Self-Service Portal** (`/sub/view/:token`) with visual quota meters and QR codes.
- 📊 **Real-Time Monitoring & Observability**:
  - Live system resource gauges (CPU %, RAM %, Disk %, and Network I/O throughput).
  - Dynamic ECharts line charts for real-time upstream & downstream traffic rates.
  - Built-in Sing-box process supervisor with WebSocket live log streaming.
  - Clash API inspection for active connections, target hosts, and source IP geolocation.
- 🤖 **Operations & Automation**:
  - Automated Telegram Bot for quota warnings, node offline alerts, and scheduled database backups.
  - SQLite one-click backup download and database restoration.
  - ACME SSL automation and configurable panel routing paths.
- 📦 **Single-Binary Zero-Dependency Packaging**:
  - Vue 3 frontend is embedded directly into the Go binary (`embed.FS`). No Node.js or external runtime needed on production servers.

---

## ⚡ Protocol Support Matrix

| Protocol | Transports | Security / Camouflage | Multi-Client | Flow / Congestion |
| :--- | :--- | :--- | :---: | :--- |
| **VLESS** | TCP, WS, gRPC, HTTPUpgrade | Reality, TLS, None | ✅ | `xtls-rprx-vision` |
| **Hysteria 2** | UDP | Salamander Obfs, TLS | ✅ | Brutal / BBR, Bandwidth Limits |
| **TUIC v5** | UDP (QUIC) | TLS | ✅ | BBR, Congestion Control |
| **Shadowsocks** | TCP, UDP | None (2022 Blake3 / AEAD) | ✅ | Multiplex |
| **Trojan** | TCP, WS, gRPC | TLS | ✅ | - |
| **VMess** | TCP, WS, gRPC, HTTPUpgrade | TLS, None | ✅ | AlterID 0 |

---

## 🚀 Quick Start

### Method 1: Linux One-Click Installer (Recommended)

Run the following command on your Linux server (Debian / Ubuntu / CentOS / Alpine / Arch):

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Gingtin/SingUI/main/scripts/install.sh)
```

Once installed, open your browser and navigate to:
- **URL**: `http://<your-server-ip>:2096`
- **Default Username**: `admin`
- **Default Password**: `admin`

> ⚠️ **Security Tip**: Please change the default admin credentials immediately after your first login in **Settings**.

---

### Method 2: Docker Compose

```bash
git clone https://github.com/Gingtin/SingUI.git
cd SingUI
docker compose up -d
```

---

### Method 3: Build from Source

#### Prerequisites
- **Go**: 1.22 or higher
- **Node.js**: 18+ & npm / pnpm

```bash
# 1. Build the frontend assets
cd frontend
npm install
npm run build

# 2. Compile the Go backend
cd ../backend
go build -ldflags="-s -w" -o ../singbox-ui ./cmd/server

# 3. Launch the panel
../singbox-ui -p 2096 -d data/singbox_ui.db
```

---

## 🖥️ User Interface Preview

- **Dashboard**: Real-time CPU, RAM, Disk, and live upstream/downstream network speed charts.
- **Inbound Management**: Quick node creation wizard with one-click Reality X25519 keypair generation and multi-client traffic allocation.
- **Subscription Center**: Export Sing-box JSON, Clash Meta YAML, Base64, and QR codes for seamless multi-platform imports.
- **Real-Time Logs**: Interactive WebSocket terminal streamer with automatic log buffering.
- **Settings & Backups**: Configurable ports, secret tokens, Telegram alerts, and database backup/restore.

---

## 📡 Subscription Endpoints

SingUI automatically delivers optimized configuration payloads based on query parameters or client User-Agents:

| Format | URL Pattern | Supported Clients |
| :--- | :--- | :--- |
| **Sing-box JSON** | `/sub/:token?flag=sing-box` | Sing-box official client, Box4, SFA |
| **Clash Meta YAML** | `/sub/:token?flag=clash` | Clash Verge Rev, Mihomo Party, Stash, Flclash |
| **Base64 URI List** | `/sub/:token?flag=base64` | Shadowrocket, v2rayN, Quantumult X, Loon |
| **Self-Service View** | `/sub/view/:token` | Web browser portal (usage meters, expiration date, QR codes) |

All subscription responses include standard `Subscription-Userinfo` headers (`upload=...; download=...; total=...; expire=...`) for client bandwidth meters.

---

## 🔧 CLI Flags

```bash
singbox-ui -h
  -p string
        Web panel listening port (default: from DB or 2096)
  -d string
        SQLite database path (default: "data/singbox_ui.db")
  -reset-admin
        Reset admin password to default 'admin'
  -v    Print version information
```

---

## 🔒 Security Best Practices

1. **Firewall Rules**: Keep port `2096` (or your custom panel port) protected or use a reverse proxy (e.g. Nginx, Caddy) with an SSL certificate.
2. **Reality Deployment**: When configuring VLESS Reality, use authoritative and widely deployed SNIs (e.g., `www.apple.com`, `addons.mozilla.org`, `www.cloudflare.com`).
3. **Database Backups**: Regularly download SQLite database backups or configure the Telegram Bot for automatic weekly backup delivery.

---

## 🤝 Contributing

Contributions are welcome! Please check out [CONTRIBUTING.md](CONTRIBUTING.md) for details on submitting issues, proposing features, and opening pull requests.

---

## 📜 License

SingUI is open-source software licensed under the [MIT License](LICENSE).

---

## ❤️ Acknowledgments & References

We would like to express our gratitude to the following outstanding open-source projects that inspired and empowered the creation of **SingUI**:

- **[Sing-box](https://github.com/SagerNet/sing-box)**: The universal proxy platform powering the high-performance core of SingUI.
- **[3x-ui](https://github.com/MHSanaei/3x-ui)**: Pioneer in multi-client Xray panel UX and workflow design.
- **[Marzban](https://github.com/Gozargah/Marzban)**: Excellent design concepts in multi-protocol subscription routing and user management.
- **[s-ui](https://github.com/alireza0/s-ui)**: Early explorations in Sing-box web user interfaces.
- **[Clash Meta / Mihomo](https://github.com/MetaCubeX/mihomo)**: Standard-setting rule provider and proxy provider specifications.
