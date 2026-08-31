<div align="center">

# SingUI

**Next-Generation Web Management Panel for Sing-box Core**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue.js-3.4+-4FC08D?style=flat-square&logo=vue.js&logoColor=white)](https://vuejs.org)
[![Sing-box](https://img.shields.io/badge/Sing--box-1.9+-blue?style=flat-square)](https://github.com/SagerNet/sing-box)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker&logoColor=white)](https://hub.docker.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)](LICENSE)

[**English**](README.md) • [**简体中文**](README_ZH.md)

</div>

---

## Overview

**SingUI** is a lightweight, high-performance web management platform engineered specifically for **Sing-box** (v1.9+). Built with Go and Vue 3, SingUI packages the backend supervisor, dynamic config engine, and frontend assets into a **single standalone binary** with zero external runtime dependencies.

SingUI brings granular multi-client traffic accounting, automated Reality X25519 keypair negotiation, dynamic rule-set routing, and adaptive multi-format subscription pipelines to modern Sing-box deployments.

---

## Architecture & Capabilities

```
                  ┌───────────────────────────────┐
                  │       SingUI Web Panel        │
                  │   (Go Backend + Vue 3 SPA)    │
                  └──────────────┬────────────────┘
                                 │
           ┌─────────────────────┼─────────────────────┐
           ▼                     ▼                     ▼
┌────────────────────┐ ┌───────────────────┐ ┌────────────────────┐
│ Inbound & Clients  │ │  Routing & DNS    │ │ Universal Subs     │
│ - VLESS Reality    │ │ - Rule-Set (SRS)  │ │ - Sing-box JSON    │
│ - Hysteria 2       │ │ - Geosite / GeoIP │ │ - Clash Meta YAML  │
│ - TUIC v5          │ │ - DoH / DoT       │ │ - Base64 URI List  │
│ - Shadowsocks 2022 │ │ - FakeIP Engine   │ │ - Web User Portal  │
└──────────┬─────────┘ └─────────┬─────────┘ └─────────┬──────────┘
           │                     │                     │
           └─────────────────────┼─────────────────────┘
                                 ▼
                     ┌───────────────────────┐
                     │ Atomic Config Check   │
                     │  (sing-box check)     │
                     └───────────┬───────────┘
                                 ▼
                     ┌───────────────────────┐
                     │    Sing-box Core      │
                     │ (Supervisor Process)  │
                     └───────────────────────┘
```

### Core Highlights

- **Native Protocol Suite**:
  - **VLESS**: Reality camouflage with automatic X25519 keypair generation, ShortID rotation, `xtls-rprx-vision` flow control, and uTLS fingerprint simulation.
  - **Hysteria 2**: UDP-based high-throughput transport with Salamander obfuscation and independent upstream/downstream bandwidth rate limits.
  - **TUIC v5**: QUIC transport with BBR congestion control and 0-RTT handshake support.
  - **Shadowsocks 2022**: Full compliance with Blake3 AEAD specifications (`2022-blake3-aes-128-gcm`, `2022-blake3-aes-256-gcm`).
  - **Trojan & VMess**: Support for TCP, WebSocket, gRPC, and HTTPUpgrade transports with TLS.
- **Multi-Client Isolation**:
  - Multiple independent client credentials per listening port.
  - Per-user traffic tracking (Upload / Download / Total Limit), expiration dates, and IP concurrency limits.
  - Automatic quota exhaustion and expiration enforcement with atomic core reloading.
- **Dynamic Rule-Set & DNS Routing**:
  - Visual route management supporting `geosite` and `geoip` rule sets.
  - Granular outbound actions (`DIRECT`, `BLOCK`, `DNS-OUT`).
  - Split DNS configuration with dedicated Remote DoH and China Direct DNS resolvers.
- **Atomic Validation & Process Supervisor**:
  - Pre-flight syntax validation via `sing-box check` before applying configuration changes to prevent core crashes.
  - Built-in process supervisor with WebSocket real-time log streaming and crash recovery.
  - Clash API controller for live connection tracing and throughput monitoring.

---

## Protocol & Transport Matrix

| Protocol | Transports | Security / Camouflage | Multi-Client | Flow / Congestion |
| :--- | :--- | :--- | :---: | :--- |
| **VLESS** | TCP, WS, gRPC, HTTPUpgrade | Reality (X25519), TLS | ✅ | `xtls-rprx-vision` |
| **Hysteria 2** | UDP | Salamander Obfs, TLS | ✅ | BBR / Brutal, Rate Limit |
| **TUIC v5** | UDP (QUIC) | TLS | ✅ | BBR, 0-RTT Handshake |
| **Shadowsocks** | TCP, UDP | Blake3-AEAD (2022) | ✅ | Multiplex |
| **Trojan** | TCP, WS, gRPC | TLS | ✅ | - |
| **VMess** | TCP, WS, gRPC, HTTPUpgrade | TLS, None | ✅ | AlterID 0 |

---

## Quick Start

### 1. One-Click Linux Installation

Deploy SingUI on any modern Linux distribution (Debian, Ubuntu, CentOS, Alpine, Arch):

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Gingtin/SingUI/main/scripts/install.sh)
```

Access the panel at:
- **URL**: `http://<server-ip>:2096`
- **Default Username**: `admin`
- **Default Password**: `admin`

---

### 2. Docker Deployment

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

# Run SingUI
../singbox-ui -p 2096 -d data/singbox_ui.db
```

---

## Subscription Endpoints

SingUI serves subscription profiles via `/sub/:token`:

| Client Type | Query Parameter | Delivered Format |
| :--- | :--- | :--- |
| **Sing-box Client** | `?flag=sing-box` | Native Sing-box Client JSON |
| **Clash Meta / Mihomo** | `?flag=clash` | Clash Meta YAML (Proxies, Proxy Groups, Rule Providers) |
| **Universal Base64** | `?flag=base64` | Standard URI list (`vless://`, `hysteria2://`, etc.) |
| **Web Self-Service** | `/sub/view/:token` | Browser portal with quota meters and QR codes |

All responses include standard `Subscription-Userinfo` headers (`upload=...; download=...; total=...; expire=...`) for client bandwidth gauges.

---

## CLI Options

```bash
singbox-ui -h
  -p string
        Panel listening port (default: 2096)
  -d string
        SQLite database file path (default: "data/singbox_ui.db")
  -reset-admin
        Reset admin credentials to default (admin/admin)
  -v    Show version
```

---

## Feedback & Issues

Encountered a bug or have a suggestion? Please open an issue on the [GitHub Issues](https://github.com/Gingtin/SingUI/issues) page.

---

## Acknowledgments

- **[SagerNet/sing-box](https://github.com/SagerNet/sing-box)**: The universal proxy platform core.
- **[XTLS/Xray-core](https://github.com/XTLS/Xray-core)**: Pioneers in Reality and XTLS specifications.
- **[MetaCubeX/mihomo](https://github.com/MetaCubeX/mihomo)**: Rule-set and Clash Meta ecosystem standards.

---

## License

This project is licensed under the [MIT License](LICENSE).
