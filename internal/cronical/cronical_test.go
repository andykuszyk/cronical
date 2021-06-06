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
		tc, err := buildTestCase(testcase.Name())
		require.NoError(t, err)
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}

func buildTestCase(dir string) (testCase, error) {
	return testCase{}, nil
}

type testCase struct {
	name string
}
