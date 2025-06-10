package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"

)

var client *resty.Client

var userExtension string = "rest/v1/user"

var gameExtension string = "rest/v1/game"

type User struct {
	Username string `json:"username"`
	Rating int `json:"rating"`
}

type Game struct {
	Gameid int `json:"gameid"`
	Whiteplayer string `json:"whiteplayer"`
	Blackplayer string `json:"blackplayer"`
	Winner string `json:"winner"`
	Opening string `json:"opening"`
	Gamemoves string `json:"gamemoves"`
	Result string `json:"result"`
}

func createClient() (*resty.Client, error) {
	if client != nil {
		return client, nil
	}
	godotenv.Load()
	apikey := os.Getenv("DBAPIKEY")
	url := os.Getenv("DBURL")

	if apikey == "" || url == "" {
		return nil, fmt.Errorf("apikey or url does not exist")
	}

	client = resty.New().
		SetBaseURL(url).
		SetHeader("apikey", apikey).
		SetHeader("Authorization", "Bearer "+apikey).
		SetHeader("Content-Type", "application/json").
		SetHeader("Prefer", "return=presentation")
	return client, nil
}

func insertUser(username string, rating int) error {
	client, err := createClient()
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	response, err := client.R(). 
		SetBody([]map[string]any{{
			"username": username,
			"rating":   rating,
		}}).
		Post(userExtension)

	if err != nil {
		return fmt.Errorf("error with request/response to supabase")
	}

	if response.IsError() {
		return fmt.Errorf("response gave error code %s", response.Status())
	}

	return nil
}

func deleteUser(username string) error {
	client, err := createClient()
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	response, err := client.R().
		SetQueryParam("username", "eq." + username).
		Delete(userExtension)
	if err != nil {
		return fmt.Errorf("error with delete method")
	}

	if response.IsError() {
		return fmt.Errorf("response gave error code %s", response.Status())
	}

	return nil
}

func getUser(username string) (*User, error) {
	client, err := createClient()
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	var result []User

	response, err := client.R().
		SetQueryParam("username", "eq." + username).
		SetResult(&result).
		Get(userExtension)
	if err != nil {
		return nil, fmt.Errorf("error with get method")
	}

	if response.IsError() {
		return nil, fmt.Errorf("response gave error code %s", response.Status())
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no user found with username: %s", username)
	}

	return &result[0], nil
}

func insertGame(whiteplayer string, blackplayer string, winner string, opening string, gamemoves string, result string) error {
	client, err := createClient()
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	response, err := client.R().
		SetBody([]map[string]any{{
			"whiteplayer": whiteplayer,
			"blackplayer": blackplayer,
			"winner": winner,
			"opening": opening,
			"gamemoves": gamemoves,
			"result": result,
		}}).
		Post(gameExtension)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	
	if response.IsError() {
		return fmt.Errorf("server response had this status: %s", response.Status())
	}

	return nil
}

func insertGames(games []Game) {
	for _, game := range games {
		insertGame(game.Whiteplayer, game.Blackplayer, game.Winner, game.Opening, game.Gamemoves, game.Result)
	}
}

func deleteGame(gameid int) error {
	client, err := createClient()
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	gameidString := strconv.Itoa(gameid)
	response, err := client.R().
		SetQueryParam("gameid", "eq."+gameidString). 
		Delete(gameExtension)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	if response.IsError() {
		return fmt.Errorf("response error: %s", response.Status())
	}

	return nil
}

func getGame(gameid int) (*Game, error) {
	client, err := createClient()
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	var result []Game
	gameidString := strconv.Itoa(gameid)
	response, err := client.R().
		SetQueryParam("gameid", "eq."+gameidString).
		SetResult(&result).
		Get(gameExtension)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	if response.IsError() {
		return nil, fmt.Errorf("response error: %s", response.Status())
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no game exists with id %s", gameidString)
	}

	return &result[0], nil
}

func getGameByPlayers(whiteplayer string, blackplayer string) ([]Game, error) {
	client, err := createClient()
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	var result []Game
	response, err := client.R(). 
		SetQueryParams(map[string]string{
			"whiteplayer": "eq."+whiteplayer,
			"blackplayer": "eq."+blackplayer,
		}).SetResult(&result).
		Get(gameExtension)
	if err != nil {
		return nil, fmt.Errorf("request gave error: %w", err)
	}
	if response.IsError() {
		return nil, fmt.Errorf("response gave error: %s", response.Status())
	}

	if len(result) == 0 {
		return result, fmt.Errorf("no game exists with whiteplayer %s and blackplayer %s", whiteplayer, blackplayer)
	}

	return result, nil
}