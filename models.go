package bop

import (
	"bytes"
	"encoding/gob"
	"log"
)

func init() {
	gob.RegisterName("message", Message{})
}

type Message struct {
	RequestID string
	Payload   []byte
	Errors    []string // TODO(kynrai) Unify errors
	Reply     string
	Values    map[string]string
}

func GobEncode(v interface{}) []byte {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		log.Printf("failed to encode %q", err)
		return nil
	}
	return b.Bytes()
}

func GobDecode(data []byte, vPtr interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(vPtr)
}
