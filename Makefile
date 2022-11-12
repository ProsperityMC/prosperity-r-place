.PHONY: all build test
SHELL := bash
VERSION := $(shell git describe --tags --dirty --always)

ifdef CI_BUILD_NUMBER
	VERSION := "build.$(CI_BUILD_NUMBER)-$(VERSION)"
endif

BUILD_DATE = $(shell date '+%Y-%m-%d %H:%M:%S')
LD_FLAGS = -s -w -X 'main.BuildVersion=${VERSION}' -X 'main.BuildDate=${BUILD_DATE}'
CC = go
OUT = ./dist/prosperity-r-place
ENTRY = ./cmd/prosperity-r-place

all: build

build:
	mkdir -p dist/
	$(CC) build -o $(OUT) $(ENTRY)

run: build
	$(OUT)

test:
	$(CC) test ./...
