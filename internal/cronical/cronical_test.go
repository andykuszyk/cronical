package cronical

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	root, err := os.Getwd()
	require.NoError(t, err)

	testdata := filepath.Join(root, "testdata")
	testcases, err := os.ReadDir(testdata)
	require.NoError(t, err)

	for _, testcase := range testcases {
		if !testcase.IsDir() {
			continue
		}
		tc := buildTestCase(t, testcase.Name(), root)
		t.Run(tc.name, func(t *testing.T) {
			webcalUrl, err := ws.addWebcal(tc.input)
			assert.NoError(t, err)
			actual, err := cronicalGetFilter(webcalUrl, tc.exclude)
			assert.NoError(t, err)
			assert.Equal(t, tc.output, actual)
		})
	}
}

func TestWebcal(t *testing.T) {
	root, err := os.Getwd()
	require.NoError(t, err)

	testdata := filepath.Join(root, "testdata")
	testcases, err := os.ReadDir(testdata)
	require.NoError(t, err)

	testcase := buildTestCase(t, testcases[0].Name(), root)
	webcalUrl, err := ws.addWebcal(testcase.input)
	require.NoError(t, err)
	actual, err := cronicalGetWebcal(webcalUrl)
	assert.NoError(t, err)
	assert.Equal(t, testcase.input, actual)
}

func cronicalGetFilter(webcalUrl, exclude string) (string, error) {
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

func cronicalGetWebcal(webcalUrl string) (string, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"http://localhost:%d/webcal?ical=%s",
			port,
			encodeFilter(webcalUrl),
		),
		nil)
	if err != nil {
		return "", err
	}
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
