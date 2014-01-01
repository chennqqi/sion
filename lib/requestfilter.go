package sion
import(
	"log"
	"io/ioutil"
	"encoding/json"
	"regexp"
)

type RequestFilter struct{
	Location regexp.Regexp
	Location_ string `json:"location"`	
	AllowMethod []string `json:"allow-method"`
	Rules []Rule `json:"rules"`
}

type Rule struct{
	Target string `json:"target"`
	Params_ []map[string]string `json:"params"`
	Params []([]RegexpPair)
}

type RegexpPair struct{
	Key regexp.Regexp
	Value regexp.Regexp
}

func LoadRequestFilters ( path string ) ([]RequestFilter, error) {
	var ufilter []RequestFilter
	jsonString, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("couldn't read %s : %s ", path, err.Error())
		return ufilter, err
	}
	err = json.Unmarshal(jsonString, &ufilter)
	if err != nil {
		log.Printf("couldn't load %s : %s ", path, err.Error())
		return ufilter, err
	}
	for f, ufelm := range ufilter {	
		ufilter[f].Location = *regexp.MustCompile(ufilter[f].Location_)
		for r, rule := range ufelm.Rules{ 
			for p, params_ := range rule.Params_{
				ufilter[f].Rules[r].Params = append(ufilter[f].Rules[r].Params,[]RegexpPair{})
				for key, value := range params_{
					ufilter[f].Rules[r].Params[p] = append(ufilter[f].Rules[r].Params[p],RegexpPair{Key:*regexp.MustCompile(key),Value:*regexp.MustCompile(value)})
				}
			}
		}
	}
	return ufilter, nil
}
