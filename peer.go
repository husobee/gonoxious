// Copyright 2015 - husobee associates, llc.  All rights reserved

// Package gonoxious - http.HandlerFuncs for gonoxious
package gonoxious

import (
	"crypto"
	"errors"
	"sync"
)

var (
	// ErrNoSuchContact - there was no such contact
	ErrNoSuchContact = errors.New("no such contact")
	// ErrContactAlreadyExists - this contact already exists
	ErrContactAlreadyExists = errors.New("contact already exists")
	// ErrContactBookUninitialized - contact book wasnt initialized
	ErrContactBookUninitialized = errors.New("contact book not initialized")

	// globalContactBook - the package global contact book
	Contacts ContactBook
	// contactsOnce - Only once do the global Contacts initialization
	contactsOnce sync.Once
)

func init() {
	contactsOnce.Do(func() {
		Contacts = NewContactBook()
	})
}

// Peer - Interface for a Peer to implement
type Peer interface {
	SendMessage(Message) error
	GetAddr() string
	GetPublicKey() crypto.PublicKey
}

// NewPeer - Create a new Peer
func NewPeer(addr string, pubKey crypto.PublicKey) Peer {
	return &peer{
		addr:   addr,
		pubKey: pubKey,
	}
}

// peer - implementation of Peer interface
type peer struct {
	addr     string
	pubKey   crypto.PublicKey
	isActive bool
}

// GetPublicKey - implementation of Peer interface
func (p *peer) GetPublicKey() crypto.PublicKey {
	return p.pubKey
}

// SendMessage - implementation of Peer interface
func (p *peer) SendMessage(m Message) error {
	// TODO: figure this out
	return nil
}

// GetAddr - implementation of Peer interface
func (p *peer) GetAddr() string {
	return p.addr
}

// ContactBook - Interface for describing how contacts are stored
type ContactBook interface {
	GetAll() ([]Peer, error)
	Get(string) (Peer, error)
	Add(Peer) error
	Remove(string) error
}

// NewContactBook - Create a new contact book
func NewContactBook() ContactBook {
	return &contactBook{
		m:         new(sync.RWMutex),
		peerIndex: make(map[string]int),
		peerList:  []Peer{},
	}
}

// contactBook - implementation of contact book, storing peers in a
// simple map of address to peer
type contactBook struct {
	m         *sync.RWMutex
	peerIndex map[string]int
	peerList  []Peer
}

// Remove - implementation of ContactBook
func (cb *contactBook) Remove(addr string) error {
	cb.m.Lock()
	defer cb.m.Unlock()
	// check if peer already exists
	if i, exists := cb.peerIndex[addr]; exists {
		cb.peerList = append(cb.peerList[:i], cb.peerList[i+1:]...)
		delete(cb.peerIndex, addr)
		return nil
	}
	return ErrNoSuchContact

}

// Add - implementation of ContactBook
func (cb *contactBook) Add(p Peer) error {
	cb.m.Lock()
	defer cb.m.Unlock()
	// check if peer already exists
	if _, exists := cb.peerIndex[p.GetAddr()]; exists {
		return ErrContactAlreadyExists
	}
	cb.peerList = append(cb.peerList, p)
	cb.peerIndex[p.GetAddr()] = len(cb.peerList) - 1
	return nil
}

// GetAll - return the full list of Peers
func (cb *contactBook) GetAll() ([]Peer, error) {
	if cb.m == nil || cb.peerIndex == nil || cb.peerList == nil {
		return nil, ErrContactBookUninitialized
	}
	cb.m.RLock()
	defer cb.m.RUnlock()
	return cb.peerList, nil
}

// Get - implementation of ContactBook interface
func (cb *contactBook) Get(addr string) (Peer, error) {
	cb.m.RLock()
	defer cb.m.RUnlock()
	if i, exists := cb.peerIndex[addr]; exists {
		return cb.peerList[i], nil
	}
	return nil, ErrNoSuchContact
}
