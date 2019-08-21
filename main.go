/*******************************************************************************
*
* Copyright 2019 laerling <laerling@posteo.de>
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
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)


/// A git repository entity for holo
type entity struct {
	fileName string
	filePath string // use actual path object
	url string
	path string
}

/// Write msg to stderr and exit with a non-zero exit code
func fail(msg string) {
	os.Stderr.Write([]byte(msg + "\n"))
	os.Exit(1)
}

/// parse a line of format key=value
func parseEntityLine(line []byte) [2]string {
	lineSplit := strings.Split(string(line), "=")

	if len(lineSplit) != 2 {
		fail("Wrong line format")
	}

	lineSplit[0] = strings.TrimSpace(lineSplit[0])
	lineSplit[1] = strings.TrimSpace(lineSplit[1])

	return [2]string{lineSplit[0], lineSplit[1]}
}

/// parse a file into an entity instance
func parseEntityFile(file io.Reader) (string, string) {
	fileReader := bufio.NewReader(file)

	// read url
	errMsg := "Error reading entity file"
	urlBytes, err := fileReader.ReadBytes('\n')
	if err != nil && err != io.EOF { fail(errMsg) }
	pathBytes, err := fileReader.ReadBytes('\n')
	if err != nil && err != io.EOF { fail(errMsg) }

	// split and clean
	url := parseEntityLine(urlBytes)
	if url[0] != "url" { fail("Erroneous key in entity file") }
	path := parseEntityLine(pathBytes)
	if path[0] != "path" { fail("Erroneous key in entity file") }

	return url[1], path[1]
}

/// parse all entities in holo resource directory
func parseEntities() []entity {
	resDirName := os.Getenv("HOLO_RESOURCE_DIR")
	if resDirName == "" {
		fail("HOLO_RESOURCE_DIR empty")
	}

	// open directory
	resDir, err := os.Open(resDirName)
	if err != nil {
		fail("Cannot open HOLO_RESOURCE_DIR")
	}

	// read files
	files, err := resDir.Readdir(0)
	if err != nil {
		fail("Cannot read files from HOLO_RESOURCE_DIR")
	}

	// parse files
	entities := make([]entity, len(files))
	for i, file := range(files) {

		// open file
		fileName := file.Name()
		filePath := resDirName + "/" + fileName
		// TODO use path joining instead of string concatenation
		file, err := os.Open(filePath)
		if err != nil {
			fail("Cannot open file " + filePath)
		}

		// read and parse file
		url, path := parseEntityFile(file)
		entities[i] = entity { fileName, filePath, url, path }
	}

	return entities
}

/// The 'scan' operation.
/// Scan $HOLO_RESOURCE_DIR for entities that can be provisioned
func holoScan() {
	for _, entity := range(parseEntities()) {
		fmt.Println("ENTITY: git-repo:" + entity.fileName)
		fmt.Println("SOURCE: " + entity.filePath)
		fmt.Println("url: " + entity.url)
		fmt.Println("clone into: " + entity.path)
	}
}

/// The 'apply' operation.
/// Apply entity with ID entityId
func holoApply(entityId string) {
	resDirName := os.Getenv("HOLO_RESOURCE_DIR")
	if resDirName == "" {
		fail("HOLO_RESOURCE_DIR empty")
	}

	// parse entity file
	entityFile, err := os.Open(resDirName + "/" + entityId)
	if err != nil { fail("Cannot open entity with ID " + entityId) }
	url, path := parseEntityFile(entityFile)

	// apply

	// delete directory (TODO: Only do this when --force'd)
	err = os.RemoveAll(path);
	if err != nil { fail("Cannot remove directory recursively: " + path) }

	// clone
	err = exec.Command("git", "clone", url, path).Run()
	if err != nil { fail("Git failed") }
}

func main() {

	// check arguments
	if len(os.Args) < 2 {
		fail("Not enough arguments")
	}

	// actions
	switch os.Args[1] {

	case "info":
		fmt.Println("MIN_API_VERSION=3")
		fmt.Println("MAX_API_VERSION=3")
		return

	case "scan":
		holoScan()
		return

	case "apply":
	case"force-apply":
		if len(os.Args) < 2 {
			fail("Not enough arguments")
		}
		holoApply(os.Args[2])
		return

	case "diff":
		// TODO (entity ID in os.Args[2])

	}
}
