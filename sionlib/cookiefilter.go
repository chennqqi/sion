package sion
import(
	"log"
	"io/ioutil"
	"encoding/json"
)

type CookieFilter struct {
	Location string `json:"location"`
	Limit int `json:"limit"`
	Rules []map[string]string `json:"rules"`
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
	return cfilter, nil
}
