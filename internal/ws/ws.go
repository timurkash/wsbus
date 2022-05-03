package ws

import (
	"github.com/gorilla/websocket"
	"github.com/timurkash/wsbus/internal/hub"
	"log"
	"net/http"
)

const BufferSize = 1024

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
	hub.SetConn <- conn
}
