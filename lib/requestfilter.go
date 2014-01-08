package sion
import(
	"log"
	"io/ioutil"
	"encoding/json"
	"regexp"
	"strings"
	"sort"
)

type RequestFilter struct{
	Location regexp.Regexp
	Location_ string `json:"location"`	
	AllowedMethod []string `json:"allowed-method"`
	Rules_ []map[string]string `json:"rules"`
	Rules []Rule
	Priority int
}

type Rule struct{
	Target string
	Params []ParamKeyValue
}

type ParamKeyValue struct{
	Key string
	Value regexp.Regexp
}


type RequestFilters []RequestFilter
func (p RequestFilters) Len() int{ return len(p) }
func (p RequestFilters) Swap(i,j int) { p[i] ,p[j] = p[j] ,p[i] }
func (p RequestFilters) Less(i,j int) bool { return p[i].Priority < p[j].Priority }

func LoadRequestFilters ( path string ) ([]RequestFilter, error) {
	var ufilter RequestFilters
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
		if ufilter[f].Location_ == "/" { 
			ufilter[f].Priority = 0 
		} else { 
			ufilter[f].Priority = strings.Count(ufilter[f].Location_,"/") 
		}
		for _, rawrule := range ufelm.Rules_{ 
			var rule Rule
			rule.Target = rawrule["[target]"]
			for key, value := range rawrule{
				rule.Params = append(rule.Params,ParamKeyValue{Key:key,Value:*regexp.MustCompile(value)})
			}
			ufilter[f].Rules = append(ufilter[f].Rules,rule)
		}
	}
	sort.Sort(ufilter)
	return ufilter, nil
}

