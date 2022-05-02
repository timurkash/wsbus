package bus

import (
	"errors"
	"github.com/timurkash/wsbus/internal/conf"
)

type Bus interface {
	WriteToBus(message []byte) error
	ReadFromBus() ([]byte, error)
}

func NewBus(confBus *conf.Bus) (Bus, error) {
	if confBus.Nats != nil {
		return NewNats(confBus.Nats)
	}
	return nil, errors.New("nats implemented only")
}
