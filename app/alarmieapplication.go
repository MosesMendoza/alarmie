package main

import "github.com/inconshreveable/log15"

// AlarmieApplication is the highest-level construct, which owns the
// configuration, slack etc
type AlarmieApplication struct {
	Config    *Config
	Logger    log15.Logger
	Connector *SlackConnection
	Context   *RtmConnectionContext
}
