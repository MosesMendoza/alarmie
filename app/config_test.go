package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

func TestLoadFromEnvironmentGetsLogfilePath(t *testing.T) {
	t.Log("Test that LoadFromEnvironment can retrieve log file path")
	os.Setenv("ALARMIE_LOGFILEPATH", "/foo/bar/baz")
	config := new(Config)
	configError := config.LoadFromEnvironment()
	if configError != nil {
		t.Errorf("Expected config.LoadFromEnvironment not to error, but got %s", configError.Error())
	}
	if config.LogFilePath != "/foo/bar/baz" {
		t.Errorf("Expected config.LogFilePath of /foo/bar/baz but got %s", config.LogFilePath)
	}
}

func TestLoadFromEnvironmentReturnsEffectiveErrorMessage(t *testing.T) {
	t.Log("Test that LoadFromEnvironment errors with a helpful message on missing value")
	os.Setenv("ALARMIE_LOGFILEPATH", "")
	config := new(Config)
	configError := config.LoadFromEnvironment()
	expectedMessage := regexp.MustCompile("Could not locate values for configuration in environment")

	if !(expectedMessage.MatchString(configError.Error())) {
		t.Errorf("Expected message to match %s but instead got %s", expectedMessage, configError.Error())
	}
}

func TestLoadFromFile(t *testing.T) {
	t.Log("Test that load from environment deserializes a config file correctly")
	var temporaryConfigFile, fileError = ioutil.TempFile("", "test")

	defer TearDown(temporaryConfigFile)

	expectedLogPath := "/baz/qux/quux"

	if fileError != nil {
		t.Errorf("Unable to create temp file for Config test: %s", fileError.Error())
	}

	temporaryConfigFile.WriteString(fmt.Sprintf("{\"LogFilePath\":\"%s\"}", expectedLogPath))

	config := new(Config)
	configError := config.LoadFromFile(temporaryConfigFile.Name())

	if configError != nil {
		t.Errorf("Expected to be able to load temp config file without error: %s", configError.Error())
	}

	if config.LogFilePath != expectedLogPath {
		t.Errorf("Expected LogFilePath to be %s, not %s", expectedLogPath, config.LogFilePath)
	}
}

func TearDown(file *os.File) {
	file.Close()
	var error = os.Remove(file.Name())
	if error != nil {
		panic(fmt.Sprintf("Could not remove temporary file %s in test teardown", file.Name()))
	}
}
