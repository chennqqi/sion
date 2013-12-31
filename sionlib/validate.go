package sion_rproxy

import (
	"net/http"
)

const (
	Safe = 0
	Doubtful
	Danger
)

func validate(req *http.Request) int {
	if validateCookies(req.Cookies()) != Safe {
	}
	if validateHeader(req.Header) != Safe{
	}
	return Safe
}

func validateHeader(headers http.Header) int {
	return Safe
}
func validateCookies(cookies [](*http.Cookie)) int {
	return Safe
}
func validateURI(url *url.URL) int {
	return Safe
}
