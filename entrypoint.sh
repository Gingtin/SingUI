#!/usr/bin/env bash

set -e

DATA_DIR="/usr/local/singbox-ui/data"
CONFIG_DIR="/usr/local/singbox-ui/config"

mkdir -p "$DATA_DIR"
mkdir -p "$CONFIG_DIR"

echo "[SingUI Docker] Starting SingUI in Container..."
exec /usr/local/singbox-ui/singbox-ui -d "${DATA_DIR}/singbox_ui.db" -p 2096
