package models

import (
	"time"
)

type RoutingRule struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Tag         string    `gorm:"size:64" json:"tag"`
	Type        string    `gorm:"size:32;default:'default'" json:"type"` // default, custom
	Outbound    string    `gorm:"size:32;not null" json:"outbound"`      // direct, block, proxy
	Domain      string    `gorm:"type:text" json:"domain"`               // JSON array of domain/geosite rules
	IP          string    `gorm:"type:text" json:"ip"`                   // JSON array of IP/geoip rules
	Protocol    string    `gorm:"size:64" json:"protocol"`               // http, tls, dns, bittorrent, etc.
	Port        string    `gorm:"size:128" json:"port"`                  // port ranges
	Network     string    `gorm:"size:32" json:"network"`                // tcp, udp
	RuleSet     string    `gorm:"type:text" json:"rule_set"`             // remote rule_set tags
	Enable      bool      `gorm:"default:true" json:"enable"`
	Order       int       `gorm:"default:0" json:"order"`
	Remark      string    `gorm:"size:255" json:"remark"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DNSSettings struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	LocalDNS       string    `gorm:"size:255;default:'local'" json:"local_dns"`
	RemoteDNS      string    `gorm:"size:255;default:'https://1.1.1.1/dns-query'" json:"remote_dns"`
	ChinaDNS       string    `gorm:"size:255;default:'https://223.5.5.5/dns-query'" json:"china_dns"`
	EnableFakeIP   bool      `gorm:"default:false" json:"enable_fakeip"`
	FakeIPInet4    string    `gorm:"size:64;default:'198.18.0.0/15'" json:"fakeip_inet4"`
	FakeIPInet6    string    `gorm:"size:64;default:'fc00::/18'" json:"fakeip_inet6"`
	Strategy       string    `gorm:"size:32;default:'prefer_ipv4'" json:"strategy"`
	CustomRules    string    `gorm:"type:text" json:"custom_rules"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
