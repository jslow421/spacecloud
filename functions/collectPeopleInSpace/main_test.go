package main

import (
	"testing"
)

func TestRead(t *testing.T) {
	value := getPeopleInSpaceFromApi()

	if value.Message != "success" {
		t.Errorf("Expected value to be not empty")
	}

	if len(value.People) < 1 {
		t.Errorf("Either the api is down or there are no people in space")
	}
}
