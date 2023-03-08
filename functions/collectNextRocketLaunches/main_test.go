package main

import (
	"testing"
)

func TestRead(t *testing.T) {
	value := getNext5Launches()

	if value.Count <= 0 {
		t.Errorf("Expected value to be greater than 0")
	}
}
