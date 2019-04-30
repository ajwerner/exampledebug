package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestExampleDebug(t *testing.T) {
	defer func(stdin, stdout *os.File) {
		os.Stdin, os.Stdout = stdin, stdout
	}(os.Stdin, os.Stdout)
	input := "testdata/input.txt"
	f, err := os.Open(input)
	if err != nil {
		t.Fatalf("failed to open %s: %v", input, err)
	}
	defer f.Close()
	os.Stdin = f
	output, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(output.Name())
	os.Stdout = output
	main()
	output.Close()
	got, err := ioutil.ReadFile(output.Name())
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	expectedOutput := "testdata/output.txt"
	want, err := ioutil.ReadFile(expectedOutput)
	if err != nil {
		t.Fatalf("failed to read expected output file: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("output\n%q\n != expected output\n%q\n", string(got), string(want))
	}
}
