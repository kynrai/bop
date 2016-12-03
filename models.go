package bop

import (
	"bytes"
	"encoding/gob"
	"log"
)

func init() {
	gob.RegisterName("message", Message{})
}

// Message represents the basic container for inter-service communications
type Message struct {
	RequestID string
	Payload   []byte
	Errors    []string // TODO(kynrai) Unify errors
	Reply     string
	Values    map[string]string
}

// GobEncode will return the given interface as a gob encoded byte array
func GobEncode(v interface{}) []byte {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		log.Printf("failed to encode %q", err)
		return nil
	}
	return b.Bytes()
}

// GobDecode will return the given data in given pointer
func GobDecode(data []byte, vPtr interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(vPtr)
}
