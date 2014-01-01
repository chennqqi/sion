package sion
import(
	"log"
	"io/ioutil"
	"encoding/json"
)
type Config struct{
	RequestFilterPath string `json:"request_filter_path"`
	CookieFilterPath string `json:"cookie_filter_path"`	
}

func LoadConfig ( path string ) (Config, error) {
	var config Config
	jsonString, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("couldn't read %s : %s ", path, err.Error())
		return config, err
	}
	err = json.Unmarshal(jsonString, &config)
	if err != nil {
		log.Printf("couldn't read %s : %s ", path, err.Error())
		return config, err
	}
	return config, nil
}
