package main

import (
	"testing"

)

func TestInsertUser(t *testing.T) {

	// Clean up any pre-existing user (ignore error)
	_ = deleteUser("testUser")

	// Insert user
	err := insertUser("testUser", 1215)
	if err != nil {
		t.Fatalf("insertUser failed: %v", err)
	}

	// Confirm user was inserted
	user, err := getUser("testUser")
	if err != nil {
		t.Fatalf("getUser failed after insert: %v", err)
	}

	if user.Username != "testUser" || user.Rating != 1215 {
		t.Errorf("user fields do not match: got %+v", user)
	}

	// Clean up
	_ = deleteUser("testUser")
}

func TestDeleteUser(t *testing.T) {

	// Ensure user exists before deletion
	err := insertUser("testUserToDelete", 1300)
	if err != nil {
		t.Fatalf("insertUser setup failed: %v", err)
	}

	// Delete user
	err = deleteUser("testUserToDelete")
	if err != nil {
		t.Fatalf("deleteUser failed: %v", err)
	}

	// Verify deletion
	_, err = getUser("testUserToDelete")
	if err == nil {
		t.Errorf("expected getUser to fail after deletion, but it succeeded")
	}
}

func TestGetUser(t *testing.T) {

	username := "testUserToGet"
	rating := 1700

	// Clean up any previous test runs
	_ = deleteUser(username)

	// Insert user to get
	err := insertUser(username, rating)
	if err != nil {
		t.Fatalf("insertUser failed: %v", err)
	}

	// Get the user
	user, err := getUser(username)
	if err != nil {
		t.Fatalf("getUser failed: %v", err)
	}

	if user.Username != username || user.Rating != rating {
		t.Errorf("user data mismatch. got %+v", user)
	}

	// Cleanup
	_ = deleteUser(username)
}

func TestInsertGame(t *testing.T) {

	white := "InsertTestWhite"
	black := "InsertTestBlack"
	err := insertGame(white, black, "white", "King's Indian", "e4 e5 Nf3", "resign")
	if err != nil {
		t.Fatalf("insertGame failed: %v", err)
	}

	// Validate insert by fetching it back
	games, err := getGameByPlayers(white, black)
	if err != nil {
		t.Fatalf("getGameByPlayers failed: %v", err)
	}

	if len(games) == 0 {
		t.Fatalf("no games found after insertion")
	}

	// Cleanup
	err = deleteGame(games[0].Gameid)
	if err != nil {
		t.Errorf("cleanup failed: %v", err)
	}
}

func TestGetGame(t *testing.T) {

	white := "GetTestWhite"
	black := "GetTestBlack"
	err := insertGame(white, black, "draw", "Italian Game", "e4 e5 Nf3 Nc6 Bc4", "resign")
	if err != nil {
		t.Fatalf("insertGame failed: %v", err)
	}

	games, err := getGameByPlayers(white, black)
	if err != nil {
		t.Fatalf("getGameByPlayers failed: %v", err)
	}

	gameID := games[0].Gameid
	game, err := getGame(gameID)
	if err != nil {
		t.Fatalf("getGame failed: %v", err)
	}

	if game.Whiteplayer != white || game.Blackplayer != black {
		t.Errorf("game data mismatch: got %s vs %s / %s vs %s", game.Whiteplayer, white, game.Blackplayer, black)
	}

	// Cleanup
	err = deleteGame(gameID)
	if err != nil {
		t.Errorf("cleanup failed: %v", err)
	}
}

func TestDeleteGame(t *testing.T) {

	white := "DeleteTestWhite"
	black := "DeleteTestBlack"
	err := insertGame(white, black, "black", "Scotch Game", "e4 e5 Nf3 Nc6 d4", "resign")
	if err != nil {
		t.Fatalf("insertGame failed: %v", err)
	}

	games, err := getGameByPlayers(white, black)
	if err != nil {
		t.Fatalf("getGameByPlayers failed: %v", err)
	}
	gameID := games[0].Gameid

	err = deleteGame(gameID)
	if err != nil {
		t.Fatalf("deleteGame failed: %v", err)
	}

	// Verify deletion
	_, err = getGame(gameID)
	if err == nil {
		t.Errorf("expected error retrieving deleted game, got none")
	}
}