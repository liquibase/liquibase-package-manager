VERSION=`cat $(PWD)/VERSION`
VEXRUN_FILE := $(PWD)/utils/vexrun.jar
VEXRUN := java -jar $(VEXRUN_FILE)

.PHONY: build darwin_amd64 darwin_arm64 windows linux_amd64 linux_arm64 s390x

release: updateVersion darwin_amd64 darwin_arm64 windows linux_amd64 linux_arm64 s390x

windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(PWD)/bin/windows/lpm.exe $(PWD)/cmd/lpm/windows.go
	cd $(PWD)/bin/windows && zip lpm-$(VERSION)-windows.zip lpm.exe

darwin_amd64:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(PWD)/bin/darwin_amd64/lpm $(PWD)/cmd/lpm/darwin.go
	cd $(PWD)/bin/darwin_amd64 && zip lpm-$(VERSION)-darwin.zip lpm

darwin_arm64:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(PWD)/bin/darwin_arm64/lpm $(PWD)/cmd/lpm/darwin.go
	cd $(PWD)/bin/darwin_arm64 && zip lpm-$(VERSION)-darwin-arm64.zip lpm

linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOARM=7 go build -ldflags="-s -w" -o $(PWD)/bin/linux_amd64/lpm $(PWD)/cmd/lpm/darwin.go
	cd $(PWD)/bin/linux_amd64 && zip lpm-$(VERSION)-linux.zip lpm

linux_arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOARM=7 go build -ldflags="-s -w" -o $(PWD)/bin/linux_arm64/lpm $(PWD)/cmd/lpm/darwin.go
	cd $(PWD)/bin/linux_arm64 && zip lpm-$(VERSION)-linux-arm64.zip lpm

s390x:
	GOOS=linux GOARCH=s390x CGO_ENABLED=0 go build -ldflags="-s -w" -o $(PWD)/bin/s390x/lpm $(PWD)/cmd/lpm/darwin.go
	cd $(PWD)/bin/s390x && zip lpm-$(VERSION)-s390x.zip lpm

updateVersion:
	cp $(PWD)/VERSION $(PWD)/internal/app/VERSION

build: updateVersion
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(PWD)/bin/lpm $(PWD)/cmd/lpm/darwin.go

generateExtensionPackages:
	go run package-manager/cmd/populator

test: test-setup
	staticcheck ./internal/app/...
	go vet ./internal/app/...
	go test -v ./internal/app/... -coverprofile=coverage.out -covermode count

cover: test
	go tool cover -html=coverage.out

test-setup:
	go install honnef.co/go/tools/cmd/staticcheck@latest
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

