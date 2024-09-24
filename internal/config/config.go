package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type ConfigYAML struct {
	KeyPath             string `yaml:"key_path"`
	CertPath            string `yaml:"cert_path"`
	DBPath              string `yaml:"db_path"`
	Port                int    `yaml:"port"`
	PebbleNotifications bool   `yaml:"pebble_notifications"`
}

type Config struct {
	Key                        []byte
	Cert                       []byte
	DBPath                     string
	Port                       int
	PebbleNotificationsEnabled bool
}

// Validate opens and processes the given yaml file, and catches errors in the process
func Validate(filePath string) (Config, error) {
	const validationErr = "config file validation failed: %w"
	config := Config{}
	configYaml, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf(validationErr, err)
	}
	c := ConfigYAML{}
	if err := yaml.Unmarshal(configYaml, &c); err != nil {
		return Config{}, fmt.Errorf(validationErr, err)
	}
	if c.CertPath == "" {
		return Config{}, fmt.Errorf(validationErr, errors.New("`cert_path` is empty"))
	}
	cert, err := os.ReadFile(c.CertPath)
	if err != nil {
		return Config{}, fmt.Errorf(validationErr, err)
	}
	if c.KeyPath == "" {
		return Config{}, fmt.Errorf(validationErr, errors.New("`key_path` is empty"))
	}
	key, err := os.ReadFile(c.KeyPath)
	if err != nil {
		return Config{}, fmt.Errorf(validationErr, err)
	}
	if c.DBPath == "" {
		return Config{}, fmt.Errorf(validationErr, errors.New("`db_path` is empty"))
	}
	dbfile, err := os.OpenFile(c.DBPath, os.O_CREATE|os.O_RDONLY, 0o644)
	if err != nil {
		return Config{}, fmt.Errorf(validationErr, err)
	}
	err = dbfile.Close()
	if err != nil {
		return Config{}, fmt.Errorf(validationErr, err)
	}
	if c.Port == 0 {
		return Config{}, fmt.Errorf(validationErr, errors.New("`port` is empty"))
	}
	if c.PebbleNotifications {
		_, err := exec.LookPath("pebble")
		if err != nil {
			return Config{}, fmt.Errorf(validationErr, errors.New("pebble binary not found"))
		}
	}

	config.Cert = cert
	config.Key = key
	config.DBPath = c.DBPath
	config.Port = c.Port
	config.PebbleNotificationsEnabled = c.PebbleNotifications
	return config, nil
}
