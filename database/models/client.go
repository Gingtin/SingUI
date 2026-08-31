package models

import (
	"time"
)

type Client struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	InboundID  uint      `gorm:"index;not null" json:"inbound_id"`
	UUID       string    `gorm:"index;size:64" json:"uuid"`
	Password   string    `gorm:"size:255" json:"password"`
	Email      string    `gorm:"index;size:128" json:"email"` // Remark / Username identifier
	Flow       string    `gorm:"size:64" json:"flow"`         // xtls-rprx-vision
	SubToken   string    `gorm:"uniqueIndex;size:64" json:"sub_token"`
	Up         int64     `gorm:"default:0" json:"up"`
	Down       int64     `gorm:"default:0" json:"down"`
	Total      int64     `gorm:"default:0" json:"total"`           // 0 = unlimited bytes
	ExpiryTime int64     `gorm:"default:0" json:"expiry_time"`     // unix timestamp in ms, 0 = unlimited
	Enable     bool      `gorm:"default:true" json:"enable"`
	LimitIP    int       `gorm:"default:0" json:"limit_ip"`        // 0 = unlimited
	ResetDay   int       `gorm:"default:0" json:"reset_day"`       // Day of month to reset traffic, 0 = disabled
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
