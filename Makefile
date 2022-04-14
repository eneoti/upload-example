GO ?= go

PACKAGE := upload-example
GITHASH := `git rev-parse HEAD`
GITTAG := `git describe --tags --always`
# LDFLAGS="-X upload-example/server.gitCommit=$(GITHASH) -X upload-example/server.gitTag=$(GITTAG)"

.PHONY: all test clean

lint:
	$(exit $(go fmt ./... | wc -l))
	go vet ./...

benchmark:
	go test -bench=. ./... -benchtime=100x -run=^# 

build:
	CGO_ENABLED=1 go build  -o $(PACKAGE) ./ 

run:
	go run ./

test:
	go test -v ./...