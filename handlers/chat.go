// Copyright 2015 - husobee associates, llc.  All rights reserved

// Package handlers - http.HandlerFuncs for gonoxious
package handlers

import (
	"crypto"
	"encoding/pem"
	"errors"
	"io"
	"log"
	"net/http"

	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"

	"github.com/husobee/gonoxious"
)

var (
	errDecodeMessage   = errors.New("Failed to Decode Payload")
	errValidateMessage = errors.New("Failed to Decode Payload")
)

func decodeAndValidate(m gonoxious.Message, reader io.Reader) error {
	if err := m.Decode(reader); err != nil {
		log.Println(err)
		return errDecodeMessage
	}
	if err := m.Validate(); err != nil {
		return errValidateMessage
	}
	return nil
}

// ChatHandler - primary handler for chat
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// each post to this handler is a new "message"
	m := gonoxious.NewMessage()

	// decode and validate the message for correctness
	defer r.Body.Close()
	if err := decodeAndValidate(m, r.Body); err != nil {
		// bad request, respond as such
		log.Println(err.Error())
		w.WriteHeader(400)
		w.Write([]byte("bad request"))
		return
	}

	t, err := m.GetType()
	if err != nil {
		// bad request, respond as such
		log.Println(err.Error())
		w.WriteHeader(400)
		w.Write([]byte("bad request"))
		return
	}
	switch t {
	case gonoxious.IntroductionContentType:
		// do the introduction message process
		// add this peer as a contact in the contact book

		// read the public key from the content
		key, err := m.GetPubPEM()
		if err != nil {
			log.Println(err.Error())
			// bad request, respond as such
			w.WriteHeader(400)
			w.Write([]byte("bad request"))
			return
		}
		log.Println(key)
		block, _ := pem.Decode([]byte(key))
		if block == nil {
			log.Println("failed to find pem block in key")
			w.WriteHeader(400)
			w.Write([]byte("bad request"))
			return
		}

		if publicKey, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
			if key, ok := publicKey.(*rsa.PublicKey); ok {
				log.Println("valid rsa public key")
				from, _ := m.GetFromAddr()

				// verify signature on content of intro message, using the key
				contentBytes, err := m.GetContentBytes()
				if err != nil {
					log.Println(err.Error())
					// bad request, respond as such
					w.WriteHeader(400)
					w.Write([]byte("bad request"))
					return
				}

				sig, err := m.GetSignature()
				if err != nil {
					log.Println(err.Error())
					// bad request, respond as such
					w.WriteHeader(400)
					w.Write([]byte("bad request"))
					return
				}

				// get the hashsum of the data in content
				hasher := sha256.New()
				hasher.Write(contentBytes)

				err = rsa.VerifyPKCS1v15(key, crypto.SHA256, hasher.Sum(nil), sig)
				if err != nil {
					log.Println("bad signature: ", err.Error())
					// bad request, respond as such
					w.WriteHeader(400)
					w.Write([]byte("bad request"))
					return
				}

				gonoxious.Contacts.Add(gonoxious.NewPeer(from, key))
				goto Success
			}
			log.Println("NOT valid rsa public key")
			w.WriteHeader(400)
			w.Write([]byte("bad request"))
			return
		}
	case gonoxious.EncryptedDataContentType:
		// do the encrypted message process
	}
Success:
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}
