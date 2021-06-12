package cronical

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/gin-gonic/gin"
	"github.com/gorhill/cronexpr"
	log "github.com/sirupsen/logrus"
)

const (
	port           int = 8080
	icalTimeFormat     = "20060102T150405Z"
)

func Run() {
	dir, err := os.Getwd()
	if err != nil {
		log.Errorf("error getting working directory: %s", err)
		os.Exit(1)
	}
	log.Infof("starting cronical with working directory: %s", dir)

	r := gin.Default()
	r.GET("/filter/", filterHandler)
	r.GET("/webcal/", webcalHandler)
	r.Static("/html/", filepath.Join(dir, "html"))

	log.Infof("running cronical on port %d", port)
	r.Run(fmt.Sprintf(":%d", port))
}

func webcalHandler(c *gin.Context) {
	req := c.Request
	resp := c.Writer

	encodedIcal := req.URL.Query().Get("ical")
	ical, err := decodeFilter(encodedIcal)
	if err != nil || len(ical) == 0 {
		log.Warnf("error decoding ical filter: %s", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	webcal, err := getWebcal(ical)
	if err != nil {
		log.Warnf("error getting webcal: %s", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	resp.Write([]byte(webcal))
}

func filterHandler(c *gin.Context) {
	req := c.Request
	resp := c.Writer

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

	log.
		WithField("ical", ical).
		WithField("cron_expression", exclude).
		WithField("url", req.URL.String()).
		Debug("handling /filter")

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
	resp, err := http.Get(webcalUrl)
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
