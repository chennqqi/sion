package sion

import(
	"net/http"
	"errors"
)

func Contains(elem string, list []string) bool { 
	for _, t := range list { if t == elem { return true } } 
	return false 
}
                                                                                                                
func (p *ReverseProxy) isSafeRequest(req *http.Request) (bool,error) {
	var enableFilters []int
	var highPriorityFilter int = 0
	for index, filter := range p.RequestFilters {
		if filter.Location.MatchString(req.URL.Path){
			enableFilters = append(enableFilters,index)
			if index > highPriorityFilter { highPriorityFilter = index }
		}
	}
	if !Contains(req.Method, p.RequestFilters[highPriorityFilter].AllowMethod){
		return false, errors.New("Method Not Allowed")
	}
	return true, nil
}


