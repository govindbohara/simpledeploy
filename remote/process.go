package remote

import (
	"fmt"
	"strings"
)

type StartResult struct {
	Pid int
}

func (c *Client) StartNodeApp(appDir string, cmd string, port int) (*StartResult, error) {
	// Writes:
	// - logs/app.log
	// - app.pid
	start := fmt.Sprintf(`
set -e
cd "%s"
mkdir -p logs
export PORT=%d
nohup %s > logs/app.log 2>&1 & echo $! > app.pid
cat app.pid
`, escape(appDir), port, cmd)

	r, err := c.Run(start)
	if err != nil {
		return nil, err
	}
	if r.ExitCode != 0 {
		return nil, fmt.Errorf("start failed (exit %d): %s", r.ExitCode, r.Stderr)
	}

	// r.Stdout contains pid
	pid, convErr := parsePid(r.Stdout)
	if convErr != nil {
		return nil, convErr
	}

	return &StartResult{Pid: pid}, nil
}

func (c *Client) StopApp(appDir string) error {
	stop := fmt.Sprintf(`
set -e
cd "%s"
if [ -f app.pid ]; then
  kill $(cat app.pid) || true
  rm -f app.pid
fi
`, escape(appDir))

	r, err := c.Run(stop)
	if err != nil {
		return err
	}
	if r.ExitCode != 0 {
		return fmt.Errorf("stop failed (exit %d): %s", r.ExitCode, r.Stderr)
	}
	return nil
}

func escape(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
