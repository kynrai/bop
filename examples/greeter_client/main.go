package main

import (
	"log"
	"time"

	"github.com/kynrai/bop"
)

func main() {
	var msg struct {
		Name string
		Age  int
	}
	msg.Name, msg.Age = "John", 25
	begin := time.Now()
	resp := bop.NewRequest("greeter-service.hi", &bop.Message{Payload: bop.GobEncode(msg)})
	log.Printf("Received: '%s'\n", string(resp.Payload))

	log.Println(time.Now().Sub(begin))
}
