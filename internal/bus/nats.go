package bus

import (
	"errors"
	"github.com/nats-io/nats.go"
	"github.com/timurkash/wsbus/internal/conf"
	"github.com/timurkash/wsbus/internal/hub"
	"log"
)

type Nats struct {
	nc  *nats.Conn
	hub *hub.Hub
}

func NewNats(confNats *conf.Nats) (Bus, error) {
	if confNats.Name == "" {
		return nil, errors.New("nats.name required")
	}
	opts := []nats.Option{nats.Name(confNats.Name)}
	if confNats.Creds != "" {
		opts = append(opts, nats.UserCredentials(confNats.Creds))
	}
	if confNats.ReconnectDelay == nil {
		return nil, errors.New("nats.reconnect_delay required")
	}
	reconnectDelayDuration := confNats.ReconnectDelay.AsDuration()
	if confNats.TotalWait == nil {
		return nil, errors.New("nats.total_wait required")
	}
	totalWaitDuration := confNats.TotalWait.AsDuration()
	opts = append(opts, nats.ReconnectWait(reconnectDelayDuration))
	opts = append(opts, nats.MaxReconnects(int(totalWaitDuration/reconnectDelayDuration)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected: will attempt reconnects for %ds", confNats.TotalWait.Seconds)
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	nc, err := nats.Connect(confNats.Url, opts...)
	if err != nil {
		return nil, err
	}
	return &Nats{
		nc: nc,
	}, nil
}

func (n *Nats) SetHub(hub *hub.Hub) {
	n.hub = hub
}

func (n *Nats) WriteTo(subject string, message []byte) error {
	if err := n.nc.Publish(subject, message); err != nil {
		return err
	}
	if err := n.nc.Flush(); err != nil {
		return err
	}
	return nil
}

func (n *Nats) Subscribe(subject string) error {
	if _, err := n.nc.Subscribe(subject, func(msg *nats.Msg) {
		if n.hub == nil {
			panic("n.hub is nil")
		}
		n.hub.FromBus <- msg.Data
	}); err != nil {
		return err
	}
	return nil
}

func (n *Nats) Close() {
	n.nc.Close()
}
