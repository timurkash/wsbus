package hub

import (
	"github.com/timurkash/wsbus/internal/bus"
	"golang.org/x/net/websocket"
)

type Hub struct {
	conn    *websocket.Conn
	bus     *bus.Bus
	toBus   chan []byte
	fromBus chan []byte
}

const BufferSize = 256

func NewHub(bus *bus.Bus) *Hub {
	return &Hub{
		bus:     bus,
		toBus:   make(chan []byte, BufferSize),
		fromBus: make(chan []byte, BufferSize),
	}
}

func (h *Hub) SetConn(conn *websocket.Conn) {
	h.conn = conn
}

func (h *Hub) Run() {
	for {
		select {
		case <-h.toBus:

		case <-h.fromBus:

		}
	}
}
