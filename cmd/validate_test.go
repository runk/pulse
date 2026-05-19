package cmd

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestValidateCmd(t *testing.T) {
	args := []string{"validate", "../example/policy-basic.json"}
	rootCmd.SetArgs(args)

	stdout := &bytes.Buffer{}
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(io.Discard)

	err := rootCmd.Execute()

	if err != nil {
		t.Fatalf("Error returned: %v", err)
	}

	if !strings.Contains(stdout.String(), "Policy is valid") {
		t.Fatalf("Stdout should indicate that policy is valid")
	}
}

func TestValidateInvalid(t *testing.T) {
	args := []string{"validate", "../example/policy-corrupted.json"}
	rootCmd.SetArgs(args)

	stdout := &bytes.Buffer{}
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(io.Discard)

	err := rootCmd.Execute()

	if err == nil {
		t.Fatal("Error expected, got nil")
	}

	if strings.Contains(err.Error(), "Policy validation failed: unexpected end of JSON input") {
		t.Fatalf("Unexpected error message returned: %v", err)
	}
}

func TestValidateNoFile(t *testing.T) {
	args := []string{"validate", "../example/doh"}
	rootCmd.SetArgs(args)

	stdout := &bytes.Buffer{}
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(io.Discard)

	err := rootCmd.Execute()

	if err == nil {
		t.Fatal("Error expected, got nil")
	}

	if !strings.Contains(err.Error(), "Cannot open policy file") {
		t.Fatalf("Unexpected error message returned: %v", err)
	}
}
