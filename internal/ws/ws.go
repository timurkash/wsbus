package ws

import (
	"github.com/gorilla/websocket"
	"github.com/timurkash/wsbus/internal/hub"
	"github.com/timurkash/wsbus/internal/message"
	"log"
	"net/http"
	"time"
)

const BufferSize = 1024

type WsClient struct {
	hub  *hub.Hub
	conn *websocket.Conn
	Send chan []byte
}

func Serve(hub *hub.Hub, w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  BufferSize,
		WriteBufferSize: BufferSize,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	wsClient := &WsClient{
		hub:  hub,
		conn: conn,
		Send: make(chan []byte, 256),
	}
	hub.SetWsClient <- wsClient
	go wsClient.readPump()
	go wsClient.writePump()
}

func (c *WsClient) writePump() {
	ticker := time.NewTicker(c.hub.ConfWs.PingPeriod.AsDuration())
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(c.hub.ConfWs.WriteWait.AsDuration()))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(c.hub.ConfWs.WriteWait.AsDuration()))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}

}

func (c *WsClient) readPump() {
	c.conn.SetReadLimit(int64(c.hub.ConfWs.MaxMessageSize))
	pongWait := c.hub.ConfWs.PongWait.AsDuration()
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error: %v", err)
			}
			break
		}
		message := &message.Message{}
		if err := message.Get(msg); err != nil {
			log.Println(err)
			continue
		}
		if err := message.Check(); err != nil {
			log.Println(err)
			continue
		}
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.ToBus <- msg
	}

}
