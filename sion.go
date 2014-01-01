package main
import (
	"net/http"
	"net/url"
	"time"
	"log"
	"./lib"
)

func main() {
	var src string
	src = ":8081"	
	u,_ := url.Parse("http://127.0.0.1:8080")
	h := sion.NewSingleHostReverseProxy(u,"config.json")
	s := &http.Server{
		Addr:           src,
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
