/*******************************************************************************
*
* Copyright 2019 lærling <laerling@posteo.de>
*
* This program is free software: you can redistribute it and/or modify it under
* the terms of the GNU General Public License as published by the Free Software
* Foundation, either version 3 of the License, or (at your option) any later
* version.
*
* This program is distributed in the hope that it will be useful, but WITHOUT
* ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
* FOR A PARTICULAR PURPOSE. See the GNU General Public License for more
* details.
*
* You should have received a copy of the GNU General Public License along with
* this program. If not, see <http://www.gnu.org/licenses/>.
*
*******************************************************************************/

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
