holo-git-repos: $(wildcard *.go)
	go build

PHONY+=run
run: holo-git-repos
	./$<

PHONY+=scan
scan: holo-git-repos
	rm -f test/*~
	HOLO_RESOURCE_DIR=./test ./holo-git-repos scan

PHONY+=apply
apply: holo-git-repos
	rm -f test/*~
	rm -rf temptest
	HOLO_RESOURCE_DIR=./test ./holo-git-repos apply foo

PHONY+=diff
diff: apply holo-git-repos
	rm -f test/*~
	HOLO_RESOURCE_DIR=./test ./holo-git-repos diff foo

PHONY+=clean
clean:
	rm -rf holo-git-repos temptest test/*~

PHONY: $(PHONY)
