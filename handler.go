package main

import (
	"encoding/json"
	"fmt"
	"github.com/zserge/webview"
	"log"
	"strconv"
	"strings"
)

func handleRPC(w webview.WebView, data string) {
	switch {
	case data == "close":
		w.Terminate()
	case data == "fullscreen":
		w.SetFullscreen(true)
	case data == "unfullscreen":
		w.SetFullscreen(false)
	case data == "open":
		log.Println("open", w.Dialog(webview.DialogTypeOpen, 0, "Open file", ""))
	case data == "opendir":
		log.Println("open", w.Dialog(webview.DialogTypeOpen, webview.DialogFlagDirectory, "Open directory", ""))
	case data == "save":
		log.Println("save", w.Dialog(webview.DialogTypeSave, 0, "Save file", ""))
	case data == "message":
		w.Dialog(webview.DialogTypeAlert, 0, "Hello", "Hello, world!")
	case data == "info":
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagInfo, "Hello", "Hello, info!")
	case data == "warning":
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagWarning, "Hello", "Hello, warning!")
	case data == "error":
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, "Hello", "Hello, error!")
	case strings.HasPrefix(data, "connectInfluxDB:"):
		message := strings.TrimPrefix(data, "connectInfluxDB:")
		s := strings.Split(message, ";")
		port, err := strconv.Atoi(s[1])
		if err != nil {
			w.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, "Connect", "Could not convert port to integer!")
			return
		}
		host := fmt.Sprintf("%s:%d", s[0], port)
		log.Println(host)
		connectionConfig = InfluxDBConnection{Host: host, Username: s[2], Password: s[3]}
		log.Println(connectionConfig)
		res := pingInfluxDB(w)
		if !res {
			return
		}
		showDatabases(w)
		w.Eval("window.toggleConnectionStatus();")
	case strings.HasPrefix(data, "createInfluxDBQuery"):
		query := strings.TrimPrefix(data, "createInfluxDBQuery:")
		queryInfo := struct {
			Query    string `json:"query"`
			Database string `json:"database"`
		}{}
		if err := json.Unmarshal([]byte(query), &queryInfo); err != nil {
			w.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, "Could not parse results!", err.Error())
			return
		}
		log.Println(queryInfo)
		log.Println(query)
		runInfluxDBQuery(w, queryInfo.Query, queryInfo.Database)
	case strings.HasPrefix(data, "changeTitle:"):
		w.SetTitle(strings.TrimPrefix(data, "changeTitle:"))
	case strings.HasPrefix(data, "changeColor:"):
		hex := strings.TrimPrefix(strings.TrimPrefix(data, "changeColor:"), "#")
		num := len(hex) / 2
		if !(num == 3 || num == 4) {
			log.Println("Color must be RRGGBB or RRGGBBAA")
			return
		}
		i, err := strconv.ParseUint(hex, 16, 64)
		if err != nil {
			log.Println(err)
			return
		}
		if num == 3 {
			r := uint8((i >> 16) & 0xFF)
			g := uint8((i >> 8) & 0xFF)
			b := uint8(i & 0xFF)
			w.SetColor(r, g, b, 255)
			return
		}
		if num == 4 {
			r := uint8((i >> 24) & 0xFF)
			g := uint8((i >> 16) & 0xFF)
			b := uint8((i >> 8) & 0xFF)
			a := uint8(i & 0xFF)
			w.SetColor(r, g, b, a)
			return
		}
	}
}
