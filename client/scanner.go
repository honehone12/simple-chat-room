package main

import (
	"bufio"
	"bytes"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

var (
	newLine = []byte{'\n'}
	space   = []byte{' '}
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
	stdin := bufio.NewScanner(os.Stdin)

	for {
		if stdin.Scan() {
			msg := bytes.TrimSpace(bytes.ReplaceAll(
				stdin.Bytes(),
				newLine,
				space,
			))
			err := s.connection.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				s.errCh <- err
			}
		} else {
			err := stdin.Err()
			if err != nil {
				s.errCh <- err
			} else {
				log.Println("reached EOF")
			}
		}
	}
}
