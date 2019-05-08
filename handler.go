package main

import (
	"encoding/json"
	"fmt"
	"github.com/zserge/webview"
	"strconv"
	"strings"
)

func handleRPC(w webview.WebView, data string) {
	switch {
	case strings.HasPrefix(data, "connectInfluxDB:"):
		message := strings.TrimPrefix(data, "connectInfluxDB:")
		s := strings.Split(message, ";")
		port, err := strconv.Atoi(s[1])
		if err != nil {
			w.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, "Connect", "Could not convert port to integer!")
			return
		}
		host := fmt.Sprintf("%s:%d", s[0], port)
		connectionConfig = InfluxDBConnection{Host: host, Username: s[2], Password: s[3]}
		res, message := pingInfluxDB()
		if !res {
			createAlertDialog(w, message, "")
			return
		}
		res, data := showDatabases()
		if !res {
			createAlertDialog(w, data, "")
			return
		}
		w.Eval(data)
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
		res, data := runInfluxDBQuery(queryInfo.Query, queryInfo.Database)
		if !res {
			createAlertDialog(w, data, "")
			return
		}
		w.Eval(data)

	}
}
