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
	"bufio"
	"io"
	"os"
	"strings"
)

// entity represents a git repository entity for holo.
type entity struct {
	fileName string
	filePath string // use actual path object
	url      string
	path     string
	revision string
}

// parseEntityLine parses a line of format 'key=value'.
func parseEntityLine(line []byte) (key string, value string) {
	lineSplit := strings.Split(string(line), "=")

	if len(lineSplit) != 2 {
		fail("Wrong line format: '" + string(line) + "'")
	}

	key = strings.TrimSpace(lineSplit[0])
	value = strings.TrimSpace(lineSplit[1])
	return
}

// parseEntityFile parses a file into an entity instance.
func parseEntityFile(file io.Reader) (url string, path string, revision string) {
	var values [3]string
	fileReader := bufio.NewReader(file)
	for i, key := range [3]string{"url", "path", "revision"} {
		lineBytes, err := fileReader.ReadBytes('\n')
		if err != io.EOF {
			failOnErr(err, "Error reading entity file")
		}
		k, v := parseEntityLine(lineBytes)
		if k != key {
			fail("Erroneous key in entity file. Expected " + key + ", found " + k)
		}
		values[i] = v
	}
	return values[0], values[1], values[2]
}

// parseEntity parses the entity with id ID.
func parseEntity(id string) (string, string, string) {

	// find resource directory
	resDirName := os.Getenv("HOLO_RESOURCE_DIR")
	if resDirName == "" {
		fail("HOLO_RESOURCE_DIR empty")
	}

	// parse entity file
	entityFile, err := os.Open(resDirName + "/" + id)
	failOnErr(err, "Cannot open entity with ID "+id)
	url, path, revision := parseEntityFile(entityFile)

	return url, path, revision
}

// parseEntities parses all entities in holo resource directory.
func parseEntities() []entity {
	resDirName := os.Getenv("HOLO_RESOURCE_DIR")
	if resDirName == "" {
		fail("HOLO_RESOURCE_DIR empty")
	}

	// open directory
	resDir, err := os.Open(resDirName)
	failOnErr(err, "Cannot open HOLO_RESOURCE_DIR")

	// read files
	files, err := resDir.Readdir(0)
	failOnErr(err, "Cannot read files from HOLO_RESOURCE_DIR")

	// parse files
	entities := make([]entity, len(files))
	for i, file := range files {

		// open file
		fileName := file.Name()
		filePath := resDirName + "/" + fileName // TODO use path joining instead of string concatenation
		file, err := os.Open(filePath)
		failOnErr(err, "Cannot open file "+filePath)

		// read and parse file
		url, path, revision := parseEntityFile(file)
		entities[i] = entity{fileName, filePath, url, path, revision}
	}

	return entities
}
