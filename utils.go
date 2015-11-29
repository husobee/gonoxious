// Copyright 2015 - husobee associates, llc.  All rights reserved

// Package handlers - http.HandlerFuncs for gonoxious
package gonoxious

import (
	"encoding/json"
	"io"
)

// Message - interface to define what a message is capable of doing
type Message interface {
	// decode request bodies
	Decode(io.Reader) error
	// encode request bodies
	Encode(io.Writer) error
	// validate message
	Validate() error
	// get the type of message
	GetType() (ContentType, error)
	// get the pubpem of message
	GetPubPEM() (string, error)
	// get the from address of message
	GetFromAddr() (string, error)
	// get the protocol of the message
	GetProtocol() (Protocol, error)
	// get signature
	GetSignature() ([]byte, error)
	// get content bytes
	GetContentBytes() ([]byte, error)
}

// validateProtocol - make sure this is a supported protocol
func validateProtocol(p Protocol) error {
	for _, v := range supportedProtocols {
		if p == v {
			return nil
		}
	}
	return ErrInvalidProtocol
}

// decodeMessage - simple generic decode message
func decodeMessage(v interface{}, reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(v)
}

// encodeMessage - simple generic encoder message
func encodeMessage(v interface{}, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(v)
}
