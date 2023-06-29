package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	fmt.Println("Your name??")
	var playerName string
	_, err := fmt.Scan(&playerName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Room's name??")
	var roomName string
	_, err = fmt.Scan(&roomName)
	if err != nil {
		log.Fatal(err)
	}

	reqURL := fmt.Sprintf(
		"ws://localhost:1323/door/%s?name=%s",
		roomName,
		playerName,
	)

	conn, res, err := websocket.DefaultDialer.Dial(reqURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusSwitchingProtocols {
		log.Fatal("connection was not switched to websocket")
	}

	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}

		log.Println(msgType)
		log.Println(string(msg))
	}
}
