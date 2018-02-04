package testUtils

import (
	"fmt"
	"net/http"

	log "github.com/inconshreveable/log15"
	"golang.org/x/net/websocket"
)

// GetTestLogger returns a logger instance that tests can use. For now just
// writes to console at INFO level.
func GetTestLogger() log.Logger {
	logger := log.New("Alarmie", "Test Logger")
	logHandler := log.LvlFilterHandler(log.LvlDebug, log.StdoutHandler)
	logger.SetHandler(logHandler)
	return logger
}

// GetTestWebsocketServer returns an instance of an HTTP Server that can serve
// websockets. Callers use the instance to stop the server. Takes two arguments:
// route <string> a path/route to mount, which will return a websocket connection
// response <string> a JSON blob that the websocket will reply to any message with,
// 	(representing an object) to reply with given a GET to that route. This is to
// 	be able to test message deserialization/handling logic in a client.
// Callers must start the server
func GetTestWebsocketServer(route string, bindPort string, response string) *http.Server {
	server := &http.Server{Addr: ":" + bindPort}

	handler := CreateWebsocketHandlerWithResponse(response)

	http.Handle(route, websocket.Handler(handler))
	return server
}

// GetTestHTTPServer returns an instance of an HTTP Server that will reply
// with a given message. Callers use the instance to stop the server. Takes two
// arguments:
// route <string> a path/route to mount
// response <string> a JSON blob that the server will reply with to any requests
// 	to the supplied route
// This is to be able to test message deserialization/handling logic in client
// The server must be started by the caller
func GetTestHTTPServer(route string, bindPort string, response string) *http.Server {
	server := &http.Server{Addr: ":" + bindPort}
	handler := CreateHTTPHandlerWithResponse(response)
	http.HandleFunc(route, handler)
	return server
}

// RunHTTPServerListenLoop is a simple helper to start an http server instance running
func RunHTTPServerListenLoop(server *http.Server) {
	go func() {
		error := server.ListenAndServe()
		if error != nil {
			fmt.Printf("Error in http server in testUtils: %s", error.Error())
		}
	}()
}

// StopTestServer gracefully shuts down the http server
func StopTestServer(server *http.Server) error {
	error := server.Shutdown(nil)
	if error != nil {
		fmt.Printf("Could not stop test websocket server in test_utils: %s", error.Error())
		return error
	}
	return nil
}

// CreateWebsocketHandlerWithResponse returns an anonymous function that can
// stand in as a websocket handler for the test server, which will reply with
// the given string
func CreateWebsocketHandlerWithResponse(response string) func(*websocket.Conn) {
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

// CreateHTTPHandlerWithResponse returns an anonymous function that can stand
// in as an HTTP request/response handler for the test server
func CreateHTTPHandlerWithResponse(response string) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, response)
	}
}
