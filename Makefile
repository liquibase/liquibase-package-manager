VERSION=`cat $(PWD)/VERSION`
VEXRUN_FILE := $(PWD)/utils/vexrun.jar
VEXRUN := java -jar $(VEXRUN_FILE)

.PHONY: build darwin windows linux s390x

release: updateVersion darwin windows linux s390x

windows:
	GOOS=windows GOARCH=amd64 go build -o $(PWD)/bin/windows/lpm.exe $(PWD)/cmd/lpm/windows.go
	cd $(PWD)/bin/windows && zip lpm-$(VERSION)-windows.zip lpm.exe

darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(PWD)/bin/darwin/lpm $(PWD)/cmd/lpm/darwin.go
	cd $(PWD)/bin/darwin && zip lpm-$(VERSION)-darwin.zip lpm

linux:
	GOOS=linux GOARCH=amd64 GOARM=7 go build -o $(PWD)/bin/linux/lpm $(PWD)/cmd/lpm/darwin.go
	cd $(PWD)/bin/linux && zip lpm-$(VERSION)-linux.zip lpm

s390x:
	GOOS=linux GOARCH=s390x go build -o $(PWD)/bin/s390x/lpm $(PWD)/cmd/lpm/darwin.go
	cd $(PWD)/bin/s390x && zip lpm-$(VERSION)-s390x.zip lpm

updateVersion:
	cp $(PWD)/VERSION $(PWD)/internal/app/VERSION

build: updateVersion
	go build -o $(PWD)/bin/lpm $(PWD)/cmd/lpm/darwin.go

generateExtensionPackages:
	go run package-manager/cmd/populator

test:
	golint ./internal/app/...
	go test -v ./internal/app/... -coverprofile=coverage.out -covermode count

cover: test
	go tool cover -html=coverage.out

test-setup:
	curl -Ls https://github.com/mcred/vexrun/releases/download/v0.0.5/vexrun-0.0.5.jar -z $(VEXRUN_FILE) -o $(VEXRUN_FILE)

e2e: test-version test-add test-completion test-help test-install test-list test-remove test-search test-cleanup

test-cleanup:
	-rm -Rf $(PWD)/liquibase_libs
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

