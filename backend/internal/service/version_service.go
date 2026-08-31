package service

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"github.com/singbox-ui/singbox-ui/internal/core"
	"github.com/singbox-ui/singbox-ui/internal/database"
	"github.com/singbox-ui/singbox-ui/internal/database/models"
)

type VersionInfo struct {
	Panel PanelVersionInfo `json:"panel"`
	Core  CoreVersionInfo  `json:"core"`
	Geo   GeoVersionInfo   `json:"geo"`
}

type PanelVersionInfo struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	HasUpdate      bool   `json:"has_update"`
	ReleaseNotes   string `json:"release_notes"`
	ReleaseURL     string `json:"release_url"`
}

type CoreVersionInfo struct {
	CurrentVersion    string   `json:"current_version"`
	LatestVersion     string   `json:"latest_version"`
	HasUpdate         bool     `json:"has_update"`
	AvailableVersions []string `json:"available_versions"`
}

type GeoVersionInfo struct {
	LastUpdated string `json:"last_updated"`
	Status      string `json:"status"`
}

var CurrentPanelVersion = "v1.0.0"

// CheckAllVersions checks latest versions for SingUI and Sing-box from GitHub
func CheckAllVersions() (*VersionInfo, error) {
	client := &http.Client{Timeout: 8 * time.Second}

	// 1. Get Sing-box releases
	coreLatest := "v1.9.7"
	var coreVersions []string
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/SagerNet/sing-box/releases?per_page=5", nil)
	req.Header.Set("User-Agent", "SingUI-Updater")
	if resp, err := client.Do(req); err == nil && resp.StatusCode == 200 {
		var releases []struct {
			TagName string `json:"tag_name"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&releases); err == nil && len(releases) > 0 {
			coreLatest = releases[0].TagName
			for _, r := range releases {
				coreVersions = append(coreVersions, r.TagName)
			}
		}
		resp.Body.Close()
	}

	if len(coreVersions) == 0 {
		coreVersions = []string{"v1.10.1", "v1.10.0", "v1.9.7", "v1.9.0"}
	}

	// Determine running core version
	currentCore := "v1.9.7"
	if core.Instance != nil {
		status := core.Instance.GetStatus()
		if status.Version != "" {
			parts := strings.Fields(status.Version)
			for i, p := range parts {
				if p == "version" && i+1 < len(parts) {
					currentCore = "v" + strings.TrimPrefix(parts[i+1], "v")
					break
				}
			}
		}
	}

	// 2. Get SingUI Panel releases
	panelLatest := CurrentPanelVersion
	releaseNotes := "已经是最新版本"
	releaseURL := "https://github.com/Gingtin/SingUI/releases"
	reqPanel, _ := http.NewRequest("GET", "https://api.github.com/repos/Gingtin/SingUI/releases/latest", nil)
	reqPanel.Header.Set("User-Agent", "SingUI-Updater")
	if resp, err := client.Do(reqPanel); err == nil && resp.StatusCode == 200 {
		var release struct {
			TagName string `json:"tag_name"`
			Body    string `json:"body"`
			HTMLURL string `json:"html_url"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&release); err == nil && release.TagName != "" {
			panelLatest = release.TagName
			releaseNotes = release.Body
			releaseURL = release.HTMLURL
		}
		resp.Body.Close()
	}

	hasPanelUpdate := panelLatest != CurrentPanelVersion && panelLatest != ""
	hasCoreUpdate := coreLatest != currentCore && coreLatest != ""

	return &VersionInfo{
		Panel: PanelVersionInfo{
			CurrentVersion: CurrentPanelVersion,
			LatestVersion:  panelLatest,
			HasUpdate:      hasPanelUpdate,
			ReleaseNotes:   releaseNotes,
			ReleaseURL:     releaseURL,
		},
		Core: CoreVersionInfo{
			CurrentVersion:    currentCore,
			LatestVersion:     coreLatest,
			HasUpdate:         hasCoreUpdate,
			AvailableVersions: coreVersions,
		},
		Geo: GeoVersionInfo{
			LastUpdated: time.Now().Format("2006-01-02 15:04"),
			Status:      "最新",
		},
	}, nil
}

// UpdateSingboxCore downloads and replaces the sing-box binary on Linux
func UpdateSingboxCore(targetVersion string) error {
	if targetVersion == "" {
		return fmt.Errorf("target version is empty")
	}

	// Format version string (strip leading v for download url)
	cleanVersion := strings.TrimPrefix(targetVersion, "v")
	arch := runtime.GOARCH
	var sbArch string
	switch arch {
	case "amd64":
		sbArch = "linux-amd64"
	case "arm64":
		sbArch = "linux-arm64"
	default:
		return fmt.Errorf("unsupported arch: %s", arch)
	}

	downloadURL := fmt.Sprintf("https://github.com/SagerNet/sing-box/releases/download/v%s/sing-box-%s-%s.tar.gz", cleanVersion, cleanVersion, sbArch)
	log.Printf("[Updater] Downloading Sing-box core: %s\n", downloadURL)

	resp, err := http.Get(downloadURL)
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("failed to download sing-box %s: %v", targetVersion, err)
	}
	defer resp.Body.Close()

	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	tempBinaryPath := "/tmp/sing-box-new"
	found := false

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if filepath.Base(header.Name) == "sing-box" && !header.FileInfo().IsDir() {
			outFile, err := os.OpenFile(tempBinaryPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("sing-box binary not found in archive")
	}

	// Stop core, replace binary, restart core
	if core.Instance != nil {
		_ = core.Instance.Stop()
	}

	targetPath := "/usr/local/bin/sing-box"
	_ = os.Remove(targetPath)
	if err := os.Rename(tempBinaryPath, targetPath); err != nil {
		_ = exec.Command("cp", tempBinaryPath, targetPath).Run()
	}
	_ = os.Chmod(targetPath, 0755)

	if core.Instance != nil {
		_ = core.Instance.Start()
	}

	log.Printf("[Updater] Sing-box core successfully updated to: %s\n", targetVersion)
	return nil
}

// UpdateGeoDatabases downloads latest GeoIP and Geosite SRS databases
func UpdateGeoDatabases() error {
	log.Println("[Updater] Updating GeoIP & Geosite databases...")
	var binPathSetting models.Setting
	database.DB.Where("key = ?", "singbox_bin_path").First(&binPathSetting)
	binPath := binPathSetting.Value
	targetDir := "/usr/local/share/sing-box/"
	if binPath != "" {
		targetDir = filepath.Dir(binPath)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", targetDir, err)
	}

	downloadFile := func(url, dest string) error {
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			return fmt.Errorf("failed to download %s: %v", url, err)
		}
		defer resp.Body.Close()

		out, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		return err
	}

	if err := downloadFile("https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip.db", filepath.Join(targetDir, "geoip.db")); err != nil {
		return err
	}

	if err := downloadFile("https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite.db", filepath.Join(targetDir, "geosite.db")); err != nil {
		return err
	}

	log.Println("[Updater] GeoIP & Geosite databases updated successfully.")
	return nil
}
