package main

import (
	"bytes"
	"container/list"
	"fmt"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"
)

var dblog *list.List

type InfluxDBConnection struct {
	Host         string
	Username     string
	Password     string
	Error        bool
	ErrorMessage string
}

type InfluxDBQuery struct {
	Query    string
	Database string
}

func getInfluxDBClient() (client.Client, error) {
	influxdbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     connectionConfig.Host,
		Username: connectionConfig.Username,
		Password: connectionConfig.Password,
	})
	return influxdbClient, err

}

func pingInfluxDB() (bool, string) {
	message := ""
	buf := bytes.NewBufferString(message)
	success := true
	influxdbClient, err := getInfluxDBClient()
	if err != nil {
		buf.WriteString("Could not connect to influxdb")
		buf.WriteString(err.Error())
		success = false
	}
	if !success {
		return success, buf.String()
	}
	message = ""
	_, _, err = influxdbClient.Ping(0)
	if err != nil {
		buf.WriteString("Could not ping to influxdb")
		buf.WriteString(err.Error())
		success = false
	}
	return success, buf.String()
}

func runInfluxDBQuery(query string, database string) (bool, string) {
	status := true
	influxdbClient, err := getInfluxDBClient()
	data := ""
	buf := bytes.NewBufferString(data)
	if err != nil {
		buf.WriteString("Could not connect to influxdb")
		buf.WriteString(err.Error())
		status = false
	}
	defer influxdbClient.Close()
	q := client.NewQuery(query, database, "")
	response, err := influxdbClient.Query(q)
	if err == nil && response.Error() == nil {
		results := "------------------\\n"
		columns := ""
		for _, serie := range response.Results[0].Series {
			for _, column := range serie.Columns {
				columns = fmt.Sprintf("%s%s\\t", columns, column)
			}
			results = fmt.Sprintf("%s%s\\n", results, columns)
			values := ""
			for _, value := range serie.Values {
				for _, val := range value {
					values = fmt.Sprintf("%s%v\\t", values, val)
				}
				values = fmt.Sprintf("%s\\n", values)
			}
			results = fmt.Sprintf("%s%s\\n", results, values)
		}

		data = createInfluxDBQueryResponse(results)
	} else if response.Error() != nil && err == nil {
		response_error := fmt.Sprintf("%v", response)
		response_error = strings.Trim(response_error, "&{[] ")
		response_error = strings.Trim(response_error, "}")
		data = createInfluxDBQueryResponse(response_error)
	}
	return status, data
}

func runInfluxDBQueryJSON(query string, database string) (bool, string) {
	status := true
	influxdbClient, err := getInfluxDBClient()
	data := ""
	buf := bytes.NewBufferString(data)
	if err != nil {
		buf.WriteString("Could not connect to influxdb")
		buf.WriteString(err.Error())
		status = false
	}
	defer influxdbClient.Close()
	q := client.NewQuery(query, database, "")
	response, err := influxdbClient.Query(q)
	if err == nil && response.Error() == nil {
		results := "{"
		columns := "'headings': ["
		for _, serie := range response.Results[0].Series {
			for _, column := range serie.Columns {
				columns = fmt.Sprintf("%s'%s',", columns, column)
			}
			values := "'data': [["
			for _, value := range serie.Values {
				for _, val := range value {
					values = fmt.Sprintf("%s'%v',", values, val)
				}
				values = fmt.Sprintf("%s],[", values)
			}
			results = fmt.Sprintf("%s%s],%s]]}", results, columns, values)
		}
		data = createInfluxDBQueryResponseJSON(results)
	} else if response.Error() != nil && err == nil {
		response_error := fmt.Sprintf("%v", response)
		response_error = strings.Trim(response_error, "&{[] ")
		response_error = strings.Trim(response_error, "}")
		data = createInfluxDBQueryResponseJSON(response_error)
	}
	return status, data
}

func showDatabases() (bool, *list.List, string) {
	status := true
	databases := list.New()
	data := ""
	buf := bytes.NewBufferString(data)
	influxdbClient, err := getInfluxDBClient()

	if err != nil {
		buf.WriteString("Could not connect to influxdb")
		buf.WriteString(err.Error())
		status = false
	}
	defer influxdbClient.Close()

	q := client.NewQuery("SHOW DATABASES;", "", "ns")
	if response, err := influxdbClient.Query(q); err == nil && response.Error() == nil {
		for _, value := range response.Results[0].Series[0].Values {
			res := value[0]
			//option := fmt.Sprintf("%s,%s,", res, res)
			databases.PushBack(res)
		}
	}
	return status, databases, data
}

func appendToLog(query string) {
	dblog.PushBack(query)
	if dblog.Len() > 10 {
		dblog.Remove(dblog.Front())
	}
}

func getLogOptions() string {
	logOptions := ""
	for value := dblog.Back(); value != nil; value = value.Prev() {
		option := fmt.Sprintf("<option value='%v'>%v</option>", value.Value, value.Value)
		logOptions = fmt.Sprintf("%s%s", logOptions, option)
	}
	return logOptions
}

func initalizeLog() {
	dblog = list.New()
	appendToLog("SHOW DATABASES;")
	appendToLog("SELECT * FROM measurement LIMIT 1")
	appendToLog("SHOW MEASUREMENTS;")
}
