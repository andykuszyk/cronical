package cronical

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInputsAndOutputs(t *testing.T) {
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

		})
	}
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
