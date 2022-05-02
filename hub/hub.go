package hub

import "github.com/timurkash/wsbus/bus"

type Hub struct {
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

func (h *Hub) Run() {
	for {
		select {
		case <-h.toBus:

		case <-h.fromBus:

		}
	}
}
