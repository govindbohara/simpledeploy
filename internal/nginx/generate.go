package nginx

import (
	"fmt"
	"strings"

	"simpledeploy/internal/config"
)

func ConfName(app string) string {
	app = sanitize(app)
	return fmt.Sprintf("simpledeploy-%s.conf", sanitize(app))
}

func GenerateStatic(cfg *config.Config) string {
	serverNames := strings.Join(cfg.Route.Hostnames, " ")
	root := fmt.Sprintf("%s/current", cfg.Static.WebRoot)

	// SPA routing for React Router
	return fmt.Sprintf(`server {
    listen 80;
    server_name %s;

    root %s;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }
}
`, serverNames, root)
}

func sanitize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	return s
}
