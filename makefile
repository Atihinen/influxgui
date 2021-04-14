GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=influxgui
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

js-lint:
	./node_modules/.bin/eslint assets/main-ko.js

go-lint:
	${GOPATH}/bin/golint ./

lint: js-lint go-lint

test:
	${GOTEST}

clean:
	rm -rf ./build

init: clean
	mkdir -p ./build

compile-assets:
	go-bindata -o assets.go -prefix assets/ assets/ assets/media/

build-dev: compile-assets
	ENV="develop" mkdir -p ./build && ${GOBUILD} -o ./build/${BINARY_NAME}_linux -v *.go

build: test build-linux
	xgo -dest ./build/ -targets windows/*,darwin/* github.com/Atihinen/influxgui

build-staging: test build-linux
	xgo -branch `git branch` -dest ./build/ -targets windows/*,darwin/* github.com/Atihinen/influxgui

build-apple: init compile-assets
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 ${GOBUILD} -o ./build/${BINARY_NAME}_macOS -v *.go

build-linux: init compile-assets
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 ${GOBUILD} -o ./build/${BINARY_NAME}_linux -v *.go

build-windows: init compile-assets
	GOOS=windows GOARCH=386 CGO_ENABLED=1 CXX=i686-w64-mingw32-g++ CC=i686-w64-mingw32-gcc ${GOBUILD} -o ./build/${BINARY_NAME}_windosX86.exe -v *.go

build-windows-dev: init compile-assets
	ENV="develop" GOOS=windows GOARCH=386 CGO_ENABLED=1 CXX=i686-w64-mingw32-g++ CC=i686-w64-mingw32-gcc ${GOBUILD} -o ./build/${BINARY_NAME}_dev_windosX86.exe -v *.go

release: init compile-assets build-apple build-windows

run: build-dev
	./build/${BINARY_NAME}_linux

#from linux to win
#sudo apt-get install gcc-multilib
#sudo apt-get install gcc-mingw-w64
#sudo apt-get install gobjc++
#sudo apt install gobjc++-mingw-w64