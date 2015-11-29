// Copyright 2015 - husobee associates, llc.  All rights reserved

// Package gonoxious - http.HandlerFuncs for gonoxious
package gonoxious

import (
	"errors"
	"sync"
)

var (
	ErrNoMessages = errors.New("No messages found")
)

type PeerStore interface {
	GetAll() ([]Peer, error)
	Get(string) (Peer, error)
	Store(Peer) error
}

type MessageStore interface {
	Get(string) ([]Message, error)
	Store(Message) error
}

type messageStore struct {
	m *sync.RWMutex
	s map[string][]Message
}

func (ms *messageStore) Get(addr string) ([]Message, error) {
	ms.m.RLock()
	defer ms.m.RUnlock()
	if _, exists := ms.s[addr]; !exists {
		return nil, ErrNoMessages
	}
	return ms.s[addr], nil
}

func (ms *messageStore) Store(m Message) error {
	ms.m.RLock()
	defer ms.m.RUnlock()
	from, err := m.GetFromAddr()
	if err != nil {
		return err
	}
	if _, exists := ms.s[from]; !exists {
		return ErrNoMessages
	}
	return nil
}

func NewMessageStore() MessageStore {
	return &messageStore{
		m: new(sync.RWMutex),
		s: make(map[string][]Message),
	}
}
