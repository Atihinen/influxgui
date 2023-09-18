package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// App struct
type App struct {
	ctx        context.Context
	connection InfluxDBConnection
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetConnections() string {
	db, err := openDB()
	if err != nil {
		return err.Error()
	}
	connections, err := getConnections(db)
	if err != nil {
		return err.Error()
	}
	var connectionsSlice []string
	for e := connections.Front(); e != nil; e = e.Next() {
		if val, ok := e.Value.(string); ok {
			connectionsSlice = append(connectionsSlice, val)
		}
	}
	jsonArray, err := json.Marshal(connectionsSlice)
	if err != nil {
		fmt.Println("Error:", err)
		return err.Error()
	}
	return string(jsonArray)

}

func (a *App) StoreConnections(host string) string {
	db, err := openDB()
	if err != nil {
		return err.Error()
	}
	rc, err := storeConnectionConfig(db, host)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%v", rc)

}

func (a *App) DeleteConnection(host string) string {
	db, err := openDB()
	if err != nil {
		return err.Error()
	}
	rc, err := deleteConnectionConfig(db, host)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%v", rc)
}

func (a *App) Connect(content string) string {
	var connectionData = InfluxDBConnection{}
	err := json.Unmarshal([]byte(content), &connectionData)
	if err != nil {
		return err.Error()
	}
	connection, err := NewClient(connectionData.Host, connectionData.Username, connectionData.Password)
	if err != nil {
		return err.Error()
	}
	rc, err := Ping(connection)
	if err != nil {
		return err.Error()
	}
	a.connection = connectionData //store connection for future queries
	return fmt.Sprintf("%v", rc)
}

func (a *App) GetDatabases() string {
	connection, err := NewClient(a.connection.Host, a.connection.Username, a.connection.Password)
	if err != nil {
		return err.Error()
	}
	status, databases, data := GetDatabases(connection)
	log.Println("databases %v", databases)
	if !status {
		return data
	}
	var slice []string
	for e := databases.Front(); e != nil; e = e.Next() {
		if val, ok := e.Value.(string); ok {
			log.Println("Db name %v", val)
			slice = append(slice, val)
		}
	}
	jsonArray, err := json.Marshal(slice)
	if err != nil {
		fmt.Println("Error:", err)
		return err.Error()
	}
	return string(jsonArray)
}

func (a *App) RunQuery(query string, database string) string {
	connection, err := NewClient(a.connection.Host, a.connection.Username, a.connection.Password)
	if err != nil {
		return err.Error()
	}
	data, err := DoQuery(connection, query, database)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf(data)
}
