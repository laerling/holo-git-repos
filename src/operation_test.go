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
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestInfo(t *testing.T) {
	expected := "MIN_API_VERSION=3\nMAX_API_VERSION=3\n"

	// since we're calling the main function we want to pass
	// arguments via os.Args, just as if we had executed the
	// binary.
	origArgs := os.Args
	os.Args = []string{"holo", "info"}
	holoOutput := getFunctionOutput(main)
	os.Args = origArgs
	assertEq(t, holoOutput, expected)
}

func TestScan(t *testing.T) {

	// create temporary directory with entity file
	tempDir, err := ioutil.TempDir(os.TempDir(), "")
	assertErrNil(t, err, "Cannot create temporary directory")
	testUrl := "TestScan_testUrl"
	testPath := "TestScan_testPath"
	testRevision := "TestScan_testRevision"
        entityFilePath1 := makeTemporaryEntityFile(t, tempDir, testUrl, testPath, testRevision)
        entityFileName1 := path.Base(entityFilePath1)
	//entityFileName2 := makeTemporaryEntityFile(t, tempDir, testUrl, testPath, testRevision)

	// call function and check output
	expected := "ENTITY: git-repo:" + entityFileName1
	expected += "\nSOURCE: " + entityFilePath1
	expected += "\nurl: " + testUrl
	expected += "\nrevision: " + testRevision
	expected += "\nclone into: " + testPath
	expected += "\n"
	os.Setenv("HOLO_RESOURCE_DIR", tempDir)
	scanOutput := getFunctionOutput(holoScan)
        assertEq(t, scanOutput, expected)
}

// TestApply does not test that the contents of the git repo are correct after cloning. It depends only on git's exit status
// for success.
func TestApply(t *testing.T) {

	// create temporary git directory for cloning
	tempGitDir, err := ioutil.TempDir(os.TempDir(), "")
	assertErrNil(t, err, "Cannot create temporary directory")
	// create arbitraty file in git repo
	_ = makeTemporaryEntityFile(t, tempGitDir, "", "", "")

	// create empty temporary directory for cloning into
	tempTargetDir, err := ioutil.TempDir(os.TempDir(), "")
	assertErrNil(t, err, "Cannot create temporary directory")

	// create temporary directory with entity file
	tempResourceDir, err := ioutil.TempDir(os.TempDir(), "")
	assertErrNil(t, err, "Cannot create temporary directory")
	entityFile := makeTemporaryEntityFile(t, tempResourceDir, tempGitDir, tempTargetDir, "tempRevision")
	entityId := path.Base(entityFile)

	// call holo and check output
	// TODO
	holoOutput := getHoloOutput(t, tempResourceDir, "apply", entityId)
	fmt.Println("TestApply: holo apply output: '" + string(holoOutput) + "'")
}

func TestApplyForce(t *testing.T) {
	t.Fatalf("unimplemented")
}

func TestDiff(t *testing.T) {
	t.Fatalf("unimplemented")
}
