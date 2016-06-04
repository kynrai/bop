package bop

import (
	"context"
	"log"
	"os"
	"os/signal"

	"fmt"

	"github.com/nats-io/nats"
	"github.com/opentracing/opentracing-go"
)

type (
	Handler interface {
		Handle(ctx context.Context, req, resp *Message) error
	}
	HandlerFunc func(ctx context.Context, req, resp *Message) error

	server struct {
		name        string
		conn        *nats.EncodedConn
		subscribers map[string]*nats.Subscription
		tracer      opentracing.Tracer
		handlers    map[string]HandlerFunc
	}
	option func(*server)
)

func (h HandlerFunc) Handle(ctx context.Context, req, resp *Message) error {
	return h(ctx, req, resp)
}

func Service(options ...option) *server {
	s := new(server)
	s.handlers = make(map[string]HandlerFunc)
	s.subscribers = make(map[string]*nats.Subscription)
	s.conn = NewConnection(true)
	for _, option := range options {
		option(s)
	}
	if s.name == "" {
		log.Fatal("Name is mandatory and must be unique per service")
	}

	log.Printf("Connected to the hive mind server %q at %q", s.conn.Conn.ConnectedServerId(), s.conn.Conn.ConnectedUrl())
	return s
}

func (s *server) Subscribe(service, hName string, handler HandlerFunc) {
	if service == "" {
		log.Fatal("Error subscribing to channel. subject cannot be empty")
	}
	var err error
	sub := fmt.Sprintf("%s.%s", service, hName)
	if s.subscribers[hName], err = s.conn.QueueSubscribe(sub, sub, func(subject, reply string, m *Message) {
		defer func() {
			if err := recover(); err != nil {
				s.respond(reply, Message{Errors: []string{err.(string)}})
				log.Print(err)
			}
		}()
		var resp Message
		if err := Trace(s.tracer, subject, handler).Handle(context.Background(), m, &resp); err != nil {
			resp.Errors = append(resp.Errors, err.Error())
		}
		s.respond(reply, resp)
	}); err != nil {
		log.Fatalf("Error trying to connect to channel %q, got error: %q", service, err)
	}
}

func (s *server) respond(reply string, resp Message) {
	if err := s.conn.Publish(reply, resp); err != nil {
		log.Print("failed to respond to request ", err)
	}
}

func (s *server) Run() {
	defer func() {
		log.Printf("ending service %q", s.name)
		s.Close()
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Kill, os.Interrupt)
	<-sig
}

func (s *server) Close() {
	for _, h := range s.subscribers {
		if err := h.Unsubscribe(); err != nil {
			log.Printf("error unsubscribing from %q, got error %q", h.Subject, err)
		}
	}
	if !s.conn.Conn.IsClosed() {
		s.conn.Close()
	}
}

func Name(n string) option {
	return func(s *server) {
		s.name = n
	}
}

func Tracer(t opentracing.Tracer) option {
	return func(s *server) {
		s.tracer = t
		opentracing.InitGlobalTracer(s.tracer)
	}
}

func Endpoint(name string, handler HandlerFunc) option {
	return func(s *server) {
		s.Subscribe(s.name, name, handler)
	}
}
