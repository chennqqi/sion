package sion

import(
	"net/http"
	"log"
)
                                                                                                                
func (p *ReverseProxy) isSafeRequest(req *http.Request) (bool,error) {
	var enableFilters []int
	for index, filter := range p.RequestFilters {
		if filter.Location.MatchString(req.URL.Path){
			enableFilters = append(enableFilters,index)
		}
	}
	log.Printf("%v",enableFilters)
	return true, nil
}

