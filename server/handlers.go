package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func Door(c echo.Context) error {
	roomName := c.Param("room")
	playerName := c.QueryParam("name")
	ctx := c.(*CRContext)

	room, err := RoomFromContext(ctx, roomName)
	if err != nil {
		return err
	}
	_, exists := room.players.Load(playerName)
	if exists {
		return ErrorPlayerExists
	}

	conn, err := ctx.protocolSwitcher.Upgrade(
		c.Response(),
		c.Request(),
		nil,
	)
	if err != nil {
		return err
	}
	player := NewPlayer(conn)
	room.players.Store(playerName, player)

	s := fmt.Sprintf("room:%s name:%s", roomName, playerName)
	room.BroadcastText(c.Logger(), s)

	return nil
}
