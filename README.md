<div align="center">

# SingUI

**Next-Generation, High-Performance Web Management Platform for Sing-box Core**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue.js-3.4+-4FC08D?style=flat-square&logo=vue.js&logoColor=white)](https://vuejs.org)
[![Sing-box](https://img.shields.io/badge/Sing--box-1.9+-blue?style=flat-square)](https://github.com/SagerNet/sing-box)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker&logoColor=white)](https://hub.docker.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)](LICENSE)

[**English**](README.md) • [**简体中文**](README_ZH.md)

</div>

---

## 📖 Overview

**SingUI** is a lightweight, modern web management platform engineered specifically for the **Sing-box** universal proxy core. Built with Go and Vue 3, SingUI packages the core supervisor, dynamic configuration engine, multi-format subscription generator, and macOS-style dashboard into a **single standalone binary** with zero external runtime dependencies.

With an ultra-low memory footprint (**< 25MB RAM on VPS**), SingUI delivers single-port multi-client traffic accounting, automated Reality & AnyTLS key negotiation, Cloudflare WARP chain detours, multi-country rule-set routing, and automated kernel BBR / UDP buffer tuning.

---

## 🏗️ Architecture & Modules

```
                  ┌───────────────────────────────────────────────┐
                  │            SingUI Web Panel (Vue 3)           │
                  │       (macOS-Inspired Minimalist Design)      │
                  └───────────────────────┬───────────────────────┘
                                          │
       ┌──────────────┬───────────────────┼───────────────────┬──────────────┐
       ▼              ▼                   ▼                   ▼              ▼
┌─────────────┐ ┌─────────────┐   ┌───────────────┐   ┌───────────────┐ ┌────────────┐
│  Inbounds   │ │  Outbounds  │   │ Rule-Sets/DNS │   │ Subscriptions │ │ Diagnostics│
│ - Reality   │ │ - Direct    │   │ - CN/IR/RU SRS│   │ - Sing-box    │ │ - Live WS  │
│ - AnyTLS    │ │ - Block     │   │ - Split-DNS   │   │ - Clash Meta  │ │ - Clash API│
│ - Hy2/TUIC  │ │ - WARP / WG │   │ - Anti-DPI    │   │ - Base64 URI  │ │ - raw JSON │
│ - SS2022    │ │ - Detours   │   │ - FakeIP      │   │ - User Portal │ │ - Backups  │
└──────┬──────┘ └──────┬──────┘   └───────┬───────┘   └───────┬───────┘ └─────┬──────┘
       │               │                  │                   │               │
       └───────────────┴──────────────────┼───────────────────┴───────────────┘
                                          ▼
                              ┌───────────────────────┐
                              │  Atomic Syntax Check  │
                              │   (sing-box check)    │
                              └───────────┬───────────┘
                                          ▼
                              ┌───────────────────────┐
                              │  Sing-box Core 1.9+   │
                              │ (Supervisor Process)  │
                              └───────────────────────┘
```

---

## ✨ Key Capabilities

### 1. Native Protocol Suite
- **VLESS Reality & Vision**: Zero-certificate deployment with automated X25519 keypair generation, ShortID rotation, `xtls-rprx-vision` zero-copy flow, and uTLS fingerprint simulation (`chrome`, `firefox`, `safari`, `ios`).
- **AnyTLS (Official Native Support)**: Mitigates TLS-in-TLS nested handshake fingerprinting using packet padding schemes and session multiplexing.
- **Hysteria 2**: Salamander obfuscation password, BBR / Brutal congestion control, and independent upstream/downstream bandwidth rate limits.
- **TUIC v5**: Native QUIC transport, BBR congestion control, and 0-RTT handshake support.
- **Shadowsocks 2022**: Blake3-AEAD 2022 specifications (`2022-blake3-aes-128-gcm`, `2022-blake3-aes-256-gcm`).
- **Trojan & VMess**: TLS, WebSocket, gRPC, and HTTPUpgrade transports.

### 2. Multi-Client per Inbound Architecture (3x-ui Paradigm)
- Multiple independent client credentials per listening port.
- Granular per-user traffic metering (Upload / Download / Total Quota), expiration timers, and IP concurrency limits.
- Sub-second quota exhaustion and expiration enforcement with atomic core reloading.

### 3. Outbounds & Cloudflare WARP Detours
- Manage `direct`, `block`, `dns-out`, and custom proxy chains.
- **Cloudflare WARP / WireGuard Detours**: Route selected traffic through Cloudflare WARP to bypass ChatGPT, Claude, Netflix, and OpenAI residential IP restrictions.

### 4. Multi-Country Rule-Set Routing & Split-DNS
- Full support for Sing-box 1.9+ binary `.srs` rule sets.
- **One-Click Regional Optimization Presets**:
  - 🇨🇳 **China Preset**: `geosite:cn` & `geoip:cn` Direct + Adblock + AliDNS / DNSPod + Remote DoH.
  - 🇮🇷 **Iran Preset**: `geosite:ir` & `geoip:ir` Direct (banking, gov) + Shecan DNS + TLS Fragmentation anti-DPI.
  - 🇷🇺 **Russia Preset**: `geosite:ru` & `geoip:ru` Direct (Gosuslugi, VK, Yandex) + Yandex / Quad9 DNS.
  - 🌐 **Global Preset**: Private LAN Direct + standard proxy routing.

### 5. Universal Subscription Pipelines
- **Sing-box Official Client JSON** (`flag=sing-box`)
- **Clash Meta / Mihomo YAML** (`flag=clash`) with auto-generated Proxy Groups and Rule Providers.
- **Universal Base64 / URI Links** (`vless://`, `hysteria2://`, `tuic://`, `ss://`, `trojan://`)
- **Web User Self-Service Portal** (`/sub/view/:token`) with visual bandwidth meters and QR codes.

### 6. Robust Reliability & Ultra-Low Memory
- **SQLite WAL Mode**: High-concurrency database engine with `PRAGMA journal_mode=WAL` preventing database locks.
- **Zero-Crash Preflight Validation**: Every configuration is validated via `sing-box check` in a sandbox before applying.
- **VPS Memory Footprint**: Less than **25MB RAM** resident memory.

---

## ⚡ Protocol & Feature Matrix

| Protocol | Transports | Security / Camouflage | Multi-Client | Flow / Congestion |
| :--- | :--- | :--- | :---: | :--- |
| **VLESS** | TCP, WS, gRPC, HTTPUpgrade | Reality (X25519), TLS | ✅ | `xtls-rprx-vision` |
| **AnyTLS** | TCP | Standard TLS, Padding | ✅ | Session Mux |
| **Hysteria 2** | UDP | Salamander Obfs, TLS | ✅ | BBR / Brutal, Rate Limit |
| **TUIC v5** | UDP (QUIC) | TLS | ✅ | BBR, 0-RTT Handshake |
| **Shadowsocks** | TCP, UDP | Blake3-AEAD (2022) | ✅ | Multiplex |
| **Trojan** | TCP, WS, gRPC | TLS | ✅ | - |
| **VMess** | TCP, WS, gRPC, HTTPUpgrade | TLS, None | ✅ | AlterID 0 |

---

## 🚀 Quick Start

### 1. Linux One-Click Installer (Recommended)

Run the following command on Debian / Ubuntu / CentOS / Alpine / Arch:

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Gingtin/SingUI/main/scripts/install.sh)
```

The installer automatically enables **Linux BBR** and **UDP kernel buffer optimizations** for peak Hysteria 2 / TUIC throughput.

- **Panel URL**: `http://<your-server-ip>:2096`
- **Default Username**: `admin`
- **Default Password**: `admin`
- **CLI Helper**: `sing-ui {start|stop|restart|status|reset-admin}`

---

### 2. Docker Compose Deployment

```bash
git clone https://github.com/Gingtin/SingUI.git
cd SingUI
docker compose up -d
```

---

### 3. Build from Source

#### Prerequisites
- Go 1.22+
- Node.js 20+

```bash
# Build frontend assets
cd frontend
npm install
npm run build

# Compile single binary
cd ../backend
go build -ldflags="-s -w" -o ../singbox-ui ./cmd/server

# Launch SingUI
../singbox-ui -p 2096 -d data/singbox_ui.db
```

---

## 🔧 CLI Options

```bash
singbox-ui -h
  -p string
        Web panel port (default: 2096)
  -d string
        SQLite database path (default: "data/singbox_ui.db")
  -reset-admin
        Reset admin password to default (admin/admin)
  -v    Show version information
```

---

## 💬 Feedback & Issues

If you encounter any bugs or have feature requests, please submit an issue via [GitHub Issues](https://github.com/Gingtin/SingUI/issues).

---

## ❤️ Acknowledgments

- **[SagerNet/sing-box](https://github.com/SagerNet/sing-box)**: The universal proxy platform powering the high-performance core of SingUI.
- **[XTLS/Xray-core](https://github.com/XTLS/Xray-core)**: Pioneers in Reality and XTLS specifications.
- **[MetaCubeX/mihomo](https://github.com/MetaCubeX/mihomo)**: Standard-setting rule provider and Clash Meta proxy provider specifications.

---

## 📜 License

SingUI is open-source software licensed under the [MIT License](LICENSE).
