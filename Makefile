# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: fmt build

build:
	$(GOBUILD) -o bin/http -v  ./app/main.go
	$(GOBUILD) -o bin/rpc -v  ./rpc/main.go
	$(GOBUILD) -o bin/cmd -v  ./cmd/main.go
	
fmt:
	gofmt -l -s -w app
	gofmt -l -s -w rpc
	gofmt -l -s -w cmd
	#gofmt -l -s -w library

lint:
	#golangci-lint run ./app/... ./logic/... ./service/... ./db/... ./library/...
	#golangci-lint run ./cmd/... ./rpc/...
	golangci-lint run
	go vet ./app/api/... ./app/service/... ./app/model/... ./app/query/... ./app/utils/... ./app/tools/...
	go vet ./cmd/... ./rpc/... ./library/...

clean:
	$(GOCLEAN)

test:
	$(GOTEST) -v -count=1 ./...

proto:
	protoc --proto_path=proto --go_out=${GOPATH}/src proto/*.proto
