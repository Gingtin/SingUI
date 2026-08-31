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

// BuildBase64Links generates standard base64 encoded node links
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
		query := url.Values{}
		query.Set("type", in.Network)
		if stream.Transport != nil {
			if stream.Transport.Type == "ws" {
				query.Set("type", "ws")
				query.Set("path", stream.Transport.Path)
			} else if stream.Transport.Type == "grpc" {
				query.Set("type", "grpc")
				query.Set("serviceName", stream.Transport.ServiceName)
			}
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
			"scy":  "auto",
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
		if stream.Transport != nil {
			if stream.Transport.Type == "ws" {
				query.Set("type", "ws")
				query.Set("path", stream.Transport.Path)
			} else if stream.Transport.Type == "grpc" {
				query.Set("type", "grpc")
				query.Set("serviceName", stream.Transport.ServiceName)
			}
		}
		return fmt.Sprintf("trojan://%s@%s:%d?%s#%s", sb.client.Password, host, in.Port, query.Encode(), remark)

	case "hysteria2":
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
		
	case "anytls":
		return fmt.Sprintf("anytls://%s@%s:%d#%s", sb.client.UUID, host, in.Port, remark)
	}

	return ""
}

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
				p["reality-opts"] = map[string]interface{}{
					"public-key": stream.Reality.PublicKey,
					"short-id":   stream.Reality.ShortIds[0],
				}
				p["client-fingerprint"] = "chrome"
			} else if in.Security == "tls" && stream.TLS != nil {
				p["tls"] = true
				if stream.TLS.ServerName != "" {
					p["servername"] = stream.TLS.ServerName
				}
			}
			if stream.Transport != nil {
				if stream.Transport.Type == "ws" {
					p["network"] = "ws"
					p["ws-opts"] = map[string]interface{}{
						"path": stream.Transport.Path,
					}
				} else if stream.Transport.Type == "grpc" {
					p["network"] = "grpc"
					p["grpc-opts"] = map[string]interface{}{
						"grpc-service-name": stream.Transport.ServiceName,
					}
				}
			}
			proxies = append(proxies, p)
			proxyNames = append(proxyNames, name)

		case "vmess":
			p := map[string]interface{}{
				"name":   name,
				"type":   "vmess",
				"server": host,
				"port":   in.Port,
				"uuid":   sb.client.UUID,
				"alterId": 0,
				"cipher": "auto",
				"network": in.Network,
				"udp":    true,
			}
			if in.Security == "tls" && stream.TLS != nil {
				p["tls"] = true
				if stream.TLS.ServerName != "" {
					p["servername"] = stream.TLS.ServerName
				}
			}
			if stream.Transport != nil {
				if stream.Transport.Type == "ws" {
					p["network"] = "ws"
					p["ws-opts"] = map[string]interface{}{
						"path": stream.Transport.Path,
					}
				}
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
			if stream.Transport != nil {
				if stream.Transport.Type == "ws" {
					p["network"] = "ws"
					p["ws-opts"] = map[string]interface{}{
						"path": stream.Transport.Path,
					}
				} else if stream.Transport.Type == "grpc" {
					p["network"] = "grpc"
					p["grpc-opts"] = map[string]interface{}{
						"grpc-service-name": stream.Transport.ServiceName,
					}
				}
			}
			proxies = append(proxies, p)
			proxyNames = append(proxyNames, name)

		case "tuic":
			p := map[string]interface{}{
				"name":     name,
				"type":     "tuic",
				"server":   host,
				"port":     in.Port,
				"uuid":     sb.client.UUID,
				"password": sb.client.Password,
				"udp":      true,
			}
			if stream.TLS != nil && stream.TLS.ServerName != "" {
				p["sni"] = stream.TLS.ServerName
			}
			proxies = append(proxies, p)
			proxyNames = append(proxyNames, name)
			
		case "shadowsocks":
			var ssSettings struct {
				Method   string `json:"method"`
				Password string `json:"password"`
			}
			_ = json.Unmarshal([]byte(in.Settings), &ssSettings)
			p := map[string]interface{}{
				"name":     name,
				"type":     "ss",
				"server":   host,
				"port":     in.Port,
				"cipher":   ssSettings.Method,
				"password": ssSettings.Password,
				"udp":      true,
			}
			proxies = append(proxies, p)
			proxyNames = append(proxyNames, name)
			
		case "anytls":
			p := map[string]interface{}{
				"name":     name,
				"type":     "vless", // Fallback
				"server":   host,
				"port":     in.Port,
				"uuid":     sb.client.UUID,
				"network":  "tcp",
				"udp":      true,
				"tls":      true,
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
				"proxies": append([]string{"AUTO", "FALLBACK"}, proxyNames...),
			},
			{
				"name":     "AUTO",
				"type":     "url-test",
				"url":      "https://www.gstatic.com/generate_204",
				"interval": 300,
				"proxies":  proxyNames,
			},
			{
				"name":     "FALLBACK",
				"type":     "fallback",
				"url":      "https://www.gstatic.com/generate_204",
				"interval": 300,
				"proxies":  proxyNames,
			},
		},
		"rule-providers": map[string]interface{}{
			"geosite-cn": map[string]interface{}{
				"type": "http",
				"behavior": "domain",
				"url": "https://cdn.jsdelivr.net/gh/Loyalsoldier/v2ray-rules-dat@release/geosite/cn.txt",
				"path": "./ruleset/geosite-cn.yaml",
				"interval": 86400,
			},
		},
		"rules": []string{
			"RULE-SET,geosite-cn,DIRECT",
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

		out := map[string]interface{}{
			"type":        in.Protocol,
			"tag":         tag,
			"server":      host,
			"server_port": in.Port,
		}

		switch in.Protocol {
		case "vless":
			out["uuid"] = sb.client.UUID
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
			} else if in.Security == "tls" && stream.TLS != nil {
				tlsMap := map[string]interface{}{
					"enabled": true,
					"utls": map[string]interface{}{
						"enabled":     true,
						"fingerprint": "chrome",
					},
				}
				if stream.TLS.ServerName != "" {
					tlsMap["server_name"] = stream.TLS.ServerName
				}
				out["tls"] = tlsMap
			}
			
			if stream.Transport != nil {
				if stream.Transport.Type == "ws" {
					out["transport"] = map[string]interface{}{
						"type": "ws",
						"path": stream.Transport.Path,
					}
				} else if stream.Transport.Type == "grpc" {
					out["transport"] = map[string]interface{}{
						"type": "grpc",
						"service_name": stream.Transport.ServiceName,
					}
				} else if stream.Transport.Type == "httpupgrade" {
					out["transport"] = map[string]interface{}{
						"type": "httpupgrade",
						"path": stream.Transport.Path,
					}
				}
			}

		case "vmess":
			out["uuid"] = sb.client.UUID
			out["security"] = "auto"
			out["alter_id"] = 0
			if in.Security == "tls" && stream.TLS != nil {
				tlsMap := map[string]interface{}{
					"enabled": true,
					"utls": map[string]interface{}{
						"enabled":     true,
						"fingerprint": "chrome",
					},
				}
				if stream.TLS.ServerName != "" {
					tlsMap["server_name"] = stream.TLS.ServerName
				}
				out["tls"] = tlsMap
			}
			if stream.Transport != nil && stream.Transport.Type == "ws" {
				out["transport"] = map[string]interface{}{
					"type": "ws",
					"path": stream.Transport.Path,
				}
			}

		case "hysteria2":
			out["password"] = sb.client.Password
			tlsMap := map[string]interface{}{
				"enabled": true,
			}
			if stream.TLS != nil && stream.TLS.ServerName != "" {
				tlsMap["server_name"] = stream.TLS.ServerName
			}
			out["tls"] = tlsMap
			var hySettings struct {
				ObfsType string `json:"obfs_type"`
				ObfsPass string `json:"obfs_password"`
			}
			_ = json.Unmarshal([]byte(in.Settings), &hySettings)
			if hySettings.ObfsType == "salamander" && hySettings.ObfsPass != "" {
				out["obfs"] = map[string]interface{}{
					"type": "salamander",
					"password": hySettings.ObfsPass,
				}
			}

		case "trojan":
			out["password"] = sb.client.Password
			tlsMap := map[string]interface{}{
				"enabled": true,
				"utls": map[string]interface{}{
					"enabled":     true,
					"fingerprint": "chrome",
				},
			}
			if stream.TLS != nil && stream.TLS.ServerName != "" {
				tlsMap["server_name"] = stream.TLS.ServerName
			}
			out["tls"] = tlsMap
			if stream.Transport != nil {
				if stream.Transport.Type == "ws" {
					out["transport"] = map[string]interface{}{
						"type": "ws",
						"path": stream.Transport.Path,
					}
				} else if stream.Transport.Type == "grpc" {
					out["transport"] = map[string]interface{}{
						"type": "grpc",
						"service_name": stream.Transport.ServiceName,
					}
				}
			}

		case "tuic":
			out["uuid"] = sb.client.UUID
			out["password"] = sb.client.Password
			tlsMap := map[string]interface{}{
				"enabled": true,
			}
			if stream.TLS != nil && stream.TLS.ServerName != "" {
				tlsMap["server_name"] = stream.TLS.ServerName
			}
			out["tls"] = tlsMap
			
		case "shadowsocks":
			var ssSettings struct {
				Method   string `json:"method"`
				Password string `json:"password"`
			}
			_ = json.Unmarshal([]byte(in.Settings), &ssSettings)
			out["method"] = ssSettings.Method
			out["password"] = ssSettings.Password
			
		case "anytls":
			out["type"] = "vless" // mapping anytls
			out["uuid"] = sb.client.UUID
			out["tls"] = map[string]interface{}{"enabled": true}
		}
		
		outbounds = append(outbounds, out)
		outboundTags = append(outboundTags, tag)
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
			"rules": []map[string]interface{}{
				{"outbound": "any", "server": "local-dns"},
				{"rule_set": "geosite-cn", "server": "local-dns"},
			},
		},
		"inbounds": []map[string]interface{}{
			{"type": "tun", "tag": "tun-in", "inet4_address": "172.19.0.1/30", "auto_route": true, "strict_route": true, "stack": "mixed", "sniff": true},
			{"type": "mixed", "tag": "mixed-in", "listen": "127.0.0.1", "listen_port": 2080, "sniff": true},
		},
		"outbounds": fullOutbounds,
		"route": map[string]interface{}{
			"rule_set": []map[string]interface{}{
				{
					"tag": "geosite-cn",
					"type": "remote",
					"format": "binary",
					"url": "https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-cn.srs",
					"download_detour": "select",
				},
				{
					"tag": "geoip-cn",
					"type": "remote",
					"format": "binary",
					"url": "https://raw.githubusercontent.com/SagerNet/sing-geoip/rule-set/geoip-cn.srs",
					"download_detour": "select",
				},
			},
			"rules": []map[string]interface{}{
				{"protocol": "dns", "outbound": "dns-out"},
				{"rule_set": "geosite-cn", "outbound": "direct"},
				{"rule_set": "geoip-cn", "outbound": "direct"},
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

func (sb *SubBuilder) BuildQuantumultX() (string, error) {
	links := []string{}
	for _, in := range sb.inbounds {
		if !in.Enable { continue }
		host := sb.serverHost
		if host == "" { host = "127.0.0.1" }
		if in.Protocol == "trojan" {
			links = append(links, fmt.Sprintf("trojan=%s:%d, password=%s, over-tls=true, tls-verification=true, fast-open=false, udp-relay=false, tag=%s", host, in.Port, sb.client.Password, in.Tag))
		} else if in.Protocol == "shadowsocks" {
			var ssSettings struct {
				Method   string `json:"method"`
				Password string `json:"password"`
			}
			_ = json.Unmarshal([]byte(in.Settings), &ssSettings)
			links = append(links, fmt.Sprintf("shadowsocks=%s:%d, method=%s, password=%s, fast-open=false, udp-relay=false, tag=%s", host, in.Port, ssSettings.Method, ssSettings.Password, in.Tag))
		}
	}
	return strings.Join(links, "\n"), nil
}

func (sb *SubBuilder) BuildSurge() (string, error) {
	links := []string{}
	for _, in := range sb.inbounds {
		if !in.Enable { continue }
		host := sb.serverHost
		if host == "" { host = "127.0.0.1" }
		if in.Protocol == "trojan" {
			links = append(links, fmt.Sprintf("%s = trojan, %s, %d, password=%s, udp-relay=true", in.Tag, host, in.Port, sb.client.Password))
		} else if in.Protocol == "shadowsocks" {
			var ssSettings struct {
				Method   string `json:"method"`
				Password string `json:"password"`
			}
			_ = json.Unmarshal([]byte(in.Settings), &ssSettings)
			links = append(links, fmt.Sprintf("%s = ss, %s, %d, encrypt-method=%s, password=%s, udp-relay=true", in.Tag, host, in.Port, ssSettings.Method, ssSettings.Password))
		}
	}
	return strings.Join(links, "\n"), nil
}

func (sb *SubBuilder) BuildLoon() (string, error) {
	links := []string{}
	for _, in := range sb.inbounds {
		if !in.Enable { continue }
		host := sb.serverHost
		if host == "" { host = "127.0.0.1" }
		if in.Protocol == "trojan" {
			links = append(links, fmt.Sprintf("%s = trojan, %s, %d, %s", in.Tag, host, in.Port, sb.client.Password))
		}
	}
	return strings.Join(links, "\n"), nil
}
