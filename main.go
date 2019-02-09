package main

import (
	"github.com/zserge/webview"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	windowWidth  = 1000
	windowHeight = 600
)

var connectionConfig InfluxDBConnection

func startServer() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer ln.Close()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(indexHTML))
		})
		log.Fatal(http.Serve(ln, nil))
	}()
	return "http://" + ln.Addr().String()
}

func getDebugFlag() bool {
	returnFlag := false
	envFlag, status := os.LookupEnv("ENV")
	if !status {
		returnFlag = false
	} else if envFlag == "develop" {
		returnFlag = true
	} else {
		returnFlag = false
	}
	return returnFlag
}

func main() {
	url := startServer()
	w := webview.New(webview.Settings{
		Width:     windowWidth,
		Height:    windowHeight,
		Debug:     true,
		Title:     "Simple window demo",
		Resizable: true,
		URL:       url,
		ExternalInvokeCallback: handleRPC,
	})
	w.SetColor(255, 255, 255, 255)
	defer w.Exit()
	w.Run()
}
