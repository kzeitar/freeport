package main

import (
	"testing"
)

func TestVersion(t *testing.T) {
	// Basic test to ensure the package can be tested
	if version == "" {
		t.Error("version should not be empty")
	}
}
