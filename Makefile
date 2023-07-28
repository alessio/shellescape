#!/usr/bin/make -f

VERSION := $(shell git describe)

all: clean escargs

escargs:
	go build \
          -ldflags="-X 'main.version=$(VERSION)'" \
          ./cmd/escargs

clean:
	rm -rf escargs

.PHONY: clean
