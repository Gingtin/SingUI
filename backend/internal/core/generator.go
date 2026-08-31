package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/singbox-ui/singbox-ui/internal/database/models"
)

type SingboxConfig struct {
	Log          LogConfig          `json:"log,omitempty"`
	DNS          DNSConfig          `json:"dns,omitempty"`
	Inbounds     []interface{}      `json:"inbounds"`
	Outbounds    []interface{}      `json:"outbounds"`
	Route        RouteConfig        `json:"route"`
	Experimental ExperimentalConfig `json:"experimental"`
}

type LogConfig struct {
	Level     string `json:"level"`
	Timestamp bool   `json:"timestamp"`
}

type DNSConfig struct {
	Servers []DNSServer `json:"servers"`
	Rules   []DNSRule   `json:"rules,omitempty"`
}

type DNSServer struct {
	Tag     string `json:"tag"`
	Address string `json:"address"`
	Detour  string `json:"detour,omitempty"`
}

type DNSRule struct {
	Outbound []string `json:"outbound,omitempty"`
	Server   string   `json:"server"`
}

type RouteConfig struct {
	Rules               []RouteRule `json:"rules"`
	AutoDetectInterface bool        `json:"auto_detect_interface"`
	Final               string      `json:"final"`
}

type RouteRule struct {
	Inbound  []string `json:"inbound,omitempty"`
	Protocol []string `json:"protocol,omitempty"`
	Outbound string   `json:"outbound"`
	RuleSet  []string `json:"rule_set,omitempty"`
}

type ExperimentalConfig struct {
	ClashAPI  *ClashAPIConfig  `json:"clash_api,omitempty"`
	CacheFile *CacheFileConfig `json:"cache_file,omitempty"`
	V2RayAPI  *V2RayAPIConfig  `json:"v2ray_api,omitempty"`
}

type ClashAPIConfig struct {
	ExternalController string `json:"external_controller"`
	Secret             string `json:"secret"`
	DefaultMode        string `json:"default_mode"`
}

type CacheFileConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

type V2RayAPIConfig struct {
	Listen string       `json:"listen,omitempty"`
	Stats  *StatsConfig `json:"stats,omitempty"`
}

type StatsConfig struct {
	Enabled  bool     `json:"enabled"`
	Inbounds []string `json:"inbounds,omitempty"`
	Users    []string `json:"users,omitempty"`
}

// GenerateConfig assembles the full Sing-box JSON config from database models
func GenerateConfig(inbounds []models.Inbound, clashPort, clashSecret string) (*SingboxConfig, error) {
	cfg := &SingboxConfig{
		Log: LogConfig{
			Level:     "info",
			Timestamp: true,
		},
		DNS: DNSConfig{
			Servers: []DNSServer{
				{Tag: "local-dns", Address: "local", Detour: "direct"},
				{Tag: "remote-dns", Address: "https://1.1.1.1/dns-query", Detour: "direct"},
			},
		},
		Inbounds: make([]interface{}, 0),
		Outbounds: []interface{}{
			map[string]interface{}{"type": "direct", "tag": "direct"},
			map[string]interface{}{"type": "block", "tag": "block"},
			map[string]interface{}{"type": "dns", "tag": "dns-out"},
		},
		Route: RouteConfig{
			AutoDetectInterface: true,
			Final:               "direct",
			Rules: []RouteRule{
				{Protocol: []string{"dns"}, Outbound: "dns-out"},
			},
		},
		Experimental: ExperimentalConfig{
			ClashAPI: &ClashAPIConfig{
				ExternalController: fmt.Sprintf("127.0.0.1:%s", clashPort),
				Secret:             clashSecret,
				DefaultMode:        "Rule",
			},
			CacheFile: &CacheFileConfig{
				Enabled: true,
				Path:    "cache.db",
			},
			V2RayAPI: &V2RayAPIConfig{
				Stats: &StatsConfig{
					Enabled: true,
				},
			},
		},
	}

	for _, in := range inbounds {
		if !in.Enable {
			continue
		}

		inboundMap := map[string]interface{}{
			"type":        in.Protocol,
			"tag":         in.Tag,
			"listen":      in.Listen,
			"listen_port": in.Port,
		}

		// Protocol specific configurations
		switch in.Protocol {
		case "vless":
			users := make([]map[string]interface{}, 0)
			for _, c := range in.Clients {
				if !c.Enable {
					continue
				}
				userObj := map[string]interface{}{
					"name": c.Email,
					"uuid": c.UUID,
				}
				if c.Flow != "" {
					userObj["flow"] = c.Flow
				}
				users = append(users, userObj)
			}
			inboundMap["users"] = users

			// TLS / Reality / Stream
			parseStreamSettings(in, inboundMap)

		case "vmess":
			users := make([]map[string]interface{}, 0)
			for _, c := range in.Clients {
				if !c.Enable {
					continue
				}
				users = append(users, map[string]interface{}{
					"name":     c.Email,
					"uuid":     c.UUID,
					"alter_id": 0,
				})
			}
			inboundMap["users"] = users
			parseStreamSettings(in, inboundMap)

		case "trojan":
			users := make([]map[string]interface{}, 0)
			for _, c := range in.Clients {
				if !c.Enable {
					continue
				}
				users = append(users, map[string]interface{}{
					"name":     c.Email,
					"password": c.Password,
				})
			}
			inboundMap["users"] = users
			parseStreamSettings(in, inboundMap)

		case "shadowsocks":
			var ssSettings struct {
				Method   string `json:"method"`
				Password string `json:"password"`
			}
			_ = json.Unmarshal([]byte(in.Settings), &ssSettings)
			if ssSettings.Method == "" {
				ssSettings.Method = "2022-blake3-aes-128-gcm"
			}
			inboundMap["method"] = ssSettings.Method
			inboundMap["password"] = ssSettings.Password
			parseStreamSettings(in, inboundMap)

		case "hysteria2":
			users := make([]map[string]interface{}, 0)
			for _, c := range in.Clients {
				if !c.Enable {
					continue
				}
				users = append(users, map[string]interface{}{
					"name":     c.Email,
					"password": c.Password,
				})
			}
			inboundMap["users"] = users

			var hySettings struct {
				UpMbps     int    `json:"up_mbps"`
				DownMbps   int    `json:"down_mbps"`
				ObfsType   string `json:"obfs_type"`
				ObfsPass   string `json:"obfs_password"`
			}
			_ = json.Unmarshal([]byte(in.Settings), &hySettings)
			if hySettings.UpMbps > 0 {
				inboundMap["up_mbps"] = hySettings.UpMbps
			}
			if hySettings.DownMbps > 0 {
				inboundMap["down_mbps"] = hySettings.DownMbps
			}
			if hySettings.ObfsType == "salamander" && hySettings.ObfsPass != "" {
				inboundMap["obfs"] = map[string]interface{}{
					"type":     "salamander",
					"password": hySettings.ObfsPass,
				}
			}
			parseStreamSettings(in, inboundMap)

		case "tuic":
			users := make([]map[string]interface{}, 0)
			for _, c := range in.Clients {
				if !c.Enable {
					continue
				}
				users = append(users, map[string]interface{}{
					"name":     c.Email,
					"uuid":     c.UUID,
					"password": c.Password,
				})
			}
			inboundMap["users"] = users
			inboundMap["congestion_control"] = "bbr"
			parseStreamSettings(in, inboundMap)
		}

		cfg.Inbounds = append(cfg.Inbounds, inboundMap)
	}

	return cfg, nil
}

func parseStreamSettings(in models.Inbound, inboundMap map[string]interface{}) {
	var stream struct {
		Network         string            `json:"network"`
		Security        string            `json:"security"`
		TLS             *TLSConfig        `json:"tls"`
		Reality         *RealityConfig    `json:"reality"`
		Transport       *TransportConfig  `json:"transport"`
		Multiplex       *MultiplexConfig  `json:"multiplex"`
	}
	_ = json.Unmarshal([]byte(in.StreamSettings), &stream)

	if in.Security == "reality" || stream.Security == "reality" {
		if stream.Reality != nil {
			realityMap := map[string]interface{}{
				"enabled":              true,
				"handshake":            map[string]interface{}{"server": stream.Reality.ServerNames[0], "server_port": 443},
				"private_key":          stream.Reality.PrivateKey,
				"short_id":             stream.Reality.ShortIds,
				"max_time_difference": "1m",
			}
			inboundMap["tls"] = realityMap
		}
	} else if in.Security == "tls" || stream.Security == "tls" {
		if stream.TLS != nil {
			tlsMap := map[string]interface{}{
				"enabled":     true,
				"server_name": stream.TLS.ServerName,
			}
			if stream.TLS.CertPath != "" && stream.TLS.KeyPath != "" {
				tlsMap["certificate_path"] = stream.TLS.CertPath
				tlsMap["key_path"] = stream.TLS.KeyPath
			}
			inboundMap["tls"] = tlsMap
		}
	}

	if stream.Transport != nil && stream.Transport.Type != "" {
		inboundMap["transport"] = map[string]interface{}{
			"type":         stream.Transport.Type,
			"path":         stream.Transport.Path,
			"service_name": stream.Transport.ServiceName,
		}
	}
}

type TLSConfig struct {
	ServerName string `json:"server_name"`
	CertPath   string `json:"cert_path"`
	KeyPath    string `json:"key_path"`
}

type RealityConfig struct {
	ServerNames []string `json:"server_names"`
	PrivateKey  string   `json:"private_key"`
	PublicKey   string   `json:"public_key"`
	ShortIds    []string `json:"short_ids"`
}

type TransportConfig struct {
	Type        string `json:"type"` // ws, grpc, httpupgrade
	Path        string `json:"path"`
	ServiceName string `json:"service_name"`
}

type MultiplexConfig struct {
	Enabled bool `json:"enabled"`
}

// WriteConfigToFile serializes config and writes it to target path
func WriteConfigToFile(cfg *SingboxConfig, filePath string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
