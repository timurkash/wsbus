package bus

import (
	"github.com/nats-io/nats.go"
	"github.com/timurkash/wsbus/internal/conf"
)

type Nats struct {
	nc *nats.Conn
}

func NewNats(confNats *conf.Nats) (Bus, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}
	return &Nats{
		nc: nc,
	}, nil
}

func (n *Nats) WriteToBus([]byte) error {

	return nil
}

func (n *Nats) ReadFromBus() ([]byte, error) {

	return nil, nil
}
