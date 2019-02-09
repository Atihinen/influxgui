GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=influxgui
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

test:
	echo "Testing..."

build:
	mkdir -p ./build && ${GOBUILD} -o ./build/${BINARY_NAME}_linux -v && xgo -dest ./build/ -targets windows/*,darwin/* github.com/Atihinen/influxgui