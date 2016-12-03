package bop

import (
	"time"

	"github.com/nats-io/nats"
)

const defaultTimeout = 100 * time.Millisecond

func init() {
	conn = NewConnection(true)
}

var conn *nats.EncodedConn

type (
	client struct {
		Service       string
		Timeout       time.Duration
		Req, Resp     *Message
		PlatformError error
	}
	// ClientOptions sets options on a given client
	ClientOptions func(*client)
)

func (c *client) Execute() {
	if c.Timeout == 0 {
		c.Timeout = defaultTimeout
	}
	c.PlatformError = conn.Request(c.Service, c.Req, c.Resp, c.Timeout)

}

// NewRequest creates a new client request to be sent to the given service via gnatsd
func NewRequest(service string, request *Message, options ...ClientOptions) *Message {
	c := new(client)
	c.Service = service
	c.Req = request
	if c.Req.Values == nil {
		c.Req.Values = make(map[string]string)
	}
	for _, option := range options {
		option(c)
	}
	c.Resp = &Message{}
	c.Execute()
	return c.Resp
}

// Timeout sets the timeout of a client
func Timeout(t time.Duration) ClientOptions {
	return func(c *client) {
		c.Timeout = t
	}
}
