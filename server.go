package bop

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/nats-io/nats"
	"github.com/opentracing/opentracing-go"
)

type (
	// Handler defines a handler for an endpoint in a service
	Handler interface {
		Handle(ctx context.Context, req, resp *Message) error
	}
	// HandlerFunc is a function adaptor for Handler
	HandlerFunc func(ctx context.Context, req, resp *Message) error

	// Server represents the server which will run the microservice
	Server struct {
		name        string
		conn        *nats.EncodedConn
		subscribers map[string]*nats.Subscription
		tracer      opentracing.Tracer
		handlers    map[string]HandlerFunc
	}
	// Option is a function which sets variables in a server
	Option func(*Server)
)

// Handle executes a handler with the given context and request
func (h HandlerFunc) Handle(ctx context.Context, req, resp *Message) error {
	return h(ctx, req, resp)
}

// Service starts a new server which runs the microservice
func Service(options ...Option) *Server {
	s := new(Server)
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

// Subscribe adds a subscriber to a topic and registers the asigned handler
func (s *Server) Subscribe(service, handlerName string, handler HandlerFunc) {
	if service == "" {
		log.Fatal("Error subscribing to channel. subject cannot be empty")
	}
	var err error
	sub := fmt.Sprintf("%s.%s", service, handlerName)
	if s.subscribers[handlerName], err = s.conn.QueueSubscribe(sub, sub, func(subject, reply string, m *Message) {
		defer func() {
			if err := recover(); err != nil {
				s.respond(reply, Message{Errors: []string{err.(string)}})
				log.Print(err)
			}
		}()
		var resp Message
		if err = Trace(s.tracer, subject, handler).Handle(context.Background(), m, &resp); err != nil {
			resp.Errors = append(resp.Errors, err.Error())
		}
		s.respond(reply, resp)
	}); err != nil {
		log.Fatalf("Error trying to connect to channel %q, got error: %q", service, err)
	}
}

func (s *Server) respond(reply string, resp Message) {
	if err := s.conn.Publish(reply, resp); err != nil {
		log.Print("failed to respond to request ", err)
	}
}

// Run starts the server which should have been setup at this point. Run will
// clean up the server on exit. This call blocks until cancelled or terminated.
func (s *Server) Run() {
	defer func() {
		log.Printf("ending service %q", s.name)
		s.Close()
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Kill, os.Interrupt)
	<-sig
}

// Close will attempt to reclaim resources on server shutdown
func (s *Server) Close() {
	for _, h := range s.subscribers {
		if err := h.Unsubscribe(); err != nil {
			log.Printf("error unsubscribing from %q, got error %q", h.Subject, err)
		}
	}
	if !s.conn.Conn.IsClosed() {
		s.conn.Close()
	}
}

// Name sets the name of the service
func Name(n string) Option {
	return func(s *Server) {
		s.name = n
	}
}

// Tracer sets the tracer engine to use with the service
func Tracer(t opentracing.Tracer) Option {
	return func(s *Server) {
		s.tracer = t
		opentracing.InitGlobalTracer(s.tracer)
	}
}

// Endpoint sets a new endpoint to register with the service
func Endpoint(name string, handler HandlerFunc) Option {
	return func(s *Server) {
		s.Subscribe(s.name, name, handler)
	}
}
