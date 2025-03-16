package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var config *Config

var configPath string;

func SetConfigPath(path string) {
	configPath = path
}

// Config struct for server, database, oauth and jwt configuration
func GetConfig() *Config {
	if config != nil {
		return config
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		log.Println("Error loading config file", err)
		log.Println(configPath)
		return nil
	}

	return config
}

func ReloadConfig() (*Config, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return config, nil
}

type Config struct {
	// Server configuration
	Server struct {
		Host string `default:"localhost" yaml:"host"`
		Port string `default:"8080" yaml:"port"`
	}
	// Database configuration
	Database struct {
		Host     string `default:"localhost" yaml:"host"`
		Port     string `default:"5432" yaml:"port"`
		User     string `default:"root" yaml:"user"`
		Password string `default:"" yaml:"password"`
		Name     string `default:"stream_bingo" yaml:"name"`
	}
	// OAuth configuration
	OAuth struct {
		Providers map[string]OAuthProvider `yaml:"providers"`
	}
	// JWT configuration
	JWT struct {
		Secret string `default:"secret" yaml:"secret"`
	}
}

// OAuthProvider config for multiple providers
type OAuthProvider struct {
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	Endpoint     struct {
		TokenURL string `yaml:"token_url"`
		AuthURL  string `yaml:"auth_url"`
		RedirectURL string `yaml:"redirect_url"`
	} `yaml:"endpoints"`
	Scopes       []string `yaml:"scopes"`	
	UserURL	  string   `yaml:"user_url"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
