package core

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// ValidateConfig uses 'sing-box check -c <tempFile>' to ensure config syntax is 100% valid before applying
func ValidateConfig(binPath string, cfg *SingboxConfig) error {
	if binPath == "" {
		binPath = "sing-box"
	}

	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, fmt.Sprintf("singbox_check_%d.json", time.Now().UnixNano()))
	defer os.Remove(tempFile)

	if err := WriteConfigToFile(cfg, tempFile); err != nil {
		return fmt.Errorf("failed to write temp check config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath, "check", "-c", tempFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("sing-box configuration check failed:\n%s", string(out))
	}

	return nil
}
