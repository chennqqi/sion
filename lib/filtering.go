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
func copyBody(body io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	var err error
    var buf bytes.Buffer
	if _, err = buf.ReadFrom(body); err != nil { return nil, nil, err }
    if err = body.Close(); err != nil {	return nil, nil, err	}
    return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewBuffer(buf.Bytes())), nil
}
func (p *ReverseProxy) ToSafeRequest(req *http.Request) (*http.Request, int, error) {	
	var (
		origin io.ReadCloser
		err error
		code int = 200
	)
	if origin, req.Body, err = copyBody(req.Body); err != nil{ return req, http.StatusInternalServerError, err }
	req.ParseForm()
	req.Body = origin
	
	filter := p.MakeFilterFromSelected(SelectEffectiveFilter(p.RequestFilters,req))
	if !Contains(req.Method, filter.AllowedMethod){
		return req, http.StatusMethodNotAllowed, errors.New("Method Not Allowed")
	}	
	for _, rule := range filter.Rules {
		var tocheck_values url.Values
		switch rule.Target {
		case "GET": tocheck_values = req.URL.Query()
		case "POST": tocheck_values = req.PostForm
		case "REGEX":
		}
		for _, param := range rule.Params{			
			if _, ok := tocheck_values[param.Key]; !ok{
				tocheck_values[param.Key] = []string{""}
			}
			for _, value := range tocheck_values[param.Key]{								
				if !param.Value.MatchString(value){
					if rule.ResponseCode != -1{
						code = rule.ResponseCode
					} else if rule.HandleTo != "" {
						req.URL.Path = rule.HandleTo
					} else if default_v, ok := rule.Defaults[param.Key]; ok {
						tocheck_values[param.Key] = []string{default_v}
						req.URL.RawQuery = tocheck_values.Encode()
					}
					return req, code, errors.New(fmt.Sprintf("Parameter Not Matched: Key=%s Value=%s Rule=%s",param.Key,req.FormValue(param.Key),param.Value.String()))
				}
			}
		}
	}
	return req, code, nil
}

