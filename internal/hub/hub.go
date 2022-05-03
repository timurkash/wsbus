package hub

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/timurkash/wsbus/internal/bus"
	"github.com/timurkash/wsbus/internal/conf"
	"github.com/timurkash/wsbus/internal/ws"
)

type Hub struct {
	ConfBus   *conf.Bus
	busClient bus.Bus

	ConfWs      *conf.Ws
	wsClient    *ws.WsClient
	SetWsClient chan *ws.WsClient

	ToBus chan []byte
	ToWs  chan []byte
}

const BufferSize = 256

func NewHub(confBus *conf.Bus, confWs *conf.Ws) (*Hub, error) {
	busClient, err := bus.NewBus(confBus)
	if err != nil {
		return nil, err
	}
	hub := &Hub{
		busClient:   busClient,
		ConfBus:     confBus,
		ConfWs:      confWs,
		SetWsClient: make(chan *ws.WsClient),
		ToBus:       make(chan []byte, BufferSize),
		ToWs:        make(chan []byte, BufferSize),
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
		case wsClient := <-h.SetWsClient:
			h.wsClient = wsClient
		case msg := <-h.ToBus:
			if err := h.busClient.WriteTo(h.ConfBus.Subject, msg); err != nil {
				log.Error(err)
			}
		case msg := <-h.ToWs:
			h.wsClient.Send <- msg
		}
	}
}
