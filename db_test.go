package main

import (
	"testing"
)

func testInsertUser(t *testing.T) {
	client, err := createClient()
	if err != nil {
		t.Errorf("error making client")
	}
	response, err := insertUser("testUser", 1215, client)
	if err != nil {
		t.Errorf("function errored")
	}

	if response == nil {
		t.Errorf("nil response")
	}
}
