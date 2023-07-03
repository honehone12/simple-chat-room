package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const (
	pingIntervalMil = 1000
)

type CRContext struct {
	echo.Context

	protocolSwitcher *websocket.Upgrader
	chatRooms        *sync.Map
}

func ping(l echo.Logger, c *sync.Map) {
	for {
		totalRooms := uint(0)
		totalPlayers := uint(0)
		emptyRoom := []interface{}{}
		// Rnage() blocks while entire loop ??
		c.Range(func(k, v any) bool {
			r, err := RoomCast(v)
			if err != nil {
				l.Error(
					err,
					"this cast error means it can not even close the room",
				)
				return true
			}
			n := r.BroadcastPing(l)
			if n == 0 {
				emptyRoom = append(emptyRoom, k)
			} else {
				totalRooms++
				totalPlayers += n
			}
			return true
		})

		if numEmpty := len(emptyRoom); numEmpty > 0 {
			for i := 0; i < numEmpty; i++ {
				c.Delete(emptyRoom[i])
			}
		}

		l.Info(fmt.Sprintf("summary: rooms[%d] players[%d]", totalRooms, totalPlayers))
		time.Sleep(time.Millisecond * pingIntervalMil)
	}
}

func main() {
	chatRooms := &sync.Map{}
	protocolSwitcher := &websocket.Upgrader{}

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := &CRContext{
				Context:          c,
				protocolSwitcher: protocolSwitcher,
				chatRooms:        chatRooms,
			}
			return next(ctx)
		}
	})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.INFO)

	e.GET("/door/:room", Door)

	go ping(e.Logger, chatRooms)
	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}
