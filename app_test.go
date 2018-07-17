package main

import (
	"testing"
)

func TestHelp(t *testing.T) {
	response, err := help(Request{})
	if response.StatusCode != 200 {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
}
