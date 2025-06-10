package main

import (
	"os"
	"strings"
	"testing"
	"path/filepath"
)

func TestParsePgnFile(t *testing.T) {
	pgnContent := `[Event "Test Game"]
[Site "Internet"]
[Date "2025.06.09"]
[Round "1"]
[White "WhitePlayer"]
[Black "BlackPlayer"]
[Result "1-0"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 1-0
`

	// Create temporary PGN file
	tmpFile, err := os.CreateTemp("", "testgame*.pgn")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(pgnContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Call the function to test
	pgns, err := ParsePgnFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("parsePgnFile returned error: %v", err)
	}

	if len(pgns) != 1 {
		t.Fatalf("Expected 1 game, got %d", len(pgns))
	}

	pgn := pgns[0]

	// Header checks
	expectedHeaders := map[string]string{
		"Event":  "Test Game",
		"Site":   "Internet",
		"Date":   "2025.06.09",
		"Round":  "1",
		"White":  "WhitePlayer",
		"Black":  "BlackPlayer",
		"Result": "1-0",
	}

	for key, expected := range expectedHeaders {
		if val, ok := pgn.Headers[key]; !ok || val != expected {
			t.Errorf("Header %s: expected '%s', got '%s'", key, expected, val)
		}
	}

	// Moves check (basic)
	if !strings.Contains(pgn.Moves, "1. e4") {
		t.Errorf("Moves string missing '1. e4': %s", pgn.Moves)
	}
	if !strings.Contains(pgn.Moves, "2. Nf3") {
		t.Errorf("Moves string missing '2. Nf3': %s", pgn.Moves)
	}
	if !strings.Contains(pgn.Moves, "3. Bb5") {
		t.Errorf("Moves string missing '3. Bb5': %s", pgn.Moves)
	}
}

func TestGetPgns(t *testing.T) {
	// Ensure the pgns directory exists
	pgnsDir := "pgns"
	if _, err := os.Stat(pgnsDir); os.IsNotExist(err) {
		t.Fatalf("pgns directory does not exist: %v", err)
	}

	// Dummy files to create
	testFiles := []string{"test1.pgn", "test2.pgn"}

	// Create the dummy files
	for _, fileName := range testFiles {
		filePath := filepath.Join(pgnsDir, fileName)
		err := os.WriteFile(filePath, []byte("dummy PGN content"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", fileName, err)
		}
	}

	// Clean up after test by removing the files we created
	defer func() {
		for _, fileName := range testFiles {
			filePath := filepath.Join(pgnsDir, fileName)
			os.Remove(filePath) // Ignoring error intentionally
		}
	}()

	// Run the function under test
	pgnFiles, err := GetPgns()
	if err != nil {
		t.Fatalf("getPgns() returned an error: %v", err)
	}

	// Check that the test files are in the output
	expectedPaths := []string{
		filepath.Join(pgnsDir, "test1.pgn"),
		filepath.Join(pgnsDir, "test2.pgn"),
	}

	found := make(map[string]bool)
	for _, path := range pgnFiles {
		found[path] = true
	}

	for _, expected := range expectedPaths {
		if !found[expected] {
			t.Errorf("expected file path %s not found in result", expected)
		}
	}
}