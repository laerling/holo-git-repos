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

// fail writes the string msg to stderr and exits with a non-zero exit code
func fail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

// failOnErr calls fail with msg and the error message if err is non-nil. Otherwise it does nothing.
func failOnErr(err error, msg string) {
	if err != nil {
		fail(msg + "\nError: " + err.Error())
	}
}

// runGit builds and runs a git command.
// If printOutput is true, the output of the command is printed to stdout.
func runGit(printOutput bool, arguments ...string) {
	// git doesn't output anything when run via exec, so no
	// output redirection is needed
	cmd := exec.Command("git", arguments...)
	if printOutput {
		cmd.Stdout = os.Stdout
	}
	err := cmd.Run()
	failOnErr(err, fmt.Sprintln("Running git with arguments", arguments, "failed"))
}

// runGitInRepo builds and runs a git command in an existing repository.
// If printOutput is true, the output of the command is printed to stdout.
func runGitInRepo(printOutput bool, repoPath string, arguments ...string) {
	// git doesn't output anything when run via exec, so no
	// output redirection is needed
	arguments = append([]string{"-C", repoPath}, arguments...)
	cmd := exec.Command("git", arguments...)
	if printOutput {
		cmd.Stdout = os.Stdout
	}
	err := cmd.Run()
	failOnErr(err, fmt.Sprintln("Running git with arguments", arguments, "failed"))
}

// isGitRepo checks whether the given path is a git repository.
// A path counts a git repository if it is a directory containing a .git directory
func isGitRepo(path string) bool {
	_, err := os.Stat(path + "/.git")
	if err == nil {
		return true
	}
	// if we get this far, there is indeed a relevant error
	if !os.IsNotExist(err) {
		failOnErr(err, fmt.Sprintln("Cannot stat", path+"/.git/"))
	}
	return false
}

// holoScan executes the 'holo scan' operation. It scans $HOLO_RESOURCE_DIR for entities that can be provisioned.
func holoScan() {
	for _, entity := range parseEntities() {
		fmt.Println("ENTITY: git-repo:" + entity.fileName)
		fmt.Println("SOURCE: " + entity.filePath)
		fmt.Println("url: " + entity.url)
		fmt.Println("revision: " + entity.revision)
		fmt.Println("clone into: " + entity.path)
	}
}

// holoApply executes the 'holo apply' operation. It applies the entity with ID entityId.
// It clones the repository and, if revision is not emptystring, checks out that revision.
// If force is true, and the git repository path of the entity already exists,
// it is recursively deleted before being cloned again.
// If force is false, and the git repository path of the entity already exists,
// the message defined in holo-plugin-interface(7) is printed to FD 3.
func holoApply(entityId string, force bool) {

	url, path, revision := parseEntity(entityId)

	// check if directory already exists
	_, err := os.Stat(path)
	exists := !os.IsNotExist(err)
	if exists && err != nil {
		failOnErr(err, fmt.Sprintln("Cannot stat path", path))
	}

	// if it exists, force is needed
	if exists {
		// If it is NOT a git repository, warn the user
		if !isGitRepo(path) {
			fmt.Fprintln(os.Stderr, "WARNING:", path, "is not a git repository")
		}

		// delete directory if forced to
		if force {
			err := os.RemoveAll(path)
			failOnErr(err, "Cannot remove recursively: "+path)
		} else {
			_, err := os.NewFile(3, "holo").Write([]byte("requires --force to overwrite\n"))
			failOnErr(err, "Can't write to file descriptor 3")
			return
		}
	}

	// if it doesn't exist or we're forced, create it
	if !exists || force {
		// We need to do clone and checkout separately, because
		// revision can be a branch/tag name or a commit ID, so it
		// can't reliably be specified to git-clone
		runGit(false, "clone", url, path)
		if revision != "" {
			runGitInRepo(false, path, "checkout", revision)
		}
	}
}

// holoDiff executes the 'holo diff' operation.
// It generates a diff of the entity with ID entityId by calling `git diff`.
// The diff is between the worktree and the revision that was checked out at clone time.
func holoDiff(entityId string) {

	_, path, revision := parseEntity(entityId)
	repo, err := filepath.EvalSymlinks(path)
	failOnErr(err, "Possibly dead symlink in path: "+path)

	// The diff is between the worktree and the revision that was checked out at clone time.
	runGitInRepo(true, repo, "diff", revision+"..HEAD")
}

func main() {

	// check arguments
	if len(os.Args) < 2 {
                fail("holo-git-repos: Not enough arguments")
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
                        fail("holo-git-repos apply: Missing entity argument")
		}
		holoApply(os.Args[2], false)
		return

	case "force-apply":
		if len(os.Args) < 3 {
                        fail("holo-git-repos force-apply: Missing entity argument")
		}
		holoApply(os.Args[2], true)
		return

	case "diff":
		if len(os.Args) < 3 {
                        fail("holo-git-repos diff: Missing entity argument")
		}
		holoDiff(os.Args[2])
		return

	}
}
