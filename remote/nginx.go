package remote

import (
	"fmt"
	"path"
)

func (c *Client) ApplyNginxConf(confFileName string, confContent string) error {
	tmp := path.Join("/tmp", "."+confFileName)
	target := path.Join("/etc/nginx/conf.d", confFileName)

	// upload temp file
	if err := c.UploadBytes(confContent, tmp); err != nil {
		return err
	}

	cmd := fmt.Sprintf(`
set -e
sudo mv "%s" "%s"
sudo nginx -t
sudo systemctl reload nginx
`, tmp, target)

	res, err := c.Run(cmd)
	if err != nil {
		return err
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("apply nginx failed (exit %d): %s", res.ExitCode, res.Stderr)
	}
	return nil
}
