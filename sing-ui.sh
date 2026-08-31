#!/usr/bin/env bash

#====================================================
# SingUI - Interactive Terminal Management Tool (TUI)
#====================================================

RED="\033[31m"
GREEN="\033[32m"
YELLOW="\033[33m"
BLUE="\033[36m"
PURPLE="\033[35m"
PLAIN="\033[0m"

INSTALL_DIR="/usr/local/singbox-ui"
CONFIG_FILE="${INSTALL_DIR}/config/singbox_config.json"
DATA_DB="${INSTALL_DIR}/data/singbox_ui.db"

# Check root
if [ "$(id -u)" != "0" ]; then
    echo -e "${RED}错误: 必须使用 root 用户运行此脚本！${PLAIN}"
    exit 1
fi

show_menu() {
    clear
    echo -e "${BLUE}====================================================${PLAIN}"
    echo -e "${BLUE}        SingUI - 现代化 Sing-box 管理面板           ${PLAIN}"
    echo -e "${BLUE}====================================================${PLAIN}"
    echo -e " ${GREEN}1.${PLAIN}  启动 SingUI 面板"
    echo -e " ${GREEN}2.${PLAIN}  停止 SingUI 面板"
    echo -e " ${GREEN}3.${PLAIN}  重启 SingUI 面板"
    echo -e " ${GREEN}4.${PLAIN}  查看运行状态与系统日志"
    echo -e "----------------------------------------------------"
    echo -e " ${GREEN}5.${PLAIN}  重置管理员账号与密码为默认值 (admin/admin)"
    echo -e " ${GREEN}6.${PLAIN}  修改面板监听端口"
    echo -e " ${GREEN}7.${PLAIN}  申请/更新 SSL 证书 (ACME Cloudflare / Let's Encrypt)"
    echo -e " ${GREEN}8.${PLAIN}  一键开启 Linux 原生 BBR 加速 & UDP 调优"
    echo -e " ${GREEN}9.${PLAIN}  备份/导出 SQLite 数据库"
    echo -e "----------------------------------------------------"
    echo -e " ${GREEN}10.${PLAIN} 检查并升级 SingUI 至最新版本"
    echo -e " ${GREEN}11.${PLAIN} 卸载 SingUI 面板"
    echo -e " ${GREEN}0.${PLAIN}  退出脚本"
    echo -e "${BLUE}====================================================${PLAIN}"
    
    # Show active status
    if systemctl is-active --quiet singbox-ui; then
        echo -e " 服务状态: ${GREEN}● 正在运行 (Running)${PLAIN}"
    else
        echo -e " 服务状态: ${RED}○ 已停止 (Stopped)${PLAIN}"
    fi
    echo -e "----------------------------------------------------"
    read -rp " 请输入选择 [0-11]: " choice

    case "$choice" in
        1)
            systemctl start singbox-ui
            echo -e "${GREEN}SingUI 面板已启动！${PLAIN}"
            sleep 2
            show_menu
            ;;
        2)
            systemctl stop singbox-ui
            echo -e "${YELLOW}SingUI 面板已停止！${PLAIN}"
            sleep 2
            show_menu
            ;;
        3)
            systemctl restart singbox-ui
            echo -e "${GREEN}SingUI 面板已重启！${PLAIN}"
            sleep 2
            show_menu
            ;;
        4)
            echo -e "${GREEN}正在获取实时日志 (按 Ctrl+C 退出)...${PLAIN}"
            journalctl -u singbox-ui -f -n 50
            ;;
        5)
            if [ -f "${INSTALL_DIR}/singbox-ui" ]; then
                "${INSTALL_DIR}/singbox-ui" -d "${DATA_DB}" -reset-admin
                echo -e "${GREEN}管理员密码已成功重置为: admin / admin${PLAIN}"
            else
                echo -e "${RED}未找到 singbox-ui 可执行程序！${PLAIN}"
            fi
            read -rp "按回车键返回菜单..."
            show_menu
            ;;
        6)
            read -rp "请输入新的面板监听端口 (1-65535): " new_port
            if [[ "$new_port" =~ ^[0-9]+$ ]] && [ "$new_port" -ge 1 ] && [ "$new_port" -le 65535 ]; then
                sed -i "s/-p [0-9]*/-p ${new_port}/g" /etc/systemd/system/singbox-ui.service
                systemctl daemon-reload
                systemctl restart singbox-ui
                echo -e "${GREEN}面板端口已修改为: ${new_port} 并已重启生效！${PLAIN}"
            else
                echo -e "${RED}输入的端口不合法！${PLAIN}"
            fi
            read -rp "按回车键返回菜单..."
            show_menu
            ;;
        7)
            echo -e "${GREEN}=== ACME 证书申请向导 ===${PLAIN}"
            read -rp "请输入你的域名 (如 node.example.com): " cert_domain
            read -rp "请输入你的邮箱 (如 admin@example.com): " cert_email
            if [ -n "$cert_domain" ]; then
                if ! command -v ~/.acme.sh/acme.sh &> /dev/null; then
                    echo -e "${YELLOW}正在安装 acme.sh 证书工具...${PLAIN}"
                    curl https://get.acme.sh | sh -s email="$cert_email"
                fi
                ~/.acme.sh/acme.sh --issue -d "$cert_domain" --standalone
                mkdir -p "${INSTALL_DIR}/cert"
                ~/.acme.sh/acme.sh --install-cert -d "$cert_domain" \
                    --key-file "${INSTALL_DIR}/cert/private.key" \
                    --fullchain-file "${INSTALL_DIR}/cert/cert.crt"
                echo -e "${GREEN}证书申请成功！保存在 ${INSTALL_DIR}/cert/${PLAIN}"
            fi
            read -rp "按回车键返回菜单..."
            show_menu
            ;;
        8)
            echo "net.core.default_qdisc=fq" >> /etc/sysctl.conf
            echo "net.ipv4.tcp_congestion_control=bbr" >> /etc/sysctl.conf
            echo "net.core.rmem_max=8388608" >> /etc/sysctl.conf
            echo "net.core.wmem_max=8388608" >> /etc/sysctl.conf
            sysctl -p
            echo -e "${GREEN}BBR 与 UDP 调优参数已成功启用！${PLAIN}"
            read -rp "按回车键返回菜单..."
            show_menu
            ;;
        9)
            backup_file="/tmp/singbox_ui_backup_$(date +%Y%m%d_%H%M%S).db"
            if [ -f "$DATA_DB" ]; then
                cp "$DATA_DB" "$backup_file"
                echo -e "${GREEN}数据库已成功备份至: ${backup_file}${PLAIN}"
            else
                echo -e "${RED}数据库文件不存在！${PLAIN}"
            fi
            read -rp "按回车键返回菜单..."
            show_menu
            ;;
        10)
            echo -e "${YELLOW}正在拉取最新安装脚本更新...${PLAIN}"
            bash <(curl -Ls https://raw.githubusercontent.com/Gingtin/SingUI/main/scripts/install.sh)
            ;;
        11)
            read -rp "确定要彻底卸载 SingUI 吗？(y/n): " confirm
            if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
                systemctl stop singbox-ui || true
                systemctl disable singbox-ui || true
                rm -f /etc/systemd/system/singbox-ui.service
                rm -f /usr/local/bin/sing-ui
                systemctl daemon-reload
                rm -rf "$INSTALL_DIR"
                echo -e "${GREEN}SingUI 面板已彻底卸载！${PLAIN}"
                exit 0
            fi
            show_menu
            ;;
        0)
            exit 0
            ;;
        *)
            echo -e "${RED}无效的选择！${PLAIN}"
            sleep 1
            show_menu
            ;;
    esac
}

show_menu
EOF
