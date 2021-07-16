.PHONY: build

build:
	go build -o $(PWD)/bin/lpm $(PWD)/cmd/lpm/darwin.go