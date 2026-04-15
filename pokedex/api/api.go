package main

import (
	"fmt"
	"os"
	"strings"
	"net/http"
	"encoding/json"
)

type LocationAreaResult struct {
	Name: string
	URL: string
}

type locationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocationAreas(offset int) ([]LocationAreaResult, error) {
	query := "https://pokeapi.co/api/v2/location-area?offset=" + fmt.Sprint(offset)

	resp, err := http.Get(query)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode JSON response
	var data locationAreaResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return []struct{}, err
	}

	return data.Results, nil
}
