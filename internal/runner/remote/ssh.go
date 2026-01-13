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
	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	c, err := ssh.Dial("tcp", addr, cfg)

	if err != nil {
		return nil, fmt.Errorf("ssh dial: %w", err)
	}

	return &Client{ssh: c}, nil

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

	res := &Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if err == nil {
		res.ExitCode = 0
		return res, nil
	}

	if ee, ok := err.(*ssh.ExitError); ok {
		res.ExitCode = ee.ExitStatus()
		return res, nil
	}

	return nil, fmt.Errorf("run remote command: %w", err)

}
