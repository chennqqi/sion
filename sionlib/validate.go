package sion_rproxy

import (
	"net/http"
	"net/url"
	"regexp"
	"container/list"
	"log"
)

const (
	Safe = 0
	Doubtful = 1
	Danger = 2
)

type Rule struct{
	target string
	rule_regexp *regexp.Regexp
}

var header_filter *list.List = list.New()
var cookie_filter *list.List = list.New()
var url_filter *list.List = list.New()

func init(){
	log.Printf("test")
	url_filter.PushBack(Rule{"all",regexp.MustCompile("aa")})
}


func validateHeader(headers http.Header) int {
	return Safe
}
func validateCookies(cookies [](*http.Cookie)) int {
	return Safe
}
func validateURL(url *url.URL) int {
	for e := url_filter.Front(); e != nil; e = e.Next() {		
		if arule, ok := e.Value.(Rule); ok && arule.rule_regexp.MatchString(url.String()) {
			return Danger
		}	
	}
	return Safe
}
