package util

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type SystemStatus struct {
	CPUPercent     float64 `json:"cpu_percent"`
	MemTotal       uint64  `json:"mem_total"`
	MemUsed        uint64  `json:"mem_used"`
	MemPercent     float64 `json:"mem_percent"`
	DiskTotal      uint64  `json:"disk_total"`
	DiskUsed       uint64  `json:"disk_used"`
	DiskPercent    float64 `json:"disk_percent"`
	NetUploadRate  uint64  `json:"net_upload_rate"`   // Bytes/s
	NetDownloadRate uint64 `json:"net_download_rate"` // Bytes/s
	NetTotalSent   uint64  `json:"net_total_sent"`
	NetTotalRecv   uint64  `json:"net_total_recv"`
	Uptime         uint64  `json:"uptime"`
	OS             string  `json:"os"`
	Platform       string  `json:"platform"`
}

var (
	lastNetSampleTime time.Time
	lastBytesSent     uint64
	lastBytesRecv     uint64
	netMutex          sync.Mutex
)

func GetSystemStatus() (*SystemStatus, error) {
	status := &SystemStatus{}

	// CPU
	cpuPercents, err := cpu.Percent(0, false)
	if err == nil && len(cpuPercents) > 0 {
		status.CPUPercent = cpuPercents[0]
	}

	// Memory
	vMem, err := mem.VirtualMemory()
	if err == nil {
		status.MemTotal = vMem.Total
		status.MemUsed = vMem.Used
		status.MemPercent = vMem.UsedPercent
	}

	// Disk
	dUsage, err := disk.Usage("/")
	if err != nil {
		dUsage, _ = disk.Usage("C:")
	}
	if dUsage != nil {
		status.DiskTotal = dUsage.Total
		status.DiskUsed = dUsage.Used
		status.DiskPercent = dUsage.UsedPercent
	}

	// Host
	hInfo, err := host.Info()
	if err == nil {
		status.Uptime = hInfo.Uptime
		status.OS = hInfo.OS
		status.Platform = hInfo.Platform + " " + hInfo.PlatformVersion
	}

	// Network
	netIOCounters, err := net.IOCounters(false)
	if err == nil && len(netIOCounters) > 0 {
		netMutex.Lock()
		now := time.Now()
		totalSent := netIOCounters[0].BytesSent
		totalRecv := netIOCounters[0].BytesRecv

		status.NetTotalSent = totalSent
		status.NetTotalRecv = totalRecv

		if !lastNetSampleTime.IsZero() {
			elapsed := now.Sub(lastNetSampleTime).Seconds()
			if elapsed > 0 {
				if totalSent >= lastBytesSent {
					status.NetUploadRate = uint64(float64(totalSent-lastBytesSent) / elapsed)
				}
				if totalRecv >= lastBytesRecv {
					status.NetDownloadRate = uint64(float64(totalRecv-lastBytesRecv) / elapsed)
				}
			}
		}

		lastNetSampleTime = now
		lastBytesSent = totalSent
		lastBytesRecv = totalRecv
		netMutex.Unlock()
	}

	return status, nil
}
