// Package config defines the config loading for the app
package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

type MailChimp struct {
	BaseURL string `yaml:"baseUrl"`
	ListID  string `yaml:"listId"`
	ApiKey  string `yaml:"apiKey"`
}

type Ometria struct {
	BaseURL string `yaml:"baseUrl"`
	ApiKey  string `yaml:"apiKey"`
}

type Ticker struct {
	Timer int `yaml:"timer"`
}

// AppConfig is main app config
type AppConfig struct {
	MailChimpAPI MailChimp `yaml:"mailchimpapi"`
	OmetriaAPI   Ometria   `yaml:"ometriaapi"`
	Ticker       Ticker    `yaml:"ticker"`
	PageLimit    int64     `yaml:"pagelimit"`
}

var (
	configFiles = []string{"config.yaml"}
	// Config contains all configuration values
	Config AppConfig
)

func searchConfig(dir string) (string, error) {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	dirPath := filepath.Join(absPath, "config")
	if fileInfo, err := os.Stat(dirPath); err == nil && fileInfo.IsDir() {
		return dirPath, nil
	}

	if absPath == "/" {
		return "", errors.New("not found")
	}

	return searchConfig(filepath.Join(absPath, ".."))
}

func init() {

	var err error
	var configDir string

	if configDir, err = searchConfig("."); err != nil {
		panic("Config dir not found")
	}

	for i, v := range configFiles {
		configFiles[i] = filepath.Join(configDir, v)
	}

	config := configor.New(&configor.Config{ENVPrefix: "-"})
	if err := config.Load(&Config, configFiles...); err != nil {
		panic("Invalid config")
	}
}
