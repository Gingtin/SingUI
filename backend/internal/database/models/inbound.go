package models

import (
	"time"
)

type Inbound struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"index" json:"user_id"`
	Tag            string    `gorm:"uniqueIndex;size:64;not null" json:"tag"`
	Protocol       string    `gorm:"size:32;not null" json:"protocol"` // vless, vmess, trojan, shadowsocks, hysteria2, tuic, wireguard, shadowtls
	Port           int       `gorm:"index;not null" json:"port"`
	Listen         string    `gorm:"size:64;default:'0.0.0.0'" json:"listen"`
	Network        string    `gorm:"size:32;default:'tcp'" json:"network"` // tcp, udp, ws, grpc, httpupgrade
	Security       string    `gorm:"size:32;default:'none'" json:"security"` // none, tls, reality
	Settings       string    `gorm:"type:text" json:"settings"`        // JSON protocol specific
	StreamSettings string    `gorm:"type:text" json:"stream_settings"` // JSON transport/TLS/Reality
	Sniffing       string    `gorm:"type:text" json:"sniffing"`        // JSON sniffing config
	Enable         bool      `gorm:"default:true" json:"enable"`
	Remark         string    `gorm:"size:255" json:"remark"`
	Clients        []Client  `gorm:"foreignKey:InboundID;constraint:OnDelete:CASCADE" json:"clients"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
