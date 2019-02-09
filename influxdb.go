package main

import (
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/zserge/webview"
	"log"
)

type InfluxDBConnection struct {
	Host     string
	Username string
	Password string
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
