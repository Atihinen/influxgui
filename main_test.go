package main

import (
	"os"
	"testing"
)

func ensureFalseDebugFlag(t *testing.T, result bool) {
	if result {
		t.Error("Return should be false was", result)
	}
}

func ensureTrueDebugFlag(t *testing.T, result bool) {
	if !result {
		t.Error("Return should be true was", result)
	}
}

func TestGetDebugFlag(t *testing.T) {
	result := getDebugFlag()
	ensureFalseDebugFlag(t, result)
	os.Setenv("ENV", "production")
	result = getDebugFlag()
	ensureFalseDebugFlag(t, result)
	os.Setenv("ENV", "develop")
	result = getDebugFlag()
	ensureTrueDebugFlag(t, result)

}
