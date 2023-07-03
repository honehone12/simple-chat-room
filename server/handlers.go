package main

import (
	"fmt"
	"simple-chat-room/common"

	"github.com/labstack/echo/v4"
)

func Door(c echo.Context) error {
	roomName := c.Param("room")
	playerName := c.QueryParam("name")

	room, err := RoomFromContext(c, roomName)
	if err != nil {
		return err
	}
	_, exists := room.players.Load(playerName)
	if exists {
		return ErrorPlayerExists
	}

	conn, err := c.(*CRContext).protocolSwitcher.Upgrade(
		c.Response(),
		c.Request(),
		nil,
	)
	if err != nil {
		return err
	}
	player := NewPlayer(
		playerName,
		conn,
		room.MsgChan(),
		room.ErrChan(),
	)
	// better key type ?? (can't use []byte here)
	room.players.Store(playerName, player)
	go player.Read(c.Logger())

	s := fmt.Sprintf("%s joined to %s!", playerName, roomName)
	room.BroadcastText(c.Logger(), common.INFOMsg(s))

	return nil
}
