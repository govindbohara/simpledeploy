package remote

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	ssh *ssh.Client
}

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func Connect(host string, user string, port int, keyPath string) (*Client, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parse key: %w", err)
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	sshClient, err := ssh.Dial("tcp", address, config)

	if err != nil {
		return nil, fmt.Errorf("ssh dial: %w", err)
	}

	return &Client{ssh: sshClient}, nil

}

func (c *Client) Close() error {
	return c.ssh.Close()
}

func (c *Client) Run(command string) (*Result, error) {
	session, err := c.ssh.NewSession()
	if err != nil {
		return nil, fmt.Errorf("new session: %w", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(command)

	result := &Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if err == nil {
		result.ExitCode = 0
		return result, nil
	}

	if ee, ok := err.(*ssh.ExitError); ok {
		result.ExitCode = ee.ExitStatus()
		return result, nil
	}

	return nil, fmt.Errorf("run remote command: %w", err)

}
