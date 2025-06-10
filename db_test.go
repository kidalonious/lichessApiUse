package main

import (
	"testing"

)

func TestInsertUser(t *testing.T) {

	// Clean up any pre-existing user (ignore error)
	_ = DeleteUser("testUser")

	// Insert user
	err := InsertUser("testUser", 1215)
	if err != nil {
		t.Fatalf("InsertUser failed: %v", err)
	}

	// Confirm user was inserted
	user, err := GetUser("testUser")
	if err != nil {
		t.Fatalf("GetUser failed after insert: %v", err)
	}

	if user.Username != "testUser" || user.Rating != 1215 {
		t.Errorf("user fields do not match: got %+v", user)
	}

	// Clean up
	_ = DeleteUser("testUser")
}

func TestDeleteUser(t *testing.T) {

	// Ensure user exists before deletion
	err := InsertUser("testUserToDelete", 1300)
	if err != nil {
		t.Fatalf("InsertUser setup failed: %v", err)
	}

	// Delete user
	err = DeleteUser("testUserToDelete")
	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	// Verify deletion
	_, err = GetUser("testUserToDelete")
	if err == nil {
		t.Errorf("expected GetUser to fail after deletion, but it succeeded")
	}
}

func TestGetUser(t *testing.T) {

	username := "testUserToGet"
	rating := 1700

	// Clean up any previous test runs
	_ = DeleteUser(username)

	// Insert user to get
	err := InsertUser(username, rating)
	if err != nil {
		t.Fatalf("InsertUser failed: %v", err)
	}

	// Get the user
	user, err := GetUser(username)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}

	if user.Username != username || user.Rating != rating {
		t.Errorf("user data mismatch. got %+v", user)
	}

	// Cleanup
	_ = DeleteUser(username)
}

func TestInsertGame(t *testing.T) {

	white := "InsertTestWhite"
	black := "InsertTestBlack"
	err := InsertGame(white, black, "white", "King's Indian", "e4 e5 Nf3", "resign")
	if err != nil {
		t.Fatalf("InsertGame failed: %v", err)
	}

	// Validate insert by fetching it back
	games, err := GetGameByPlayers(white, black)
	if err != nil {
		t.Fatalf("GetGameByPlayers failed: %v", err)
	}

	if len(games) == 0 {
		t.Fatalf("no games found after insertion")
	}

	// Cleanup
	err = DeleteGame(games[0].Gameid)
	if err != nil {
		t.Errorf("cleanup failed: %v", err)
	}
}

func TestGetGame(t *testing.T) {

	white := "GetTestWhite"
	black := "GetTestBlack"
	err := InsertGame(white, black, "draw", "Italian Game", "e4 e5 Nf3 Nc6 Bc4", "resign")
	if err != nil {
		t.Fatalf("InsertGame failed: %v", err)
	}

	games, err := GetGameByPlayers(white, black)
	if err != nil {
		t.Fatalf("GetGameByPlayers failed: %v", err)
	}

	gameID := games[0].Gameid
	game, err := GetGame(gameID)
	if err != nil {
		t.Fatalf("GetGame failed: %v", err)
	}

	if game.Whiteplayer != white || game.Blackplayer != black {
		t.Errorf("game data mismatch: got %s vs %s / %s vs %s", game.Whiteplayer, white, game.Blackplayer, black)
	}

	// Cleanup
	err = DeleteGame(gameID)
	if err != nil {
		t.Errorf("cleanup failed: %v", err)
	}
}

func TestDeleteGame(t *testing.T) {

	white := "DeleteTestWhite"
	black := "DeleteTestBlack"
	err := InsertGame(white, black, "black", "Scotch Game", "e4 e5 Nf3 Nc6 d4", "resign")
	if err != nil {
		t.Fatalf("InsertGame failed: %v", err)
	}

	games, err := GetGameByPlayers(white, black)
	if err != nil {
		t.Fatalf("GetGameByPlayers failed: %v", err)
	}
	gameID := games[0].Gameid

	err = DeleteGame(gameID)
	if err != nil {
		t.Fatalf("DeleteGame failed: %v", err)
	}

	// Verify deletion
	_, err = GetGame(gameID)
	if err == nil {
		t.Errorf("expected error retrieving deleted game, got none")
	}
}