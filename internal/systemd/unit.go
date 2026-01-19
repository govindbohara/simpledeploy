package systemd

import (
	"fmt"
	"simpledeploy/internal/config"
)

func ServiceName(app string) string {
	return fmt.Sprintf("simpledeploy-%s.service", app)
}

func GenerateUnit(cfg *config.Config, workDir string, logPath string) string {
	return fmt.Sprintf(`[Unit]
Description=SimpleDeploy %s
After=network.target

[Service]
Type=simple
User=%s
WorkingDirectory=%s
Environment=NODE_ENV=production
Environment=PORT=%d
ExecStart=%s
Restart=always
RestartSec=2
StandardOutput=append:%s
StandardError=append:%s

[Install]
WantedBy=multi-user.target
`, cfg.App, cfg.Target.User, workDir, cfg.Node.Port, cfg.Node.Start, logPath, logPath)
}
