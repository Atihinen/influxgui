package main

import (
	"bytes"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tidwall/buntdb"
	"github.com/webview/webview"
)

const (
	windowWidth  = 1050
	windowHeight = 600
)

var connectionConfig InfluxDBConnection
var databaseHandler *buntdb.DB

func startServer() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer ln.Close()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if len(path) > 0 && path[0] == '/' {
				path = path[1:]
			}
			if path == "" {
				path = "index.html"
			}
			if bs, err := Asset(path); err != nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.Header().Add("Content-Type", mime.TypeByExtension(filepath.Ext(path)))
				io.Copy(w, bytes.NewBuffer(bs))
			}
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
	initalizeLog()

	url := startServer()
	w := webview.New(webview.Settings{
		Width:                  windowWidth,
		Height:                 windowHeight,
		Debug:                  true,
		Title:                  "InfluxGUI",
		Resizable:              true,
		URL:                    url,
		ExternalInvokeCallback: handleRPC,
	})
	w.SetColor(255, 255, 255, 255)
	defer w.Exit()
	db, _ := setupDB(w)
	databaseHandler = db
	w.Run()
}
