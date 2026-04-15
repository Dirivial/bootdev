package main

import (
	"fmt"
	"bufio"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := config{}

	for true {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		usrInput := cleanInput(scanner.Text())

		firstWord := ""
		if len(usrInput) > 0 {
			firstWord = usrInput[0]
		}

		cmd, ok := commands[firstWord]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		cmd.callback(&cfg)
	}
}

