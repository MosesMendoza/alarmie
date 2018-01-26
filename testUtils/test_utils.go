package testUtils

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

// StartTestWebsocketServer returns an instance of an HTTP Server that can serve
// websockets. Callers use the instance to stop the server. Takes two arguments:
// route <string> a path/route to mount, which will return a websocket connection
// response <string> a JSON blob that the websocket will reply to any message with,
// 	(representing an object) to reply with given a GET to that route. This is to
// 	be able to test message deserialization/handling logic in a client.
func StartTestWebsocketServer(route string, response string) *http.Server {
	const bindPort = ":9999"

	server := &http.Server{Addr: bindPort}

	handler := createHandlerWithResponse(response)

	http.Handle("/", websocket.Handler(handler))

	go func() {
		error := server.ListenAndServe()
		if error != nil {
			fmt.Printf("Error in http server in testUtils: %s", error.Error())
		}
	}()

	return server
}

// StopTestWebsocketServer gracefully shuts down the http server
func StopTestWebsocketServer(server *http.Server) error {
	error := server.Shutdown(nil)
	if error != nil {
		fmt.Printf("Could not stop test websocket server in test_utils: %s", error.Error())
		return error
	}
	return nil
}

// createHandlerWithResponse returns an anonymous function that can stand in as
// a handler for the test server, which will reply with the given string
func createHandlerWithResponse(response string) func(*websocket.Conn) {
	return func(socket *websocket.Conn) {
		for {
			var receivedMessage string
			receiveError := websocket.Message.Receive(socket, &receivedMessage)
			if receiveError != nil {
				fmt.Printf("Could not receive message from websocket in testUtils: %s", receiveError.Error())
			}
			sendError := websocket.Message.Send(socket, response)
			if sendError != nil {
				fmt.Printf("Could not send reply message over websocket in testUtils: %s", sendError.Error())
			}
		}
	}
}
