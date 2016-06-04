package bop

import (
	"time"

	"github.com/nats-io/nats"
)

const DEFAULT_TIMEOUT = 100 * time.Millisecond

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
	clientOptions func(*client)
)

func (c *client) Execute() {
	if c.Timeout == 0 {
		c.Timeout = DEFAULT_TIMEOUT
	}
	c.PlatformError = conn.Request(c.Service, c.Req, c.Resp, c.Timeout)

}

func NewRequest(service string, request *Message, options ...clientOptions) *Message {
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

func Timeout(t time.Duration) clientOptions {
	return func(c *client) {
		c.Timeout = t
	}
}
