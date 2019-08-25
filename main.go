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
	"os"
	"os/exec"
	"path/filepath"
)

/// fail writes the string msg to stderr and exits with a non-zero exit code.
func fail(msg string) {
	os.Stderr.Write([]byte(msg + "\n"))
	os.Exit(1)
}

/// holoScan executes the 'holo scan' operation. It scans $HOLO_RESOURCE_DIR for entities that can be provisioned.
func holoScan() {
	for _, entity := range(parseEntities()) {
		fmt.Println("ENTITY: git-repo:" + entity.fileName)
		fmt.Println("SOURCE: " + entity.filePath)
		fmt.Println("url: " + entity.url)
		fmt.Println("clone into: " + entity.path)
	}
}

/// holoApply executes the 'holo apply' operation. It applies the entity with ID entityId.
/// If force is true, and the git repository path of the entity already exists, it is recursively deleted before being cloned again.
func holoApply(entityId string, force bool) {

	url, path := parseEntity(entityId)

	// delete directory
	err := os.RemoveAll(path);
	if err != nil { fail("Cannot remove directory recursively: " + path) }

	// clone
	// git doesn't output anything when run via exec, so no output redirection is needed
	err = exec.Command("git", "clone", url, path).Run()
	if err != nil { fail("Git failed") }
}

/// holoDiff executes the 'holo diff' operation. It generates a diff of the entity with ID entityId by calling `git diff`.
func holoDiff(entityId string) {

	_, path := parseEntity(entityId)

	// git fetch
	cmd := exec.Command("git", "fetch")
	cmdDir, err := filepath.EvalSymlinks(path)
	if err != nil { fail("Possibly dead symlink in path: " + path) }
	cmd.Dir = cmdDir
	err = cmd.Run()
	if err != nil { fail("Git fetch failed: " + err.Error()) }

	// diff
	cmd = exec.Command("git", "diff", "HEAD", "origin/master")
	cmdDir, err = filepath.EvalSymlinks(path)
	if err != nil { fail("Possibly dead symlink in path: " + path) }
	cmd.Dir = cmdDir
	// TODO cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil { fail("Git diff failed") }
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
		if len(os.Args) < 3 {
			fail("Not enough arguments")
		}
		holoApply(os.Args[2], false)
		return

	case "force-apply":
		if len(os.Args) < 3 {
			fail("Not enough arguments")
		}
		holoApply(os.Args[2], true)
		return

	case "diff":
		if len(os.Args) < 3 {
			fail("Not enough arguments")
		}
		holoDiff(os.Args[2])
		return

	}
}
