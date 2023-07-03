package main

import (
	"fmt"
	"simple-chat-room/common"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Player struct {
	name       string
	connection *websocket.Conn

	msgCh chan<- string
	errCh chan<- error
}

func NewPlayer(
	name string,
	conn *websocket.Conn,
	ch chan<- string,
	eCh chan<- error,
) Player {
	p := Player{
		name:       name,
		connection: conn,

		msgCh: ch,
		errCh: eCh,
	}
	return p
}

func PlayerCast(i interface{}) (Player, error) {
	p, ok := i.(Player)
	if !ok {
		return Player{"", nil, nil, nil}, ErrorCastFailed
	}
	return p, nil
}

func (p Player) Read(l echo.Logger) {
	var err error
	for {
		var msgType int
		var msg []byte
		// closed client will return error
		// then this goroutine will be closed
		msgType, msg, err = p.connection.ReadMessage()
		if err != nil {
			break
		}

		if msgType == websocket.TextMessage {
			p.msgCh <- common.PlayerMsg(p.name, string(msg))
		}
	}
	p.errCh <- fmt.Errorf("player: %s, error %s", p.name, err)
}
