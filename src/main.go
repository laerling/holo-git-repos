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
func runGit(printOutput bool, arguments ...string) error {
	// git doesn't output anything when run via exec, so no
	// output redirection is needed
	cmd := exec.Command("git", arguments...)
	if printOutput {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runGitInDir builds and runs a git command in an existing repository.
// If printOutput is true, the output of the command is printed to stdout.
func runGitInDir(printOutput bool, repoPath string, arguments ...string) error {
	arguments = append([]string{"-C", repoPath}, arguments...)
	return runGit(printOutput, arguments...)
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

// clone clones the git repo from url to path, then checks out the
// given revision if it is not emptystring.
func clone(url string, path string, revision string) error {
	// We need to do clone and checkout separately, because
	// revision can be a branch/tag name or a commit ID, so it
	// can't reliably be specified to git-clone

	// clone
	err := runGit(false, "clone", url, path)
	if err != nil {
		return err
	}

	// checkout
	if revision != "" {
		return checkout(path, revision)
	}

	return nil
}

// checkout checks out the given revision in the git repository
// denoted by path.
func checkout(path string, revision string) error {
	return runGitInDir(false, path, "checkout", revision)
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
// If the target already exists, the behavior depends on a few things:
// - If force is false, return control to holo with the corresponding message
// - If the target is a git repo, try checking out the revision
// - If the target is not a git repo or the checkout failed (supposedly because
//   the revision does not exist), delete it before clone and checkout is done
func holoApply(entityId string, force bool) {

	url, path, revision := parseEntity(entityId)

	// check if directory already exists
	_, err := os.Stat(path)
	exists := !os.IsNotExist(err)

	// if the target already exists, the behavior depends on a few things
	if exists {
		// fail if we encountered an error
		failOnErr(err, fmt.Sprintln("Cannot stat path", path))

		// check if it's even a repo
		isRepo := isGitRepo(path)

		// if we're not forced, return with a warning
		if !force {
			// if it's not even a repo, the user might want to know what they're doing
			if !isRepo {
				fmt.Fprintln(os.Stderr, "WARNING:", path, "is not a git repository")
			}
			_, err := os.NewFile(3, "holo").Write([]byte("requires --force to overwrite\n"))
			failOnErr(err, "Can't write to file descriptor 3")
			return
		}

		// if it is a repo, let's try a simple checkout first
		err = nil
		if isRepo {
			err = checkout(path, revision)
		}

		// if it's not a repo or checkout failed, delete and reclone it
		if !isRepo || err != nil {
			err := os.RemoveAll(path)
			failOnErr(err, "Cannot remove recursively: "+path)
			exists = false
		}
	}

	// if the target does not yet exist, clone it
	// we cannot use an else branch, since exists might have been reassigned above
	if !exists {
		err = clone(url, path, revision)
		failOnErr(err, "Cannot clone repository "+url+" into "+path+" with revision "+revision)
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
	runGitInDir(true, repo, "diff", revision+"..HEAD")
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
