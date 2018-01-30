package main

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"

	"math/rand"

	"github.com/MosesMendoza/alarmie/testUtils"
	"golang.org/x/net/websocket"
)

func TestInitiateWebsocketConnectionCreatesWebsocket(t *testing.T) {
	// setup
	reply := "foo"
	const route = "/"
	const listenAddress = "127.0.0.1"
	const scheme = "ws://"
	bindPort := strconv.Itoa(rand.Intn(9999))

	websocketURL := scheme + listenAddress + ":" + bindPort + route
	server := testUtils.StartTestWebsocketServer(route, bindPort, reply)

	logger := testUtils.GetTestLogger()
	slackConnection := SlackConnection{logger: logger}
	// teardown
	defer testUtils.StopTestWebsocketServer(server)

	connection, error := slackConnection.InitiateWebsocketConnection(websocketURL)
	if error != nil {
		t.Errorf("Could not dial websocket connection in slack_test: %s", error.Error())
		t.FailNow()
	}

	if reflect.TypeOf(connection) != reflect.TypeOf(new(websocket.Conn)) {
		t.Errorf("Expected InitiateWebsocketConnection to return a *websocket.Conn, but received %s", reflect.TypeOf(connection))
		t.FailNow()
	}
}

func TestGetSecureRtmConnectionInfoDeserializes(t *testing.T) {
	// Setup
	logger := testUtils.GetTestLogger()
	slackConnection := SlackConnection{logger: logger}

	const route = "/"
	const scheme = "http://"
	const listenAddress = "127.0.0.1"
	bindPort := strconv.Itoa(rand.Intn(9999))
	rtmConnectURL := scheme + listenAddress + ":" + bindPort + route

	// These are the objects the HTTP server will reply with
	expectedTeam := RtmConnectTeam{Domain: "fooTeam.com", ID: "teamID", Name: "teamName"}
	expectedSelf := RtmConnectSelf{ID: "selfID", Name: "selfName"}
	expectedConnectResponse := RtmConnectResponse{Ok: true, Self: expectedSelf, Team: expectedTeam, WsURL: "https://aSecureUrl.foo/websocket"}

	replyObject, error := json.Marshal(expectedConnectResponse)
	if error != nil {
		t.Errorf("Could not serialize expected response in test setup: %s", error.Error())
		t.FailNow()
	}
	server := testUtils.StartTestHTTPServer(route, bindPort, string(replyObject))

	rtmConnectionInfo, error := slackConnection.GetSecureRtmConnectionInfo("aPretentTokne", rtmConnectURL)

	if error != nil {
		t.Errorf("Unexpected error from GetSecureRtmConnectionInfo: %s", error.Error())
		t.FailNow()
	}

	// Confirm type & characteristics
	if rtmConnectionInfo.Ok != true {
		t.Errorf("Expected response property 'Ok' to be true, not %t", rtmConnectionInfo.Ok)
		t.FailNow()
	}

	if rtmConnectionInfo.Team.Name != "teamName" {
		t.Errorf("Expected response property Team to have a Name of teamName, not %s", rtmConnectionInfo.Team.Name)
		t.FailNow()
	}

	if rtmConnectionInfo.WsURL != "https://aSecureUrl.foo/websocket" {
		t.Errorf("Expected response property WsUrl to be https://aSecureUrl.foo/websocket, not %s", rtmConnectionInfo.WsURL)
		t.FailNow()
	}
	// teardown
	defer testUtils.StopTestWebsocketServer(server)
}

/*
	message := "Test Message"
	var reply string

	sendError := websocket.Message.Send(connection, message)
	if sendError != nil {
		fmt.Printf("Could not send test message over websocket in slack_test: %s", sendError.Error())
		t.FailNow()
	}

	receiveError := websocket.Message.Receive(connection, &reply)
	if receiveError != nil {
		fmt.Printf("Could not receive reply message over websocket in slack_test: %s", receiveError.Error())
		t.FailNow()
	}

	fmt.Printf("Server replied with %s", reply)
*/
