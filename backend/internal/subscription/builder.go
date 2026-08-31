package subscription

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/singbox-ui/singbox-ui/internal/core"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
	"gopkg.in/yaml.v3"
)

type SubBuilder struct {
	serverHost string
	client     *models.Client
	inbounds   []models.Inbound
}

func NewSubBuilder(serverHost string, client *models.Client, inbounds []models.Inbound) *SubBuilder {
	return &SubBuilder{
		serverHost: serverHost,
		client:     client,
		inbounds:   inbounds,
	}
}

// BuildBase64Links generates standard base64 encoded node links (vless://, vmess://, hysteria2://, tuic://, ss://, trojan://)
func (sb *SubBuilder) BuildBase64Links() string {
	var links []string

	for _, in := range sb.inbounds {
		if !in.Enable {
			continue
		}

		link := sb.buildSingleLink(in)
		if link != "" {
			links = append(links, link)
		}
	}

	raw := strings.Join(links, "\n")
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func (sb *SubBuilder) buildSingleLink(in models.Inbound) string {
	host := sb.serverHost
	if host == "" {
		host = "127.0.0.1"
	}
	remark := url.QueryEscape(fmt.Sprintf("%s-%s", in.Tag, sb.client.Email))

	var stream struct {
		Network   string                 `json:"network"`
		Security  string                 `json:"security"`
		TLS       *core.TLSConfig        `json:"tls"`
		Reality   *core.RealityConfig    `json:"reality"`
		Transport *core.TransportConfig  `json:"transport"`
	}
	_ = json.Unmarshal([]byte(in.StreamSettings), &stream)

	switch in.Protocol {
	case "vless":
		// vless://uuid@host:port?security=reality&sni=...&pbk=...&sid=...&type=tcp&flow=xtls-rprx-vision#remark
		query := url.Values{}
		query.Set("type", in.Network)
		if stream.Transport != nil && stream.Transport.Type == "ws" {
			query.Set("type", "ws")
			query.Set("path", stream.Transport.Path)
		} else if stream.Transport != nil && stream.Transport.Type == "grpc" {
			query.Set("type", "grpc")
			query.Set("serviceName", stream.Transport.ServiceName)
		}

		if in.Security == "reality" && stream.Reality != nil {
			query.Set("security", "reality")
			if len(stream.Reality.ServerNames) > 0 {
				query.Set("sni", stream.Reality.ServerNames[0])
			}
			query.Set("pbk", stream.Reality.PublicKey)
			if len(stream.Reality.ShortIds) > 0 {
				query.Set("sid", stream.Reality.ShortIds[0])
			}
			if sb.client.Flow != "" {
				query.Set("flow", sb.client.Flow)
			}
		} else if in.Security == "tls" {
			query.Set("security", "tls")
			if stream.TLS != nil && stream.TLS.ServerName != "" {
				query.Set("sni", stream.TLS.ServerName)
			}
		}

		return fmt.Sprintf("vless://%s@%s:%d?%s#%s", sb.client.UUID, host, in.Port, query.Encode(), remark)

	case "vmess":
		vmessObj := map[string]interface{}{
			"v":    "2",
			"ps":   fmt.Sprintf("%s-%s", in.Tag, sb.client.Email),
			"add":  host,
			"port": in.Port,
			"id":   sb.client.UUID,
			"aid":  0,
			"net":  in.Network,
			"type": "none",
			"tls":  in.Security,
		}
		if stream.Transport != nil && stream.Transport.Type == "ws" {
			vmessObj["path"] = stream.Transport.Path
		}
		data, _ := json.Marshal(vmessObj)
		return "vmess://" + base64.StdEncoding.EncodeToString(data)

	case "trojan":
		query := url.Values{}
		if stream.TLS != nil && stream.TLS.ServerName != "" {
			query.Set("sni", stream.TLS.ServerName)
		}
		return fmt.Sprintf("trojan://%s@%s:%d?%s#%s", sb.client.Password, host, in.Port, query.Encode(), remark)

	case "hysteria2":
		// hysteria2://password@host:port/?sni=...#remark
		query := url.Values{}
		if stream.TLS != nil && stream.TLS.ServerName != "" {
			query.Set("sni", stream.TLS.ServerName)
		}
		var hySettings struct {
			ObfsType string `json:"obfs_type"`
			ObfsPass string `json:"obfs_password"`
		}
		_ = json.Unmarshal([]byte(in.Settings), &hySettings)
		if hySettings.ObfsType == "salamander" && hySettings.ObfsPass != "" {
			query.Set("obfs", "salamander")
			query.Set("obfs-password", hySettings.ObfsPass)
		}
		return fmt.Sprintf("hysteria2://%s@%s:%d/?%s#%s", sb.client.Password, host, in.Port, query.Encode(), remark)

	case "tuic":
		query := url.Values{}
		query.Set("congestion_control", "bbr")
		if stream.TLS != nil && stream.TLS.ServerName != "" {
			query.Set("sni", stream.TLS.ServerName)
		}
		return fmt.Sprintf("tuic://%s:%s@%s:%d/?%s#%s", sb.client.UUID, sb.client.Password, host, in.Port, query.Encode(), remark)

	case "shadowsocks":
		var ssSettings struct {
			Method   string `json:"method"`
			Password string `json:"password"`
		}
		_ = json.Unmarshal([]byte(in.Settings), &ssSettings)
		userInfo := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", ssSettings.Method, ssSettings.Password)))
		return fmt.Sprintf("ss://%s@%s:%d#%s", userInfo, host, in.Port, remark)
	}

	return ""
}

// BuildClashMeta generates Mihomo / Clash Meta YAML configuration
func (sb *SubBuilder) BuildClashMeta() (string, error) {
	proxies := make([]map[string]interface{}, 0)
	proxyNames := make([]string, 0)

	host := sb.serverHost
	if host == "" {
		host = "127.0.0.1"
	}

	for _, in := range sb.inbounds {
		if !in.Enable {
			continue
		}

		name := fmt.Sprintf("%s - %s", in.Tag, in.Protocol)
		var stream struct {
			Network   string                 `json:"network"`
			Security  string                 `json:"security"`
			TLS       *core.TLSConfig        `json:"tls"`
			Reality   *core.RealityConfig    `json:"reality"`
			Transport *core.TransportConfig  `json:"transport"`
		}
		_ = json.Unmarshal([]byte(in.StreamSettings), &stream)

		switch in.Protocol {
		case "vless":
			p := map[string]interface{}{
				"name":    name,
				"type":    "vless",
				"server":  host,
				"port":    in.Port,
				"uuid":    sb.client.UUID,
				"network": in.Network,
				"udp":     true,
			}
			if sb.client.Flow != "" {
				p["flow"] = sb.client.Flow
			}
			if in.Security == "reality" && stream.Reality != nil {
				p["tls"] = true
				if len(stream.Reality.ServerNames) > 0 {
					p["servername"] = stream.Reality.ServerNames[0]
				}
				realityOpts := map[string]interface{}{
					"public-key": stream.Reality.PublicKey,
				}
				if len(stream.Reality.ShortIds) > 0 {
					realityOpts["short-id"] = stream.Reality.ShortIds[0]
				}
				p["reality-opts"] = realityOpts
				p["client-fingerprint"] = "chrome"
			}
			proxies = append(proxies, p)
			proxyNames = append(proxyNames, name)

		case "hysteria2":
			p := map[string]interface{}{
				"name":     name,
				"type":     "hysteria2",
				"server":   host,
				"port":     in.Port,
				"password": sb.client.Password,
				"tls":      true,
			}
			if stream.TLS != nil && stream.TLS.ServerName != "" {
				p["sni"] = stream.TLS.ServerName
			}
			var hySettings struct {
				ObfsType string `json:"obfs_type"`
				ObfsPass string `json:"obfs_password"`
			}
			_ = json.Unmarshal([]byte(in.Settings), &hySettings)
			if hySettings.ObfsType == "salamander" && hySettings.ObfsPass != "" {
				p["obfs"] = "salamander"
				p["obfs-password"] = hySettings.ObfsPass
			}
			proxies = append(proxies, p)
			proxyNames = append(proxyNames, name)

		case "trojan":
			p := map[string]interface{}{
				"name":     name,
				"type":     "trojan",
				"server":   host,
				"port":     in.Port,
				"password": sb.client.Password,
				"udp":      true,
			}
			if stream.TLS != nil && stream.TLS.ServerName != "" {
				p["sni"] = stream.TLS.ServerName
			}
			proxies = append(proxies, p)
			proxyNames = append(proxyNames, name)
		}
	}

	if len(proxyNames) == 0 {
		proxyNames = append(proxyNames, "DIRECT")
	}

	clashConfig := map[string]interface{}{
		"port":                7890,
		"socks-port":          7891,
		"allow-lan":           true,
		"mode":                "rule",
		"log-level":           "info",
		"external-controller": "127.0.0.1:9090",
		"proxies":             proxies,
		"proxy-groups": []map[string]interface{}{
			{
				"name":    "PROXIES",
				"type":    "select",
				"proxies": append([]string{"AUTO"}, proxyNames...),
			},
			{
				"name":     "AUTO",
				"type":     "url-test",
				"url":      "https://www.gstatic.com/generate_204",
				"interval": 300,
				"proxies":  proxyNames,
			},
		},
		"rules": []string{
			"GEOIP,LAN,DIRECT,no-resolve",
			"GEOIP,CN,DIRECT",
			"MATCH,PROXIES",
		},
	}

	yamlData, err := yaml.Marshal(clashConfig)
	if err != nil {
		return "", err
	}
	return string(yamlData), nil
}

// BuildSingboxJSON generates official Sing-box client JSON configuration
func (sb *SubBuilder) BuildSingboxJSON() (string, error) {
	outbounds := make([]map[string]interface{}, 0)
	outboundTags := make([]string, 0)

	host := sb.serverHost
	if host == "" {
		host = "127.0.0.1"
	}

	for _, in := range sb.inbounds {
		if !in.Enable {
			continue
		}

		tag := fmt.Sprintf("%s-%s", in.Tag, in.Protocol)
		var stream struct {
			Network   string                 `json:"network"`
			Security  string                 `json:"security"`
			TLS       *core.TLSConfig        `json:"tls"`
			Reality   *core.RealityConfig    `json:"reality"`
			Transport *core.TransportConfig  `json:"transport"`
		}
		_ = json.Unmarshal([]byte(in.StreamSettings), &stream)

		switch in.Protocol {
		case "vless":
			out := map[string]interface{}{
				"type":        "vless",
				"tag":         tag,
				"server":      host,
				"server_port": in.Port,
				"uuid":        sb.client.UUID,
			}
			if sb.client.Flow != "" {
				out["flow"] = sb.client.Flow
			}
			if in.Security == "reality" && stream.Reality != nil {
				realityMap := map[string]interface{}{
					"enabled":    true,
					"public_key": stream.Reality.PublicKey,
					"utls": map[string]interface{}{
						"enabled":     true,
						"fingerprint": "chrome",
					},
				}
				if len(stream.Reality.ServerNames) > 0 {
					realityMap["server_name"] = stream.Reality.ServerNames[0]
				}
				if len(stream.Reality.ShortIds) > 0 {
					realityMap["short_id"] = stream.Reality.ShortIds[0]
				}
				out["tls"] = realityMap
			}
			outbounds = append(outbounds, out)
			outboundTags = append(outboundTags, tag)

		case "hysteria2":
			out := map[string]interface{}{
				"type":        "hysteria2",
				"tag":         tag,
				"server":      host,
				"server_port": in.Port,
				"password":    sb.client.Password,
				"tls": map[string]interface{}{
					"enabled": true,
				},
			}
			if stream.TLS != nil && stream.TLS.ServerName != "" {
				out["tls"].(map[string]interface{})["server_name"] = stream.TLS.ServerName
			}
			outbounds = append(outbounds, out)
			outboundTags = append(outboundTags, tag)
		}
	}

	if len(outboundTags) == 0 {
		outboundTags = append(outboundTags, "direct")
	}

	fullOutbounds := []interface{}{
		map[string]interface{}{
			"type":      "selector",
			"tag":       "select",
			"outbounds": append([]string{"auto"}, outboundTags...),
			"default":   "auto",
		},
		map[string]interface{}{
			"type":      "urltest",
			"tag":       "auto",
			"outbounds": outboundTags,
			"url":       "https://www.gstatic.com/generate_204",
			"interval":  "3m",
		},
	}

	for _, o := range outbounds {
		fullOutbounds = append(fullOutbounds, o)
	}

	fullOutbounds = append(fullOutbounds,
		map[string]interface{}{"type": "direct", "tag": "direct"},
		map[string]interface{}{"type": "block", "tag": "block"},
		map[string]interface{}{"type": "dns", "tag": "dns-out"},
	)

	config := map[string]interface{}{
		"log": map[string]interface{}{"level": "info"},
		"dns": map[string]interface{}{
			"servers": []map[string]interface{}{
				{"tag": "remote-dns", "address": "https://1.1.1.1/dns-query", "detour": "select"},
				{"tag": "local-dns", "address": "local", "detour": "direct"},
			},
		},
		"inbounds": []map[string]interface{}{
			{"type": "tun", "tag": "tun-in", "inet4_address": "172.19.0.1/30", "auto_route": true, "strict_route": true, "stack": "mixed", "sniff": true},
			{"type": "mixed", "tag": "mixed-in", "listen": "127.0.0.1", "listen_port": 2080, "sniff": true},
		},
		"outbounds": fullOutbounds,
		"route": map[string]interface{}{
			"rules": []map[string]interface{}{
				{"protocol": "dns", "outbound": "dns-out"},
				{"geosite": "cn", "outbound": "direct"},
				{"geoip": "cn", "outbound": "direct"},
				{"geoip": "private", "outbound": "direct"},
			},
			"auto_detect_interface": true,
			"final":                 "select",
		},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
