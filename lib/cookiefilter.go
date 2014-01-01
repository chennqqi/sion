package sion
import(
	"log"
	"io/ioutil"
	"encoding/json"
	"regexp"
)

type CookieFilter struct {
	Location string `json:"location"`
	Limit int `json:"limit"`
	Rules_ []map[string]string `json:"rules"`
	Rules []([]RegexpPair)
}

func LoadCookieFilters ( path string ) ([]CookieFilter, error) {
	var cfilter []CookieFilter
	jsonString, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("couldn't read %s : %s ", path, err.Error())
		return cfilter, err
	}
	err = json.Unmarshal(jsonString, &cfilter)
	if err != nil {
		log.Printf("couldn't load %s : %s ", path, err.Error())
		return cfilter, err
	}
	for f, filter := range cfilter {
		for r, rule := range filter.Rules_{
			cfilter[f].Rules = append(cfilter[f].Rules,[]RegexpPair{})
			for key, value := range rule {
				cfilter[f].Rules[r] = append(cfilter[f].Rules[r],RegexpPair{Key:*regexp.MustCompile(key), Value:*regexp.MustCompile(value)})
			}
		}
	}
	return cfilter, nil
}
