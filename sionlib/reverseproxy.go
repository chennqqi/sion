// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// HTTP reverse proxy handler

// From: https://github.com/methane/rproxy

package sion

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"regexp"
	"encoding/json"
	"io/ioutil"
)

// onExitFlushLoop is a callback set by tests to detect the state of the
// flushLoop() goroutine.
var onExitFlushLoop func()

// ReverseProxy is an HTTP Handler that takes an incoming request and
// sends it to another server, proxying the response back to the
// client.
type ReverseProxy struct {
	// Director must be a function which modifies
	// the request into a new request to be sent
	// using Transport. Its response is then copied
	// back to the original client unmodified.
	Director func(*http.Request)

	// The transport used to perform proxy requests.
	// If nil, http.DefaultTransport is used.
	Transport http.RoundTripper

	// FlushInterval specifies the flush interval
	// to flush to the client while copying the
	// response body.
	// If zero, no periodic flushing is done.
	FlushInterval time.Duration
	
	Config Config

	// Filters
	HeaderFilter Filter

	CookieFilter Filter

	UrlFilter Filter

}

type Rule struct{
	Target string `json:"target"`
	Regexp_ string `json:"Regexp"`
	Regexp *regexp.Regexp 
	Level int `json:"Level"`
}
type Filter struct{
	Rules []Rule `json:"rules"` 
}
type Config struct{
	HeaderFilterPath string `json:"header_filter_path"`
	UrlFilterPath string `json:"url_filter_path"`
	CookieFilterPath string `json:"cookie_filter_path"`
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// NewSingleHostReverseProxy returns a new ReverseProxy that rewrites
// URLs to the scheme, host, and base path provided in target. If the
// target's path is "/base" and the incoming request was for "/dir",
// the target request will be for /base/dir.
func NewSingleHostReverseProxy(target *url.URL,cfgpath string) *ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
	config,err := parseGeneralConfig(cfgpath)
	if err != nil{
		log.Printf("%v",err.Error())
	}
	headerFilter,err := parseFilterConfig(config.HeaderFilterPath)
	if err != nil{
		//TODO: default configration 
		log.Printf("%v",err.Error())
	}
	cookieFilter,err := parseFilterConfig(config.CookieFilterPath)
	if err != nil{
		//TODO: default configration 
		log.Printf("%v",err.Error())
	}
	urlFilter,err := parseFilterConfig(config.UrlFilterPath)
	if err != nil{
		//TODO: default configration 
		log.Printf("%v",err.Error())
	}
	return &ReverseProxy{Director: director, Config:config, HeaderFilter:headerFilter,CookieFilter:cookieFilter,UrlFilter:urlFilter}
}
func parseGeneralConfig(path string) (Config, error) {
	var config Config
	json_string, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("failed to read %s",path)
		return config, err
	}
	err = json.Unmarshal(json_string,&config)
	if err != nil {
		log.Printf("failed to load %s",path)
		return config, err
	}
	log.Printf("loaded %s",path)
	return config, nil
}

func parseFilterConfig(path string)(Filter,error){
	var filter Filter
	json_string, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("failed to read %s",path)
		return filter, err
	}
	err = json.Unmarshal(json_string,&filter)
	if err != nil {
		log.Printf("failed to load %s",path)
		return filter, err
	}
	for index,rule := range filter.Rules {
		filter.Rules[index].Regexp = regexp.MustCompile(rule.Regexp_)
	}
	log.Printf("loaded %s",path)
	return filter, nil
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func tcpProxy(rw http.ResponseWriter, outreq *http.Request) {
	clientConn, _, err := rw.(http.Hijacker).Hijack()
	if err != nil {
		panic("Fail to hijack.")
	}
	defer clientConn.Close()

	host := outreq.URL.Host
	if !strings.ContainsRune(host, ':') {
		host = host + ":80"
	}
	serverConn, err := net.Dial("tcp", host)
	if err != nil {
		panic("Can't connect to " + host)
	}
	defer serverConn.Close()

	// pass request
	outreq.Write(serverConn)
	go io.Copy(serverConn, clientConn)

	// pass response
	io.Copy(clientConn, serverConn)
}

func (p *ReverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	transport := p.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	outreq := new(http.Request)
	*outreq = *req // includes shallow copies of maps, but okay
	
	p.Director(outreq)
	outreq.Proto = "HTTP/1.1"
	outreq.ProtoMajor = 1
	outreq.ProtoMinor = 1
	outreq.Close = false

	validate_result_header := p.validateHeader(outreq.Header)
	if validate_result_header  != Safe{
	}
	validate_result_cookies := p.validateCookies(outreq.Cookies())
	if validate_result_cookies != Safe{
	}
	validate_result_url := p.validateURL(outreq.URL)
	if validate_result_url != Safe{
		log.Printf("It's danger! haha")
	}
	
	upgrading := outreq.Header.Get("Upgrade") == "websocket"

	if !upgrading && outreq.Header.Get("Connection") != "" {
		// Remove the connection header to the backend.  We want a
		// persistent connection, regardless of what the client sent
		// to us.  This is modifying the same underlying map from req
		// (shallow copied above) so we only copy it if necessary.
		outreq.Header = make(http.Header)
		copyHeader(outreq.Header, req.Header)
		outreq.Header.Del("Connection")
	}

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// If we aren't the first proxy retain prior
		// X-Forwarded-For information as a comma+space
		// separated list and fold multiple headers into one.
		if prior, ok := outreq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outreq.Header.Set("X-Forwarded-For", clientIP)
	}

	if upgrading {
		log.Println("hijacking:", outreq.URL)
		tcpProxy(rw, outreq)
		return
	}

	res, err := transport.RoundTrip(outreq)
	if err != nil {
		log.Printf("http: proxy error: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	copyHeader(rw.Header(), res.Header)

	rw.WriteHeader(res.StatusCode)
	p.copyResponse(rw, res.Body)
}

func (p *ReverseProxy) copyResponse(dst io.Writer, src io.Reader) {
	if p.FlushInterval != 0 {
		if wf, ok := dst.(writeFlusher); ok {
			mlw := &maxLatencyWriter{
				dst:     wf,
				latency: p.FlushInterval,
				done:    make(chan bool),
			}
			go mlw.flushLoop()
			defer mlw.stop()
			dst = mlw
		}
	}

	io.Copy(dst, src)
}

type writeFlusher interface {
	io.Writer
	http.Flusher
}

type maxLatencyWriter struct {
	dst     writeFlusher
	latency time.Duration

	lk   sync.Mutex // protects Write + Flush
	done chan bool
}

func (m *maxLatencyWriter) Write(p []byte) (int, error) {
	m.lk.Lock()
	defer m.lk.Unlock()
	return m.dst.Write(p)
}

func (m *maxLatencyWriter) flushLoop() {
	t := time.NewTicker(m.latency)
	defer t.Stop()
	for {
		select {
		case <-m.done:
			if onExitFlushLoop != nil {
				onExitFlushLoop()
			}
			return
		case <-t.C:
			m.lk.Lock()
			m.dst.Flush()
			m.lk.Unlock()
		}
	}
	panic("unreached")
}

func (m *maxLatencyWriter) stop() { m.done <- true }
