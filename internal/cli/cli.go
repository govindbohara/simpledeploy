package cli

import (
	"fmt"
	"os"
	"path"
	"simpledeploy/internal/config"
	"simpledeploy/internal/nginx"
	"simpledeploy/internal/packageutil"
	"simpledeploy/internal/runner"
	"simpledeploy/internal/systemd"
	"simpledeploy/remote"
	"time"
)

func Deploy() error {
	config, err := config.Load("simpledeploy.yaml")
	if err != nil {
		return err
	}
	if err := config.Validate(); err != nil {
		return err
	}

	if len(config.Build.Local) > 0 {
		if err := runner.RunLocalCommands(config.Build.Local, ".", nil); err != nil {
			return err
		}
	}

	archiveFile := ".simpledeploy-release.tar.gz"
	if err := packageutil.CreateTarGz(archiveFile, ".", config.Package.Include); err != nil {
		return err
	}
	defer os.Remove(archiveFile)

	// 3) SSH connect
	fmt.Println("Connecting to server via SSH...")
	client, err := remote.Connect(config.Target.Host, config.Target.User, config.Target.Port, config.Target.KeyPath)
	if err != nil {
		return err
	}
	defer client.Close()

	if config.Type == "static" {
		if err := deployStatic(config, client, archiveFile); err != nil {
			return err
		}

		conf := nginx.GenerateStatic(config)
		confName := nginx.ConfigurationName(config.App)
		if err := client.ApplyNginxConfiguration(confName, conf); err != nil {
			return err
		}

		return nil
	} else if config.Type == "node" {
		if err := deployNode(config, client, archiveFile); err != nil {
			return err
		}
	}
	return nil
}

func deployStatic(config *config.Config, client *remote.Client, archive string) error {
	releaseID := fmt.Sprintf("%d", time.Now().Unix())
	webRoot := config.Static.WebRoot
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
	distPath := path.Join(releaseDir, config.Static.DistDir)

	// ln -sfn makes the symlink
	_, _ = client.Run(fmt.Sprintf(`ln -sfn "%s" "%s"`, distPath, current))

	_, _ = client.Run(`sudo systemctl reload nginx`)

	fmt.Println("Static site deployed.")
	fmt.Printf("Serving from: %s\n", current)
	return nil
}

func deployNode(config *config.Config, client *remote.Client, archive string) error {
	releaseID := fmt.Sprintf("%d", time.Now().Unix())

	base := path.Join("/home", config.Target.User, "simpledeploy", "apps", config.App)
	releasesDir := path.Join(base, "releases")
	releaseDir := path.Join(releasesDir, releaseID)
	current := path.Join(base, "current")
	logsDir := path.Join(base, "logs")
	logFile := path.Join(logsDir, "app.log")

	if err := client.MkdirAll(releaseDir); err != nil {
		return err
	}
	if err := client.MkdirAll(logsDir); err != nil {
		return err
	}
	// Give ubuntu ownership for extract/install steps
	_, _ = client.Run(fmt.Sprintf(`sudo chown -R %s:%s "%s"`, config.Target.User, config.Target.User, base))

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
	res, err = client.Run(fmt.Sprintf(`cd "%s" && %s`, current, config.Node.Install))
	if err != nil {
		return err
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("npm install failed (exit %d): %s", res.ExitCode, res.Stderr)
	}

	// Write systemd unit for process manager
	unit := systemd.GenerateUnit(config, current, logFile)
	unitName := systemd.ServiceName(config.App)
	unitPath := path.Join("/etc/systemd/system", unitName)

	fmt.Println("Configuring systemd...")
	if err := client.WriteFileSudo(unit, unitPath); err != nil {
		return err
	}
	if err := client.RestartSystemdService(unitName); err != nil {
		return err
	}

	// Apply nginx proxy config
	configuration := nginx.GenerateNodeProxy(config)
	configurationName := nginx.ConfigurationName(config.App)
	fmt.Println("Configuring nginx...")
	if err := client.ApplyNginxConfiguration(configurationName, configuration); err != nil {
		return err
	}

	fmt.Println("Node app deployed.")
	for _, h := range config.Route.Hostnames {
		fmt.Printf("Public URL: http://%s/\n", h)
	}
	return nil
}
