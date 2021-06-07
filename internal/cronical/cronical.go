package cronical

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
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
		logrus.Warnf("error decoding ical filter: %s", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	encodedExclude := req.URL.Query().Get("exclude")
	exclude, err := decodeFilter(encodedExclude)
	if err != nil || len(exclude) == 0 {
		logrus.Warnf("error decoding exclude filter: %s", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	logrus.Infof("handle with ical: %s and exclude: %s and url: %s", ical, exclude, req.URL.String())
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
