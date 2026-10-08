package taco

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config mirrors taco.yaml.
type Config struct {
	Tacos         int                 `mapstructure:"tacos"`
	NumberOfTacos int                 `mapstructure:"numberOfTacos"`
	Message       string              `mapstructure:"message"`
	Emoji         string              `mapstructure:"emoji"`
	Template      string              `mapstructure:"template"`
	Teams         map[string][]string `mapstructure:"teams"`
}

// TacoCount returns the configured taco count, falling back to fallback.
func (c Config) TacoCount(fallback int) int {
	if c.NumberOfTacos > 0 {
		return c.NumberOfTacos
	}
	if c.Tacos > 0 {
		return c.Tacos
	}
	return fallback
}

func LoadConfig(path string, explicit bool) (Config, error) {
	cfg := Config{}

	configPath, err := resolveConfigPath(path, explicit)
	if err != nil {
		return cfg, err
	}
	if configPath == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}

	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

func resolveConfigPath(path string, explicit bool) (string, error) {
	if path == "" {
		path = "taco.yaml"
	}

	if filepath.IsAbs(path) || strings.Contains(path, string(os.PathSeparator)) {
		if _, err := os.Stat(path); err != nil {
			if errors.Is(err, os.ErrNotExist) && !explicit {
				return "", nil
			}
			return "", fmt.Errorf("read config: %w", err)
		}
		return path, nil
	}

	candidates := []string{path}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, "."+path), filepath.Join(home, path))
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("read config: %w", err)
		}
	}

	if explicit {
		return "", fmt.Errorf("read config: %s was not found", path)
	}
	return "", nil
}
