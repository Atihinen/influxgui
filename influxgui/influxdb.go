package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
)

type InfluxDBConnection struct {
	Host     string
	Username string
	Password string
}

func NewClient(url string, username string, password string) (client.Client, error) {
	if username != "" {
		c, err := client.NewHTTPClient(client.HTTPConfig{
			Addr:     url,
			Username: username,
			Password: password,
		})
		return c, err
	}
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: url,
	})
	return c, err
}

func Ping(connection client.Client) (int, error) {
	_, _, err := connection.Ping(0)
	if err != nil {
		return 400, err
	}
	return 200, nil
}

func GetDatabases(connection client.Client) (bool, *list.List, string) {
	status := true
	databases := list.New()
	data := ""
	q := client.NewQuery("SHOW DATABASES;", "", "ns")
	log.Println("Query %v", q)
	if response, err := connection.Query(q); err == nil && response.Error() == nil {
		for _, value := range response.Results[0].Series[0].Values {
			res := value[0]
			log.Println("Result %v", value[0])
			//option := fmt.Sprintf("%s,%s,", res, res)
			databases.PushBack(res)
		}
	}
	return status, databases, data
}

/*
func DoQuery(connection client.Client, query string, database string) (string, error) {
	q := client.NewQuery(query, database, "")
	response, err := connection.Query(q)
	if err == nil && response.Error() == nil {
		results := "{"
		columns := "\"headings\": ["
		for _, serie := range response.Results[0].Series {
			for _, column := range serie.Columns {
				columns = fmt.Sprintf("%s\"%s\",", columns, column)
			}
			values := "\"data\": [["
			for _, value := range serie.Values {
				for _, val := range value {
					values = fmt.Sprintf("%s\"%v\",", values, val)
				}
				values = fmt.Sprintf("%s],[", values)
			}
			results = fmt.Sprintf("%s%s],%s]]}", results, columns, values)
			return results, nil
		}
	} else if response.Error() != nil && err == nil {
		response_error := fmt.Sprintf("%v", response)
		response_error = strings.Trim(response_error, "&{[] ")
		response_error = strings.Trim(response_error, "}")
		return "", response.Error()
	}
	return "", err
}*/

func DoQuery(connection client.Client, query string, database string) (string, error) {
	q := client.NewQuery(query, database, "")
	response, err := connection.Query(q)
	if err == nil && response.Error() == nil {
		results := map[string]interface{}{} // Use a map to construct the JSON

		// Construct the "headings" array
		var headings []string
		for _, serie := range response.Results[0].Series {
			for _, column := range serie.Columns {
				headings = append(headings, column)
			}
		}
		results["headings"] = headings

		// Construct the "data" array
		var data [][]interface{}
		for _, serie := range response.Results[0].Series {
			for _, value := range serie.Values {
				data = append(data, value)
			}
		}
		results["data"] = data

		// Marshal the map to JSON
		jsonData, err := json.Marshal(results)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	} else if response.Error() != nil && err == nil {
		responseError := fmt.Sprintf("%v", response)
		responseError = strings.Trim(responseError, "&{[] ")
		responseError = strings.Trim(responseError, "}")
		return "", response.Error()
	}
	return "", err
}
