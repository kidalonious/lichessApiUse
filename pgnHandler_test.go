package main

import (
	"os"
	"strings"
	"testing"
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
	pgns, err := parsePgnFile(tmpFile.Name())
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