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

func contains(elem string, list []string) bool { 
	for _, t := range list { if t == elem { return true } } 
	return false 
}
func copyReader(body io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	var err error
    var buf bytes.Buffer
	if _, err = buf.ReadFrom(body); err != nil { return nil, nil, err }
    if err = body.Close(); err != nil {	return nil, nil, err	}
    return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewBuffer(buf.Bytes())), nil
}       
func copyBody(req *http.Request) (f url.Values, err error){
	var	origin io.ReadCloser
	if origin, req.Body, err = copyReader(req.Body); err != nil{ return nil, err }
	req.ParseForm()
	req.Body = origin
	return req.PostForm, nil
}
func modifyBody(req *http.Request, values url.Values){
	req.Body = ioutil.NopCloser(bytes.NewBufferString(values.Encode()))
	req.ContentLength = int64(len(values.Encode()))
}
func SelectEffectiveFilter(filters []RequestFilter,req *http.Request) (enable_filters []int) {
	for index, filter := range filters {
		if filter.Location.MatchString(req.URL.Path){
			enable_filters = append(enable_filters,index)
		}
	}
	return	
}
func (p *ReverseProxy) MakeFilterFromSelected(enableFilters []int) (filter RequestFilter) {
	// TODO:THIS IS MAKESHIFT 
	for _, index := range enableFilters {
		filter.Location = p.RequestFilters[index].Location
		filter.AllowedMethod = p.RequestFilters[index].AllowedMethod
		filter.Rules = p.RequestFilters[index].Rules
	}	
	return
}
func (p *ReverseProxy) IsMethodAllowed(req *http.Request, filter RequestFilter) (int, error) {
	if !contains(req.Method, filter.AllowedMethod){
		return http.StatusMethodNotAllowed, errors.New("Method Not Allowed")
	}
	return 200, nil
}
// by blacklist
func (p *ReverseProxy) CheckSafetyRequest(req *http.Request) (code int, err error) {	
	return 200, nil
}
func filterByRules(req *http.Request,values map[string]url.Values, rule Rule) (code int, err error){
	for _, param := range rule.Params{			
		for _, value := range values[rule.Target][param.Key]{								
			if param.Value.MatchString(value){ continue }
			_, ok := rule.Defaults[param.Key]
			switch {
			case ok : 
				values[rule.Target][param.Key] = []string{rule.Defaults[param.Key]}						
			case rule.HandleTo != "" : 
				req.URL.Path = rule.HandleTo
			default:
				return rule.ResponseCode, errors.New(fmt.Sprintf("Parameter Not Matched: Key=%s Value=%s Rule=%s",param.Key,req.FormValue(param.Key),param.Value.String()))
			}
		}
	}
	return 200, nil
}

// by whitelist
func (p *ReverseProxy) ToValidRequest(req *http.Request, filter RequestFilter) (code int, err error) {	
	var values = map[string]url.Values{ "POST":url.Values{}, "GET":req.URL.Query() }
	if values["POST"], err = copyBody(req); err != nil { return http.StatusInternalServerError, err }
	for _, rule := range filter.Rules {
		code, err = filterByRules(req, values, rule)
		if err != nil { return }
	}
	req.URL.RawQuery = values["GET"].Encode()
	modifyBody(req, values["POST"])	
	return code, nil
}

