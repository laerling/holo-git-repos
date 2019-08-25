main: $(wildcard *.go)
	go build main.go

PHONY+=run
run: main
	./$<

PHONY+=scan
scan: main
	rm -f test/*~
	HOLO_RESOURCE_DIR=./test ./$< scan

PHONY+=clean
clean:
	rm -rf temptest test/*~

PHONY: $(PHONY)
