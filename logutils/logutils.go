// Copyright 2015 - husobee associates, llc.  All rights reserved

// Package logutils - log utilities for gonoxious
package logutils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/golang/glog"
)

const (
	// Fatal glog level
	Fatal glog.Level = 0
	// Error glog level
	Error = 2
	// Warning glog level
	Warning = 3
	// Info glog level
	Info = 4
	// Debug glog level
	Debug = 5
)

// LogMiddleware - Negroni Logging middleware
type LogMiddleware struct {
	level glog.Level
}

// NewLogMiddleware - Create a new LogMiddleware and initialize with a log level
func NewLogMiddleware(level glog.Level) *LogMiddleware {
	return &LogMiddleware{
		level: level,
	}
}

// ServeHTTP - Implementation of negroni middleware handler
func (lm *LogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	s := time.Now()
	remoteAddr := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		remoteAddr = realIP
	}

	next(w, r)
	res := w.(negroni.ResponseWriter)

	d := time.Since(s)
	var duration string

	if d > 1*time.Second {
		duration = fmt.Sprintf("%d s", d/(1*time.Second))
	} else if d > 1*time.Millisecond {
		duration = fmt.Sprintf("%d ms", d/(1*time.Millisecond))
	} else if d > 1*time.Microsecond {
		duration = fmt.Sprintf("%d μs", d/(1*time.Microsecond))
	} else {
		duration = fmt.Sprintf("%d ns", d)
	}

	var color string
	if res.Status() < 400 && res.Status() >= 200 {
		color = string(0x1b) + "[32m"
	} else if res.Status() < 500 && res.Status() >= 400 {
		color = string(0x1b) + "[33m"
	} else if res.Status() < 0 && res.Status() >= 100 {
		color = string(0x1b) + "[34m"
	} else if res.Status() > 500 {
		color = string(0x1b) + "[31m"
	}
	clear := string(0x1b) + "[0m"

	if glog.V(lm.level) {
		// TODO: get the user information from the request
		glog.Infof(`%s%s %s %s (remote=%s) (status=%d "%s") (size=%d) (took=%s)%s`,
			color,
			r.Method,
			r.URL.Path,
			r.Proto,
			remoteAddr,
			res.Status(),
			http.StatusText(res.Status()),
			res.Size(),
			duration,
			clear,
		)
	}

}
