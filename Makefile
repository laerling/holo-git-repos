holo-git-repos: $(wildcard *.go)
	go build

PHONY+=run
run: holo-git-repos
	./$<

PHONY+=scan
scan: holo-git-repos
	rm -f test/*~
	HOLO_RESOURCE_DIR=./test ./$< scan

PHONY+=clean
clean:
	rm -rf holo-git-repos temptest test/*~

PHONY: $(PHONY)
