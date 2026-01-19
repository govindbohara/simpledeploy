package remote

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/pkg/sftp"
)

func (client *Client) MkdirAll(remoteDir string) error {
	remoteDir = strings.ReplaceAll(remoteDir, `"`, `\"`)
	r, err := client.Run(fmt.Sprintf(`mkdir -p "%s"`, remoteDir))
	if err != nil {
		return err
	}
	if r.ExitCode != 0 {
		return fmt.Errorf("mkdir failed (exit %d): %s", r.ExitCode, r.Stderr)
	}
	return nil
}

func (client *Client) UploadFile(localPath string, remotePath string) error {
	sftpClient, err := sftp.NewClient(client.ssh)
	if err != nil {
		return fmt.Errorf("sftp client: %w", err)
	}
	defer sftpClient.Close()

	if err := client.MkdirAll(path.Dir(remotePath)); err != nil {
		return err
	}

	sourceFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("open local: %w", err)
	}
	defer sourceFile.Close()

	destinationFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("create remote: %w", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(sourceFile, sourceFile)
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	return nil
}

func (client *Client) UploadBytes(content string, remotePath string) error {
	sshClient, err := sftp.NewClient(client.ssh)
	if err != nil {
		return fmt.Errorf("sftp client: %w", err)
	}
	defer sshClient.Close()

	if err := client.MkdirAll(path.Dir(remotePath)); err != nil {
		return err
	}

	sftpFile, err := sshClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("create remote file: %w", err)
	}
	defer sftpFile.Close()

	_, err = sftpFile.Write([]byte(content))
	return err
}
