package cronical

import (
	"fmt"
	"log"
	"net/http"
)

const port int = 8080

func Run() {
	http.HandleFunc("/filter", handler)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func handler(resp http.ResponseWriter, req *http.Request) {
	ical := req.URL.Query().Get("ical")
	exclude := req.URL.Query().Get("exclude")
	log.Printf("handle with ical: %s and exclude: %s and url: %s", ical, exclude, req.URL.String())
}
