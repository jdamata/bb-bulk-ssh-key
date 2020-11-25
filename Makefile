export GO111MODULE=on
VERSION=$(shell git describe --tags --candidates=1 --dirty)
BUILD_FLAGS=-ldflags="-X main.version=$(VERSION)"
SRC=$(shell find . -name '*.go')

.PHONY: all clean release install

all: linux darwin

clean:
	rm -f bb-bulk-ssh-key linux darwin

linux: $(SRC)
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o bb-bulk-ssh-key-linux .

darwin: $(SRC)
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o bb-bulk-ssh-key-darwin .

install:
	rm -f bb-bulk-ssh-key
	go build $(BUILD_FLAGS) .
	mv bb-bulk-ssh-key ~/bin/
