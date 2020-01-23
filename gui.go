package main

import (
	"fmt"
	"github.com/zserge/webview"
)

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

			function toggleConnectionStatus(){
				var connectionStatus = document.getElementById('connection_status');
				if(connectionStatus.className == "not_connected"){
					connectionStatus.classList.remove("not_connected");
					connectionStatus.classList.add("connected");
				}
				else {
					connectionStatus.classList.remove("connected");
					connectionStatus.classList.add("not_connected");	
				}
			}

			function toggleHistory(self){
				var historyContainer = document.getElementById('history_container');
				if(historyContainer.className == "hidden"){
					historyContainer.classList.remove("hidden");
					self.classList.remove("no_history");
					self.classList.add("history");
				}
				else {
					historyContainer.classList.add("hidden");
					self.classList.remove("history");
					self.classList.add("no_history");
				}
			}

			function setQuery(query) {
				document.getElementById('influxdb_query').value = query;
			}
		</script>
		<style>
			html, body {
				margin: 0;
				padding: 0;
				overflow: hidden;
			}
			header {
				width: 100vw;
				height: 3em;
				background: #ccc;
				padding-top: 0.5em;
				padding-left: 1em;
			}

			#connect_container form {
				width: 70%
				float: left;
			}

			#influxdb_host {
				width: 20%;
			}

			#connection_status {
				float: right;
				margin-right: 2em;
				position: relative;
				top: -1.7em;
			}

			#show_history {
				float: right;
				margin-right: 0.5em;
				cursor: pointer;
				font-size: 2em;
				position: relative;
				top: -1.2em;
			}

			#query_input_container {
				width: 100vw;
				height: 100%;
				padding-left: 1em;
			}

			#influxdb_query {
				width: 70%;
			}

			#query_content_container {
				width: 100vw;
				height: 100%;
			}

			#query_content {
				width: 96%;
				margin-left: 2%;
				height: 40em;
			}

			.connected, .history {
				color: #4CAF50;
			}

			.not_connected, .no_history {
				color: #f44336;
			}

			.hidden {
				display: none;
			}


			input[type=submit] {
				background-color: #008CBA;
				border: none;
				color: white;
				padding: 8px 25px;
				text-align: center;
				text-decoration: none;
				display: inline-block;
				font-size: 12px;
				font-weight: bold;
				-moz-border-radius: 4px;
				 -webkit-border-radius: 4px;
				 border-radius: 4px;
			}

			select, input[type=text], input[type=password] {
				border: 1px solid #80e5ff;
				-moz-border-radius: 4px;
				 -webkit-border-radius: 4px;
				 border-radius: 4px;
				 padding: 3px;
			}

			#history_container {
				width: 100vw;
				height: 100%;
				padding-left: 1em;
				padding-right: 1em;
			}

			#history_content_dropdown {
				width: 95%;
			}
		</style>
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
					<span id="connection_status" title="Connection status"class="not_connected">●</span>
					<span onClick='toggleHistory(this);' id="show_history" title="Show history" class="no_history">⏍</span>
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
			<div id="history_container" class="hidden">
				<select onchange="setQuery(this.options[this.selectedIndex].value)" id="history_content_dropdown">
					<option value="SHOW MEASUREMENTS;">SHOW MEASUREMENTS;</option>
					<option value="SELECT * FROM measurement LIMIT 1;">SELECT * FROM measurement LIMIT 1;</option>
					<option value="SHOW DATABASES;">SHOW DATABASES;</option>
				</select>
			</div>
			<div id="query_content_container">
				<textarea id="query_content"></textarea>
			</div>
		</div>
		
	</body>
</html>
`

func createAlertDialog(w webview.WebView, message string, errMessage string) {
	w.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, message, errMessage)
}

func writeHistoryLogs(w webview.WebView, content string) {
	data := fmt.Sprintf("document.getElementById('history_content_dropdown').innerHTML = \"%s\";", content)
	w.Eval(data)
}

func createInfluxDBQueryResponse(data string) (string){
	return `document.getElementById('query_content').value = "` + data + `";`
}