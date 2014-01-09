package sion
import(
	"log"
	"io/ioutil"
	"encoding/json"
	"regexp"
	"strings"
	"sort"
	"strconv"
)

type RequestFilterRaw struct{
	Location string `json:"location"`	
	AllowedMethod []string `json:"allowed_method"`
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
	HandleTo string
	ResponseCode int
	Params []ParamKeyValue
	Options map[string]string	
	Defaults map[string]string
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
			rule.Params, rule.Options, rule.Defaults = []ParamKeyValue{}, map[string]string{}, map[string]string{}
			for key, value := range rawrule{
				if strings.HasPrefix(key,":") {
					rule.Options[strings.TrimPrefix(key, ":")] = value
				} else if strings.HasPrefix(key,"@") {
					rule.Defaults[strings.TrimPrefix(key, "@")] = value
				} else {
					rule.Params = append(rule.Params,ParamKeyValue{Key:key,Value:*regexp.MustCompile(value)})
				}
			}
			rule.Target = rule.Options["target"]
			rule.HandleTo = rule.Options["handle_to"]
			rule.ResponseCode, err = strconv.Atoi(rule.Options["response_code"])
			if err != nil { rule.ResponseCode = -1 } 
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
			log.Printf("Target: %s",rule.Target)
			for key, value := range rule.Options{
				log.Printf("Option %s : %s", key, value)
			}
			for key, value := range rule.Defaults{
				log.Printf("Default %s : %s", key, value)
			}			
			for _, param := range rule.Params{
				log.Printf("Rule %s : %s",param.Key,param.Value.String())
			}			
		}
		log.Printf("-----")
	}
}

