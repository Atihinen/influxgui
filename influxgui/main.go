package main

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/zserge/webview"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

//"./guihtml"
const (
	windowWidth  = 600
	windowHeight = 480
)

type InfluxDBConnection struct {
	Host     string
	Username string
	Password string
}

var connectionConfig InfluxDBConnection

var indexHTML = `
<!doctype html>
<html>
	<head>
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<!--<script type="text/javascript" src="https://getfirebug.com/firebug-lite.js"></script>-->
		<script type="text/javascript">
			function connectForm(self){
				var host = document.getElementById('influxdb_host').value;
				var port = document.getElementById('influxdb_port').value;
				var username = document.getElementById('influxdb_username').value;
				var password = document.getElementById('influxdb_password').value;
				var msg = host+";"+port+";"+username+";"+password;
				external.invoke('connectInfluxDB:'+msg);
				return false;
			}
			function createQuery(){
				var query = document.getElementById('influxdb_query').value;
				var database = document.getElementById('inluxdb_db').value;
				var data = {
					query: query,
					database: database
				};
				external.invoke('createInfluxDBQuery:'+JSON.stringify(data));
				return false;
			}
		</script>
	</head>
	<body>
		<div id="app">
			<header>
				<div id="connect_container">
					<form onsubmit="return connectForm(this)">
						<input type="text" placeholder="Hostname" id="influxdb_host" name="influxdb_host" />
						<input type="text" placeholder="Port" id="influxdb_port" name="influxdb_port" />
						<input type="text" placeholder="Username" id="influxdb_username" name="influxdb_username" />
						<input type="password" placeholder="Password" id="influxdb_password" name="influxdb_password" />
						<input type="submit" value="Connect" />
					</form>
					<span id="connection_status" title="Connection status">Not connected</span>
				</div>
			</header>
			<div id="query_input_container">
				<form onsubmit="return createQuery()">
					<input type="text" placeholder="Query" id="influxdb_query" name="influxdb_query" />
					<select name="inluxdb_db" id="inluxdb_db">
						<option value="">Database</option>
					</select>
					<input type="submit" value="Send query" />
				</form>
			</div>
			<div id="query_content_container">
				<textarea id="query_content"></textarea>
			</div>
		</div>
		
	</body>
</html>
`

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

func pingInfluxDB(w webview.WebView) bool {
	success := true
	influxdbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     connectionConfig.Host,
		Username: connectionConfig.Username,
		Password: connectionConfig.Password,
	})
	if err != nil {
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, "Could not connect to influxdb", err.Error())
		success = false
	}
	_, _, err = influxdbClient.Ping(0)
	if err != nil {
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, "Could not ping to influxdb", err.Error())
		success = false
	}
	return success
}

func runInfluxDBQuery(w webview.WebView, query string, database string) {
	// Make client
	influxdbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     connectionConfig.Host,
		Username: connectionConfig.Username,
		Password: connectionConfig.Password,
	})
	if err != nil {
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, "Could not connect to influxdb", err.Error())
	}
	defer influxdbClient.Close()

	q := client.NewQuery(query, database, "")
	if response, err := influxdbClient.Query(q); err == nil && response.Error() == nil {
		log.Println(response.Results)
		results := "------------------\\n"
		columns := ""
		for _, serie := range response.Results[0].Series {
			log.Println(serie.Columns)
			for _, column := range serie.Columns {
				columns = fmt.Sprintf("%s%s\\t", columns, column)
			}
			results = fmt.Sprintf("%s%s\\n", results, columns)
			values := ""
			log.Println(serie.Values)
			for _, value := range serie.Values {
				log.Println(value)
				for _, val := range value {
					log.Println(val)
					values = fmt.Sprintf("%s%s\\t", values, val)
				}
				values = fmt.Sprintf("%s\\n", values)
			}
			results = fmt.Sprintf("%s%s\\n", results, values)
		}

		jsCmd := `document.getElementById('query_content').value = "` + results + `";`
		log.Println(jsCmd)
		w.Eval(jsCmd)
	}
}

func showDatabases(w webview.WebView) {
	// Make client
	influxdbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     connectionConfig.Host,
		Username: connectionConfig.Username,
		Password: connectionConfig.Password,
	})
	if err != nil {
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, "Could not connect to influxdb", err.Error())
	}
	defer influxdbClient.Close()

	q := client.NewQuery("SHOW DATABASES;", "", "ns")
	if response, err := influxdbClient.Query(q); err == nil && response.Error() == nil {
		log.Println(response.Results[0].Series[0].Values)
		dbs := "<option value=''>Database</option>"
		for _, value := range response.Results[0].Series[0].Values {
			res := value[0]
			log.Println(res)
			option := fmt.Sprintf("<option value='%s'>%s</option>", res, res)
			dbs = fmt.Sprintf("%s%s", dbs, option)
		}
		jsCmd := fmt.Sprintf("document.getElementById('inluxdb_db').innerHTML = \"%s\";", dbs)
		log.Println(jsCmd)
		w.Eval(jsCmd)

	}
}

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
		w.Eval("document.getElementById('connection_status').innerHTML = 'Connected';")
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

/*<button onclick="external.invoke('close')">Close</button>
<button onclick="external.invoke('fullscreen')">Fullscreen</button>
<button onclick="external.invoke('unfullscreen')">Unfullscreen</button>
<button onclick="external.invoke('open')">Open</button>
<button onclick="external.invoke('opendir')">Open directory</button>
<button onclick="external.invoke('save')">Save</button>
<button onclick="external.invoke('message')">Message</button>
<button onclick="external.invoke('info')">Info</button>
<button onclick="external.invoke('warning')">Warning</button>
<button onclick="external.invoke('error')">Error</button>
<button onclick="external.invoke('changeTitle:'+document.getElementById('new-title').value)">
	Change title
</button>
<input id="new-title" type="text" />
<button onclick="external.invoke('changeColor:'+document.getElementById('new-color').value)">
	Change color
</button>
<input id="new-color" value="#e91e63" type="color" />*/
