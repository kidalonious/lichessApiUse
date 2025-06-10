package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/notnil/chess"
)
// from the notnil/chess documentation on scanning a large pgn file:
/* 
	f, err := os.Open("lichess_db_standard_rated_2013-01.pgn")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := chess.NewScanner(f)
	for scanner.Scan() {
		game := scanner.Next()
		fmt.Println(game.GetTagPair("Site"))
		// Output &{Site https://lichess.org/8jb5kiqw}
	}
*/

type Pgn struct {
	Headers map[string]string
	Moves string
}

func parsePgnFile(pgnPath string) ([]Pgn, error) {
	var filePgns []Pgn
	f, err := os.Open(pgnPath)
	if err != nil {
		return nil, fmt.Errorf("error opening pgn file: %s", pgnPath)
	}

	defer f.Close()

	notation := chess.AlgebraicNotation{}

	scanner := chess.NewScanner(f)
	for scanner.Scan() {
		game := scanner.Next()
		var pgn Pgn
		headers := make(map[string]string)
		for _, tag := range game.TagPairs() {
			headers[tag.Key] = tag.Value
		}
		pgn.Headers = headers

		moves := game.Moves()
		moveStr := ""
		for i, move := range moves {
			correctNotation := notation.Encode(game.Positions()[i], move)
			if i % 2 == 0 {
				moveStr += fmt.Sprintf("%d. %s ", i/2 + 1, correctNotation)
			} else {
				moveStr += fmt.Sprintf("%s ", correctNotation)
			}
		}
		pgn.Moves = moveStr

		filePgns = append(filePgns, pgn)
	}

	return filePgns, nil
}

func getPgns() ([]string, error) {
	var pgnFilePaths []string
	pathBegin := "pgns"
	dirEntries, err := os.ReadDir("pgns")
	if err != nil {
		return nil, fmt.Errorf("read directory returned error: %w", err)
	}

	for _, entry := range dirEntries {
		filepath := filepath.Join(pathBegin, entry.Name())
		pgnFilePaths = append(pgnFilePaths, filepath)
	}

	return pgnFilePaths, nil
}
