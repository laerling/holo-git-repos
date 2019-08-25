package main

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

// These tests test functions that exit with a non-zero exit code.
// Therefore, a hack is required. See https://stackoverflow.com/questions/26225513/how-to-test-os-exit-scenarios-in-go

func TestFail(t *testing.T) {

	// base case of (one-stepped) recursion
	if os.Getenv("HOLO_GIT_REPOS_FAIL") == "1" {
		fail("fail")
		return
	}

	// rerun test with HOLO_GIT_REPOS_FAIL set
	cmd := exec.Command(os.Args[0], "-test.run=TestFail")
	cmd.Env = append(os.Environ(), "HOLO_GIT_REPOS_FAIL=1")
	err := cmd.Run()

	// check exit code
	if err, ok := err.(*exec.ExitError); ok && !err.Success() {
		return // exit code was non-zero -> test succeeds
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestFailOnErr(t *testing.T) {

	// base case of (one-stepped) recursion
	if os.Getenv("HOLO_GIT_REPOS_FAIL") == "1" {
		failOnErr(errors.New("non-nil error"), "should fail")
		return
	}

	// invocation that shouldn't fail
	failOnErr(nil, "shouldn't fail")

	// rerun test with HOLO_GIT_REPOS_FAIL set
	cmd := exec.Command(os.Args[0], "-test.run=TestFailOnErr")
	cmd.Env = append(os.Environ(), "HOLO_GIT_REPOS_FAIL=1")
	err := cmd.Run()

	// check exit code
	if err, ok := err.(*exec.ExitError); ok && !err.Success() {
		return // exit code was non-zero -> test succeeds
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
