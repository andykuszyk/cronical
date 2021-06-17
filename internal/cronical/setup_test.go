package cronical

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var ws *webcalServer

func TestMain(m *testing.M) {
	logrus.SetLevel(logrus.DebugLevel)
	ws = newWebcalServer()
	go Run()
	m.Run()
}

type webcalServer struct {
	port    int
	webcals map[string]string
}

func newWebcalServer() *webcalServer {
	ws := &webcalServer{
		port:    8081,
		webcals: make(map[string]string),
	}
	go ws.start()
	return ws
}

func (ws *webcalServer) start() {
	http.HandleFunc("/webcalmock", ws.handler)
	http.ListenAndServe(fmt.Sprintf(":%d", ws.port), nil)
}

func (ws *webcalServer) handler(resp http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if webcal, ok := ws.webcals[id]; ok {
		resp.Write([]byte(webcal))
		return
	}
	resp.WriteHeader(http.StatusNotFound)
}

func (ws *webcalServer) addWebcal(webcal string) (string, error) {
	id := uuid.New().String()
	ws.webcals[id] = webcal
	u, err := url.Parse(fmt.Sprintf("webcal://localhost:%d/webcalmock?id=%s", ws.port, id))
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
