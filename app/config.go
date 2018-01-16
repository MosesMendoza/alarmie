package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config is a basic struct we'll populate with configuration properties
type Config struct {
	LogFilePath string
}

// LoadFromEnvironment will populate the values of the given Config pointer from
// the pre-known keys in the environment, and return non-nil error on failure
func (config *Config) LoadFromEnvironment() error {
	const logFilePath = "ALARMIE_LOGFILEPATH"
	config.LogFilePath = os.Getenv("ALARMIE_LOGFILEPATH")
	if config.LogFilePath == "" {
		return fmt.Errorf("Could not locate values for configuration in environment: %s", logFilePath)
	}
	return nil
}

// LoadFromFile will populate the values of the given Config pointer from the
// JSON file located at the supplied path. Will return a non-nil error on
// failure.
func (config *Config) LoadFromFile(configFilePath string) error {
	var handle, error = ioutil.ReadFile(configFilePath)
	if error != nil {
		fmt.Printf("Could not load configuration from file %s: %s", configFilePath, error.Error())
		return error
	}

	var unmarshalError = json.Unmarshal(handle, config)
	if unmarshalError != nil {
		fmt.Printf("Could not deserialize JSON config file %s: %s", configFilePath, error.Error())
	}
	return unmarshalError
}
