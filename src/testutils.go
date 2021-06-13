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
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

// makeTemporaryEntityFile creates a valid temporary entity file. It
// returns the whole path of the file (without expanded symlinks,
// though).
func makeTemporaryEntityFile(t *testing.T, baseDir string, url string, targetDir string, revision string) string {

	// create temporary entity file
	tempFile, err := ioutil.TempFile(baseDir, "")
	assertErrNil(t, err, "Cannot open temporary file")

	// write contents
	_, err = tempFile.WriteString("url=" + url + "\npath=" + targetDir + "\nrevision=" + revision)
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

// assertEq fails if value and expected are not equal.
func assertEq(t *testing.T, value interface{}, expected interface{}) {
	if value != expected {
		t.Fatalf("\nExpected: %v\nFound %v\n", expected, value)
	}
}

// getFunctionOutput calls a function and returns whatever that
// function printed to stdout as a string. If you want to pass
// arguments to that function, let the caller define a lambda
// expression that calls the function witht the desired arguments and
// pass that expression to getFunctionOutput.
// For example:
//
// func greet(person string) {
//         fmt.Println("Hello, " + person)
// }
// func TestGreet(t *testing.T) {
//         f := func(){greet("you!")}
//         assertEq(t, getFunctionOutput(f), "Hello, you!\n")
// }
func getFunctionOutput(f func()) string {
	r, w, err := os.Pipe()
	failOnErr(err, "Cannot syscall pipe")

	// call function with changed stdout
	oldStdout := os.Stdout
	os.Stdout = w
	f()
	os.Stdout = oldStdout
	w.Close()

	// read function output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
