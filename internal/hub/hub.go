package hub

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"
	"github.com/timurkash/wsbus/internal/bus"
	"github.com/timurkash/wsbus/internal/conf"
)

type Hub struct {
	confBus *conf.Bus
	confWs  *conf.Ws

	conn    *websocket.Conn
	SetConn chan *websocket.Conn

	busClient bus.Bus

	ToBus   chan []byte
	FromBus chan []byte
}

const BufferSize = 256

func NewHub(confBus *conf.Bus, confWs *conf.Ws) (*Hub, error) {
	busClient, err := bus.NewBus(confBus)
	if err != nil {
		return nil, err
	}
	hub := &Hub{
		busClient: busClient,
		confBus:   confBus,
		confWs:    confWs,
		SetConn:   make(chan *websocket.Conn),
		ToBus:     make(chan []byte, BufferSize),
		FromBus:   make(chan []byte, BufferSize),
	}
	busClient.SetHub(hub)
	return hub, nil
}

func (h *Hub) CloseBus() {
	h.busClient.Close()
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.SetConn:
			h.conn = conn
		case msg := <-h.ToBus:
			if err := h.busClient.WriteTo(h.confBus.Subject, msg); err != nil {
				log.Error(err)
			}
		case msg := <-h.FromBus:
			if err := h.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Error(err)
			}
		}
	}
}
