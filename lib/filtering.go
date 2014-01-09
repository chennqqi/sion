package sion

import(
	"net/http"
	"errors"
	"fmt"
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
		filter.Rules = append(filter.Rules, p.RequestFilters[index].Rules...)
		//TODO: priority !!
	}
	return filter
}
func (p *ReverseProxy) IsSafeRequest(req *http.Request) (bool,error) {
	filter := p.MakeFilterFromSelected(SelectEffectiveFilter(p.RequestFilters,req))
	if !Contains(req.Method, filter.AllowedMethod){
		return false, errors.New("Method Not Allowed")
	}
	req.ParseForm()
	for _, rule := range filter.Rules {
		var tocheck_values map[string][]string
		switch rule.Target {
		case "GET": tocheck_values = req.URL.Query()
		case "POST": tocheck_values = req.PostForm
		case "REGEX":
		}
		for _, param := range rule.Params{
			for _, value := range tocheck_values[param.Key]{
				if !param.Value.MatchString(value){
					req.URL.Path = rule.HandleTo
					return false, errors.New(fmt.Sprintf("Parameter Not Matched: Key=%s Value=%s Rule=%s",param.Key,req.FormValue(param.Key),param.Value.String()))
				}
			}
		}
	}
	return true, nil
}


