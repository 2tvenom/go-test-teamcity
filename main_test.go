package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
)

var (
	inDir  = "testdata/input/"
	outDir = "testdata/output/"

	timestampRegexp      = regexp.MustCompile(`timestamp='.*?'`)
	timestampReplacement = `timestamp=''`
)

func TestProcessReader(t *testing.T) {
	files, err := ioutil.ReadDir(inDir)
	if err != nil {
		t.Error(err)
	}
	for _, file := range files {
		f, err := os.Open(inDir + file.Name())
		if err != nil {
			t.Error(err)
		}
		in := bufio.NewReader(f)

		out := &bytes.Buffer{}
		processReader(in, out)
		actual := out.String()
		actual = timestampRegexp.ReplaceAllString(actual, timestampReplacement)

		expectedBytes, err := ioutil.ReadFile(outDir + file.Name())
		if err != nil {
			t.Error(err)
		}
		expected := string(expectedBytes)
		expected = timestampRegexp.ReplaceAllString(expected, timestampReplacement)

		if strings.Compare(expected, actual) != 0 {
			t.Errorf("expected %s, but got %s", expected, actual)
		}
	}
}
