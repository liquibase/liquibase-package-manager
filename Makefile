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

genrateExentionPackages:
	go run cmd/populator/Populator.go liquibase-cache >> temp_packages.json