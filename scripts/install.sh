#!/usr/bin/env bash

#====================================================
# System Request: Debian / Ubuntu / CentOS / Alpine / Arch
# Description: One-click installer for Sing-box UI Panel
#====================================================

set -e

RED="\033[31m"
GREEN="\033[32m"
YELLOW="\033[33m"
PLAIN="\033[0m"

# 请将此处修改为你自己的 GitHub 仓库用户名和仓库名
GITHUB_REPO="Gingtin/SingUI"

INSTALL_DIR="/usr/local/singbox-ui"
CONFIG_DIR="/usr/local/singbox-ui/config"
DATA_DIR="/usr/local/singbox-ui/data"

echo -e "${GREEN}=======================================${PLAIN}"
echo -e "${GREEN}      Sing-box UI 一键安装/更新脚本      ${PLAIN}"
echo -e "${GREEN}=======================================${PLAIN}"

# Check root
if [ "$(id -u)" != "0" ]; then
    echo -e "${RED}错误: 必须使用 root 用户运行此脚本！${PLAIN}"
    exit 1
fi

# Detect Architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64)
        SB_ARCH="linux-amd64"
        UI_ARCH="linux-amd64"
        ;;
    aarch64|arm64)
        SB_ARCH="linux-arm64"
        UI_ARCH="linux-arm64"
        ;;
    *)
        echo -e "${RED}不支持的系统架构: $ARCH${PLAIN}"
        exit 1
        ;;
esac

# Create dirs
mkdir -p "$INSTALL_DIR"
mkdir -p "$CONFIG_DIR"
mkdir -p "$DATA_DIR"

echo -e "${GREEN}[1/4] 检查并下载 Sing-box 核心...${PLAIN}"
if [ ! -f "/usr/local/bin/sing-box" ]; then
    SINGBOX_VERSION=$(curl -s "https://api.github.com/repos/SagerNet/sing-box/releases/latest" | grep -Po '"tag_name": "\K.*?(?=")') || SINGBOX_VERSION="v1.9.7"
    SINGBOX_URL="https://github.com/SagerNet/sing-box/releases/download/${SINGBOX_VERSION}/sing-box-${SINGBOX_VERSION#v}-${SB_ARCH}.tar.gz"

    curl -L -o /tmp/sing-box.tar.gz "$SINGBOX_URL"
    tar -xzf /tmp/sing-box.tar.gz -C /tmp/
    cp /tmp/sing-box-*/sing-box /usr/local/bin/sing-box
    chmod +x /usr/local/bin/sing-box
    rm -rf /tmp/sing-box*
fi

echo -e "${GREEN}[2/4] 获取 Sing-box UI 面板程序...${PLAIN}"
if [ -f "./singbox-ui" ]; then
    echo -e "${YELLOW}使用本地编译的 singbox-ui 二进制文件...${PLAIN}"
    cp ./singbox-ui "$INSTALL_DIR/singbox-ui"
else
    echo -e "${YELLOW}正在从 GitHub (${GITHUB_REPO}) 下载最新发行版...${PLAIN}"
    LATEST_TAG=$(curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep -Po '"tag_name": "\K.*?(?=")' || echo "")
    if [ -n "$LATEST_TAG" ]; then
        DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${LATEST_TAG}/singbox-ui-${UI_ARCH}.tar.gz"
        curl -L -o /tmp/singbox-ui.tar.gz "$DOWNLOAD_URL"
        tar -xzf /tmp/singbox-ui.tar.gz -C "$INSTALL_DIR/"
        mv "$INSTALL_DIR/singbox-ui-${UI_ARCH}" "$INSTALL_DIR/singbox-ui" 2>/dev/null || true
        rm -f /tmp/singbox-ui.tar.gz
    else
        echo -e "${YELLOW}未检测到远程 GitHub Release，如果这是开发环境，请将编译好的 singbox-ui 二进制放到当前目录运行。${PLAIN}"
    fi
fi

if [ -f "$INSTALL_DIR/singbox-ui" ]; then
    chmod +x "$INSTALL_DIR/singbox-ui"
fi

echo -e "${GREEN}[3/4] 配置 Systemd 守护进程...${PLAIN}"
cat > /etc/systemd/system/singbox-ui.service <<EOF
[Unit]
Description=Sing-box UI Management Panel
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=${INSTALL_DIR}
ExecStart=${INSTALL_DIR}/singbox-ui -d ${DATA_DIR}/singbox_ui.db -p 2096
Restart=on-failure
RestartSec=5s
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable singbox-ui
systemctl restart singbox-ui

echo -e "${GREEN}[4/4] 部署/更新完成！${PLAIN}"
IP=$(curl -s4 ifconfig.me || echo "你的服务器IP")
echo -e "${GREEN}===============================================${PLAIN}"
echo -e "${GREEN} Sing-box UI 服务运行状态：${PLAIN}"
systemctl status singbox-ui --no-pager -l || true
echo -e "-----------------------------------------------"
echo -e " 访问地址: ${YELLOW}http://${IP}:2096${PLAIN}"
echo -e " 默认用户: ${YELLOW}admin${PLAIN}"
echo -e " 默认密码: ${YELLOW}admin${PLAIN}"
echo -e " 管理指令: ${YELLOW}systemctl restart|status|stop singbox-ui${PLAIN}"
echo -e "${GREEN}===============================================${PLAIN}"
