package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	TEAMCITY_TIMESTAMP_FORMAT = "2006-01-02T15:04:05.000"
)

type Test struct {
	Name, Output           string
	Race, Fail, Skip, Pass bool
}

var (
	input  = os.Stdin
	output = os.Stdout

	additionalTestName = ""

	run   = regexp.MustCompile("=== RUN\\s+([a-zA-Z_]\\S*)")
	end   = regexp.MustCompile("--- (PASS|SKIP|FAIL):\\s+([a-zA-Z_]\\S*) \\((-?[\\.\\d]+)")
	suite = regexp.MustCompile("^(ok|FAIL)\\s+([^\\s]+)\\s+(-?[\\.\\d]+)s")
	race  = regexp.MustCompile("^WARNING: DATA RACE")
)

func init() {
	flag.StringVar(&additionalTestName, "name", "", "Add prefix to test name")
}

func escapeOutput(outputLines []string) string {
	newOutput := strings.Join(outputLines, "\n")
	newOutput = strings.Replace(newOutput, "|", "||", -1)
	newOutput = strings.Replace(newOutput, "\n", "|n", -1)
	newOutput = strings.Replace(newOutput, "\r", "|n", -1)
	newOutput = strings.Replace(newOutput, "'", "|'", -1)
	newOutput = strings.Replace(newOutput, "]", "|]", -1)
	newOutput = strings.Replace(newOutput, "[", "|[", -1)
	return newOutput
}

func outputTest(w io.Writer, test *Test, out []string) {
	now := time.Now().Format(TEAMCITY_TIMESTAMP_FORMAT)
	var testName = additionalTestName + test.Name
	if test.Fail {
		fmt.Fprintf(w, "##teamcity[testFailed timestamp='%s' name='%s' details='%s']\n", now,
			testName, escapeOutput(out))
		fmt.Fprintf(w, "##teamcity[testFinished timestamp='%s' name='%s']\n", now, testName)
	} else if test.Race {
		fmt.Fprintf(w, "##teamcity[testFailed timestamp='%s' name='%s' message='Race detected!' details='%s']\n", now,
			testName, test.Output)
		fmt.Fprintf(w, "##teamcity[testFinished timestamp='%s' name='%s']\n", now, testName)
	} else if test.Skip {
		fmt.Fprintf(w, "##teamcity[testIgnored timestamp='%s' name='%s']\n", now, testName)
	} else if test.Pass {
		fmt.Fprintf(w, "##teamcity[testFinished timestamp='%s' name='%s']\n", now, testName)
	} else {
		fmt.Fprintf(w, "##teamcity[testFailed timestamp='%s' name='%s' message='Test ended in panic.' details='%s']\n", now,
			test.Name, escapeOutput(out))
		fmt.Fprintf(w, "##teamcity[testFinished timestamp='%s' name='%s']\n", now, test.Name)
	}
}

func processReader(r *bufio.Reader, w io.Writer) {
	var out []string
	var test *Test
	for {
		line, err := r.ReadString('\n')

		if err != nil {
			break
		}

		now := time.Now().Format(TEAMCITY_TIMESTAMP_FORMAT)

		runOut := run.FindStringSubmatch(line)
		if runOut != nil {
			if test != nil {
				if strings.HasPrefix(runOut[1], test.Name+"/") {
					// Just ignore subtests.
					continue
				}
				outputTest(w, test, out)
			}
			fmt.Fprintf(w, "##teamcity[testStarted timestamp='%s' name='%s']\n", now,
				additionalTestName+runOut[1])

			test = &Test{Name: runOut[1]}
			out = []string{}
			continue
		}

		endOut := end.FindStringSubmatch(line)
		if endOut != nil && test != nil {
			if endOut[1] == "FAIL" {
				test.Fail = true
			} else if endOut[1] == "SKIP" {
				test.Skip = true
			} else if endOut[1] == "PASS" {
				test.Pass = true
			}
			test.Name = endOut[2]
			test.Output = escapeOutput(out)
			if test.Pass || test.Skip {
				outputTest(w, test, out)
				test = nil
			}
			out = []string{}
			continue
		}

		suiteOut := suite.FindStringSubmatch(line)
		if suiteOut != nil {
			if test != nil {
				outputTest(w, test, out)
				out = []string{}
				test = nil
				continue
			}
		}

		if race.MatchString(line) {
			if test != nil {
				test.Race = true
			}
		}
		out = append(out, line[:len(line)-1])

		fmt.Fprint(w, line)
	}
	if test != nil {
		outputTest(w, test, out)
	}
}

func main() {
	flag.Parse()

	if len(additionalTestName) > 0 {
		additionalTestName += " "
	}

	reader := bufio.NewReader(input)

	processReader(reader, output)
}
