package bus

import (
	"errors"
	"github.com/timurkash/wsbus/internal/conf"
	"github.com/timurkash/wsbus/internal/hub"
)

type Bus interface {
	WriteTo(subject string, message []byte) error
	Subscribe(subject string) error
	Close()
	SetHub(hub *hub.Hub)
}

func NewBus(confBus *conf.Bus) (Bus, error) {
	if confBus.Nats != nil {
		return NewNats(confBus.Nats)
	}
	return nil, errors.New("nats implemented only")
}
