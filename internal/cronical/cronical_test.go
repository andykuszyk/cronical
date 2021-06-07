package cronical

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInputsAndOutputs(t *testing.T) {
	root, err := os.Getwd()
	require.NoError(t, err)

	testdata := filepath.Join(root, "testdata")
	testcases, err := os.ReadDir(testdata)
	require.NoError(t, err)

	ws := newWebcalServer()

	for _, testcase := range testcases {
		if !testcase.IsDir() {
			continue
		}
		tc := buildTestCase(t, testcase.Name(), root)
		t.Run(tc.name, func(t *testing.T) {
			webcalUrl, err := ws.addWebcal(tc.input)
			assert.NoError(t, err)
			actual, err := cronicalGet(webcalUrl, tc.exclude)
			assert.NoError(t, err)
			assert.Equal(t, tc.output, actual)
		})
	}
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
	http.HandleFunc("/", ws.handler)
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
	u, err := url.Parse(fmt.Sprintf("http://localhost:%d?id=%s", ws.port, id))
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func cronicalGet(webcalUrl, exclude string) (string, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"http://localhost:%d/filter?ical=%s&exclude=%s",
			port,
			encodeFilter(webcalUrl),
			encodeFilter(exclude),
		),
		nil)
	if err != nil {
		return "", err
	}
	logrus.Infof("getting croncial at: %s", request.URL.String())
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func buildTestCase(t *testing.T, dir string, rootDir string) testCase {
	input, err := os.ReadFile(filepath.Join(rootDir, "testdata", dir, "input.ical"))
	require.NoError(t, err)
	output, err := os.ReadFile(filepath.Join(rootDir, "testdata", dir, "output.ical"))
	require.NoError(t, err)
	exclude, err := os.ReadFile(filepath.Join(rootDir, "testdata", dir, "exclude.cron"))
	require.NoError(t, err)
	return testCase{
		input:   string(input),
		output:  string(output),
		exclude: string(exclude),
		name:    filepath.Base(dir),
	}
}

type testCase struct {
	input   string
	output  string
	exclude string
	name    string
}
