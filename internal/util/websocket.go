package util

import (
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

var (
	ErrWSUpgrade      = errors.New("Websocket: ErrWSUpgrade")
	ErrWSParseMessage = errors.New("Websocket: ErrParseMessage")
	ErrWSWriteMessage = errors.New("Websocket: ErrWriteMessage")
	ErrWSReadMessage  = errors.New("Websocket: ErrReadMessage")
)

// Reads messages of a given websocket connection.
func ReadWS(conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println(ErrWSReadMessage)
			return
		}
	}
}

// Write messages to a given websocket connection.
// Second argument is a function, which acts as a
// getter for the data that should be sent via the
// websocket.
func WriteWS(
	conn *websocket.Conn,
	getMessage func() ([]byte, error),
) {
	for {
		message, err := getMessage()
		if err != nil {
			log.Println(ErrWSParseMessage)
			return
		}
		err = conn.WriteMessage(1, message)
		if err != nil {
			log.Println(ErrWSWriteMessage)
			return
		}
	}
}
