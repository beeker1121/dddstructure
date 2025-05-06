package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
)

// APIEnvironment defines the API environment.
type APIEnvironment string

const (
	APIEnvironmentDevelop    APIEnvironment = "DEVELOP"
	APIEnvironmentProduction APIEnvironment = "PRODUCTION"
)

// Config defines the Go Todo API settings.
type Config struct {
	DBHost         string         `json:"db_host"`
	DBPort         string         `json:"db_port"`
	DBName         string         `json:"db_name"`
	DBUser         string         `json:"db_user"`
	DBPass         string         `json:"db_pass"`
	APIHost        string         `json:"api_host"`
	APIPort        string         `json:"api_port"`
	APIEnvironment APIEnvironment `json:"api_environment"`
	LogFile        string         `json:"log_file"`
	JWTSecret      string         `json:"jwt_secret"`
	JWTExpiryTime  time.Duration  `json:"jwt_expiry_time"`
	LimitDefault   uint           `json:"limit_default"`
	LimitMax       uint           `json:"limit_max"`
}

// ParseConfigFile parses the API configuration file.
func ParseConfigFile(filepath string) (*Config, error) {
	config := &Config{}

	// Read the config file.
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not find config.json file at the given path of %s", filepath)
	}

	// Try to unmarshal config file JSON into Config struct.
	if err := json.Unmarshal(file, config); err != nil {
		return nil, errors.New("failed to unmarshal JSON into struct")
	}

	return config, nil
}
