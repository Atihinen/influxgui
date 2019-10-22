GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=influxgui
BINARY_UNIX=$(BINARY_NAME)_unix

clear:
	rm -rf ./build

all: test build

test:
	${GOTEST}

build-dev:
	ENV="develop" mkdir -p ./build && ${GOBUILD} -o ./build/${BINARY_NAME}_linux -v *.go

build-linux:
	mkdir -p ./build && ${GOBUILD} -o ./build/${BINARY_NAME}_linux -v *.go

build: test clear build-linux
	xgo -dest ./build/ -targets windows/*,darwin/* github.com/Atihinen/influxgui

run: build-dev
	./build/${BINARY_NAME}_linux