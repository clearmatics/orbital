package main

import (
	"bytes"
	"fmt"
)

type Ring struct {
	PubKeys []PubKey `json:"pubkeys"`
}

func (r Ring) Bytes() []byte {
	var xbytes, ybytes []byte
	for _, key := range r.PubKeys {
		xbytes = append(xbytes, key.X.Bytes()...)
		ybytes = append(ybytes, key.Y.Bytes()...)
	}

	return append(xbytes, ybytes...)
}

func (r Ring) String() string {
	var buffer bytes.Buffer

	for i := 0; i < len(r.PubKeys); i++ {
		buffer.WriteString(fmt.Sprintf("%s\n", r.PubKeys[i]))
	}

	return buffer.String()
}
