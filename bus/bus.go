package bus

type Bus interface {
	WriteToBus(message []byte) error
	ReadFromBus() ([]byte, error)
}
