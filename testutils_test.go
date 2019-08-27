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
	"io/ioutil"
	"os"
	"testing"
)


/// makeTemporaryFile creates a valid temporary entity file.
/// It returns the whole path of the file (without expanded symlinks, though).
func makeTemporaryFile(t *testing.T, baseDir string) string {

	// create temporary entity file
	tempFile, err := ioutil.TempFile(baseDir, "")
	assertErrNil(t, err, "Cannot open temporary file")
	defer tempFile.Close()

	// write contents
	_, err = tempFile.WriteString("url=a\npath=b")
	assertErrNil(t, err, "Cannot write to temporary file")

	// get path
	tempFileInfo, err := tempFile.Stat()
	assertErrNil(t, err, "Cannot stat temporary file")
	return os.TempDir() + "/" + tempFileInfo.Name() // TODO use path joining instead of string concatenation
}

func assertErrNil(t *testing.T, err error, msg string) {
	if err != nil {
		t.Fatalf(msg)
	}
}

func assertEq(t *testing.T, value interface{}, expected interface{}) {
	if value != expected {
		t.Fatalf("expected %v, found %v", value, expected)
	}
}
