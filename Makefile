APP_NAME=lakitu
TOOLS = $(shell go list -f '{{range .Imports}}{{.}} {{end}}' tools.go)
VERSION = $(shell cat VERSION)

all: clean mac

build:
	env GOOS=linux GOARCH=amd64 go build -o bin/linux_amd64/$(APP_NAME)

mac:
	env GOOS=darwin GOARCH=amd64 go build -o bin/darwin_amd64/$(APP_NAME)

test:
	go test -v ./...

clean:
	go clean ./...
	rm -rf bin/*

format:
	@gofmt -s -l .

tools:
	go install $(TOOLS)

docker:
	docker build . -t $(APP_NAME):$(VERSION)
