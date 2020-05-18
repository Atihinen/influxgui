package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/webview/webview"
)

func storeConnection(w webview.WebView, host string) {
	err := addInfluxDBConfig(w, databaseHandler, host)
	if err != nil {
		createAlertDialog(w, "could not save", err.Error())
	}
	//createAlertDialog(w, "juu", "tesmos")
	connections, _ := getInfluxDBConfigs(w, databaseHandler)
	/*fmt.Println("connections: %v", connections)
	connectionsJSON, _ := json.Marshal(connections)
	data = fmt.Sprintf("window.connections=%s", connectionsJSON)
	w.Eval(data)*/
	popluateConnections(w, connections)
	//createAlertDialog(w, "Connections", "jee")
}

func deleteConnection(w webview.WebView, host string) {
	err := deleteInfluxDBConfig(w, databaseHandler, host)
	if err != nil {
		createAlertDialog(w, "Could not delete", err.Error())
	}
	connections, _ := getInfluxDBConfigs(w, databaseHandler)
	popluateConnections(w, connections)
}
func handleRPC(w webview.WebView, data string) {
	cmd := struct {
		Name string `json:"cmd"`
	}{}
	if err := json.Unmarshal([]byte(data), &cmd); err != nil {
		log.Println(err)
		return
	}
	switch cmd.Name {
	case "connections":
		storeConnection(w, "http://localhost:8086")
	case "addCon":
		connection := InfluxDBConnection{}
		if err := json.Unmarshal([]byte(data), &connection); err != nil {
			createAlertDialog(w, "Could not decrypt connection info", err.Error())
		} else if len(connection.Host) > 0 {
			storeConnection(w, connection.Host)
		} else {
			createAlertDialog(w, "Host not given", "Host needs to be http://<domain>:<port>")
		}
	case "delCon":
		connection := InfluxDBConnection{}
		if err := json.Unmarshal([]byte(data), &connection); err != nil {
			createAlertDialog(w, "Could not decrypt delete connection info", err.Error())
		} else {
			deleteConnection(w, connection.Host)
		}
	case "setHost":
		if err := json.Unmarshal([]byte(data), &connectionConfig); err != nil {
			createAlertDialog(w, "Could not set host", err.Error())
		}
	case "connectInflux":
		if err := json.Unmarshal([]byte(data), &connectionConfig); err != nil {
			createAlertDialog(w, "Could not connect influxdb", err.Error())
		}
		fmt.Printf("*****, %v, %v, %v", connectionConfig.Host, connectionConfig.Username, connectionConfig.Password)
		res, msg := pingInfluxDB()
		if !res {
			createAlertDialog(w, connectionConfig.Host, msg)
		}
		res, databases, data := showDatabases()
		if !res {
			createAlertDialog(w, data, "")
			return
		}
		updateConnectionStatus(w, res)
		populateDatabases(w, databases)
	case "sendQuery":
		query := InfluxDBQuery{}
		if err := json.Unmarshal([]byte(data), &query); err != nil {
			createAlertDialog(w, "Could not send query", err.Error())
		}
		fmt.Printf("*****, %v, %v", query.Query, query.Database)
		res, data := runInfluxDBQueryJSON(query.Query, query.Database)
		if !res {
			createAlertDialog(w, data, "")
			return
		}
		w.Eval(data)

	default:
		createAlertDialog(w, "Default", "nothing happened")
		/*case strings.HasPrefix(data, "connectInfluxDB:"):
		//	message := strings.TrimPrefix(data, "connectInfluxDB:")
		//	s := strings.Split(message, ";")
		//	port, err := strconv.Atoi(s[1])
		//	if err != nil {
		/		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, "Connect", "Could not convert port to integer!")
		/		return
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
			appendToLog(queryInfo.Query)
			updateHistoryLogs(w)
			res, data := runInfluxDBQuery(queryInfo.Query, queryInfo.Database)
			if !res {
				createAlertDialog(w, data, "")
				return
			}
			w.Eval(data)
		*/
	}
}

func updateHistoryLogs(w webview.WebView) {
	writeHistoryLogs(w, getLogOptions())
}
