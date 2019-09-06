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
	"os/exec"
	"io/ioutil"
	"os"
	"path"
	"testing"
)


/// makeTemporaryEntityFile creates a valid temporary entity file. It
/// returns the whole path of the file (without expanded symlinks,
/// though).
func makeTemporaryEntityFile(t *testing.T, baseDir string, url string, targetDir string) string {

	// create temporary entity file
	tempFile, err := ioutil.TempFile(baseDir, "")
	assertErrNil(t, err, "Cannot open temporary file")

	// write contents
	_, err = tempFile.WriteString("url=" + url + "\npath=" + targetDir)
	assertErrNil(t, err, "Cannot write to temporary file")

	// get file name
	tempFileInfo, err := tempFile.Stat()
	assertErrNil(t, err, "Cannot stat temporary file")
	tempFilePath := path.Join(baseDir, tempFileInfo.Name())
	return tempFilePath
}

func assertErrNil(t *testing.T, err error, msg string) {
	if err != nil {
		t.Fatalf(msg)
	}
}

/// assertEq fails if value and expected are not equal.
func assertEq(t *testing.T, value interface{}, expected interface{}) {
	if value != expected {
		t.Fatalf("\nExpected: %v\nFound %v\n", expected, value)
	}
}

/// getHoloOutput calls the main function with HOLO_RESOURCE_DIR set
/// and args as arguments and returns its stdout output as a byte
/// slice. args must not contain the program's name, it is inserted by
/// getHoloOutput. This approach to calling holo from within tests is
/// used, because during test execution the binary does not yet
/// exist. If an entity file must exist before holo is called, is has
/// to be created before invoking getHoloOutput. Use
/// makeTemporaryEntityFile.
func getHoloOutput(t *testing.T, holoResourceDir string, args ...string) []byte {
	fmt.Println("getHoloOutput: Calling " + os.Args[0])

	// redirect stdout to temporary file
	tempFile, err := ioutil.TempFile("", "")
	assertErrNil(t, err, "Cannot open temporary file")
	oldStdout := os.Stdout // save old stdout fd
	stdoutFile := os.NewFile(tempFile.Fd(), "/dev/stdout")
	os.Stdout = stdoutFile

	// call holo
	os.Args = []string{"holo-git-repos"}
	for _, arg := range(args) {
		os.Args = append(os.Args, arg)
	}
	os.Setenv("HOLO_RESOURCE_DIR", holoResourceDir)
	// FIXME: What if main calls exit?
	cmd := exec.Command(os.Args[0], "apply")
	cmd.Env = append(os.Environ(), "HOLO_RESOURCE_DIR=" + holoResourceDir)
	cmd.Run()
	//main()

	// end redirection and return output
	os.Stdout = oldStdout
	mainOutput, err := ioutil.ReadFile(tempFile.Name())
	assertErrNil(t, err, "Cannot read temporary file")
	return mainOutput
}
