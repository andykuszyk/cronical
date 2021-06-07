package cronical

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	ics "github.com/arran4/golang-ical"
	log "github.com/sirupsen/logrus"
)

const port int = 8080

func Run() {
	http.HandleFunc("/filter/", handler)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func handler(resp http.ResponseWriter, req *http.Request) {
	encodedIcal := req.URL.Query().Get("ical")
	ical, err := decodeFilter(encodedIcal)
	if err != nil || len(ical) == 0 {
		log.Warnf("error decoding ical filter: %s", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	encodedExclude := req.URL.Query().Get("exclude")
	exclude, err := decodeFilter(encodedExclude)
	if err != nil || len(exclude) == 0 {
		log.Warnf("error decoding exclude filter: %s", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Infof("handle with ical: %s and exclude: %s and url: %s", ical, exclude, req.URL.String())

	webcal, err := getWebcal(ical)
	if err != nil {
		log.Warnf("error getting webcal: %s", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	filteredWebcal, err := filterWebcal(webcal, exclude)
	if err != nil {
		log.Warnf("error filtering webcal: %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.Header().Add("content-type", "text/calendar")
	resp.Write([]byte(filteredWebcal))
}

func filterWebcal(webcal, exclude string) (string, error) {
	ical, err := ics.ParseCalendar(strings.NewReader(webcal))
	if err != nil {
		return "", err
	}
	filteredIcal, err := ics.ParseCalendar(strings.NewReader(webcal))
	if err != nil {
		return "", err
	}
	filteredIcal.ClearEvents()

	for _, event := range ical.Events() {
		// TODO skip filtered events
		filteredIcal.AddEntireEvent(event)
	}
	return filteredIcal.Serialize(), nil
}

func getWebcal(webcalUrl string) (string, error) {
	resp, err := http.Get(webcalUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func encodeFilter(filter string) string {
	return base64.StdEncoding.EncodeToString([]byte(filter))
}

func decodeFilter(filter string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(filter)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
