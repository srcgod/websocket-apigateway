package types

type ClientInterface interface {
	ID() string

	UserID() int64

	SendMessage(data []byte)

	Close()
}
