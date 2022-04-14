GO ?= go

PACKAGE := upload-example
GITHASH := `git rev-parse HEAD`
GITTAG := `git describe --tags --always`
# LDFLAGS="-X upload-example/server.gitCommit=$(GITHASH) -X upload-example/server.gitTag=$(GITTAG)"

.PHONY: build, test, lint-pkgs

# run supported packages
# lint-pkgs:
# 	GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell

lint:
	$(exit $(go fmt ./... | wc -l))
	go vet ./...

test:
	go test ./...

benchmark:
	go test -bench=. ./...

build:
	CGO_ENABLED=1 go build  -o $(PACKAGE) ./

run:
	go run ./
