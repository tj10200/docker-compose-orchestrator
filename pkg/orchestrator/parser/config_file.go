package parser

import "time"

type Command struct {
	Type    string            `yaml:"type"`
	Image   string            `yaml:"image"`
	Env     map[string]string `yaml:"environment"`
	Cmd     string            `yaml:"command"`
	Timeout time.Duration     `yaml:"timeout"`
	NetName string            `yaml:"network_name"`

	Dir  string   `yaml:"dir"`
	Tool string   `yaml:"tool"`
	Args []string `yaml:"args"`
}

type ConfigService struct {
	Name     string    `yaml:"name"`
	Commands []Command `yaml:"commands"`
	AfterRun []Command `yaml:"after_run"`
}

type ConfigFile struct {
	Services map[string]ConfigService `yaml:"services"`
}
