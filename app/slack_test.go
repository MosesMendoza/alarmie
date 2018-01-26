package main

import "testing"
import "github.com/MosesMendoza/alarmie/testUtils"
import "golang.org/x/net/websocket"
import "fmt"
import "encoding/json"

func TestInitiatWebSocketConnectionCreatesWebsocket(t *testing.T) {
	/*
		Test that, given a URL, uses websocket to "dial" that URL and initiate
		a socket connection. then returns the socket connection

		Test by creating an http server that serves websocket connections, and testing
		that the connection is returned
	*/

	replyObjectAsString, error := json.Marshal(Message{Type: "test", Channel: "channel-foo", User: "coder", Text: "test text", Timestamp: "<timestamp>"})

	server := testUtils.StartTestWebsocketServer("/", string(replyObjectAsString))

	defer testUtils.StopTestWebsocketServer(server)

	fmt.Printf("Dialing " + "ws://127.0.0.1" + server.Addr + "\n")
	connection, error := websocket.Dial("ws://127.0.0.1"+server.Addr, "", "http://localhost")

	if error != nil {
		fmt.Printf("Could not dial websocket connection in slack_test: %s", error.Error())
		t.FailNow()
	}
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
}
