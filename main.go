package main

import (
	"os"
	"bufio"
	"regexp"
	"fmt"
	"time"
	"flag"
)

const (
	TEAMCITY_TIMESTAMP_FORMAT = "2006-01-02T15:04:05.000"
)

var (
	input = os.Stdin
	output = os.Stdout

	additionalTestName = ""

	run = regexp.MustCompile("=== RUN\\s+(\\w+)")
	pass = regexp.MustCompile("--- PASS:\\s+(\\w+) \\(([\\.\\d]+)s\\)")
	skip = regexp.MustCompile("--- SKIP:\\s+(\\w+)\\s+\\(([\\.\\d]+)s\\)")
	fail = regexp.MustCompile("--- FAIL:\\s+(\\w+)\\s+\\(([\\.\\d]+)s\\)")
)
func init() {
	flag.StringVar(&additionalTestName, "name", "", "Add prefix to test name")
	if len(additionalTestName) > 0 {
		additionalTestName += " "
	}
}

func main() {
	flag.Parse()

	reader := bufio.NewReader(input)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		now := time.Now().Format(TEAMCITY_TIMESTAMP_FORMAT)

		runOut := run.FindStringSubmatch(line)
		if runOut != nil {
			fmt.Fprintf(output, "##teamcity[testStarted timestamp='%s' name='%s']\n", additionalTestName + runOut[1],
				now)
			continue
		}

		passOut := pass.FindStringSubmatch(line)
		if passOut != nil {
			fmt.Fprintf(output, "##teamcity[testFinished timestamp='%s' name='%s']\n", additionalTestName + passOut[1],
				now)
			continue
		}

		skipOut := skip.FindStringSubmatch(line)
		if skipOut != nil {
			fmt.Fprintf(output, "##teamcity[testIgnored timestamp='%s' name='%s']\n", additionalTestName + skipOut[1],
				now)
			continue
		}

		failOut := fail.FindStringSubmatch(line)
		if failOut != nil {
			fmt.Fprintf(output, "##teamcity[testFailed timestamp='%s' name='%s']\n", additionalTestName + failOut[1],
				now)
			continue
		}

		fmt.Fprint(output, line)
	}
}
