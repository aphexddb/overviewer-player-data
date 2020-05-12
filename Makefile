ifeq ($(OS),Windows_NT)
BINARY=odp.exe
else
BINARY=odp
endif

.PHONY: all test clean build

all: clean test build

clean:
	rm ${BINARY} || true

test:
	go test ./...

build:
	go build -mod vendor -o ${BINARY} main.go rcon.go
