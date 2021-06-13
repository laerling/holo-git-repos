PHONY+=build
build: holo-git-repos
holo-git-repos: format src/*.go
	go build -o "$@" $(WILDCARD src/*.go)

PHONY+=run
run: holo-git-repos
	./$<

PHONY+=test
test: format
	cd src; go test; cd -

PHONY+=format
format: src/*.go
	hash goimports && goimports -w $+
	hash gofmt && gofmt -w $+

PHONY+=clean
clean:
	rm -rf holo-git-repos temptest
	find src -name '*~' -delete

.PHONY: $(PHONY)
