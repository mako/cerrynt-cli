// Package config handles loading and parsing of the cerrynt configuration file.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the top-level structure of the cerrynt config file.
type Config struct {
	APIBaseURL string `yaml:"api_base_url"`
	AuthToken  string `yaml:"auth_token"`
	Feeds      []Feed `yaml:"feeds"`
}

// Feed represents a single RSS feed subscription as stored in config.
type Feed struct {
	ID    string `yaml:"id"`
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
}

// Load reads and parses the YAML config file at the given path.
// It returns an error for any problem, including when the file does not exist
// (os.IsNotExist will be true in that case). The caller is responsible for
// deciding what to do when the file is missing.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	dec := yaml.NewDecoder(f)
	dec.KnownFields(true) // surface unknown keys as errors
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: parse %q: %w", path, err)
	}

	return &cfg, nil
}
