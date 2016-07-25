package main

import (
	"bufio"
	"flag"
	"fmt"
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

	run   = regexp.MustCompile("=== RUN\\s+(\\w+)")
	end   = regexp.MustCompile("--- (PASS|SKIP|FAIL):\\s+(\\w+) \\(([\\.\\d]+)s\\)")
	suite = regexp.MustCompile("^(ok|FAIL)\\w+(\\s+)\\w+([\\.\\d]+)s")
	race  = regexp.MustCompile("^WARNING: DATA RACE")
)

func init() {
	flag.StringVar(&additionalTestName, "name", "", "Add prefix to test name")
}

func main() {
	flag.Parse()

	if len(additionalTestName) > 0 {
		additionalTestName += " "
	}

	reader := bufio.NewReader(input)

	var out []string
	var test *Test
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		now := time.Now().Format(TEAMCITY_TIMESTAMP_FORMAT)

		runOut := run.FindStringSubmatch(line)
		if runOut != nil {
			if test != nil {
				var testName = additionalTestName + test.Name
				if test.Fail {
					fmt.Fprintf(output, "##teamcity[testFailed timestamp='%s' name='%s' details='%s']\n", now,
						testName, strings.Join(out, "|n"))
					fmt.Fprintf(output, "##teamcity[testFinished timestamp='%s' name='%s']\n", now, testName)
				} else if test.Race {
					fmt.Fprintf(output, "##teamcity[testFailed timestamp='%s' name='%s' message='Race detected!' details='%s']\n", now,
						testName, test.Output)
					fmt.Fprintf(output, "##teamcity[testFinished timestamp='%s' name='%s']\n", now, testName)
				} else if test.Skip {
					fmt.Fprintf(output, "##teamcity[testIgnored timestamp='%s' name='%s']\n", now, testName)
				} else if test.Pass {
					fmt.Fprintf(output, "##teamcity[testFinished timestamp='%s' name='%s']\n", now, testName)
				}
			}
			fmt.Fprintf(output, "##teamcity[testStarted timestamp='%s' name='%s']\n", now,
				additionalTestName+runOut[1])

			test = &Test{Name: runOut[1]}
			out = []string{}
			continue
		}

		endOut := end.FindStringSubmatch(line)
		if endOut != nil {
			if endOut[1] == "FAIL" {
				test.Fail = true
			} else if endOut[1] == "SKIP" {
				test.Skip = true
			} else if endOut[1] == "PASS" {
				test.Pass = true
			}
			test.Name = endOut[2]
			test.Output = strings.Join(out, "|n")
			out = []string{}
			continue
		}

		suiteOut := suite.FindStringSubmatch(line)
		if suiteOut != nil {
			if test != nil {
				fmt.Fprintf(output, "##teamcity[testFailed timestamp='%s' name='%s' message='Test ended in panic.' details='%s']\n", now,
					testName, test.Output)
				fmt.Fprintf(output, "##teamcity[testFinished timestamp='%s' name='%s']\n", now, testName)
			}
		}

		if race.MatchString(line) {
			if test != nil {
				test.Race = true
			}
		}
		out = append(out, line[:len(line)-1])

		fmt.Fprint(output, line)
	}
}
