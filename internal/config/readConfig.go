package config

type Config struct {
	App     string `yaml:"app"`
	Target  TargetConfig
	Build   BuildConfig
	Package PackageConfig
	Start   StartConfig
	Env     map[string]string `yaml:"env"`
}

type TargetConfig struct {
	Host    string `yaml:"host"`
	User    string `yaml:"user"`
	Port    int    `yaml:"port"`
	KeyPath string `yaml:"keyPath"`
}

type BuildConfig struct {
	Local []string `yaml:"local"`
}

type PackageConfig struct {
	Include []string `yaml:"include"`
}

type StartConfig struct {
	Cmd string `yaml:"cmd"`
}
