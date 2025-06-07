package main

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"

)

var client *resty.Client

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
        SetHeader("Content-Type", "application/json")
	return client, nil
}



func insertUser(username string, rating int, client *resty.Client) ([]map[string]any, error) {
	var result []map[string]any

	_, err := client.R(). // _ is the response
		SetBody([]map[string]any{{
			"username" : username,
			"rating" : rating,
		}}).
		SetResult(&result).
		Post("/rest/v1/user")

	if err != nil {
		return nil, fmt.Errorf("error with request/response to supabase")
	}

	fmt.Println("succesful insert")
	return result, nil
}


