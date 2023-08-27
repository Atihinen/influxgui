package main

import (
	"github.com/tidwall/buntdb"
)

func openDB() (*buntdb.DB, error) {
	db, err := buntdb.Open("influxConfigs.db")
	return db, err
}
