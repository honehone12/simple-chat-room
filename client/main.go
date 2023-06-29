package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	fmt.Println("[SETUP] Your name??")
	var playerName string
	_, err := fmt.Scan(&playerName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[SETUP] Room's name??")
	var roomName string
	_, err = fmt.Scan(&roomName)
	if err != nil {
		log.Fatal(err)
	}

	conn, res, err := websocket.DefaultDialer.Dial(
		fmt.Sprintf(
			"ws://localhost:1323/door/%s?name=%s",
			roomName,
			playerName,
		),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusSwitchingProtocols {
		log.Fatal("connection was not switched to websocket")
	}
	defer conn.Close()

	p := NewPrinter(conn)
	pe := p.ErrorChan()
	s := NewScanner(conn)
	se := s.ErrChan()
	go p.PrintMessages()
	go s.ScanInputs()

	select {
	case err = <-pe:
		break
	case err = <-se:
		break
	}
	log.Fatal(err)
}
