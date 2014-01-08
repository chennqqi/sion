package sion
import(
	"log"
	"io/ioutil"
	"encoding/json"
	"regexp"
	"strings"
	"sort"
)

type RequestFilterRaw struct{
	Location string `json:"location"`	
	AllowedMethod []string `json:"allowed-method"`
	Rules []map[string]string `json:"rules"`

}
type RequestFilter struct{
	Location regexp.Regexp
	AllowedMethod []string
	Rules []Rule	
	Priority int
}

type Rule struct{
	Target string
	Params []ParamKeyValue
	Options []([]string)
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
	var ufilterraw []RequestFilterRaw
	var ufilter RequestFilters
	jsonString, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("couldn't read %s : %s ", path, err.Error())
		return ufilter, err
	}
	err = json.Unmarshal(jsonString, &ufilterraw)
	if err != nil {
		log.Printf("couldn't load %s : %s ", path, err.Error())
		return ufilter, err
	}
	for f, ufelm := range ufilterraw {	
		ufilter = append(ufilter, RequestFilter{})		
		ufilter[f].Location = *regexp.MustCompile(ufelm.Location)
		ufilter[f].AllowedMethod = ufelm.AllowedMethod		
		if ufelm.Location == "/" { 
			ufilter[f].Priority = 0 
		} else { 
			ufilter[f].Priority = strings.Count(ufelm.Location,"/") 
		}
		for _, rawrule := range ufelm.Rules{ 
			var rule Rule			
			rule.Target = rawrule["[target]"]
			for key, value := range rawrule{
				if strings.HasPrefix(key,"[") && strings.HasSuffix(key,"]"){
					rule.Options = append(rule.Options,[]string{strings.TrimPrefix(strings.TrimSuffix(key,"]"),"["),value})
				}
				rule.Params = append(rule.Params,ParamKeyValue{Key:key,Value:*regexp.MustCompile(value)})
			}
			ufilter[f].Rules = append(ufilter[f].Rules,rule)
		}
	}
	sort.Sort(ufilter)
	debug(ufilter)
	return ufilter, nil
}

func debug(filters RequestFilters){
	for _, filter := range filters{
		log.Printf("####")
		log.Printf("Location: %s",filter.Location.String())
		log.Printf("Priority: %d",filter.Priority)
		log.Printf("AllowedMethod: %v",filter.AllowedMethod)		
		for _, rule := range filter.Rules{
			log.Printf("---")
			log.Printf("Rule Target: %s",rule.Target)
			for _, param := range rule.Params{
				log.Printf("Rule %s : %s",param.Key,param.Value.String())
			}			
			for _, option := range rule.Options{
				log.Printf("Option %s : %s",option[0],option[1])
			}
		}
		log.Printf("-----")
	}
}

