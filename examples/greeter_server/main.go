package main

import (
	"context"
	"fmt"

	"github.com/kynrai/bop"
)

func main() {
	tracer := bop.StartAppDash()
	s := bop.Service(
		bop.Name("greeter-service"),
		bop.Tracer(tracer),
		bop.Endpoint("hi", Hello),
		bop.Endpoint("bye", Bye),
	)
	s.Run()
}

func Hello(ctx context.Context, req, resp *bop.Message) error {
	var msg struct {
		Name string
		Age  int
	}
	bop.GobDecode(req.Payload, &msg)
	fmt.Printf("Received a greeting message from: Name %q, Age %d\n", msg.Name, msg.Age)
	byeResp := bop.NewRequest("greeter-service.bye", req)
	resp.Payload = []byte("hello " + string(byeResp.Payload))
	return nil
}

func Bye(ctx context.Context, req, resp *bop.Message) error {
	resp.Payload = []byte("bye response payload")
	return nil
}
