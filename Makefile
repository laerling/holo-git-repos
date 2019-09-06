PHONY+=build
build: holo-git-repos
holo-git-repos: $(wildcard src/*.go)
	hash gofmt && gofmt -w $+
	go build -o "$@" $+

PHONY+=run
run: holo-git-repos
	./$<

PHONY+=test
test:
	cd src; go test; cd -

PHONY+=clean
clean:
	rm -rf holo-git-repos temptest
	find src -name '*~' -delete

.PHONY: $(PHONY)
