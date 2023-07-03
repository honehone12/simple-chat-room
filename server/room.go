package main

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	pingBytes = []byte{0x70, 0x69, 0x6e, 0x67} // "ping"
)

type Room struct {
	players *sync.Map

	msgCh chan string
	errCh chan error
}

func NewRoom() Room {
	return Room{
		players: &sync.Map{},
		msgCh:   make(chan string),
		errCh:   make(chan error),
	}
}

func RoomCast(i interface{}) (Room, error) {
	r, ok := i.(Room)
	if !ok {
		return Room{nil, nil, nil}, ErrorCastFailed
	}
	return r, nil
}

func RoomFromContext(c echo.Context, roomName string) (Room, error) {
	ctx := c.(*CRContext)
	i, exists := ctx.chatRooms.Load(roomName)
	if exists {
		room, err := RoomCast(i)
		if err != nil {
			return Room{nil, nil, nil}, err
		}
		return room, nil
	} else {
		room := NewRoom()
		go room.Hub(c.Logger())
		ctx.chatRooms.Store(roomName, room)
		return room, nil
	}
}

func (r Room) MsgChan() chan<- string {
	return r.msgCh
}

func (r Room) ErrChan() chan<- error {
	return r.errCh
}

// returns count of active players
func (r Room) BroadcastText(l echo.Logger, text string) uint {
	return r.broadcastInternal(l, websocket.TextMessage, []byte(text))
}

// returns count of active players
func (r Room) BroadcastPing(l echo.Logger) uint {
	return r.broadcastInternal(l, websocket.PingMessage, pingBytes)
}

// returns count of active players
func (r Room) broadcastInternal(l echo.Logger, msgType int, msg []byte) uint {
	badConns := []interface{}{}
	numPlayers := uint(0)
	r.players.Range(func(k, v any) bool {
		p, err := PlayerCast(v)
		if err != nil {
			l.Error(
				err,
				"this cast error means it can not even close the connection",
			)
			badConns = append(badConns, k)
			// want iterate all anyway
			return true
		}
		d := time.Now().Add(time.Millisecond * 1000)
		p.connection.SetWriteDeadline(d)
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

func (r Room) Hub(l echo.Logger) {
	for {
		select {
		case msg := <-r.msgCh:
			r.BroadcastText(l, msg)
		case err := <-r.errCh:
			// catching this ch means the goroutine was closed
			// that's enough all here
			// actual errored connection handling is done in ping loop
			l.Warn(err)
		}
	}
}
