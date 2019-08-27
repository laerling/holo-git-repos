holo-git-repos: $(wildcard *.go)
	go build

PHONY+=run
run: holo-git-repos
	./$<

PHONY+=clean
clean:
	rm -rf holo-git-repos temptest test/*~

PHONY: $(PHONY)
