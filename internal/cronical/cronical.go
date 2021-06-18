package cronical

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/gorhill/cronexpr"
	log "github.com/sirupsen/logrus"
)

const (
	port           int = 8080
	icalTimeFormat     = "20060102T150405Z"
)

func Run() {
	log.Infof("running cronical on port %d", port)
	mux := http.NewServeMux()
	mux.HandleFunc("/webcal", webcalHandler)
	mux.Handle("/", http.FileServer(http.Dir("./html/")))
	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func webcalHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	encodedIcal := req.URL.Query().Get("ical")
	ical, err := decodeFilter(encodedIcal)
	if err != nil || len(ical) == 0 {
		log.Warnf("error decoding ical filter: %s", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	encodedExclude := req.URL.Query().Get("exclude")
	exclude := ""
	if encodedExclude != "" {
		exclude, err = decodeFilter(encodedExclude)
		if err != nil || len(exclude) == 0 {
			log.Warnf("error decoding exclude filter: %s", err)
			resp.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	log.
		WithField("ical", ical).
		WithField("cron_expression", exclude).
		WithField("url", req.URL.String()).
		Debug("handling /webcal")

	webcal, err := getWebcal(ical)
	if err != nil {
		log.Warnf("error getting webcal: %s", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	filteredWebcal := webcal
	if exclude != "" {
		filteredWebcal, err = filterWebcal(webcal, exclude)
		if err != nil {
			log.Warnf("error filtering webcal: %s", err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	resp.Header().Add("content-type", "text/calendar")
	resp.Write([]byte(filteredWebcal))
}

func filterWebcal(webcal, exclude string) (string, error) {
	cron, err := cronexpr.Parse(exclude)
	if err != nil {
		return "", err
	}
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
		start, err := time.Parse(icalTimeFormat, event.GetProperty(ics.ComponentPropertyDtStart).Value)
		if err != nil {
			return "", err
		}
		end, err := time.Parse(icalTimeFormat, event.GetProperty(ics.ComponentPropertyDtEnd).Value)
		if err != nil {
			return "", err
		}

		cronTime := cron.Next(start)
		log.
			WithField("start_time", start).
			WithField("end_time", end).
			WithField("cron_time", cronTime).
			WithField("cron_expression", exclude).
			Debug("evaluating event")
		if cronTime.After(start) && cronTime.Before(end) {
			continue
		}
		filteredIcal.AddEntireEvent(event)
	}
	return strings.Replace(filteredIcal.Serialize(), "\r", "", -1), nil
}

func getWebcal(webcalUrl string) (string, error) {
	resp, err := http.Get(strings.Replace(webcalUrl, "webcal", "http", 1))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code getting ical: %d", resp.StatusCode)
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
