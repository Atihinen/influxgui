window.rpc = {
    askConnections: function() { 
      this.databases = [];
      this.databases.push("New connection");
      window.external.invoke(JSON.stringify({cmd : 'connections'})); 
    },
    setConnections: function(connections) {
      this.databases = connections;
    },
    getConnections: function() { 
      return this.connections 
    },
    storeConnection: function(host){
      console.log("sending... "+host)
      window.external.invoke(JSON.stringify({cmd : 'addCon', host: host}));
    },
    deleteConnection: function(host){
        window.external.invoke(JSON.stringify({cmd: 'delCon', host: host}));
    },
    setHost: function(host){
      window.external.invoke(JSON.stringify({cmd : 'setHost', host: host}));
    },
    connectInfluxDB: function(username, password){
      window.external.invoke(JSON.stringify({cmd: 'connectInflux', username: username, password: password}));
    },
    sendQuery: function(query, database){
        //alert("Database: "+database+", query: "+query);
        window.external.invoke(JSON.stringify({cmd: 'sendQuery', query: query, database: database}));
    },
    results: [],
    databases: ["New connection", "http://localhost:8086"],
    connected: false,
};

window.ready = function (fn) {
    if (document.readyState != 'loading') {
        fn();
    } else if (document.addEventListener) {
        document.addEventListener('DOMContentLoaded', fn);
    } else {
        document.attachEvent('onreadystatechange', function() {
        if (document.readyState != 'loading')
            fn();
        });
    }
}

function InfluxDBGUI(){
    var self = this;
    self.Rpc = window.rpc;
    self.connections = ko.observableArray(self.Rpc.databases);
    self.historyContent = ko.observableArray(["SHOW MEASUREMENTS;",
                        "SELECT * FROM measurement LIMIT 1;",
                        "SHOW DATABASES;"]);
    self.databases = ko.observableArray([]);
    self.connection = ko.observable();
    self.database = ko.observable();
    self.query = ko.observable();
    self.modalClass = ko.observable(false);
    self.connectedClass = ko.observable(false);
    self.username = ko.observable();
    self.password = ko.observable();
    self.newConnectionValue = ko.observable();
    self.queryResult = ko.observable();
    self.history = ko.observable("Show history");
    self.showHistory = ko.observable(false);
    self.closeModal = function(){
        self.modalClass(false);
    };
    self.setConnection = function(isConnected){
        self.connectedClass(isConnected);
    }
    self.toggleHistory = function(){
        self.showHistory(!self.showHistory());
        if(self.showHistory()){
            self.history("Hide history");
        }
        else {
            self.history("Show history")
        }
    }
    self.setQueryFromHistory = function(query){
        console.log("Adding... "+query);
        self.query(query);
        console.log("Added.");
    };
    self.deleteConnection = function(host){
        console.log("Deleting..."+host);
        self.Rpc.deleteConnection(host);
        return false;
    };
    /* eslint-disable no-unused-vars */
    self.newConnection = function(formElement){
        self.Rpc.storeConnection(self.newConnectionValue());
        self.modalClass(false);
        return false;
    }
    self.connectToInfluxDB = function(formElement){
        self.Rpc.connectInfluxDB(self.username(), self.password());
        return false;
    }

    self.createQuery = function(formElement){
        console.log("Database: "+self.database());
        console.log("Query: "+self.query());
        self.Rpc.sendQuery(self.query(), self.database());
        self.historyContent.push(self.query());
        return false;
    }
    /* eslint-enable no-unused-vars */
    self.setConnections = function(newConnections){
        self.connections.removeAll();
        for(var i=0; i<newConnections.length; i++){
            self.connections.push(newConnections[i]);
        }
        self.connections.push("New connection");
    }
    self.setDatabases = function(databases){
        self.databases.removeAll();
        for(var i=0; i<databases.length; i++){
            self.databases.push(databases[i]);
        }
    }

}
window.mv = { viewmodel: new InfluxDBGUI() };
window.mv.viewmodel.connection.subscribe(function(newValue){
    if(newValue == "New connection"){
        window.mv.viewmodel.modalClass(true);
    }
    else {
        window.mv.viewmodel.Rpc.setHost(newValue);
        window.mv.viewmodel.setConnection(false);
    }
});
window.mv.viewmodel.toggleModal = ko.pureComputed(function(){
    return window.mv.viewmodel.modalClass() == false ? "hidden" : "block";
});

window.mv.viewmodel.connected = ko.pureComputed(function(){
    return window.mv.viewmodel.connectedClass() == false ? "not_connected" : "connected";
});

window.ready(function() {
    ko.applyBindings(window.mv.viewmodel);
    window.rpc.askConnections();
});
