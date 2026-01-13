package cli

import (
	"fmt"
	"os"
	"simpledeploy/internal/config"
	"simpledeploy/internal/daemon"
	"simpledeploy/internal/runner"
	"simpledeploy/internal/runner/remote"
)

func Deploy() error {
	cfg, err := config.Load("simpledeploy.yaml")

	if err != nil {
		fmt.Println("Failed to read File")
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	// 1) Local build
	if len(cfg.Build.Local) > 0 {
		if err := runner.RunLocalCommands(cfg.Build.Local, ".", nil); err != nil {
			return err
		}
	}

	client, err := remote.Connect(cfg.Target.Host, cfg.Target.User, cfg.Target.Port, cfg.Target.KeyPath)

	if err != nil {
		return err
	}
	defer client.Close()
	fmt.Println("Connecting to server via SSH...")

	r, err := client.Run("uname -a")

	if err != nil {
		return err
	}
	if r.ExitCode != 0 {
		return fmt.Errorf("remote command failed (exit %d): %s", r.ExitCode, r.Stderr)
	}
	fmt.Println("Remote OK:", r.Stdout)

	return nil

}

func Status() error {
	fmt.Println("Status command not implemented yet")
	return nil
}

func Ping() error {
	if err := daemon.Ping(); err != nil {
		return err
	}

	fmt.Println("daemon is alive")
	return nil
}
