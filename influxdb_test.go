package main

import (
	"container/list"
	"fmt"
	"testing"
)

func testSetup() {
	dblog = list.New()
}

func TestAppendToLog(t *testing.T) {
	testSetup()
	dbloglen := dblog.Len()
	if dbloglen != 0 {
		t.Error("Len should be 0, was", dbloglen)
	}
	appendToLog("SHOW MEASUREMENTS")
	dbloglen = dblog.Len()
	if dbloglen != 1 {
		t.Error("Len should not be 1, was", dbloglen)
	}
}

func TestAppendToLogTakesOnly10Elements(t *testing.T) {
	testSetup()
	for i := 1; i < 12; i++ {
		appendToLog(fmt.Sprintf("SELECT * FROM data LIMIT %d", i))
	}
	dbloglen := dblog.Len()
	if dbloglen != 10 {
		t.Error("Len should not be 10, was", dbloglen)
	}
	expected_first := "SELECT * FROM data LIMIT 2"
	observed_first := dblog.Front().Value
	if expected_first != observed_first {
		t.Error(fmt.Sprintf("First element should be %s, was %s", expected_first, observed_first))
	}
	expected_last := "SELECT * FROM data LIMIT 11"
	observed_last := dblog.Back().Value
	if expected_last != observed_last {
		t.Error(fmt.Sprintf("First element should be %s, was %s", expected_last, observed_last))
	}
}

func TestGetLogOptions(t *testing.T) {
	testSetup()
	appendToLog("SELECT * FROM data LIMIT 1;")
	expected := "<option value='SELECT * FROM data LIMIT 1;'>SELECT * FROM data LIMIT 1;</option>"
	observed := getLogOptions()
	if expected != observed {
		t.Error(fmt.Sprintf("Observed %s was not expected: %s", observed, expected))
	}
}
