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
	"path"
	"testing"
)

func TestEntityParseLine(t *testing.T) {

	// url
	result := parseEntityLine([]byte("url=ugh..."))
	assertEq(t, result[0], "url")
	assertEq(t, result[1], "ugh...")

	// path
	result = parseEntityLine([]byte("path=whatever"))
	assertEq(t, result[0], "path")
	assertEq(t, result[1], "whatever")
}

// "url=a\npath=b" is the only format currently accepted, therefore there's only one test for now
func TestEntityParseFile(t *testing.T) {

	// create temporary entity file
	testUrl := "testUrl"
	testPath := "testPath"
	filePath := makeTemporaryEntityFile(t, os.TempDir(), testUrl, testPath)

	// call function
	file, err := os.Open(filePath)
	assertErrNil(t, err, "Cannot re-open temporary file")
	url, path := parseEntityFile(file)
	assertEq(t, url, testUrl)
	assertEq(t, path, testPath)
}

func TestEntityParse(t *testing.T) {

	// create temporary entity file
	testUrl := "testUrl"
	testPath := "testPath"
	filePath := makeTemporaryEntityFile(t, os.TempDir(), testUrl, testPath)
	entityId := path.Base(filePath)

	// call function
	os.Setenv("HOLO_RESOURCE_DIR", path.Dir(filePath))
	url, path := parseEntity(entityId)
	assertEq(t, url, testUrl)
	assertEq(t, path, testPath)
}

func TestEntities(t *testing.T) {

	// create temporary directory with entity file
	tempDir, err := ioutil.TempDir(os.TempDir(), "")
	assertErrNil(t, err, "Cannot create temporary directory")
	testUrl := "testUrl"
	testPath := "testPath"
	_ = makeTemporaryEntityFile(t, tempDir, testUrl, testPath)

	// call function
	os.Setenv("HOLO_RESOURCE_DIR", tempDir)
	entities := parseEntities()
	assertEq(t, len(entities), 1)
	assertEq(t, entities[0].url, testUrl)
	assertEq(t, entities[0].path, testPath)
}
