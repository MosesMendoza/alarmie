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
	logHandler := log.LvlFilterHandler(log.LvlInfo, log.StdoutHandler)
	logger.SetHandler(logHandler)
	return logger
}

// StartTestWebsocketServer returns an instance of an HTTP Server that can serve
// websockets. Callers use the instance to stop the server. Takes two arguments:
// route <string> a path/route to mount, which will return a websocket connection
// response <string> a JSON blob that the websocket will reply to any message with,
// 	(representing an object) to reply with given a GET to that route. This is to
// 	be able to test message deserialization/handling logic in a client.
func StartTestWebsocketServer(route string, bindPort string, response string) *http.Server {
	server := &http.Server{Addr: ":" + bindPort}

	handler := createWebsocketHandlerWithResponse(response)

	http.Handle(route, websocket.Handler(handler))

	runHTTPServerListenLoop(server)
	return server
}

// StartTestHTTPServer returns an instance of an HTTP Server that will reply
// with a given message. Callers use the instance to stop the server. Takes two
// arguments:
// route <string> a path/route to mount
// response <string> a JSON blob that the server will reply with to any requests
// 	to the supplied route
// This is to be able to test message deserialization/handling logic in client
func StartTestHTTPServer(route string, bindPort string, response string) *http.Server {
	server := &http.Server{Addr: ":" + bindPort}
	handler := createHTTPHandlerWithResponse(response)
	http.HandleFunc(route, handler)
	runHTTPServerListenLoop(server)
	return server
}

func runHTTPServerListenLoop(server *http.Server) {
	go func() {
		error := server.ListenAndServe()
		if error != nil {
			fmt.Printf("Error in http server in testUtils: %s", error.Error())
		}
	}()
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

// createWebsocketHandlerWithResponse returns an anonymous function that can
// stand in as a websocket handler for the test server, which will reply with
// the given string
func createWebsocketHandlerWithResponse(response string) func(*websocket.Conn) {
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

// createHTTPHandlerWithResponse returns an anonymous function that can stand
// in as an HTTP request/response handler for the test server
func createHTTPHandlerWithResponse(response string) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, response)
	}
}
