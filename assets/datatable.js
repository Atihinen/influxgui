window.ready(function() {
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
});


//{headings: ["time", "host", "region", "value"], data: [["2020-01-15T19:57:29.149121976Z"], ["serverA"], ["us_west"], ["0.64"]]}

//window.populateDataTable({headings: ["host", "region", "value"], data: [["serverA"], ["us_west"], ["0.64"]]})