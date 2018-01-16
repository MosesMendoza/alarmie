package main

import (
	"errors"
	"fmt"
	"os"

	log "github.com/inconshreveable/log15"
)

func main() {
	var application = new(AlarmieApplication)

	// set up configuration
	config := new(Config)
	var configError = config.LoadFromEnvironment()
	if configError != nil {
		panic(errors.New("Could not locate values for configuration in environment"))
	}

	application.Config = config

	// set up logger
	logger := log.New("Alarmie", "application startup")
	// TODO: Have to actually make the log file on first run
	logHandler, error := log.FileHandler(config.LogFilePath, log.LogfmtFormat())
	if error != nil {
		panic(fmt.Errorf("Could not open log file handle: %s", error.Error()))
	} else {
		logger.SetHandler(logHandler)
		application.Logger = logger
	}
	application.Logger.Warn("Started application")
	// initiate connection to slack
}

// Retrieve authorization credentials from environment variable
// SLACK_ALARMIE_TOKEN
func getAuthCredentialsFromEnv() string {
	// will add error handling, logging. WIP
	return os.Getenv("SLACK_ALARMIE_TOKEN")
}
