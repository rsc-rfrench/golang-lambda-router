package main

import (
	"router"
	"testing"
)

func TestHelp(t *testing.T) {
	response, err := help(router.Request{})
	if response.StatusCode != 200 {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
}
