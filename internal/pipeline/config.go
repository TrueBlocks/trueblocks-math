package pipeline

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	API       APIConfig       `yaml:"api"`
	Models    ModelsConfig    `yaml:"models"`
	Pipeline  PipelineConfig  `yaml:"pipeline"`
	Dashboard DashboardConfig `yaml:"dashboard"`
}

type APIConfig struct {
	AnthropicKey string `yaml:"anthropic_key"`
}

type ModelsConfig struct {
	Research  string `yaml:"research"`
	Outline   string `yaml:"outline"`
	Draft     string `yaml:"draft"`
	Factcheck string `yaml:"factcheck"`
	Draft2    string `yaml:"draft2"`
}

type PipelineConfig struct {
	MaxPerCycle   int     `yaml:"max_per_cycle"`
	NewPerCycle   int     `yaml:"new_per_cycle"`
	DryRun        bool    `yaml:"dry_run"`
	CycleInterval int     `yaml:"cycle_interval"`
	APITimeout    int     `yaml:"api_timeout"`
	Concurrency   int     `yaml:"concurrency"`
	Verbose       bool    `yaml:"verbose"`
	ReadMean      float64 `yaml:"read_mean"`
	ReadSpread    float64 `yaml:"read_spread"`
}

type DashboardConfig struct {
	Port int `yaml:"port"`
}

func DefaultConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "trueblocks", "math", "config.yaml")
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}

	cfg := &Config{
		Models: ModelsConfig{
			Research:  "claude-sonnet-4-20250514",
			Outline:   "claude-sonnet-4-20250514",
			Draft:     "claude-sonnet-4-20250514",
			Factcheck: "claude-sonnet-4-20250514",
			Draft2:    "claude-sonnet-4-20250514",
		},
		Pipeline: PipelineConfig{
			MaxPerCycle:   6,
			NewPerCycle:   4,
			DryRun:        true,
			CycleInterval: 15,
			APITimeout:    300,
			Concurrency:   3,
			ReadMean:      5.0,
			ReadSpread:    1.0,
		},
		Dashboard: DashboardConfig{
			Port: 8787,
		},
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}

	return cfg, nil
}

func SaveConfig(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}
