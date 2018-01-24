package testUtils

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

// StartTestWebsocketServer returns an instance of an HTTP Server that can serve
// websockets. Callers use the instance to stop the server.
func StartTestWebsocketServer() *http.Server {
	const bindPort = ":9999"

	server := &http.Server{Addr: bindPort}

	http.Handle("/", websocket.Handler(replyForverHandler))

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

// replyForeverHandler just replies back with the message that was received.
// Used to validate connection only.
func replyForverHandler(socket *websocket.Conn) {
	// FYI the websocket package has a number of "static" methods on the package that
	// take a websocket connection to act on.
	fmt.Println("In websocket handler")
	for {
		// drop the message into a string
		var receivedMessage string
		receiveError := websocket.Message.Receive(socket, &receivedMessage)
		if receiveError != nil {
			fmt.Printf("Could not receive message from websocket in testUtils: %s", receiveError.Error())
		}
		sendError := websocket.Message.Send(socket, receivedMessage)
		if sendError != nil {
			fmt.Printf("Could not send reply message over websocket in testUtils: %s", sendError.Error())
		}
	}
}
