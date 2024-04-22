package testrunner

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sigs.k8s.io/yaml"
)

type TestRunner struct {
	T    *testing.T
	Path string
}

func NewTestRunner(t *testing.T, path string) *TestRunner {
	return &TestRunner{
		T:    t,
		Path: path,
	}
}

func (tr *TestRunner) Run() {
	var t = tr.T

	testCases, err := tr.loadFromPath(tr.Path)
	if err != nil {
		t.Fatalf(err.Error())
	}

	hasFocus := false
	for _, testCase := range testCases {
		if testCase.Focus {
			if hasFocus {
				t.Fatal("multiple focus test cases")
			}

			hasFocus = true
			break
		}
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.FullName(), func(t *testing.T) {
			t.Parallel()

			if hasFocus && !tc.Focus {
				t.SkipNow()
				return
			}

			tc.Run(t)
		})
	}
}

func (tr *TestRunner) loadFromPath(path string) ([]TestCase, error) {
	var t = tr.T
	var testCases []TestCase
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, "test.yaml") == false {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("failed to read file %s: %v", path, err)
		}

		documents := bytes.SplitN(data, []byte("\n---"), -1)

		for i, testcaseRaw := range documents {
			var testCase TestCase

			jsonData, err := yaml.YAMLToJSON(testcaseRaw)
			if err != nil {
				t.Errorf("failed to convert yaml to json at %s#%d%v", path, i+1, err)
				continue
			}

			var invalidSchemaMessage = fmt.Sprintf("invalid schema at %s#%d", path, i+1)
			err = TestCaseSchema.Validator().ValidateWithMessage(bytes.NewReader(jsonData), invalidSchemaMessage)
			if err != nil {
				t.Errorf(err.Error())
				continue
			}

			testCase.Path = path
			testCases = append(testCases, testCase)
		}

		return nil
	})

	return testCases, err
}
