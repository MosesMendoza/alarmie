package main

import (
	"fmt"
	"os"

	log "github.com/inconshreveable/log15"
)

func main() {
	var application = new(AlarmieApplication)

	// set up configuration
	config := new(Config)
	configError := config.LoadFromEnvironment()
	if configError != nil {
		fmt.Printf(configError.Error())
		os.Exit(1)
	}
	application.Config = config

	logger := log.New("Alarmie", "application startup")

	// TODO: Something has to actually make the log file on first run
	// Note: this is the insertion point for where we'll set log level
	logHandler := log.LvlFilterHandler(log.LvlInfo, log.FailoverHandler(
		// Try to open log file and use it, but fall back to stdout on error
		log.Must.FileHandler(config.LogFilePath, log.LogfmtFormat()),
		log.StdoutHandler))

	logger.SetHandler(logHandler)

	application.Logger = logger

	application.Logger.Debug("Initialized logging via log handler %v", logHandler)

	// initiate connection to slack
	// Retrieve authorization credentials from environment variable
	// SLACK_ALARMIE_TOKEN
	//Note: This token is NOT maintained in any internal state
	// anywhere in the application. It is passed directly to the slack connection.
	token := os.Getenv("SLACK_ALARMIE_TOKEN")
	if token == "" {
		application.Logger.Crit("Could not find SLACK_ALARMIE_TOKEN in environment. Do `export SLACK_ALARMIE_TOKEN=<token>` and run again.")
		os.Exit(1)
	}

	connector := &SlackConnection{logger}
	connectionContext, error := connector.Connect(token, "foo")

	if error != nil {
		application.Logger.Crit("Alarmie could not initialize its connection to Slack")
		os.Exit(1)
	}

	application.Connector = connector
	application.Context = connectionContext

	fmt.Printf("context: %v\nerror: %s", connectionContext, error)
}
