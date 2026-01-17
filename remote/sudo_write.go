package remote

import (
	"fmt"
)

func (c *Client) WriteFileSudo(content string, targetPath string) error {
	tmp := "/tmp/simpledeploy.tmp"
	if err := c.UploadBytes(content, tmp); err != nil {
		return err
	}

	cmd := fmt.Sprintf(`
set -e
sudo mv "%s" "%s"
`, tmp, targetPath)

	res, err := c.Run(cmd)
	if err != nil {
		return err
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("write file failed (exit %d): %s", res.ExitCode, res.Stderr)
	}
	return nil
}

func (c *Client) RestartSystemdService(serviceFile string) error {
	cmd := fmt.Sprintf(`
set -e
sudo systemctl daemon-reload
sudo systemctl enable "%s" >/dev/null 2>&1 || true
sudo systemctl restart "%s"
`, serviceFile, serviceFile)

	res, err := c.Run(cmd)
	if err != nil {
		return err
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("systemd restart failed (exit %d): %s", res.ExitCode, res.Stderr)
	}
	return nil
}
