package main

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	pingBytes = []byte{0x70, 0x69, 0x6e, 0x67}
)

type Room struct {
	players *sync.Map
}

func NewRoom() *Room {
	return &Room{
		players: &sync.Map{},
	}
}

func RoomCast(i interface{}) (*Room, error) {
	r, ok := i.(*Room)
	if !ok {
		return nil, ErrorCastFailed
	}
	return r, nil
}

func RoomFromContext(ctx *CRContext, roomName string) (*Room, error) {
	i, exists := ctx.chatRooms.Load(roomName)
	if exists {
		room, err := RoomCast(i)
		if err != nil {
			return nil, err
		}
		return room, nil
	} else {
		room := NewRoom()
		ctx.chatRooms.Store(roomName, room)
		return room, nil
	}
}

func (r *Room) BroadcastText(l echo.Logger, text string) uint {
	return r.broadcastInternal(l, websocket.TextMessage, []byte(text))
}

func (r *Room) BroadcastPing(l echo.Logger) uint {
	return r.broadcastInternal(l, websocket.PingMessage, pingBytes)
}

func (r *Room) broadcastInternal(l echo.Logger, msgType int, msg []byte) uint {
	badConns := []interface{}{}
	numPlayers := uint(0)
	r.players.Range(func(k, v any) bool {
		p, err := PlayerCast(v)
		if err != nil {
			l.Error(
				err,
				"[Leak] this cast error means it can not even close the connection",
			)
			badConns = append(badConns, k)
			// want iterate all anyway
			return true
		}
		err = p.connection.WriteMessage(msgType, msg)
		if err != nil {
			defer p.connection.Close()
			l.Warn(err)
			badConns = append(badConns, k)
		}
		numPlayers++
		return true
	})

	if numClose := len(badConns); numClose > 0 {
		for i := 0; i < numClose; i++ {
			r.players.Delete(badConns[i])
			numPlayers--
		}
	}

	return numPlayers
}
