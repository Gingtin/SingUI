package models

import "time"

type Outbound struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Tag       string    `gorm:"size:64;uniqueIndex;not null" json:"tag"`
	Type      string    `gorm:"size:32;not null" json:"type"` // direct, block, dns, warp, wireguard, socks, http, custom
	Server    string    `gorm:"size:255" json:"server"`
	Port      int       `gorm:"default:0" json:"port"`
	Settings  string    `gorm:"type:text" json:"settings"` // Extra JSON parameters (e.g. WARP private key, peer IP, reserved, etc.)
	Enable    bool      `gorm:"default:true" json:"enable"`
	Remark    string    `gorm:"size:255" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
