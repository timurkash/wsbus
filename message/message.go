package message

import (
	"encoding/json"
	"errors"
	"time"
)

type Message struct {
	MessageType string     `json:"messageType"`
	Time        *time.Time `json:"time"`
	Message     string     `json:"message"`
}

func (m *Message) Get(bytes []byte) error {
	if err := json.Unmarshal(bytes, m); err != nil {
		return err
	}
	if m.MessageType == "" {
		return errors.New("messageType required")
	}
	if m.Message == "" {
		return errors.New("message required")
	}
	return nil
}

func (m *Message) Bytes() ([]byte, error) {
	return json.Marshal(m)
}
