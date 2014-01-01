package sion
import(
	"log"
	"io/ioutil"
	"encoding/json"
	"regexp"
)

type UrlFilter struct{
	Location string
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

func LoadUrlFilters ( path string ) ([]UrlFilter, error) {
	var ufilter []UrlFilter
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
