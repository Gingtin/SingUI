package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/singbox-ui/singbox-ui/internal/database"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
	"github.com/singbox-ui/singbox-ui/internal/subscription"
)

func HandleSubscription(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.String(http.StatusBadRequest, "Invalid subscription token")
		return
	}

	var client models.Client
	if err := database.DB.Where("sub_token = ?", token).First(&client).Error; err != nil {
		c.String(http.StatusNotFound, "Subscription not found or expired")
		return
	}

	if !client.Enable {
		c.String(http.StatusForbidden, "Subscription is disabled")
		return
	}

	expireSec := client.ExpiryTime / 1000
	userinfo := fmt.Sprintf("upload=%d; download=%d; total=%d; expire=%d", client.Up, client.Down, client.Total, expireSec)
	c.Header("Subscription-Userinfo", userinfo)
	c.Header("Profile-Update-Interval", "6")

	var subDomainSetting models.Setting
	database.DB.Where("key = ?", "sub_domain").First(&subDomainSetting)
	serverHost := subDomainSetting.Value
	if serverHost == "" {
		serverHost = c.Request.Host
		if strings.Contains(serverHost, ":") {
			serverHost = strings.Split(serverHost, ":")[0]
		}
	}

	var inbounds []models.Inbound
	database.DB.Preload("Clients").Where("enable = ?", true).Find(&inbounds)

	builder := subscription.NewSubBuilder(serverHost, &client, inbounds)

	flag := strings.ToLower(c.Query("flag"))
	ua := strings.ToLower(c.GetHeader("User-Agent"))

	if flag == "sing-box" || flag == "sbox" || strings.Contains(ua, "sing-box") || strings.Contains(ua, "s-box") {
		jsonStr, err := builder.BuildSingboxJSON()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(http.StatusOK, jsonStr)
		return
	}

	if flag == "clash" || flag == "meta" || flag == "mihomo" || strings.Contains(ua, "clash") || strings.Contains(ua, "meta") || strings.Contains(ua, "stash") || strings.Contains(ua, "verge") {
		yamlStr, err := builder.BuildClashMeta()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Header("Content-Type", "text/yaml; charset=utf-8")
		c.String(http.StatusOK, yamlStr)
		return
	}

	if flag == "quanx" || strings.Contains(ua, "quantumult%20x") || strings.Contains(ua, "quanx") {
		str, err := builder.BuildQuantumultX()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(http.StatusOK, str)
		return
	}

	if flag == "surge" || strings.Contains(ua, "surge") {
		str, err := builder.BuildSurge()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(http.StatusOK, str)
		return
	}

	if flag == "loon" || strings.Contains(ua, "loon") {
		str, err := builder.BuildLoon()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(http.StatusOK, str)
		return
	}

	base64Links := builder.BuildBase64Links()
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, base64Links)
}

func HandleSubView(c *gin.Context) {
	token := c.Param("token")
	var client models.Client
	if err := database.DB.Where("sub_token = ?", token).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	var subDomainSetting models.Setting
	database.DB.Where("key = ?", "sub_domain").First(&subDomainSetting)
	serverHost := subDomainSetting.Value
	if serverHost == "" {
		serverHost = c.Request.Host
		if strings.Contains(serverHost, ":") {
			serverHost = strings.Split(serverHost, ":")[0]
		}
	}

	var inbounds []models.Inbound
	database.DB.Preload("Clients").Where("enable = ?", true).Find(&inbounds)

	now := time.Now().UnixMilli()
	isExpired := client.ExpiryTime > 0 && now >= client.ExpiryTime
	isOverQuota := client.Total > 0 && (client.Up+client.Down) >= client.Total

	if strings.Contains(c.GetHeader("Accept"), "application/json") {
		c.JSON(http.StatusOK, gin.H{
			"email":         client.Email,
			"up":            client.Up,
			"down":          client.Down,
			"total":         client.Total,
			"expiry_time":   client.ExpiryTime,
			"is_expired":    isExpired,
			"is_over_quota": isOverQuota,
			"sub_token":     client.SubToken,
			"server_host":   serverHost,
		})
		return
	}

	// HTML View
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	subURL := fmt.Sprintf("%s://%s/sub/%s", scheme, c.Request.Host, token)
	
	used := client.Up + client.Down
	total := client.Total
	percentage := 0.0
	if total > 0 {
		percentage = float64(used) / float64(total) * 100
	}
	if percentage > 100 { percentage = 100 }

	statusHTML := "🟢 正常 (Active)"
	if isExpired {
		statusHTML = "🔴 已过期 (Expired)"
	} else if isOverQuota {
		statusHTML = "🟡 流量耗尽 (Out of Data)"
	}

	expDate := "永不 (Never)"
	if client.ExpiryTime > 0 {
		expDate = time.UnixMilli(client.ExpiryTime).Format("2006-01-02 15:04:05")
	}
	
	nodesHTML := ""
	builder := subscription.NewSubBuilder(serverHost, &client, inbounds)
	links := strings.Split(builder.BuildBase64Links(), "\n")
	// For node list rendering we need actual links, we might decode them or we can just show the encoded list.
	var plainLinks []string
	for _, in := range inbounds {
		if !in.Enable { continue }
		b := subscription.NewSubBuilder(serverHost, &client, []models.Inbound{in})
		plainLinks = append(plainLinks, b.BuildBase64Links())
	}
	
	// Better to decode them in JS or just show in table.
	nodesJSON, _ := json.Marshal(inbounds)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Subscription Portal</title>
    <style>
        :root {
            --bg: #f5f5f7;
            --card-bg: rgba(255, 255, 255, 0.6);
            --text: #1d1d1f;
            --accent: #0071e3;
        }
        @media (prefers-color-scheme: dark) {
            :root {
                --bg: #000000;
                --card-bg: rgba(28, 28, 30, 0.6);
                --text: #f5f5f7;
                --accent: #2997ff;
            }
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: var(--bg);
            color: var(--text);
            margin: 0;
            padding: 2rem;
            display: flex;
            flex-direction: column;
            align-items: center;
        }
        .container {
            max-width: 800px;
            width: 100%%;
        }
        .card {
            background: var(--card-bg);
            backdrop-filter: blur(20px);
            -webkit-backdrop-filter: blur(20px);
            border-radius: 20px;
            padding: 2rem;
            box-shadow: 0 8px 32px rgba(0,0,0,0.1);
            margin-bottom: 2rem;
        }
        .header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            border-bottom: 1px solid rgba(128,128,128,0.2);
            padding-bottom: 1rem;
            margin-bottom: 2rem;
        }
        .traffic-lights {
            display: flex;
            gap: 8px;
        }
        .traffic-light {
            width: 12px;
            height: 12px;
            border-radius: 50%%;
        }
        .red { background: #ff5f56; }
        .yellow { background: #ffbd2e; }
        .green { background: #27c93f; }
        
        .gauge-container {
            display: flex;
            justify-content: space-around;
            align-items: center;
            flex-wrap: wrap;
        }
        .stats {
            text-align: center;
        }
        
        svg {
            transform: rotate(-90deg);
        }
        .circle-bg {
            fill: none;
            stroke: rgba(128,128,128,0.2);
            stroke-width: 8;
        }
        .circle-progress {
            fill: none;
            stroke: var(--accent);
            stroke-width: 8;
            stroke-linecap: round;
            transition: stroke-dasharray 1s ease;
        }
        
        .buttons {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1rem;
            margin-top: 2rem;
        }
        .btn {
            display: inline-block;
            background: var(--accent);
            color: #fff;
            text-decoration: none;
            padding: 12px 20px;
            border-radius: 12px;
            text-align: center;
            font-weight: 600;
            transition: opacity 0.2s;
        }
        .btn:hover {
            opacity: 0.9;
        }
        table {
            width: 100%%;
            border-collapse: collapse;
        }
        th, td {
            text-align: left;
            padding: 12px;
            border-bottom: 1px solid rgba(128,128,128,0.2);
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <div class="header">
                <div class="traffic-lights">
                    <div class="traffic-light red"></div>
                    <div class="traffic-light yellow"></div>
                    <div class="traffic-light green"></div>
                </div>
                <h2>Subscription Portal</h2>
            </div>
            
            <div class="gauge-container">
                <div style="position: relative; width: 150px; height: 150px;">
                    <svg width="150" height="150" viewBox="0 0 100 100">
                        <circle class="circle-bg" cx="50" cy="50" r="40"></circle>
                        <circle class="circle-progress" cx="50" cy="50" r="40" stroke-dasharray="%f, 251.2"></circle>
                    </svg>
                    <div style="position: absolute; top: 50%%; left: 50%%; transform: translate(-50%%, -50%%); font-size: 1.2rem; font-weight: bold;">
                        %.1f%%
                    </div>
                </div>
                <div class="stats">
                    <p><strong>Status:</strong> %s</p>
                    <p><strong>Used:</strong> %.2f GB</p>
                    <p><strong>Total:</strong> %.2f GB</p>
                    <p><strong>Expires:</strong> %s</p>
                </div>
            </div>

            <div class="buttons">
                <a class="btn" href="sing-box://import-remote-profile?url=%s">Import to Sing-box</a>
                <a class="btn" href="clash://install-config?url=%s">Import to Clash</a>
                <a class="btn" href="shadowrocket://add/sub://%s">Import to Shadowrocket</a>
            </div>
        </div>
        
        <div class="card">
            <h3>Node List</h3>
            <table>
                <thead>
                    <tr>
                        <th>Tag</th>
                        <th>Protocol</th>
                    </tr>
                </thead>
                <tbody id="nodes">
                </tbody>
            </table>
        </div>
    </div>
    
    <script>
        const nodes = %s;
        const tbody = document.getElementById('nodes');
        nodes.forEach(n => {
            if(!n.enable) return;
            const tr = document.createElement('tr');
            tr.innerHTML = '<td>' + n.tag + '</td><td>' + n.protocol + '</td>';
            tbody.appendChild(tr);
        });
    </script>
</body>
</html>`,
		percentage*2.512, percentage, statusHTML, float64(used)/1073741824, float64(total)/1073741824, expDate, 
		url.QueryEscape(subURL+"?flag=sing-box"), url.QueryEscape(subURL+"?flag=clash"), 
		base64.StdEncoding.EncodeToString([]byte(subURL)), string(nodesJSON))

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}
