package config

import "errors"

func (c *Config) Validate() error {
	if c.App == "" {
		return errors.New("app is required")
	}
	if c.Type == "" {
		c.Type = "static"
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
	if len(c.Route.Hostnames) == 0 {
		return errors.New("route.hostnames is required for static deploy")
	}
	if len(c.Package.Include) == 0 {
		c.Package.Include = []string{c.Static.DistDir}
	}

	if c.Type == "static" {
		if c.Static.WebRoot == "" {
			return errors.New("static.webRoot is required")
		}
		if c.Static.DistDir == "" {
			return errors.New("static.distDir is required")
		}

		return nil
	}
	if c.Type == "node" {
		if c.Node.Port == 0 {
			return errors.New("node.port is required")
		}
		if c.Node.Install == "" {
			c.Node.Install = "npm ci --omit=dev"
		}
		if c.Node.Start == "" {
			return errors.New("node.start is required")
		}

		return nil
	}

	return errors.New("Please provide a valid type: " + c.Type)
}
