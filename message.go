// Copyright 2015 - husobee associates, llc.  All rights reserved

// Package gonoxious - http.HandlerFuncs for gonoxious
package gonoxious

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
)

// ContentType - the type of the content
type ContentType string

// Protocol - the type of the content
type Protocol string

const (
	// IntroductionContentType - This is the introduction content type
	IntroductionContentType = ContentType("introduction")
	// EncryptedDataContentType - This is the introduction content type
	EncryptedDataContentType = ContentType("encryptedData")
	// Protocolv1 - version 1 of the protocol
	ProtocolV1 = Protocol("1.0")
)

var (
	// supportedProtocols - private variable to hold allowable protocols
	supportedProtocols = []Protocol{
		ProtocolV1,
	}
	ErrInvalidProtocol    = errors.New("Invalid Message Protocol")
	ErrInvalidContentType = errors.New("Invalid Message ContentType")
	ErrInvalidFromAddr    = errors.New("Invalid From Addr")
	ErrInvalidPublicKey   = errors.New("Invalid PublicKey")
	ErrInvalidSignature   = errors.New("Invalid Signature")
)

func NewMessage() Message {
	return &message{}
}

type messageContent struct {
	From      string      `json:"from,omitempty"`
	PubPEM    string      `json:"pubPem,omitempty"`
	To        string      `json:"to,omitempty"`
	Type      ContentType `json:"type,omitempty"`
	ClearFrom string      `json:"clearFrom,omitempty"`
	Data      []byte      `json:"data,omitempty"`
}
type message struct {
	Content   messageContent `json:"content,omitempty"`
	Protocol  Protocol       `json:"protocol,omitempty"`
	Signature string         `json:"signature,omitempty"`
}

// GetFromAddr - get the from address from the message
func (im *message) GetFromAddr() (string, error) {
	var err error
	if im.Content.From == "" {
		err = ErrInvalidFromAddr
	}
	return im.Content.From, err
}

// Decode - introduction message implementation
func (im *message) Decode(reader io.Reader) error {
	return decodeMessage(im, reader)
}

// Encode - introduction message implementation
func (im *message) Encode(writer io.Writer) error {
	return encodeMessage(im, writer)
}

// Validate - introduction message implementation
func (im *message) Validate() error {
	// validate protocol is supported
	if err := validateProtocol(im.Protocol); err != nil {
		return err
	}
	return nil
}

// GetContentBytes - get the content's data
func (im *message) GetContentBytes() ([]byte, error) {
	var err error
	b, err := json.Marshal(im.Content)
	return b, err
}

// GetSignature - get the content's signature
func (im *message) GetSignature() ([]byte, error) {
	var err error
	if im.Signature == "" && im.Content.Type == IntroductionContentType {
		err = ErrInvalidSignature
	}
	data, err := base64.StdEncoding.DecodeString(im.Signature)
	if err != nil {
		err = ErrInvalidSignature
	}
	return data, err
}

// GetPubPem - get the content's public pem
func (im *message) GetPubPEM() (string, error) {
	var err error
	if im.Content.PubPEM == "" && im.Content.Type == IntroductionContentType {
		err = ErrInvalidPublicKey
	}
	return im.Content.PubPEM, err
}

// GetType - get the content's type
func (im *message) GetType() (ContentType, error) {
	var err error
	if im.Content.Type != EncryptedDataContentType && im.Content.Type != IntroductionContentType {
		err = ErrInvalidContentType
	}
	return im.Content.Type, err
}

// GetProtocol - get the protocol version of the intro message
func (im *message) GetProtocol() (Protocol, error) {
	return im.Protocol, validateProtocol(im.Protocol)
}
