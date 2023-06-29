package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Scanner struct {
	connection *websocket.Conn
	errCh      chan error
}

func NewScanner(conn *websocket.Conn) Scanner {
	return Scanner{
		connection: conn,
		errCh:      make(chan error),
	}
}

func (s Scanner) ErrChan() <-chan error {
	return s.errCh
}

func (s Scanner) ScanInputs() {
	for {
		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			s.errCh <- err
		}

		err = s.connection.WriteMessage(websocket.TextMessage, []byte(input))
		if err != nil {
			s.errCh <- err
		}
	}
}
