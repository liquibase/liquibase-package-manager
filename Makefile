VERSION=`cat VERSION`

.PHONY: build darwin windows

release: updateVersion darwin windows

windows:
	GOOS=windows GOARCH=amd64 go build -o $(PWD)/bin/windows/lpm.exe $(PWD)/cmd/lpm/windows.go
	zip ./bin/windows/lpm-$(VERSION)-windows.zip ./bin/windows/lpm.exe

darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(PWD)/bin/darwin/lpm $(PWD)/cmd/lpm/darwin.go
	zip ./bin/darwin/lpm-$(VERSION)-darwin.zip ./bin/lpm

updateVersion:
	cp $(PWD)/VERSION $(PWD)/internal/app/VERSION

build: updateVersion
	go build -o $(PWD)/bin/lpm $(PWD)/cmd/lpm/darwin.go

generateExtensionPackages:
	go run cmd/populator/Populator.go >> temp_packages.json

test:
	golint ./internal/app/...
	go test -v ./internal/app/... -coverprofile=coverage.out -covermode count

cover: test
	go tool cover -html=coverage.out