package sion

import (
	"net/http"
	"net/url"
)

const Safe = 0

func (p *ReverseProxy) validateHeader(headers http.Header) int {
	return Safe
}
func (p *ReverseProxy) validateCookies(cookies []*(http.Cookie)) int {
	return Safe
}
func (p *ReverseProxy) validateURL(url *url.URL) int {
	for _,arule := range p.UrlFilter.Rules {
		if arule.Regexp.MatchString(url.String()){
			return arule.Level
		}
	}
	return Safe
}
