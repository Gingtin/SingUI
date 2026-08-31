#!/usr/bin/env bash

#====================================================
# System: Debian / Ubuntu / CentOS / Alpine / Arch
# Description: Next-Generation One-Click Installer for SingUI
#====================================================

set -e

RED="\033[31m"
GREEN="\033[32m"
YELLOW="\033[33m"
BLUE="\033[36m"
PLAIN="\033[0m"

GITHUB_REPO="Gingtin/SingUI"
INSTALL_DIR="/usr/local/singbox-ui"
CONFIG_DIR="/usr/local/singbox-ui/config"
DATA_DIR="/usr/local/singbox-ui/data"

echo -e "${BLUE}====================================================${PLAIN}"
echo -e "${BLUE}      SingUI - Next-Gen Sing-box Management Panel    ${PLAIN}"
echo -e "${BLUE}====================================================${PLAIN}"

# Check root
if [ "$(id -u)" != "0" ]; then
    echo -e "${RED}错误: 必须使用 root 用户运行此脚本！${PLAIN}"
    exit 1
fi

# 1. Optimize Linux Kernel (BBR & UDP Buffers for Hysteria 2 / TUIC)
echo -e "${GREEN}[1/5] 检查并开启 Linux 原生 BBR 与 UDP 缓冲区调优...${PLAIN}"
sysctl_file="/etc/sysctl.conf"
if ! grep -q "net.core.default_qdisc=fq" "$sysctl_file"; then
    echo "net.core.default_qdisc=fq" >> "$sysctl_file"
fi
if ! grep -q "net.ipv4.tcp_congestion_control=bbr" "$sysctl_file"; then
    echo "net.ipv4.tcp_congestion_control=bbr" >> "$sysctl_file"
fi
if ! grep -q "net.core.rmem_max=8388608" "$sysctl_file"; then
    echo "net.core.rmem_max=8388608" >> "$sysctl_file"
    echo "net.core.wmem_max=8388608" >> "$sysctl_file"
fi
sysctl -p >/dev/null 2>&1 || true

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

echo -e "${GREEN}[2/5] 检查并安装 Sing-box 核心...${PLAIN}"
if [ ! -f "/usr/local/bin/sing-box" ]; then
    SINGBOX_VERSION=$(curl -s "https://api.github.com/repos/SagerNet/sing-box/releases/latest" | grep -Po '"tag_name": "\K.*?(?=")') || SINGBOX_VERSION="v1.9.7"
    SINGBOX_URL="https://github.com/SagerNet/sing-box/releases/download/${SINGBOX_VERSION}/sing-box-${SINGBOX_VERSION#v}-${SB_ARCH}.tar.gz"

    echo -e "${YELLOW}下载 Sing-box ${SINGBOX_VERSION} (${SB_ARCH})...${PLAIN}"
    curl -L -o /tmp/sing-box.tar.gz "$SINGBOX_URL"
    tar -xzf /tmp/sing-box.tar.gz -C /tmp/
    cp /tmp/sing-box-*/sing-box /usr/local/bin/sing-box
    chmod +x /usr/local/bin/sing-box
    rm -rf /tmp/sing-box*
fi

echo -e "${GREEN}[3/5] 获取 SingUI 面板程序...${PLAIN}"
if [ -f "./singbox-ui" ]; then
    echo -e "${YELLOW}使用本地编译的 singbox-ui 二进制文件...${PLAIN}"
    cp ./singbox-ui "$INSTALL_DIR/singbox-ui"
else
    echo -e "${YELLOW}正在从 GitHub (${GITHUB_REPO}) 获取最新发布版本...${PLAIN}"
    LATEST_TAG=$(curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep -Po '"tag_name": "\K.*?(?=")' || echo "")
    if [ -n "$LATEST_TAG" ]; then
        DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${LATEST_TAG}/singbox-ui-${UI_ARCH}.tar.gz"
        curl -L -o /tmp/singbox-ui.tar.gz "$DOWNLOAD_URL"
        tar -xzf /tmp/singbox-ui.tar.gz -C "$INSTALL_DIR/"
        mv "$INSTALL_DIR/singbox-ui-${UI_ARCH}" "$INSTALL_DIR/singbox-ui" 2>/dev/null || true
        rm -f /tmp/singbox-ui.tar.gz
    else
        echo -e "${YELLOW}提示: 未检测到远程 GitHub Release。若为本地源码编译部署，请运行 make 或 go build 后重新执行。${PLAIN}"
    fi
fi

if [ -f "$INSTALL_DIR/singbox-ui" ]; then
    chmod +x "$INSTALL_DIR/singbox-ui"
fi

echo -e "${GREEN}[4/5] 配置 Systemd 守护进程与快捷命令...${PLAIN}"
cat > /etc/systemd/system/singbox-ui.service <<EOF
[Unit]
Description=SingUI Management Panel
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

# Create CLI helper script
cat > /usr/local/bin/sing-ui <<'EOF'
#!/usr/bin/env bash
case "$1" in
    start)
        systemctl start singbox-ui
        echo "SingUI 已启动"
        ;;
    stop)
        systemctl stop singbox-ui
        echo "SingUI 已停止"
        ;;
    restart)
        systemctl restart singbox-ui
        echo "SingUI 已重启"
        ;;
    status)
        systemctl status singbox-ui
        ;;
    reset-admin)
        /usr/local/singbox-ui/singbox-ui -d /usr/local/singbox-ui/data/singbox_ui.db -reset-admin
        echo "管理员密码已重置为: admin / admin"
        ;;
    *)
        echo "用法: sing-ui {start|stop|restart|status|reset-admin}"
        ;;
esac
EOF
chmod +x /usr/local/bin/sing-ui

systemctl daemon-reload
systemctl enable singbox-ui
systemctl restart singbox-ui

echo -e "${GREEN}[5/5] 部署/更新完成！${PLAIN}"
IP=$(curl -s4 ifconfig.me || echo "服务器公网IP")
echo -e "${BLUE}====================================================${PLAIN}"
echo -e "${GREEN} SingUI 运行状态：${PLAIN}"
systemctl status singbox-ui --no-pager -l || true
echo -e "----------------------------------------------------"
echo -e " 面板地址: ${YELLOW}http://${IP}:2096${PLAIN}"
echo -e " 默认用户: ${YELLOW}admin${PLAIN}"
echo -e " 默认密码: ${YELLOW}admin${PLAIN}"
echo -e " 快捷指令: ${BLUE}sing-ui {start|stop|restart|status|reset-admin}${PLAIN}"
echo -e "${BLUE}====================================================${PLAIN}"
