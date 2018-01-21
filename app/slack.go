// Our basic API/interface to slack
package main

// Slack is the basic connection interface
type Slack interface {
	Connect() bool
	SendMessage(string) bool
	GetMessage() string
}

// SlackConnection is the concrete implementation of the basic connection
// interface
type SlackConnection struct{}

// Connect instantiates the slack connection instance and maintains the
// websocket reference internally
func (s SlackConnection) Connect() bool {
	return true
}

// SendMessage takes a string and sends it over the websocket connection
func (s SlackConnection) SendMessage(message string) bool {
	return true
}

// GetMessage retrieves the next message from the internal websocket connection.
// This function acts like a blocking queue/channel - retrieving a message when
// the queue is empty will block
func (s SlackConnection) GetMessage() string {
	return "foo"
}

// Models, particularly Slack API

// Message corresponds 1:1 to the slack Message object
type Message struct {
	Type      string `json:"type"`    // I think this is always "message"?
	Channel   string `json:"channel"` // name of slack channel message was on
	User      string `json:"user"`    // user sending message
	Text      string `json:"text"`    // content of message
	Timestamp string `json:"ts"`      // timestamp
}

// RtmConnectResponse corresponds to the API response from rtm.connect
type RtmConnectResponse struct {
	Ok    bool           `json:"ok"`   // did we successfully connect
	Self  RtmConnectSelf `json:"self"` // contains id, name response
	Team  RtmConnectTeam `json:"team"` // info on team
	WsURL string         `json:"url"`  // websocket url to connect to - in <= 30 seconds
}

// RtmConnectSelf is one of the member fields on RtmConnectResponse
type RtmConnectSelf struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RtmConnectTeam is one of the member fields on RtmConnectReponse
type RtmConnectTeam struct {
	Domain string `json:"domain"`
	ID     string `json:"id"`
	Name   string `json:"name"`
}
