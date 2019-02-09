package main

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

			.connected {
				color: #4CAF50;
			}

			.not_connected {
				color: #f44336;
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
					<span id="connection_status" title="Connection status" title="Connection status" class="not_connected">●</span>
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
