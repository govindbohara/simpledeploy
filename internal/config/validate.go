package config

import "errors"

func (c *Config) Validate() error {
	if c.App == "" {
		return errors.New("app is required")
	}

	if c.Target.Host == "" {
		return errors.New("target.host is required")
	}
	if c.Target.User == "" {
		return errors.New("target.user is required")
	}
	if c.Target.KeyPath == "" {
		return errors.New("target.keyPath is required")
	}
	if c.Target.Port == 0 {
		c.Target.Port = 22
	}

	// Decide mode:
	isStatic := c.Static.WebRoot != "" || c.Static.DistDir != ""

	if isStatic {
		// Static deploy needs these:
		if c.Static.WebRoot == "" {
			return errors.New("static.webRoot is required for static deploy")
		}
		if c.Static.DistDir == "" {
			return errors.New("static.distDir is required for static deploy")
		}
		if len(c.Package.Include) == 0 {
			// default to distDir
			c.Package.Include = []string{c.Static.DistDir}
		}
		return nil
	}

	// Process deploy mode:
	if c.Start.Cmd == "" {
		return errors.New("start.cmd is required")
	}
	if c.Start.Port == 0 {
		c.Start.Port = 3000
	}
	if c.Start.Workdir == "" {
		c.Start.Workdir = "."
	}
	if len(c.Package.Include) == 0 {
		return errors.New("package.include is required")
	}

	return nil
}
