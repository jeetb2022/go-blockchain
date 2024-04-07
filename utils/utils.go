package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Serialize serializes a data structure into bytes
func Serialize(data interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		log.Fatal("encode error:", err)
	}
	return buf.Bytes()
}
