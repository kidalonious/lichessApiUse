package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/notnil/chess"
)

type Pgn struct {
	Headers map[string]string
	Moves string
}

func ParsePgnFile(pgnPath string) ([]Pgn, error) {
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

func GetPgns() ([]string, error) {
	var pgnFilepaths []string
	pathBegin := "pgns"
	dirEntries, err := os.ReadDir("pgns")
	if err != nil {
		return nil, fmt.Errorf("read directory returned error: %w", err)
	}

	for _, entry := range dirEntries {
		filepath := filepath.Join(pathBegin, entry.Name())
		pgnFilepaths = append(pgnFilepaths, filepath)
	}

	return pgnFilepaths, nil
}

func PgnToGame(pgn Pgn) Game {
	var game Game
	game.Blackplayer = pgn.Headers["Black"]
	game.Whiteplayer = pgn.Headers["White"]
	game.Opening = pgn.Headers["Opening"]
	game.Result = pgn.Headers["Termination"]
	game.Gamemoves = pgn.Moves
	if pgn.Headers["Result"] == "1-0" {
		game.Winner = pgn.Headers["White"]
	} else if pgn.Headers["Result"] == "0-1" {
		game.Winner = pgn.Headers["Black"]
	} else {
		game.Winner = ""
	}
	return game
}

func PgnsToGames(pgns []Pgn) []Game {
	var games []Game
	for _, pgn := range pgns {
		game := PgnToGame(pgn)
		games = append(games, game)
	}
	return games
}

func PgnToUser(pgn Pgn) (User, User) {
	var whiteUser User
	var blackUser User
	whiteUser.Username = pgn.Headers["White"]
	whiteUser.Rating, _ = strconv.Atoi(pgn.Headers["WhiteElo"])
	blackUser.Username = pgn.Headers["Black"]
	blackUser.Rating, _ = strconv.Atoi(pgn.Headers["BlackElo"])
	return whiteUser, blackUser
}

func PgnsToUsers(pgns []Pgn) []User {
	userMap := make(map[User]struct{})
	for _, pgn := range pgns {
		whiteUser, blackUser := PgnToUser(pgn)
		userMap[whiteUser] = struct{}{}
		userMap[blackUser] = struct{}{}
	}
	users := make([]User, 0, len(userMap))
	for user := range userMap {
		users = append(users, user)
	}
	return users
}

func PgnsToStructs(pgns []Pgn) ([]Game, []User) {
	return PgnsToGames(pgns), PgnsToUsers(pgns)
}