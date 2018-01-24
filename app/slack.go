// Our basic API/interface to slack
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/inconshreveable/log15"
	"golang.org/x/net/websocket"
)

// Slack is the basic connection interface
type Slack interface {
	Connect() bool
	SendMessage(string) bool
	GetMessage() string
}

// SlackConnection is the concrete implementation of the basic connection
// interface
type SlackConnection struct {
	logger log15.Logger
}

// Connect instantiates the slack connection instance and returns an
// object representing the connection.
func (s SlackConnection) Connect(token string, slackAPIURL string) (*RtmConnectionContext, error) {
	// TODO move this URL out to config after rewriting config class

	rtmConnectResponse, error := s.GetSecureRtmConnectionInfo(token, slackAPIURL)
	if error != nil {
		return nil, error
	}

	var temporaryURL = rtmConnectResponse.WsURL
	var connectedAsID = rtmConnectResponse.Self.ID
	var domain = rtmConnectResponse.Team.Domain
	var teamName = rtmConnectResponse.Team.Name
	var name = rtmConnectResponse.Self.Name

	socket, error := s.InitiateWebsocketConnection(temporaryURL)
	if error != nil {
		return nil, error
	}

	s.logger.Debug("Successfully initiated connection to Slack")
	s.logger.Debug("Connected as ID: %s", connectedAsID)
	s.logger.Debug("My Name: %s", name)
	s.logger.Debug("Domain: %s", domain)
	s.logger.Debug("Team Name: %s", teamName)

	context := &RtmConnectionContext{ID: connectedAsID,
		SocketConnection: socket,
		Domain:           domain,
		TeamName:         teamName,
		Name:             name}

	return context, nil
}

// GetSecureRtmConnectionInfo obtains the Rtm info needed to generate a websocket connection
func (s SlackConnection) GetSecureRtmConnectionInfo(token string, slackAPIURL string) (*RtmConnectResponse, error) {
	s.logger.Info("Initiating connection to Slack")
	url := fmt.Sprintf("%s?token=%s", slackAPIURL, token)
	response, error := http.Get(url)

	if error != nil || response.StatusCode != http.StatusOK {
		s.logger.Crit("Could not initiate RTM connection to Slack: %s", error.Error())
		s.logger.Crit("HTTP response code: %d", response.StatusCode)
		return nil, error
	}
	s.logger.Debug("%s: %d", response.Status, response.StatusCode)

	body, error := s.readHTTPBody(response)
	if error != nil {
		return nil, error
	}

	rtmConnectResponse := new(RtmConnectResponse)
	unmarshallError := json.Unmarshal(body, rtmConnectResponse)

	if error != nil {
		s.logger.Crit("Could not deserialize connection response object: %s", unmarshallError.Error())
		return nil, unmarshallError
	}
	return rtmConnectResponse, nil
}

// ReadHTTPBody is a simple helper to return the byte-array content of an http
// response object
func (s SlackConnection) readHTTPBody(response *http.Response) ([]byte, error) {
	defer response.Body.Close()

	body, error := ioutil.ReadAll(response.Body)

	if error != nil {
		s.logger.Crit("Could not read body from RTM connection response, %s", error.Error())
		return nil, error
	}

	return body, nil
}

// InitiateWebsocketConnection creats a websocket connection given a URL
func (s SlackConnection) InitiateWebsocketConnection(url string) (*websocket.Conn, error) {
	socket, error := websocket.Dial(url, "", "https://slack.com")
	if error != nil {
		s.logger.Crit("Could not initiate websocket connection: %s", error.Error())
		return nil, error
	}

	return socket, error
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

// Models, including representing Slack API

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

// RtmConnectionContext has the information we need to identify ourself and
// send/receive messages
type RtmConnectionContext struct {
	ID               string
	Name             string
	TeamName         string
	Domain           string
	SocketConnection *websocket.Conn
}
