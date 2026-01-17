package cli

import (
	"fmt"
	"os"
	"path"
	"simpledeploy/internal/config"
	"simpledeploy/internal/daemon"
	"simpledeploy/internal/nginx"
	"simpledeploy/internal/packageutil"
	"simpledeploy/internal/runner"
	"simpledeploy/internal/systemd"
	"simpledeploy/remote"
	"time"
)

func Deploy() error {
	cfg, err := config.Load("simpledeploy.yaml")
	if err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return err
	}

	if len(cfg.Build.Local) > 0 {
		if err := runner.RunLocalCommands(cfg.Build.Local, ".", nil); err != nil {
			return err
		}
	}

	// 2) Package build output
	archive := ".simpledeploy-release.tar.gz"
	if err := packageutil.CreateTarGz(archive, ".", cfg.Package.Include); err != nil {
		return err
	}
	defer os.Remove(archive)

	// 3) SSH connect
	fmt.Println("Connecting to server via SSH...")
	client, err := remote.Connect(cfg.Target.Host, cfg.Target.User, cfg.Target.Port, cfg.Target.KeyPath)
	if err != nil {
		return err
	}
	defer client.Close()

	// Choose mode

	if cfg.Type == "static" {
		// 1) Deploy the static artifact first (creates releases/<id>, switches current symlink)
		if err := deployStatic(cfg, client, archive); err != nil {
			return err
		}

		// 2) Then apply nginx config for this app (server_name -> /var/www/<app>/current)
		conf := nginx.GenerateStatic(cfg)
		confName := nginx.ConfName(cfg.App)
		if err := client.ApplyNginxConf(confName, conf); err != nil {
			return err
		}

		return nil
	} else if cfg.Type == "node" {
		if err := deployNode(cfg, client, archive); err != nil {
			return err
		}
	}
	return nil
}

func deployStatic(cfg *config.Config, client *remote.Client, archive string) error {
	// Release paths
	releaseID := fmt.Sprintf("%d", time.Now().Unix())
	webRoot := cfg.Static.WebRoot
	releasesDir := path.Join(webRoot, "releases")
	releaseDir := path.Join(releasesDir, releaseID)

	// Upload archive to release dir
	if err := client.MkdirAll(releaseDir); err != nil {
		return err
	}

	remoteArchive := path.Join(releaseDir, "release.tar.gz")
	fmt.Println("Uploading artifact...")
	if err := client.UploadFile(archive, remoteArchive); err != nil {
		return err
	}

	// Extract on server
	fmt.Println("Extracting...")
	r, err := client.Run(fmt.Sprintf(`cd "%s" && tar -xzf release.tar.gz`, releaseDir))
	if err != nil {
		return err
	}
	if r.ExitCode != 0 {
		return fmt.Errorf("extract failed (exit %d): %s", r.ExitCode, r.Stderr)
	}

	// Point current -> releaseDir/dist
	current := path.Join(webRoot, "current")
	distPath := path.Join(releaseDir, cfg.Static.DistDir)

	// ln -sfn makes the symlink atomically switch
	_, _ = client.Run(fmt.Sprintf(`ln -sfn "%s" "%s"`, distPath, current))

	_, _ = client.Run(`sudo systemctl reload nginx`)

	fmt.Println("✅ Static site deployed.")
	fmt.Printf("Serving from: %s\n", current)
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

func deployNode(cfg *config.Config, client *remote.Client, archive string) error {
	releaseID := fmt.Sprintf("%d", time.Now().Unix())

	base := path.Join("/home", cfg.Target.User, "simpledeploy", "apps", cfg.App)
	releasesDir := path.Join(base, "releases")
	releaseDir := path.Join(releasesDir, releaseID)
	current := path.Join(base, "current")
	logsDir := path.Join(base, "logs")
	logFile := path.Join(logsDir, "app.log")

	// Prepare dirs
	if err := client.MkdirAll(releaseDir); err != nil {
		return err
	}
	if err := client.MkdirAll(logsDir); err != nil {
		return err
	}
	// Give ubuntu ownership for extract/install steps
	_, _ = client.Run(fmt.Sprintf(`sudo chown -R %s:%s "%s"`, cfg.Target.User, cfg.Target.User, base))

	// Upload artifact to release dir
	remoteArchive := path.Join(releaseDir, "release.tar.gz")
	fmt.Println("Uploading artifact...")
	if err := client.UploadFile(archive, remoteArchive); err != nil {
		return err
	}

	// Extract
	fmt.Println("Extracting...")
	res, err := client.Run(fmt.Sprintf(`cd "%s" && tar -xzf release.tar.gz`, releaseDir))
	if err != nil {
		return err
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("extract failed (exit %d): %s", res.ExitCode, res.Stderr)
	}

	// Switch current symlink
	_, _ = client.Run(fmt.Sprintf(`sudo ln -sfn "%s" "%s"`, releaseDir, current))

	// Install prod deps on server
	fmt.Println("Installing dependencies on server...")
	res, err = client.Run(fmt.Sprintf(`cd "%s" && %s`, current, cfg.Node.Install))
	if err != nil {
		return err
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("npm install failed (exit %d): %s", res.ExitCode, res.Stderr)
	}

	// Write systemd unit
	unit := systemd.GenerateUnit(cfg, current, logFile)
	unitName := systemd.ServiceName(cfg.App)
	unitPath := path.Join("/etc/systemd/system", unitName)

	fmt.Println("Configuring systemd...")
	if err := client.WriteFileSudo(unit, unitPath); err != nil {
		return err
	}
	if err := client.RestartSystemdService(unitName); err != nil {
		return err
	}

	// Apply nginx proxy config (per app)
	conf := nginx.GenerateNodeProxy(cfg)
	confName := nginx.ConfName(cfg.App)
	fmt.Println("Configuring nginx...")
	if err := client.ApplyNginxConf(confName, conf); err != nil {
		return err
	}

	fmt.Println("✅ Node app deployed.")
	for _, h := range cfg.Route.Hostnames {
		fmt.Printf("Public URL: http://%s/\n", h)
	}
	return nil
}
