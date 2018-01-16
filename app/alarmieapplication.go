package main

import "github.com/inconshreveable/log15"

// AlarmieApplication is the highest-level construct, which owns the
// configuration, slack etc
type AlarmieApplication struct {
	Config *Config
	// this doesn't work right
	Logger log15.Logger
}
