package main

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"

)

var client *resty.Client

var userExtension string = "rest/v1/user"

// var gameExtension string = "rest/v1/game"

type User struct {
	Username string `json:"username"`
	Rating int `json:"rating"`
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

func insertUser(username string, rating int, client *resty.Client) error {
	response, err := client.R(). // _ is the response
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

func deleteUser(username string, client *resty.Client) error {
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

func getUser(username string, client *resty.Client) (*User, error) {
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
