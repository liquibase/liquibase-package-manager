VERSION=`cat VERSION`
VEXRUN_FILE := $(PWD)/utils/vexrun.jar
VEXRUN := java -jar $(VEXRUN_FILE)

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

test-setup:
	curl -Ls https://github.com/mcred/vexrun/releases/download/v0.0.5/vexrun-0.0.5.jar -z $(VEXRUN_FILE) -o $(VEXRUN_FILE)

e2e: test-version test-add test-completion test-help test-install test-list test-remove test-search test-cleanup

test-cleanup:
	-rm -Rf $(PWD)/liquibase_modules
	-rm $(PWD)/liquibase.json

test-version:
	$(VEXRUN) -f $(PWD)/tests/endtoend/version.yml -p VERSION $(VERSION)

test-add:
	$(VEXRUN) -f $(PWD)/tests/endtoend/add.yml

test-completion:
	$(VEXRUN) -f $(PWD)/tests/endtoend/completion.yml

test-help:
	$(VEXRUN) -f $(PWD)/tests/endtoend/help.yml

test-install:
	$(VEXRUN) -f $(PWD)/tests/endtoend/install.yml

test-list:
	$(VEXRUN) -f $(PWD)/tests/endtoend/list.yml

test-remove:
	$(VEXRUN) -f $(PWD)/tests/endtoend/remove.yml

test-search:
	$(VEXRUN) -f $(PWD)/tests/endtoend/search.yml

