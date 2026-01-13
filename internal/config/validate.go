package config

import "errors"

func (c *Config) Validate() error {
	if c.App == "" {
		return errors.New("app name is required")
	}
	if c.Target.Host == "" {
		return errors.New("target.host is required")
	}
	if c.Target.User == "" {
		return errors.New("target.user is required")
	}
	if c.Target.Port == 0 {
		c.Target.Port = 22
	}
	if c.Target.KeyPath == "" {
		return errors.New("target.keyPath is required")
	}
	if c.Start.Cmd == "" {
		return errors.New("start.cmd is required")
	}
	return nil
}
