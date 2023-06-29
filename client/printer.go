package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Printer struct {
	connection *websocket.Conn
	errCh      chan error
}

func NewPrinter(conn *websocket.Conn) Printer {
	return Printer{
		connection: conn,
		errCh:      make(chan error),
	}
}

func (p Printer) ErrorChan() <-chan error {
	return p.errCh
}

func (p Printer) PrintMessages() {
	for {
		msgType, msg, err := p.connection.ReadMessage()
		if err != nil {
			p.errCh <- err
		}

		if msgType == websocket.TextMessage {
			fmt.Println(string(msg))
		}
	}
}
