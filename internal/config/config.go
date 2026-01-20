package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Temporal TemporalConfig `yaml:"temporal"`
	Nats     NatsConfig     `yaml:"nats"`
	Activity ActivityConfig `yaml:"activity"`
}

type TemporalConfig struct {
	Host      string `yaml:"host"`
	TaskQueue string `yaml:"task_queue"`
}

type NatsConfig struct {
	URL     string        `yaml:"url"`
	Timeout time.Duration `yaml:"timeout"`
}

type ActivityConfig struct {
	Timeout time.Duration `yaml:"timeout"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	cfg.applyEnvOverrides()
	return &cfg, nil
}

func (c *Config) applyEnvOverrides() {
	if v := os.Getenv("TEMPORAL_HOST"); v != "" {
		c.Temporal.Host = v
	}
	if v := os.Getenv("NATS_URL"); v != "" {
		c.Nats.URL = v
	}
}
