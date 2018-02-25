package main

import (
	"encoding/json"
	"net/http"
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
	server := testUtils.GetTestWebsocketServer(route, bindPort, reply)
	testUtils.RunHTTPServerListenLoop(server)

	logger := testUtils.GetTestLogger()
	slackConnection := SlackConnection{logger: logger}
	// teardown
	defer testUtils.StopTestServer(server)

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
	server := testUtils.GetTestHTTPServer(route, bindPort, string(replyObject))
	testUtils.RunHTTPServerListenLoop(server)
	// teardown
	defer testUtils.StopTestServer(server)

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
}

func TestConnectCreatesConnectionContext(t *testing.T) {
	logger := testUtils.GetTestLogger()
	slackConnection := SlackConnection{logger: logger}

	const listenAddress = "127.0.0.1"
	bindPort := strconv.Itoa(rand.Intn(9999))

	// ============ HTTP Server Setup ===================================
	const route = "/"
	const scheme = "http://"
	rtmConnectURL := scheme + listenAddress + ":" + bindPort + route

	// =========== Websocket Server Setup ==================================
	const websocketReply = "foo"
	const websocketRoute = "/getwebsocket"
	const websocketScheme = "ws://"
	websocketURL := websocketScheme + listenAddress + ":" + bindPort + websocketRoute

	// These are the objects the HTTP server will reply with to enable us to set
	// up the websocket connection
	expectedTeam := RtmConnectTeam{Domain: "fooTeam.com", ID: "teamID", Name: "teamName"}
	expectedSelf := RtmConnectSelf{ID: "selfID", Name: "selfName"}
	expectedConnectResponse := RtmConnectResponse{Ok: true, Self: expectedSelf, Team: expectedTeam, WsURL: websocketURL}

	replyObject, error := json.Marshal(expectedConnectResponse)
	if error != nil {
		t.Errorf("Could not serialize expected response in test setup: %s", error.Error())
		t.FailNow()
	}

	server := testUtils.GetTestHTTPServer(route, bindPort, string(replyObject))
	websocketHandler := testUtils.CreateWebsocketHandlerWithResponse(websocketReply)
	http.Handle(websocketRoute, websocket.Handler(websocketHandler))

	defer testUtils.StopTestServer(server)
	testUtils.RunHTTPServerListenLoop(server)

	//  now the actual testing begins
	rtmConnectionContext, error := slackConnection.Connect("aPretendToken", rtmConnectURL)

	if rtmConnectionContext.TeamName != "teamName" {
		t.Errorf("Expected response property Team to have a Name of teamName, not %s", rtmConnectionContext.TeamName)
		t.FailNow()
	}

	// validate that the socket connection is active
	socket := rtmConnectionContext.SocketConnection
	message := "Test Message"
	var reply string

	sendError := websocket.Message.Send(socket, message)
	if sendError != nil {
		t.Errorf("Could not send test message over websocket in slack_test: %s", sendError.Error())
		t.FailNow()
	}

	receiveError := websocket.Message.Receive(socket, &reply)
	if receiveError != nil {
		t.Errorf("Could not receive reply message over websocket in slack_test: %s", receiveError.Error())
		t.FailNow()
	}

	if reply != websocketReply {
		t.Errorf("Expected websocket to reply with %s, not %s", websocketReply, reply)
		t.FailNow()
	}
}
