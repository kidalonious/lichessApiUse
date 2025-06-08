package main

import (
	"testing"

)

func TestInsertUser(t *testing.T) {
	client, err := createClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Clean up any pre-existing user (ignore error)
	_ = deleteUser("testUser", client)

	// Insert user
	err = insertUser("testUser", 1215, client)
	if err != nil {
		t.Fatalf("insertUser failed: %v", err)
	}

	// Confirm user was inserted
	user, err := getUser("testUser", client)
	if err != nil {
		t.Fatalf("getUser failed after insert: %v", err)
	}

	if user.Username != "testUser" || user.Rating != 1215 {
		t.Errorf("user fields do not match: got %+v", user)
	}

	// Clean up
	_ = deleteUser("testUser", client)
}

func TestDeleteUser(t *testing.T) {
	client, err := createClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Ensure user exists before deletion
	err = insertUser("testUser", 1300, client)
	if err != nil {
		t.Fatalf("insertUser setup failed: %v", err)
	}

	// Delete user
	err = deleteUser("testUser", client)
	if err != nil {
		t.Fatalf("deleteUser failed: %v", err)
	}

	// Verify deletion
	_, err = getUser("testUser", client)
	if err == nil {
		t.Errorf("expected getUser to fail after deletion, but it succeeded")
	}
}

func TestGetUser(t *testing.T) {
	client, err := createClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	username := "testUserToGet"
	rating := 1700

	// Clean up any previous test runs
	_ = deleteUser(username, client)

	// Insert user to get
	err = insertUser(username, rating, client)
	if err != nil {
		t.Fatalf("insertUser failed: %v", err)
	}

	// Get the user
	user, err := getUser(username, client)
	if err != nil {
		t.Fatalf("getUser failed: %v", err)
	}

	if user.Username != username || user.Rating != rating {
		t.Errorf("user data mismatch. got %+v", user)
	}

	// Cleanup
	_ = deleteUser(username, client)
}