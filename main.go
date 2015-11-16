package main

import (
	"os"
	"bufio"
	"regexp"
	"fmt"
	"time"
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

func main() {
	reader := bufio.NewReader(input)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		now := time.Now().Format(TEAMCITY_TIMESTAMP_FORMAT)

		runOut := run.FindStringSubmatch(line)
		if runOut != nil {
			fmt.Fprintf(output, "##teamcity[timestamp='%s' testStarted name='%s']\n", now,
				additionalTestName + runOut[1])
			continue
		}

		passOut := pass.FindStringSubmatch(line)
		if passOut != nil {
			fmt.Fprintf(output, "##teamcity[timestamp='%s' testFinished name='%s']\n", now,
				additionalTestName + passOut[1])
			continue
		}

		skipOut := skip.FindStringSubmatch(line)
		if skipOut != nil {
			fmt.Fprintf(output, "##teamcity[timestamp='%s' testIgnored name='%s']\n", now,
				additionalTestName + skipOut[1])
			continue
		}

		failOut := fail.FindStringSubmatch(line)
		if failOut != nil {
			fmt.Fprintf(output, "##teamcity[timestamp='%s' testFailed name='%s']\n", now,
				additionalTestName + failOut[1])
			continue
		}

		fmt.Fprint(output, line)
	}
}
