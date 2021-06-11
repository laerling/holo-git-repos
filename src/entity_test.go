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
	k, v := parseEntityLine([]byte("key=value"))
	assertEq(t, k, "key")
	assertEq(t, v, "value")
}

// "url=a\npath=b\nrevision=c" is the only format currently accepted, therefore there's only one test for now
func TestEntityParseFile(t *testing.T) {

	// create temporary entity file
	testUrl := "TestEntityParseFile_testUrl"
	testPath := "TestEntityParseFile_testPath"
	testRevision := "TestEntityParseFile_testRevision"
	filePath := makeTemporaryEntityFile(t, os.TempDir(), testUrl, testPath, testRevision)

	// call function
	file, err := os.Open(filePath)
	assertErrNil(t, err, "Cannot re-open temporary file")
	url, path, revision := parseEntityFile(file)
	assertEq(t, url, testUrl)
	assertEq(t, path, testPath)
	assertEq(t, revision, testRevision)
}

func TestEntityParse(t *testing.T) {

	// create temporary entity file
	testUrl := "TestEntityParse_testUrl"
	testPath := "TestEntityParse_testPath"
	testRevision := "TestEntityParse_testRevision"
	filePath := makeTemporaryEntityFile(t, os.TempDir(), testUrl, testPath, testRevision)
	entityId := path.Base(filePath)

	// call function
	os.Setenv("HOLO_RESOURCE_DIR", path.Dir(filePath))
	url, path, revision := parseEntity(entityId)
	assertEq(t, url, testUrl)
	assertEq(t, path, testPath)
	assertEq(t, revision, testRevision)
}

func TestEntities(t *testing.T) {

	// create temporary directory with entity file
	tempDir, err := ioutil.TempDir(os.TempDir(), "")
	assertErrNil(t, err, "Cannot create temporary directory")
	testUrl := "TestEntities_testUrl"
	testPath := "TestEntities_testPath"
	testRevision := "TestEntities_testRevision"
	_ = makeTemporaryEntityFile(t, tempDir, testUrl, testPath, testRevision)

	// call function
	os.Setenv("HOLO_RESOURCE_DIR", tempDir)
	entities := parseEntities()
	assertEq(t, len(entities), 1)
	assertEq(t, entities[0].url, testUrl)
	assertEq(t, entities[0].path, testPath)
	assertEq(t, entities[0].revision, testRevision)
}
