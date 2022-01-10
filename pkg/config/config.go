package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Proxy struct {
		Port               string   `yaml:"port"`
		TrustedRootCaCerts []string `yaml:"trustedRootCaCerts"`
		ForgingRootCa      struct {
			Cert string `yaml:"cert"`
			Key  string `yaml:"key"`
		} `yaml:"forgingRootCa"`
	} `yaml:"proxy"`
	Paths struct {
		Doh     string `yaml:"doh"`
		Http    string `yaml:"http"`
		Session string `yaml:"session"`
	} `yaml:"paths"`
	Tamper struct {
		TypeA map[string][4]byte `yaml:"typeA"`
	} `yaml:"tamper"`
}

func ReadFile(configFile string) (*Config, error) {
	var cfg Config

	f, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", configFile, err)
	}
	defer func() { _ = f.Close() }()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("decodig config %w", err)
	}

	return &cfg, nil
}
