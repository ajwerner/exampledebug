// Command exampledebug takes a test output for a golang Example and computes a
// diff.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"syscall"
)

var re = regexp.MustCompile("(?s).*?--- FAIL: (\\w+).*?\ngot:\n(.*?)\nwant:\n(.*?)\n(--- FAIL|FAIL)")

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	exitIfError(err, "failed to read input from stdin")
	for len(input) > 0 {
		submatches := re.FindSubmatch(input)
		if submatches == nil {
			return
		}
		//fmt.Printf("%q\n%q\n%q\n", string(submatches[1]), string(submatches[2]), string(submatches[3]))
		test := string(submatches[1])
		got := submatches[2]
		want := submatches[3]
		diff, err := computeDiff("got", "want", got, want)
		exitIfError(err, fmt.Sprintf("failed to compute diff for test %s", test))
		fmt.Printf("%s:\n%s\n---\n", test, string(diff))
		input = input[len(submatches[0])-len(submatches[4]):]
	}
}

func exitIfError(err error, message string) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}

func computeDiff(labelA, labelB string, a, b []byte) ([]byte, error) {
	aFile, err := makeTempFile("a", a)
	exitIfError(err, "failed to create temp file")
	defer os.Remove(aFile)
	bFile, err := makeTempFile("b", b)
	exitIfError(err, "failed to create temp file")
	defer os.Remove(bFile)
	cmd := exec.Command("diff", "-u", "--label", labelA, "--label", labelB, aFile, bFile)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start diff: %v", err)
	}
	if ps, err := cmd.Process.Wait(); err != nil || ps.Sys().(syscall.WaitStatus).ExitStatus() > 1 {
		return nil, fmt.Errorf("failed to run diff: %v %v", ps.Sys().(syscall.WaitStatus).ExitStatus(), err)
	}
	return stdout.Bytes(), nil
}

func makeTempFile(pattern string, data []byte) (path string, err error) {
	f, err := ioutil.TempFile("", pattern)
	if err != nil {
		return "", err
	}
	f.Close()
	if err = ioutil.WriteFile(f.Name(), data, 0755); err != nil {
		os.Remove(f.Name())
		return "", err
	}
	return f.Name(), nil
}
