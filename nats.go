package bop

import (
	"log"
	"os"
	"time"

	"github.com/nats-io/nats"
)

// NewConnection creates a new nats encoded connection, If reconnect is true,
// the connection will attempt ot repair itself on disconnect
func NewConnection(reconnect bool) *nats.EncodedConn {
	var address string
	if host := os.Getenv("NATS_HOST"); host == "" {
		address = nats.DefaultURL
	} else {
		address = host
	}
	options := []nats.Option{
		nats.MaxReconnects(10),
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectHandler(func(_ *nats.Conn) {
			log.Print("Got disconnected!\n")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("Got reconnected to %v!\n", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Printf("Connection closed. Reason: %q\n", nc.LastError())
		}),
	}

	var c *nats.Conn
	var err error
	switch reconnect {
	case true:
		c, err = nats.Connect(address, options...)
	case false:
		c, err = nats.Connect(address)
	}
	if err != nil {
		log.Fatalf("Error trying to connect to nats server at %q, got error: %q", address, err)
	}

	nc, err := nats.NewEncodedConn(c, nats.GOB_ENCODER)
	if err != nil {
		log.Fatalf("Error trying to connect to encoded connection at %q, got error: %q", address, err)
	}

	return nc
}
