package sion
import(
	"log"
	"io/ioutil"
	"encoding/json"
)

type UrlFilter struct{
	Location string `json:"location"`
	AllowMethod []string `json:"allow-method"`
	Rules []Rule `json:"rules"`
}

type Rule struct{
	Target string `json:"target"`
	Params []map[string]string `json:"params"`
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
	return ufilter, nil
}
