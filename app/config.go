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
func (c Config) LoadFromEnvironment() error {
	const logFilePath = "ALARMIE_LOGFILEPATH"
	c.LogFilePath = os.Getenv("ALARMIE_LOGFILEPATH")
	if c.LogFilePath == "" {
		return fmt.Errorf("Could not locate values for configuration in environment: %s", logFilePath)
	}
	return nil
}

// LoadFromFile will populate the values of the given Config pointer from the
// JSON file located at the supplied path. Will return a non-nil error on
// failure.
func (c Config) LoadFromFile(configFilePath string) error {
	var handle, error = ioutil.ReadFile(configFilePath)
	if error != nil {
		fmt.Printf("Could not load configuration from file %s: %s", configFilePath, error.Error())
		return error
	}

	var unmarshalError = json.Unmarshal(handle, &c)
	if unmarshalError != nil {
		fmt.Printf("Could not deserialize JSON config file %s: %s", configFilePath, error.Error())
	}
	return unmarshalError
}
