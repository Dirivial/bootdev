package main

import (
	"fmt"
	"os"
	"strings"
)

type config struct {
	offset int32
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config) error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Show this help",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name: "map",
			description: "Map locations",
			callback: commandMap,
		},
	}
}

func cleanInput(text string) []string {
	words := make([]string, 0)
	for _, word := range strings.Split(text, " ") {
		if len(word) > 0 {
			word = strings.ToLower(word)
			words = append(words, word)
		}
	}
	return words
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("usage:\n")

	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	locationAreas := api.GetLocationAreas(cfg.offset)

	// Print the results
	for _, location := range locaitonAreas {
		fmt.Println(location.Name)
	}

	cfg.offset += 20
	return nil
}
