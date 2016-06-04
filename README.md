# BOP Microservice framework

[![Build Status](https://travis-ci.org/kynrai/bop.svg?branch=master)](https://travis-ci.org/kynrai/bop)
[![godoc](https://godoc.org/github.com/kynrai/bop?status.svg)](http://godoc.org/github.com/kynrai/bop)

BOP is a light weight micro-service framework written in Go. It relies on RPC patterns using gnatsd as an rpc broker.

## Development

In order to play with the code you will need to have [golang](https://golang.org/doc/install) installed and [gnatsd](http://nats.io/download/)

## Architecture

The framework is built around the [gnatsd](http://nats.io/download/) message broker. This is a lightweight simple and blazing fast broker. We take advantage of its native queueing features as well as its request-reply model.

## Trying it out

Download the repo then run [gnatsd](http://nats.io/download/), now you can start the greeter server and try the simple client.

If you run gnatsd on any other address or port then set the env var NATS_HOST

Starting the server:
```
cd $GOPATH/src/github.com/kynrai/bop/examples/greeter_server
go run main.go
```

Run the client:
```
cd $GOPATH/src/github.com/kynrai/bop/examples/greeter_client
go run main.go
```

## Tracing
BOP comes with tracing out the box, if no app server is declared then no traces will be visible as it will default to the opentracing noop server. The examples show how traving might work with appdash. See the opentracing documentation for more information

## Acknowledgments

This project has been inspired by my own work and works of others across the open source community. Notable sources of inspiration cames from the following:
 
[go-mciro](http://micro.mu) A microservice ecosystem

[gnatsd](https://nats.io) The messaging system

[Micro on Nats](https://nats.io/blog/microonnats/) Talk from the creator of go-micro on using nats as a transport