package main

import (
	"os"
)

func main() {
	// set up logger
	// parse args for auth info, log level?
	// create an instance of app logic
	// initiate connection to slack
	var s = new(SlackConnection)
	testSlackInterface(s)
}

// Retrieve authorization credentials from environment variable
// SLACK_ALARMIE_TOKEN
func getAuthCredentialsFromEnv() string {
	// will add error handling, logging. WIP
	return os.Getenv("SLACK_ALARMIE_TOKEN")
}
