package runner

import (
	"fmt"
	"os"
	"os/exec"
)

func RunLocalCommands(commands []string, workDir string, envs map[string]string) error {
	for _, command := range commands {
		if err := RunOneCommand(command, workDir, envs); err != nil {
			return err
		}
	}
	return nil
}

func RunOneCommand(command string, workDir string, envs map[string]string) error {
	fmt.Println("Start: ", command)
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = workDir
	cmd.Env = os.Environ()

	for k, v := range envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("command failed: %q (exit code %d)", command, ee.ExitCode())
		}
		return fmt.Errorf("failed to run command: %q (%v)", command, err)
	}

	return nil
}
