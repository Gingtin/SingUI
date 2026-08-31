package network

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// PingTCP tests real TCP handshake latency to target host:port
func PingTCP(host string, port int, timeout time.Duration) (int64, error) {
	if timeout <= 0 {
		timeout = 3 * time.Second
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return -1, err
	}
	defer conn.Close()

	latency := time.Since(start).Milliseconds()
	return latency, nil
}

// GetPublicIP returns the server public IP address
func GetPublicIP() string {
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return "127.0.0.1"
	}
	defer resp.Body.Close()

	var ip string
	_, _ = fmt.Fscan(resp.Body, &ip)
	if ip == "" {
		return "127.0.0.1"
	}
	return ip
}
