package remote

import (
	"fmt"
)

func (client *Client) WriteFileSudo(content string, targetPath string) error {
	tmp := "/tmp/simpledeploy.tmp"
	if err := client.UploadBytes(content, tmp); err != nil {
		return err
	}

	cmd := fmt.Sprintf(`
set -e
sudo mv "%s" "%s"
`, tmp, targetPath)

	result, err := client.Run(cmd)
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("write file failed (exit %d): %s", result.ExitCode, result.Stderr)
	}
	return nil
}

func (client *Client) RestartSystemdService(serviceFile string) error {
	cmd := fmt.Sprintf(`
set -e
sudo systemctl daemon-reload
sudo systemctl enable "%s" >/dev/null 2>&1 || true
sudo systemctl restart "%s"
`, serviceFile, serviceFile)

	result, err := client.Run(cmd)
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("systemd restart failed (exit %d): %s", result.ExitCode, result.Stderr)
	}
	return nil
}
