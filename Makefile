# creates static binaries
CC := CGO_ENABLED=0 go build -ldflags "-w -s" -trimpath -a -installsuffix cgo

SOURCES := $(shell find . -type f -name '*.go')
BIN := dns2mdns

.PHONY: all fmt docker clean

all: dns2mdns

deps: go.mod
	GOPROXY=direct go mod download
	GOPROXY=direct go get -u all
	go mod tidy

docker: Dockerfile
	docker build -t="lanrat/dns2mdns" .

$(BIN): $(SOURCES) go.mod go.sum
	$(CC) -o $@ 

check: | lint check1 check2

check1:
	golangci-lint run

check2:
	staticcheck -f stylish -checks all ./...

lint:
	golint ./...

clean:
	rm $(BIN)

fmt:
	gofmt -s -w -l .