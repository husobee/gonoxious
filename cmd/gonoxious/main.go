// Copyright 2015 - husobee associates, llc.  All rights reserved

// Package main - entry point for gonoxious client
package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/golang/glog"
	"github.com/husobee/gonoxious/handlers"
	"github.com/husobee/gonoxious/logutils"
	"github.com/husobee/vestigo"
	"github.com/tylerb/graceful"
)

var (
	// Flags

	//ListenAddr - address to listen on
	ListenAddr string
	// PrivKeyFilename - filename for the PEM encoded Private Key File
	PrivKeyFileName string
	// PubKeyFilename - filename for the PEM encoded Public Key File
	PubKeyFileName string
)

func init() {
	flag.StringVar(&PrivKeyFileName, "priv_key", "~/.noxious/privkey.pem", "Private Key for signing noxious messages")
	flag.StringVar(&PubKeyFileName, "pub_key", "~/.noxious/pubkey.pem", "Public Key for receiving encrypted noxious messages")
	flag.StringVar(&ListenAddr, "addr", ":1111", "Listen Address for web server")
	// parse command line flags
	flag.Parse()
	// Load Keys
}

func main() {
	// logger setup
	defer glog.Flush()
	if glog.V(logutils.Info) {
		glog.Infof("Starting gonoxious client, listening on %s", ListenAddr)
		if glog.V(logutils.Debug) {
			glog.Infof("start time: %v, listen: %s, privkey: %s, pubkey: %s",
				time.Now(),
				ListenAddr,
				PrivKeyFileName,
				PubKeyFileName)
		}
	}
	// load up keypair, if no keypair, generate keypair
	// TODO: get all the crypto in place

	// negroni middleware setup, with custom logger
	n := negroni.New(
		negroni.NewRecovery(),
		logutils.NewLogMiddleware(logutils.Info),
	)

	// setup URL router
	router := vestigo.NewRouter()
	// chat message endpoint:
	router.Post("/", handlers.ChatHandler)

	n.UseHandler(router)
	// graceful start/stop server
	srv := &graceful.Server{
		Timeout: 5 * time.Second,
		Server:  &http.Server{Addr: ListenAddr, Handler: n},
	}
	srv.ListenAndServe()
	glog.Fatal("ending gonoxious")
}
