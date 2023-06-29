package main

import "github.com/gorilla/websocket"

type Player struct {
	connection *websocket.Conn
}

func NewPlayer(conn *websocket.Conn) *Player {
	return &Player{
		connection: conn,
	}
}

func PlayerCast(i interface{}) (*Player, error) {
	p, ok := i.(*Player)
	if !ok {
		return nil, ErrorCastFailed
	}
	return p, nil
}
