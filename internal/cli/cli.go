package cli

import (
	"fmt"
	"os"
	"path"
	"simpledeploy/internal/config"
	"simpledeploy/internal/daemon"
	"simpledeploy/internal/packageutil"
	"simpledeploy/internal/runner"
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
	isStatic := cfg.Static.WebRoot != "" && cfg.Static.DistDir != ""
	if isStatic {
		return deployStatic(cfg, client, archive)
	}

	// Otherwise: process deploy (Node backend etc.)
	return fmt.Errorf("process deploy not implemented in this version (static deploy configured)")
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
