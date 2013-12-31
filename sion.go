package main
import (
	"net/http"
	"net/url"
	"time"
	"log"
	"./sionlib"
)

func main() {
	var src string
	src = ":8081"	
	u, e := url.Parse("http://127.0.0.1:8080")
	h := sion_rproxy.NewSingleHostReverseProxy(u)
	s := &http.Server{
		Addr:           src,
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
