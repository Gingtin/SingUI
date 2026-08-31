# Changelog

All notable changes to **SingUI** will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2026-08-31

### Added
- **Core Engine**: Complete integration with Sing-box 1.9+ / 1.10+ universal proxy core.
- **Protocol Suite**:
  - VLESS with Reality (`X25519` keypair generator, `ShortID`, `xtls-rprx-vision` flow), TCP, WebSocket, gRPC, HTTPUpgrade.
  - Hysteria 2 with Salamander obfuscation and bandwidth limits.
  - TUIC v5 with QUIC / BBR congestion control.
  - Shadowsocks 2022 specifications (`2022-blake3-aes-128-gcm`, etc.).
  - Trojan & VMess transports.
- **Multi-Client Architecture**: Support for single inbound port with multiple independent users, per-client traffic quotas, and expiration timers.
- **Traffic & Limits Enforcement**: Automatic background quota check and core hot-reloading on expiration/quota exhaustion.
- **Universal Subscriptions**:
  - Sing-box Official Client JSON format (`flag=sing-box`).
  - Clash Meta / Mihomo YAML format (`flag=clash`).
  - Standard Base64 / URI node links (`vless://`, `hysteria2://`, `tuic://`, `ss://`, `trojan://`).
  - Self-service web portal with visual progress bar and QR code importer.
- **Dashboard & Diagnostics**: Real-time system gauges (CPU, RAM, Disk), ECharts network rate line graphs, and WebSocket live log streaming.
- **Operations & Security**: Telegram Bot notification alerts, SQLite one-click backup download/restore, JWT authentication, and 2FA TOTP support.
- **Packaging & CI/CD**: Single-binary Go `embed.FS` packaging and GitHub Actions multi-arch cross-compilation pipeline.
