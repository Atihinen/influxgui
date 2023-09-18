import './style.css';
import './app.css';

import logo from './assets/images/logo.png';
import {Greet} from '../wailsjs/go/main/App';
import {GetConnections} from '../wailsjs/go/main/App';
import {StoreConnections} from "../wailsjs/go/main/App";
import {DeleteConnection} from "../wailsjs/go/main/App";
import {Connect} from "../wailsjs/go/main/App";
import {GetDatabases} from "../wailsjs/go/main/App";
import {RunQuery} from "../wailsjs/go/main/App";
document.querySelector('#app').innerHTML = `
    <header id="header">
    <div id="tools">
        <div id="logo"><img id="logo_img" alt="InfluxGUI" /> InfluxGUI</div>
        <ul class="ul inline right">
            <li>Connections
                <select id="iconnections">
                    <option disabled value="">Please select one</option>
                    <option value="new">Modify connections</option>
                </select>
            </li>
            <li><button class="secondary-input" id="history">Show history</button></li>
            <li>Help</li></ul>
        </div>
    </header>
    <div id="alert_container" class="modal show">
        <!-- Modal content -->
        <div class="modal-content">
            <div id="alert_header">
                <div id="alert_topic">Alert</div>
                <span class="close" onclick="window.toggleAlertDialog(false)">&times;</span>
            </div>
            <div id="alert_message">
                Some message
            </div>
        </div>
    </div>
    <div id="new-database" class="modal show">
        <!-- Modal content -->
        <div class="modal-content">
            <span class="close" onclick="window.toggleConnectionsDialog(false)">&times;</span>
            <div>
                <div id="manage_connections">
                    <ul id="existing_connections">
                    </ul>
                </div>
                <form id="store_connection">
                    <input id="new_connetion" type="text" placeholder="http://domain:port" /> <input type="submit" value="Create new connection" />
                </form>
            </div>
        </div>
    </div>
    <div class="container mx-auto p-4">
        <div id="connection_details">
            <div id="connect_container">
                <form id="influxdb_connection_form">
                    <input type="text" placeholder="Username" id="influxdb_username" name="influxdb_username" />
                    <input type="password" placeholder="Password" id="influxdb_password" name="influxdb_password" />
                    <input type="submit" value="Connect" />
                </form>
                <span id="connection_status">●</span>
            </div>
            <div id="query_input_container">
                <form  id="send_query_form">
                    <input type="text" placeholder="Query" id="influxdb_query" name="influxdb_query" />
                    <select id="database_selection">
                    </select>
                    <input type="submit" value="Send query" />
                </form>
            </div>
        </div>
        <div id="query_content_container">
            <textarea class="hidden" id="query_content"></textarea>
            <div id="query_content_table">

            </div>
        </div>
    </div>
    <!--<img id="logo" class="logo">
      <div class="result" id="result">Please enter your name below 👇</div>
      <div class="input-box" id="input">
        <input class="input" id="name" type="text" autocomplete="off" />
        <button class="btn" onclick="greet()">Greet</button>
      </div>
    </div>-->
`;
document.getElementById('logo_img').src = logo;

//let nameElement = document.getElementById("name");
let resultElement = document.getElementById("result");
let alertDialog = document.getElementById("alert_container");
let connectionsDialog = document.getElementById("new-database");
let alertTopic = document.getElementById("alert_topic");
let alertMessage = document.getElementById("alert_message");
let selectConnections = document.getElementById("iconnections");
let managedConnections = document.getElementById("existing_connections");
let storeHostInput = document.getElementById("new_connetion");
let storeHostForm = document.getElementById("store_connection");
let connectionStatus = document.getElementById("connection_status");
let influxUser = document.getElementById("influxdb_username");
let influxPw = document.getElementById("influxdb_password");
let connectionForm = document.getElementById("influxdb_connection_form");
let queryForm = document.getElementById("send_query_form");
let queryInput = document.getElementById("influxdb_query");
let influxDBSelection = document.getElementById("database_selection");

function toggleState(element, state){
    if(state == true){
        element.classList.remove("hidden");
        element.classList.add("show");
    }
    else {
        element.classList.remove("show");
        element.classList.add("hidden");
    }
}

function removeConnection(url){
    window.toggleConnectionsDialog(false);
    console.log("This is going to be deleted "+url);
    window.deleteHost(url);
}

function popuplateConnections(connections){
    managedConnections.innerHTML='';
    for (var i = 0; i < connections.length; i++) {
        if(connections[i] == "http://localhost:8086"){
            continue;
        }
        console.log(connections[i]);
        var item = document.createElement("li");
        var connectionTxt = document.createElement("span");
        connectionTxt.textContent = connections[i]
        var removeBtn = document.createElement("button");
        removeBtn.classList.add("seconday-button");
        removeBtn.textContent = "Remove";
        removeBtn.setAttribute("connection-url", connections[i]);
        removeBtn.addEventListener("click", function(evt){
            var url = evt.target.getAttribute("connection-url");
            removeConnection(url);
        });
        item.appendChild(connectionTxt);
        item.appendChild(removeBtn);
        managedConnections.appendChild(item);
        console.log("added li element");
    }
};

window.toggleAlertDialog = function(state){
    //var dialogState = state == true ? "show" : "hidden";
    toggleState(alertDialog, state);
    
};

window.toggleConnectionsDialog = function(state){
    if(state == false){
        selectConnections.selectedIndex = 0;
    }
    toggleState(connectionsDialog, state);
};

window.setAlertMessage = function(message, topic){
    var aTopic = "Alert"
    if(topic != null){
        aTopic = topic;
    }
    alertMessage.textContent = message;
    alertTopic.textContent = aTopic;
};

selectConnections.onchange = function() {
    var value = selectConnections.value;
    console.log("Value was changed: "+selectConnections.value);
    if(value == "Manage connections") {
        window.toggleConnectionsDialog(true);
    }
}

window.toggleConnectionStatus = function(state){
    if(state == true){
        connectionStatus.classList.remove("not_connected");
        connectionStatus.classList.add("connected");
    }
    else {
        connectionStatus.classList.remove("connected");
        connectionStatus.classList.add("not_connected");
    }
}


// backend calls
// Setup the greet function
window.greet = function () {
    // Get name
    let name = nameElement.value;

    // Check if the input is empty
    if (name === "") return;

    // Call App.Greet(name)
    try {
        Greet(name)
            .then((result) => {
                // Update result with data back from App.Greet()
                resultElement.innerText = result;
            })
            .catch((err) => {
                console.error(err);
            });
    } catch (err) {
        console.error(err);
    }
};

window.getConnections = function() {
    let connections = document.getElementById("iconnections");
    try {
        GetConnections().then((result) => {
            console.log(result);
            try {
                var data = JSON.parse(result);
            }
            catch(err) {
                window.setAlertMessage(result);
                window.toggleAlertDialog(true);
                data = [];
            }
            console.log(Array.isArray(data));
            popuplateConnections(data);
            data.push("Manage connections");
            connections.innerHTML='';
            for (var i = 0; i < data.length; i++) {
                var option = document.createElement("option");
                option.value = data[i];
                option.text = data[i];
                connections.appendChild(option);
            }
            
        })
        .catch((err) => {
            console.log("error: "+err);
            window.setAlertMessage(err);
            window.toggleAlertDialog(true);
        });
    } catch (err) {
        console.log("Error from catch: "+err);
        window.setAlertMessage(err);
        window.toggleAlertDialog(true);

    };
};

window.storeHost = function(){
    let host = storeHostInput.value;
    if (host == "") return;
    try {
        StoreConnections(host)
            .then((result) => {
                if(result == 200){
                    window.setAlertMessage(host+" stored", "Success");    
                }
                else {
                    window.setAlertMessage(result);
                }
                window.toggleAlertDialog(true);
            })
            .catch((err) => {
                window.setAlertMessage(err);
                window.toggleAlertDialog(true);
            })
    } catch (err) {
        window.setAlertMessage(err);
        window.toggleAlertDialog(true);
    };
    window.getConnections();
    return false;
};

window.deleteHost = function(host) {
    if(host == "") return;
    try {
        DeleteConnection(host)
            .then((result) => {
                if(result == 200){
                    window.setAlertMessage(host+" deleted", "Success");
                }
                else {
                    window.setAlertMessage(result);
                }
                window.toggleAlertDialog(true);
            })
            .catch((err) => {
                window.setAlertMessage(err);
                window.toggleAlertDialog(true);
            });
    } catch (err) {
        window.setAlertMessage(err);
        window.toggleAlertDialog(true);
    }
    window.getConnections();
}

window.connectInfluxDB = function(){
    var connecion = {
        "Host": selectConnections.value,
        "Usename": influxUser.value,
        "Password": influxPw.value,
    };
    try {
        Connect(JSON.stringify(connecion))
            .then((result)=>{
                if(result==200){
                    toggleConnectionStatus(true);
                    window.getDatabases();
                }
                else {
                    window.setAlertMessage(result);
                    window.toggleAlertDialog(true);
                }
            })
            .catch((err) => {
                window.setAlertMessage(err);
                window.toggleAlertDialog(true);
            });
    } catch (err) {
        window.setAlertMessage(err);
        window.toggleAlertDialog(true);
    }
};

window.getDatabases = function(){
    console.log("Fetching databases");
    try {
        GetDatabases()
            .then((result) => {
                // Update result with data back from App.Greet()
                console.log("Data: "+result);
                try {
                    var data = JSON.parse(result);
                }catch (err) {
                    window.setAlertMessage(err);
                    window.toggleAlertDialog(true);
                }
                influxDBSelection.innerHTML='';
                for(var i=0; i<data.length; i++){
                    console.log("DB name: "+data[i]);
                    var option = document.createElement("option");
                    option.value = data[i];
                    option.text = data[i];
                    influxDBSelection.appendChild(option);
                }
            })
            .catch((err) => {
                window.setAlertMessage(err);
                window.toggleAlertDialog(true);
            });
    } catch (err) {
        window.setAlertMessage(err);
        window.toggleAlertDialog(true);
    }
    console.log("Done fetching databases");
};

window.doQuery = function(){
    console.log("Start query");
    var query = queryInput.value;
    var database = influxDBSelection.value;
    if(query == ""){
        window.alertMessage("Query can't be empty");
        window.toggleAlertDialog(true);
        return
    }
    try {
        RunQuery(query, database)
            .then((result) => {
                console.log("query result");
                console.log(result);
                if(result == ""){
                    window.setAlertMessage("Query "+query+" done", "Success");
                    window.toggleAlertDialog(true);
                }
                else {
                    try {
                        var data=JSON.parse(result);
                        data["data"].push('');
                        window.populateDataTable(data);
                    }
                    catch(err) {
                        window.setAlertMessage(result+"<br />"+err);
                        window.toggleAlertDialog(true);        
                    }
                }
            })
            .catch((err) => {
                window.setAlertMessage(err);
                window.toggleAlertDialog(true);
            });
    } catch(err){
        window.setAlertMessage(err);
        window.toggleAlertDialog(true);
    }
    console.log("end query");
}

// when loaded
storeHostForm.addEventListener("submit", function(evt){
    evt.preventDefault();
    window.toggleConnectionsDialog(false);
    return window.storeHost();
});
connectionForm.addEventListener("submit", function(evt){
    evt.preventDefault();
    window.connectInfluxDB();
    return false;
});
queryForm.addEventListener("submit", function(evt){
    evt.preventDefault();
    window.doQuery();
    return false;
});
window.getConnections();
window.toggleAlertDialog(false);
window.toggleConnectionsDialog(false);
window.toggleConnectionStatus(false);

//Datatable

window.roott = document.getElementById('query_content_table');
window.t= document.createElement('table');
window.roott.appendChild(window.t);
var data = {
    "headings": [
        "Query run"
    ],
    "data": [
        ["Need to run query first"],
        [],
    ]
};
window.populateDataTable = function(jsonData){
    console.log(jsonData);
    jsonData["data"].pop();
    var children = window.roott.children;
    for (var i = 0; i < children.length; i++) {
        var tableChild = children[i];
        window.roott.removeChild(tableChild);
    }
    var dt = document.createElement('table')
    window.roott.appendChild(dt);   
    var dataTable = new DataTable(dt);
    dataTable.insert(jsonData);
    console.log("updated");
};
window.populateDataTable(data);
console.log("Datatable popuplated");