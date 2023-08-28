package main

import (
	"container/list"
	"fmt"

	"github.com/tidwall/buntdb"
)

func getExists(connections *list.List, host string) bool {
	exists := false
	for connection := connections.Front(); connection != nil; connection = connection.Next() {
		if host == connection.Value {
			exists = true
		}
	}
	return exists
}

func getConnections(db *buntdb.DB) (*list.List, error) {

	connections := list.New()
	connections.PushBack("http://localhost:8086")
	err := db.View(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(key, val string) bool {
			connections.PushBack(val)
			fmt.Printf("%s %s\n", key, val)
			return true
		})
		return nil
	})
	return connections, err
}

func storeConnectionConfig(db *buntdb.DB, host string) (int, error) {
	connections, _ := getConnections(db)
	exists := false
	for connection := connections.Front(); connection != nil; connection = connection.Next() {
		if host == connection.Value {
			exists = true
		}
	}
	if !exists {
		err := db.Update(func(tx *buntdb.Tx) error {
			_, _, err := tx.Set(host, host, nil)
			return err
		})
		if err != nil {
			return 401, err
		}
	}
	return 200, nil

}

func deleteConnectionConfig(db *buntdb.DB, host string) (int, error) {
	connections, _ := getConnections(db)
	exists := getExists(connections, host)
	if exists {
		err := db.Update(func(tx *buntdb.Tx) error {
			_, err := tx.Delete(host)
			return err
		})
		if err != nil {
			return 400, err
		}
	}
	return 200, nil
}
