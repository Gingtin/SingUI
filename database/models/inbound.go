package models

import "time"

type Inbound struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Tag            string    `gorm:"size:64;uniqueIndex;not null" json:"tag"`
	Protocol       string    `gorm:"size:32;not null" json:"protocol"` // vless, vmess, trojan, shadowsocks, hysteria2, tuic, anytls, shadowtls, direct, socks, http
	Port           int       `gorm:"not null" json:"port"`
	Listen         string    `gorm:"size:64;default:'0.0.0.0'" json:"listen"`
	Network        string    `gorm:"size:32;default:'tcp'" json:"network"`    // tcp, udp, ws, grpc, httpupgrade
	Security       string    `gorm:"size:32;default:'none'" json:"security"`  // none, tls, reality
	Settings       string    `gorm:"type:text" json:"settings"`               // JSON for protocol-specific settings
	StreamSettings string    `gorm:"type:text" json:"stream_settings"`        // JSON for TLS/Reality/Transport/uTLS settings
	Sniffing       string    `gorm:"type:text" json:"sniffing"`               // JSON for packet sniffing settings
	Total          int64     `gorm:"default:0" json:"total"`                  // Total quota in bytes (0 = unlimited)
	ExpiryTime     int64     `gorm:"default:0" json:"expiry_time"`             // Expiration timestamp in ms
	Enable         bool      `gorm:"default:true" json:"enable"`
	Remark         string    `gorm:"size:255" json:"remark"`
	Clients        []Client  `gorm:"foreignKey:InboundID;constraint:OnDelete:CASCADE" json:"clients"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
