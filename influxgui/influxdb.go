package main

import (
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
			Addr:     "http://localhost:8086",
			Username: username,
			Password: password,
		})
		return c, err
	}
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
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
