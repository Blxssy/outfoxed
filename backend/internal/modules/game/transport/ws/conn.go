package ws

import (
	"context"
	"encoding/json"
	"sync"

	cws "github.com/coder/websocket"
)

type Conn struct {
	ws   *cws.Conn
	send chan []byte

	closeOnce sync.Once
}

func NewConn(ws *cws.Conn) *Conn {
	return &Conn{
		ws:   ws,
		send: make(chan []byte, 32),
	}
}

// WriteLoop — единственная горутина, которая пишет в сокет.
func (c *Conn) WriteLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			_ = c.ws.Write(ctx, cws.MessageText, msg)
		}
	}
}

func (c *Conn) SendJSON(v any) {
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	select {
	case c.send <- b:
	default:
	}
}

func (c *Conn) Close() {
	c.closeOnce.Do(func() {
		close(c.send)
		_ = c.ws.Close(cws.StatusNormalClosure, "bye")
	})
}
