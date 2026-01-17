package nginx

import (
	"fmt"
	"strings"

	"simpledeploy/internal/config"
)

func GenerateNodeProxy(cfg *config.Config) string {
	serverNames := strings.Join(cfg.Route.Hostnames, " ")

	return fmt.Sprintf(`server {
    listen 80;
    server_name %s;

    location / {
        proxy_pass http://127.0.0.1:%d;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
`, serverNames, cfg.Node.Port)
}
