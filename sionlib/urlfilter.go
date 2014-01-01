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
	Params []map[string]regexp.Regexp
}

type RagexpPair struct{
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
	for findex, ufelm := range ufilter {				
		for rindex, rule := range ufelm.Rules{
			for pindex, params_ := range rule.Params_{
				ufilter[findex].Rules[rindex].Params = append(ufilter[findex].Rules[rindex].Params,map[string]regexp.Regexp{})
				for key, value := range params_{
					ufilter[findex].Rules[rindex].Params[pindex][key] = *regexp.MustCompile(value)
				}
			}
		}
	}
	return ufilter, nil
}
