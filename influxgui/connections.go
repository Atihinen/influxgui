package main

import (
	"container/list"
)

func getConnections() (*list.List, error) {
	connections := list.New()
	connections.PushBack("http://localhost:8086")
	return connections, nil
}
