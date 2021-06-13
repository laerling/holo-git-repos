/*******************************************************************************
*
* Copyright 2021 laerling <laerling@posteo.de>
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
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestInfo(t *testing.T) {
	// since we're calling the main function we want to pass
	// arguments via os.Args, just as if we had executed the
	// binary.
	origArgs := os.Args
	os.Args = []string{"holo", "info"}
	mainOutput := getFunctionOutput(main)
	os.Args = origArgs

	// check output
	expected := "MIN_API_VERSION=3\nMAX_API_VERSION=3\n"
	assertEq(t, mainOutput, expected)
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

// holoApply: Target does not exist
// => clone, checkout
// basically like first-time provisioning
func TestApplyNotexistentTarget(t *testing.T) {
	t.Fatalf("unimplemented")
}

// holoApply: Target exists and not forced
// => "needs force"
func TestApplyExistNoForce(t *testing.T) {
	t.Fatalf("unimplemented")
}

// holoApply: Target exists and forced and target is repo
// => checkout
func TestApplyForceRepo(t *testing.T) {
	t.Fatalf("unimplemented")
}

// holoApply: Target exists and forced and target is repo and revision non-existent
// => checkout fails, delete, clone, checkout
func TestApplyForceRepoNonexistentRevision(t *testing.T) {
	t.Fatalf("unimplemented")
}

// holoApply: Target exists and forced and target is no repo
// => delete, clone, checkout
func TestApplyForceNoRepo(t *testing.T) {
	t.Fatalf("unimplemented")
}

// TODO: Remove this general test in favor of the specific tests above
func TestApply(t *testing.T) {

	// create git repo with content
	tempGitDir, err := ioutil.TempDir(os.TempDir(), "")
	assertErrNil(t, err, "Cannot create temporary directory for git repo")
	t.Log("tempGitDir (where to clone from):", tempGitDir)
	runGitInDir(false, tempGitDir, "init", "-b", "main")
	_ = makeTemporaryEntityFile(t, tempGitDir, "", "", "")
	runGitInDir(false, tempGitDir, "add", "-A")
	runGitInDir(false, tempGitDir, "commit", "-m", "TestApply")

	// create empty temporary directory for cloning into
	tempTargetDir, err := ioutil.TempDir(os.TempDir(), "")
	assertErrNil(t, err, "Cannot create temporary directory for cloning into")
	tempTargetDir += "/repo" // name of the cloned repo
	t.Log("tempTargetDir (where to clone to):", tempTargetDir)

	// create temporary directory with entity file
	tempResourceDir, err := ioutil.TempDir(os.TempDir(), "")
	assertErrNil(t, err, "Cannot create temporary directory for entity file")
	t.Log("tempResourceDir (where the entity file lies):", tempResourceDir)
	entityFile := makeTemporaryEntityFile(t, tempResourceDir, tempGitDir, tempTargetDir, "")
	entityId := path.Base(entityFile)
	t.Log("entityFile / entityId:", entityFile, "/", entityId)

	// call main function as if binary had been called
	// Positive case - expecting clone to succeed
	os.Setenv("HOLO_RESOURCE_DIR", tempResourceDir)
	origArgs := os.Args
	os.Args = []string{"holo", "apply", entityId}
	mainOutput := getFunctionOutput(main)
	os.Args = origArgs

	// check output
	assertEq(t, mainOutput, "")
}

func TestDiff(t *testing.T) {
	t.Fatalf("unimplemented")
}
