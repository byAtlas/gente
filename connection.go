package gente

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Connection struct {
	id uuid.UUID

	msgPipe MessagePipeline
	log     logrus.Logger

	ws       *websocket.Conn
	inbound  chan []byte
	outbound chan []byte
}

//Connection won't actually start listening until ServeHTTP is called
func NewConnection(p MessagePipeline, log logrus.Logger) *Connection {
	return &Connection{id: uuid.NewUUID(), msgPipe: p, log: log}
}

func (c *Connection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.log.Error(err)
		return
	}

	c.outbound = make(chan []byte, 256)
	c.inbound = make(chan []byte, 256)
	c.ws = ws

	c.msgPipe.Register(c.inbound, c.outbound)

	go c.writePump()
	c.readPump()
}

// readPump pumps messages from the websocket connection to the messagePipline.
func (c *Connection) readPump() {
	defer func() {
		c.ws.Close()
	}()
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			c.log.Error(err)
			break
		}
		c.inbound <- message
	}
}

//actually send a message over the socket.
func (c *Connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

//send messages from the outbound channel
func (c *Connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.outbound:
			if !ok {
				c.log.Info("Outbound channel closed, sending close msg.")
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
