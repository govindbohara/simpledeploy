package remote

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/pkg/sftp"
)

func (c *Client) MkdirAll(remoteDir string) error {
	remoteDir = strings.ReplaceAll(remoteDir, `"`, `\"`)
	r, err := c.Run(fmt.Sprintf(`mkdir -p "%s"`, remoteDir))
	if err != nil {
		return err
	}
	if r.ExitCode != 0 {
		return fmt.Errorf("mkdir failed (exit %d): %s", r.ExitCode, r.Stderr)
	}
	return nil
}

func (c *Client) UploadFile(localPath string, remotePath string) error {
	s, err := sftp.NewClient(c.ssh)
	if err != nil {
		return fmt.Errorf("sftp client: %w", err)
	}
	defer s.Close()

	if err := c.MkdirAll(path.Dir(remotePath)); err != nil {
		return err
	}

	src, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("open local: %w", err)
	}
	defer src.Close()

	dst, err := s.Create(remotePath)
	if err != nil {
		return fmt.Errorf("create remote: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	return nil
}

func (c *Client) UploadBytes(content string, remotePath string) error {
	s, err := sftp.NewClient(c.ssh)
	if err != nil {
		return fmt.Errorf("sftp client: %w", err)
	}
	defer s.Close()

	if err := c.MkdirAll(path.Dir(remotePath)); err != nil {
		return err
	}

	f, err := s.Create(remotePath)
	if err != nil {
		return fmt.Errorf("create remote file: %w", err)
	}
	defer f.Close()

	_, err = f.Write([]byte(content))
	return err
}
