window.rpc = {
  askConnections: function() { 
    this.connections = [];
    window.external.invoke(JSON.stringify({cmd : 'connections'})); 
  },
  setConnections: function(connections) {
    this.connections = connections;
  },
  getConnections: function() { 
    return this.connections 
  },
  storeConnection: function(host){
    console.log("sending... "+host)
    window.external.invoke(JSON.stringify({cmd : 'addCon', host: host}))
  },
  setHost: function(host){
    window.external.invoke(JSON.stringify({cmd : 'setHost', host: host}));
  },
  connectInfluxDB: function(username, password){
    window.external.invoke(JSON.stringify({cmd: 'connectInflux', username: username, password, password}));
  },
  results: [],
  databases: [],
  connected: false,
};

import App from './components/App.js';
Vue.prototype.Rpc = window.rpc;
Vue.prototype.influxDBConnections = window.connections;
window.vue = new Vue({
  render: h => h(App),
}).$mount(`#app`);

