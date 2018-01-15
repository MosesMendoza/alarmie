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
