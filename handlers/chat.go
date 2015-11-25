// Copyright 2015 - husobee associates, llc.  All rights reserved

// Package handlers - http.HandlerFuncs for gonoxious
package handlers

import "net/http"

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
	supportedProtocols = []Protocol{
		ProtocolV1,
	}
)

// validateProtocol - make sure this is a supported protocol
func validateProtocol(p Protocol) bool {
	for _, v := range supportedProtocols {
		if p == v {
			return true
		}
	}
	return false
}

type Message interface {
	GetType() ContentType
	GetProtocol() Protocol
}

// IntroductionContent - This is the introduction message structure
type IntroductionContent struct {
	From   string      `json:"from"`
	PubPEM string      `json:"pubPem"`
	To     string      `json:"to"`
	Type   ContentType `json:"type"`
}

// IntroductionMessage - this is the structure of the Introduction Message
type IntroductionMessage struct {
	Content   IntroductionContent `json:"content"`
	Protocol  Protocol            `json:"protocol"`
	Signature string              `json:"signature"`
}

// GetType - get the content's type
func (im *IntroductionMessage) GetType() ContentType {
	return im.Content.Type
}

// GetProtocol - get the protocol version of the intro message
func (im *IntroductionMessage) GetProtocol() Protocol {
	return im.Protocol
}

// EncryptedDataContent - this is the content of an encrypted data message
type EncryptedDataContent struct {
	Type      ContentType `json:"type"`
	ClearFrom string      `json:"clearFrom"`
	Data      []byte      `json:"data"`
}

// EncryptedDataMessage - this is the structure of the encrypted data message
type EncryptedDataMessage struct {
	Content  EncryptedDataContent `json:"content"`
	Protocol Protocol             `json:"protocol"`
}

// GetType - get the content's type
func (edm *EncryptedDataMessage) GetType() ContentType {
	return edm.Content.Type
}

// GetProtocol - get the protocol version of the encrypted data  message
func (edm *EncryptedDataMessage) GetProtocol() Protocol {
	return edm.Protocol
}

// ChatHandler - primary handler for chat
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// two types of messages possible, Introduction messages,
	// and encrypted messages
	w.WriteHeader(200)
	w.Write([]byte(""))
}
