package sion

import(
	"net/http"
	"errors"
	"fmt"
	"bytes"
	"io"
	"io/ioutil"
	"net/url"
)

func Contains(elem string, list []string) bool { 
	for _, t := range list { if t == elem { return true } } 
	return false 
}
func copyBody(body io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	var err error
    var buf bytes.Buffer
	if _, err = buf.ReadFrom(body); err != nil { return nil, nil, err }
    if err = body.Close(); err != nil {	return nil, nil, err	}
    return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewBuffer(buf.Bytes())), nil
}       
func SelectEffectiveFilter(filters []RequestFilter,req *http.Request) []int {
	var enableFilters []int
	for index, filter := range filters {
		if filter.Location.MatchString(req.URL.Path){
			enableFilters = append(enableFilters,index)
		}
	}
	return enableFilters	
}
func (p *ReverseProxy) MakeFilterFromSelected(enableFilters []int) RequestFilter {
	var filter RequestFilter
	for _, index := range enableFilters {
		filter.Location = p.RequestFilters[index].Location
		filter.AllowedMethod = p.RequestFilters[index].AllowedMethod
		filter.Rules = p.RequestFilters[index].Rules
	}
	return filter
}
func (p *ReverseProxy) IsMethodAllowed(req *http.Request, filter RequestFilter) (int, error) {
	if !Contains(req.Method, filter.AllowedMethod){
		return http.StatusMethodNotAllowed, errors.New("Method Not Allowed")
	}
	return 200, nil
}
func (p *ReverseProxy) ToValidRequest(req *http.Request, filter RequestFilter) (code int, err error) {	
	var	origin io.ReadCloser
	if origin, req.Body, err = copyBody(req.Body); err != nil{ return http.StatusInternalServerError, err }
	req.ParseForm()
	req.Body = origin	
	for _, rule := range filter.Rules {
		var tocheck_values url.Values
		switch rule.Target {
		case "GET": tocheck_values = req.URL.Query()
		case "POST": tocheck_values = req.PostForm
		case "REGEX":
		}
		for _, param := range rule.Params{			
			if _, ok := tocheck_values[param.Key]; !ok {
				tocheck_values[param.Key] = []string{""}
			}
			for _, value := range tocheck_values[param.Key]{								
				if param.Value.MatchString(value){ continue }
				_, ok := rule.Defaults[param.Key]
				switch {
				case rule.HandleTo != "" : 
					req.URL.Path = rule.HandleTo
				case ok : 
					tocheck_values[param.Key] = []string{rule.Defaults[param.Key]}
					req.URL.RawQuery = tocheck_values.Encode()
					continue
				case true:
					return rule.ResponseCode, errors.New(fmt.Sprintf("Parameter Not Matched: Key=%s Value=%s Rule=%s",param.Key,req.FormValue(param.Key),param.Value.String()))
				}
			}
		}
	}
	return code, nil
}

