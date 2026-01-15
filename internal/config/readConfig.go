package config

type Config struct {
	App     string        `yaml:"app"`
	Target  TargetConfig  `yaml:"target"`
	Build   BuildConfig   `yaml:"build"`
	Package PackageConfig `yaml:"package"`
	Static  StaticConfig  `yaml:"static"`

	Start StartConfig `yaml:"start"`
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
	Cmd     string `yaml:"cmd"`
	Workdir string `yaml:"workdir"`
	Port    int    `yaml:"port"`
}

type StaticConfig struct {
	WebRoot string `yaml:"webRoot"`
	DistDir string `yaml:"distDir"`
}
