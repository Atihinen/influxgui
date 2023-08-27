package main

import (
	"context"
	"encoding/json"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
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
